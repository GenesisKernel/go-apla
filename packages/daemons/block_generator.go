// Copyright (C) 2017, 2018, 2019 EGAAS S.A.
//
// This program is free software; you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation; either version 2 of the License, or (at
// your option) any later version.
//
// This program is distributed in the hope that it will be useful, but
// WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301, USA.

package daemons

import (
	"bytes"
	"context"
	"time"

	"github.com/AplaProject/go-apla/packages/block"
	"github.com/AplaProject/go-apla/packages/conf"
	"github.com/AplaProject/go-apla/packages/conf/syspar"
	"github.com/AplaProject/go-apla/packages/consts"
	"github.com/AplaProject/go-apla/packages/model"
	"github.com/AplaProject/go-apla/packages/notificator"
	"github.com/AplaProject/go-apla/packages/protocols"
	"github.com/AplaProject/go-apla/packages/service"
	"github.com/AplaProject/go-apla/packages/transaction"
	"github.com/AplaProject/go-apla/packages/utils"

	log "github.com/sirupsen/logrus"
)

// BlockGenerator is daemon that generates blocks
func BlockGenerator(ctx context.Context, d *daemon) error {
	d.sleepTime = time.Second
	if service.IsNodePaused() {
		return nil
	}

	nodePosition, err := syspar.GetNodePositionByKeyID(conf.Config.KeyID)
	if err != nil {
		// we are not full node and can't generate new blocks
		d.sleepTime = 10 * time.Second
		d.logger.WithFields(log.Fields{"type": consts.JustWaiting, "error": err}).Debug("we are not full node, sleep for 10 seconds")
		return nil
	}

	DBLock()
	defer DBUnlock()

	// wee need fresh myNodePosition after locking
	nodePosition, err = syspar.GetNodePositionByKeyID(conf.Config.KeyID)
	if err != nil {
		d.logger.WithFields(log.Fields{"type": consts.DBError, "error": err}).Error("getting node position by key id")
		return err
	}

	btc := protocols.NewBlockTimeCounter()
	at := time.Now()

	if exists, err := btc.BlockForTimeExists(at, int(nodePosition)); exists || err != nil {
		return nil
	}

	timeToGenerate, err := btc.TimeToGenerate(at, int(nodePosition))
	if err != nil {
		d.logger.WithFields(log.Fields{"type": consts.BlockError, "error": err, "position": nodePosition}).Debug("calculating block time")
		return err
	}

	if !timeToGenerate {
		d.logger.WithFields(log.Fields{"type": consts.JustWaiting}).Debug("not my generation time")
		return nil
	}

	_, endTime, err := btc.RangeByTime(time.Now())
	if err != nil {
		log.WithFields(log.Fields{"type": consts.TimeCalcError, "error": err}).Error("on getting end time of generation")
	}

	done := time.After(endTime.Sub(time.Now()))
	prevBlock := &model.InfoBlock{}
	_, err = prevBlock.Get()
	if err != nil {
		d.logger.WithFields(log.Fields{"type": consts.DBError, "error": err}).Error("getting previous block")
		return err
	}

	NodePrivateKey, NodePublicKey, err := utils.GetNodeKeys()
	if err != nil || len(NodePrivateKey) < 1 {
		if err == nil {
			d.logger.WithFields(log.Fields{"type": consts.EmptyObject}).Error("node private key is empty")
		}
		return err
	}

	dtx := DelayedTx{
		privateKey: NodePrivateKey,
		publicKey:  NodePublicKey,
		logger:     d.logger,
	}

	dtx.RunForBlockID(prevBlock.BlockID + 1)

	trs, err := processTransactions(d.logger, done)
	if err != nil {
		return err
	}

	// Block generation will be started only if we have transactions
	if len(trs) == 0 {
		return nil
	}

	header := &utils.BlockData{
		BlockID:      prevBlock.BlockID + 1,
		Time:         time.Now().Unix(),
		EcosystemID:  0,
		KeyID:        conf.Config.KeyID,
		NodePosition: nodePosition,
		Version:      consts.BLOCK_VERSION,
	}

	pb := &utils.BlockData{
		BlockID:       prevBlock.BlockID,
		Hash:          prevBlock.Hash,
		RollbacksHash: prevBlock.RollbacksHash,
	}

	blockBin, err := generateNextBlock(header, trs, NodePrivateKey, pb)
	if err != nil {
		return err
	}

	err = block.InsertBlockWOForks(blockBin, true, false)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("on inserting new block")
		return err
	}
	log.WithFields(log.Fields{"Block": header.String(), "type": consts.SyncProcess}).Debug("Generated block ID")

	go notificator.CheckTokenMovementLimits(nil, conf.Config.TokenMovement, header.BlockID)
	return nil
}

