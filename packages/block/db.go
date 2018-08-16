package block

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/GenesisKernel/go-genesis/packages/blockchain"
	"github.com/GenesisKernel/go-genesis/packages/consts"
	"github.com/GenesisKernel/go-genesis/packages/converter"
	"github.com/GenesisKernel/go-genesis/packages/crypto"
	"github.com/GenesisKernel/go-genesis/packages/model"
	"github.com/GenesisKernel/go-genesis/packages/utils"

	log "github.com/sirupsen/logrus"
	msgpack "gopkg.in/vmihailenco/msgpack.v2"
)

// UpdBlockInfo updates info_block table
func UpdBlockInfo(dbTransaction *model.DbTransaction, block *PlayableBlock) error {
	blockID := block.Header.BlockID
	// for the local tests
	forSha := fmt.Sprintf("%d,%x,%s,%d,%d,%d,%d", blockID, block.PrevHeader.Hash, block.MrklRoot,
		block.Header.Time, block.Header.EcosystemID, block.Header.KeyID, block.Header.NodePosition)

	hash, err := crypto.DoubleHash([]byte(forSha))
	if err != nil {
		log.WithFields(log.Fields{"type": consts.CryptoError, "error": err}).Fatal("double hashing block")
	}

	block.Header.Hash = hash
	if block.Header.BlockID == 1 {
		ib := &model.InfoBlock{
			Hash:           hash,
			BlockID:        blockID,
			Time:           block.Header.Time,
			EcosystemID:    block.Header.EcosystemID,
			KeyID:          block.Header.KeyID,
			NodePosition:   converter.Int64ToStr(block.Header.NodePosition),
			CurrentVersion: fmt.Sprintf("%d", block.Header.Version),
		}
		err := ib.Create(dbTransaction)
		if err != nil {
			log.WithFields(log.Fields{"type": consts.DBError, "error": err}).Error("creating info block")
			return fmt.Errorf("error insert into info_block %s", err)
		}
	} else {
		ibUpdate := &model.InfoBlock{
			Hash:         hash,
			BlockID:      blockID,
			Time:         block.Header.Time,
			EcosystemID:  block.Header.EcosystemID,
			KeyID:        block.Header.KeyID,
			NodePosition: converter.Int64ToStr(block.Header.NodePosition),
			Sent:         0,
		}
		if err := ibUpdate.Update(dbTransaction); err != nil {
			log.WithFields(log.Fields{"type": consts.DBError, "error": err}).Error("creating info block")
			return fmt.Errorf("error while updating info_block: %s", err)
		}
	}

	return nil
}

// InsertIntoBlockchain inserts a block into the blockchain
func InsertIntoBlockchain(transaction *model.DbTransaction, block *PlayableBlock) error {
	// for local tests
	blockID := block.Header.BlockID

	// record into the block chain
	bl := &model.Block{}
	err := bl.DeleteById(transaction, blockID)
	if err != nil {
		log.WithFields(log.Fields{"type": consts.DBError, "error": err}).Error("deleting block by id")
		return err
	}
	rollbackTx := &model.RollbackTx{}
	blockRollbackTxs, err := rollbackTx.GetBlockRollbackTransactions(transaction, blockID)
	if err != nil {
		log.WithFields(log.Fields{"type": consts.DBError, "error": err}).Error("getting block rollback txs")
		return err
	}
	buffer := bytes.Buffer{}
	for _, rollbackTx := range blockRollbackTxs {
		rollbackTxBytes, err := json.Marshal(rollbackTx)
		if err != nil {
			log.WithFields(log.Fields{"type": consts.JSONMarshallError, "error": err}).Error("marshalling rollback_tx to json")
			return err
		}

		buffer.Write(rollbackTxBytes)
	}
	rollbackTxsHash, err := crypto.Hash(buffer.Bytes())
	if err != nil {
		log.WithFields(log.Fields{"type": consts.CryptoError, "error": err}).Error("hashing block rollback_txs")
		return err
	}
	b := &model.Block{
		ID:            blockID,
		Hash:          block.Header.Hash,
		Data:          block.BinData,
		EcosystemID:   block.Header.EcosystemID,
		KeyID:         block.Header.KeyID,
		NodePosition:  block.Header.NodePosition,
		Time:          block.Header.Time,
		RollbacksHash: rollbackTxsHash,
		Tx:            int32(len(block.Transactions)),
	}
	blockTimeCalculator, err := utils.BuildBlockTimeCalculator(nil)
	if err != nil {
		log.WithFields(log.Fields{"type": consts.DBError, "error": err}).Error("creating block")
		return err
	}
	validBlockTime := true
	if blockID > 1 {
		validBlockTime, err = blockTimeCalculator.ValidateBlock(b.NodePosition, time.Unix(b.Time, 0))
		if err != nil {
			log.WithFields(log.Fields{"type": consts.BlockError, "error": err}).Error("block validation")
			return err
		}
	}
	if validBlockTime {
		err = b.Create(transaction)
		if err != nil {
			log.WithFields(log.Fields{"type": consts.DBError, "error": err}).Error("creating block")
			return err
		}
	} else {
		err := fmt.Errorf("Invalid block time: %d", block.Header.Time)
		log.WithFields(log.Fields{"type": consts.BlockError, "error": err}).Error("invalid block time")
		return err
	}

	return nil
}

// GetBlockDataFromBlockChain is retrieving block data from blockchain
func GetBlockDataFromBlockChain(hash []byte) (*utils.BlockData, error) {
	block, found, err := blockchain.GetBlock(hash)
	if err != nil {
		log.WithFields(log.Fields{"type": consts.DBError, "error": err}).Error("Getting block by hash")
		return nil, utils.ErrInfo(err)
	}
	if !found {
		return nil, nil
	}
	return block.Header, nil
}

// GetDataFromFirstBlock returns data of first block
func GetDataFromFirstBlock() (data *consts.FirstBlock, ok bool) {
	bBlock, found, err := blockchain.GetFirstBlock()
	if !found {
		return
	}
	if err != nil {
		log.WithFields(log.Fields{"type": consts.DBError, "error": err}).Error("getting record of first block")
		return
	}
	if len(bBlock.Transactions) == 0 {
		log.WithFields(log.Fields{"type": consts.ParserError}).Error("list of parsers is empty")
		return
	}

	t := bBlock.Transactions[0]
	fb := &consts.FirstBlock{}
	if err := msgpack.Unmarshal(t, fb); err != nil {
		log.WithFields(log.Fields{"type": consts.UnmarshallingError, "error": err}).Error("getting data of first block")
		return
	} else {
		data = fb
	}

	return
}
