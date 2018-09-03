package registry

import (
	"encoding/json"
	"fmt"

	"sync"

	"github.com/GenesisKernel/go-genesis/packages/consts"
	"github.com/GenesisKernel/go-genesis/packages/storage/kv"
	"github.com/GenesisKernel/go-genesis/packages/types"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/match"
	"github.com/yddmat/memdb"
)

const keyConvention = "%s.%s.%s"

var (
	ErrUnknownContext   = errors.New("unknown writing operation context (block o/or hash empty)")
	ErrWrongRegistry    = errors.New("wrong registry")
	ErrRollbackDisabled = errors.New("rollback is disabled")
)

// metadataTx must be closed by calling Commit() or Rollback() when done
type metadataTx struct {
	db      kv.Database
	tx      kv.Transaction
	durable bool

	rollback *metadataRollback
	indexer  *indexer

	currentBlockHash []byte
	currentTxHash    []byte
	stateMu          sync.RWMutex
}

func (m *metadataTx) Insert(registry *types.Registry, pkValue string, value interface{}) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return errors.Wrapf(err, "marshalling struct to json")
	}

	key, err := m.formatKey(registry, pkValue)
	if err != nil {
		return err
	}

	err = m.tx.Set(key, string(jsonValue))
	if err != nil {
		return errors.Wrapf(err, "inserting value %s to %s registry", value, registry.Name)
	}

	if m.rollback != nil {
		m.stateMu.RLock()
		block := m.currentBlockHash
		tx := m.currentTxHash
		m.stateMu.RUnlock()

		if len(block) == 0 || len(tx) == 0 {
			return ErrUnknownContext
		}

		err = m.rollback.saveState(block, tx, registry, pkValue, "")
		if err != nil {
			return errors.Wrapf(err, "saving rollback info")
		}
	}

	return nil
}

func (m *metadataTx) Update(registry *types.Registry, pkValue string, newValue interface{}) error {
	jsonValue, err := json.Marshal(newValue)
	if err != nil {
		return errors.Wrapf(err, "marshalling struct to json")
	}

	key, err := m.formatKey(registry, pkValue)
	if err != nil {
		return err
	}

	old, err := m.tx.Update(key, string(jsonValue))
	if err != nil {
		return errors.Wrapf(err, "inserting value %s to %s registry", pkValue, registry.Name)
	}

	m.stateMu.RLock()
	block := m.currentBlockHash
	tx := m.currentTxHash
	m.stateMu.RUnlock()

	if len(block) == 0 || len(tx) == 0 {
		return ErrUnknownContext
	}

	err = m.rollback.saveState(block, tx, registry, pkValue, old)
	if err != nil {
		return errors.Wrapf(err, "saving rollback info")
	}

	return nil
}

func (m *metadataTx) Get(registry *types.Registry, pkValue string, out interface{}) error {
	err := m.refreshTx()
	if err != nil {
		return err
	}
	defer m.endRead()

	key, err := m.formatKey(registry, pkValue)
	if err != nil {
		return err
	}

	value, err := m.tx.Get(key)
	if err != nil {
		return errors.Wrapf(err, "retrieving %s from databse", key)
	}

	err = json.Unmarshal([]byte(value), out)
	if err != nil {
		return errors.Wrapf(err, "unmarshalling value %s to struct", value)
	}

	return nil
}

func (m *metadataTx) Walk(registry *types.Registry, field string, fn func(value string) bool) error {
	err := m.refreshTx()
	if err != nil {
		return err
	}
	defer m.endRead()

	prefix := fmt.Sprintf("%s.*", registry.Name)

	return m.tx.Ascend(m.indexer.formatIndexName(registry, field), func(key, value string) bool {
		if match.Match(key, prefix) {
			return fn(value)
		}

		return true
	})
}

func (m *metadataTx) Rollback() error {
	tx := m.tx
	m.tx = nil
	return tx.Rollback()
}

func (m *metadataTx) Commit() error {
	err := m.tx.Commit()
	m.tx = nil
	return err
}

func (m *metadataTx) SetTxHash(txHash []byte) {
	m.stateMu.Lock()
	m.currentTxHash = txHash
	m.stateMu.Unlock()
}

func (m *metadataTx) SetBlockHash(blockHash []byte) {
	m.stateMu.Lock()
	m.currentBlockHash = blockHash
	m.stateMu.Unlock()
}

func (m *metadataTx) AddIndex(indexes ...types.Index) error {
	return m.addIndex(false, indexes...)
}

func (m *metadataTx) addIndex(init bool, indexes ...types.Index) error {
	if init {
		if err := m.indexer.init(indexes); err != nil {
			return err
		}
	}

	return m.indexer.addIndexes(indexes...)
}

func (m *metadataTx) refreshTx() error {
	if m.durable {
		if m.tx == nil {
			return memdb.ErrTxClosed
		}

		return nil
	}

	// Non-durable transaction can only be readable. All writable transaction called directly from Begin()
	// and must be committed/rollback manually. So here we create readonly tx without providing any choice
	m.tx = m.db.Begin(false)

	return nil
}

func (m *metadataTx) formatKey(reg *types.Registry, pk string) (string, error) {
	if reg.Name == "ecosystem" {
		return fmt.Sprintf("%s.%s", reg.Name, pk), nil
	}

	if reg.Ecosystem == nil {
		return "", ErrWrongRegistry
	}

	return fmt.Sprintf(keyConvention, reg.Name, reg.Ecosystem.Name, pk), nil
}

func (m *metadataTx) endRead() error {
	if !m.durable {
		if err := m.tx.Commit(); err != nil {
			return errors.Wrapf(err, "ending read transaction")
		}

		m.tx = nil
	}

	return nil
}

type metadataStorage struct {
	db       kv.Database
	rollback bool
}

func NewMetadataStorage(db kv.Database, indexes []types.Index, rollback bool) (types.MetadataRegistryStorage, error) {
	ms := &metadataStorage{
		db:       db,
		rollback: rollback,
	}

	mtx := ms.Begin()
	tx := mtx.(*metadataTx)
	if err := tx.addIndex(true, indexes...); err != nil {
		return nil, err
	}
	mtx.Commit()

	return ms, nil
}

func (m *metadataStorage) Begin() types.MetadataRegistryReaderWriter {
	databaseTx := m.db.Begin(true)
	tx := &metadataTx{tx: databaseTx, durable: true, indexer: &indexer{tx: databaseTx}}

	if m.rollback {
		tx.rollback = &metadataRollback{tx: databaseTx, txCounter: make(map[string]uint64)}
	}

	return tx
}

func (m *metadataStorage) Rollback(block []byte) error {
	if !m.rollback {
		return ErrRollbackDisabled
	}

	databaseTx := m.db.Begin(true)
	rollback := &metadataRollback{tx: databaseTx, txCounter: make(map[string]uint64)}

	err := rollback.rollbackState(block)
	if err != nil {
		rbErr := databaseTx.Rollback()
		log.WithFields(log.Fields{"type": consts.DBError, "error": rbErr}).Error("rollback metadata db")
		return err
	}

	err = databaseTx.Commit()
	if err != nil {
		log.WithFields(log.Fields{"type": consts.DBError, "error": err}).Error("commiting metadata db")
		return err
	}

	return nil
}

func (m *metadataStorage) Reader() types.MetadataRegistryReader {
	return &metadataTx{db: m.db}
}