func generateNextBlock(blockHeader *utils.BlockData, trs []*model.Transaction, key string, prevBlock *utils.BlockData) ([]byte, error) {
	trData := make([][]byte, 0, len(trs))
	for _, tr := range trs {
		trData = append(trData, tr.Data)
	}

	return block.MarshallBlock(blockHeader, trData, prevBlock, key)
}

func processTransactions(logger *log.Entry, done <-chan time.Time) ([]*model.Transaction, error) {
	p := new(transaction.Transaction)

	// verify transactions
	err := transaction.ProcessTransactionsQueue(p.DbTransaction)
	if err != nil {
		return nil, err
	}

	trs, err := model.GetAllUnusedTransactions(syspar.GetMaxTxCount())
	if err != nil {
		logger.WithFields(log.Fields{"type": consts.DBError, "error": err}).Error("getting all unused transactions")
		return nil, err
	}

	limits := block.NewLimits(nil)

	type badTxStruct struct {
		hash  []byte
		msg   string
		keyID int64
	}

	processBadTx := func(dbTx *model.DbTransaction) chan badTxStruct {
		ch := make(chan badTxStruct)

		go func() {
			for badTxItem := range ch {
				block.BadTxForBan(badTxItem.keyID)
				transaction.MarkTransactionBad(p.DbTransaction, badTxItem.hash, badTxItem.msg)
			}
		}()

		return ch
	}

	processIncAttemptCnt := func() chan []byte {
		ch := make(chan []byte)
		go func() {
			for tx := range ch {
				model.IncrementTxAttemptCount(tx)
			}
		}()

		return ch
	}

	txBadChan := processBadTx(p.DbTransaction)
	attemptCountChan := processIncAttemptCnt()

	defer func() {
		close(txBadChan)
		close(attemptCountChan)
	}()

	// Checks preprocessing count limits
	txList := make([]*model.Transaction, 0, len(trs))
	for i, txItem := range trs {
		select {
		case <-done:
			return txList, err
		default:
			bufTransaction := bytes.NewBuffer(txItem.Data)
			p, err := transaction.UnmarshallTransaction(bufTransaction, true)
			if err != nil {
				if p != nil {
					txBadChan <- badTxStruct{hash: p.TxHash, msg: err.Error(), keyID: p.TxHeader.KeyID}
				}
				continue
			}

			if err := p.Check(time.Now().Unix(), false); err != nil {
				txBadChan <- badTxStruct{hash: p.TxHash, msg: err.Error(), keyID: p.TxHeader.KeyID}
				continue
			}

			if p.TxSmart != nil {
				err = limits.CheckLimit(p)
				if err == block.ErrLimitStop && i > 0 {
					attemptCountChan <- p.TxHash
					break
				} else if err != nil {
					if err == block.ErrLimitSkip {
						attemptCountChan <- p.TxHash
					} else {
						txBadChan <- badTxStruct{hash: p.TxHash, msg: err.Error(), keyID: p.TxHeader.KeyID}
					}
					continue
				}
			}
			txList = append(txList, trs[i])
		}
	}
	return txList, nil
}
