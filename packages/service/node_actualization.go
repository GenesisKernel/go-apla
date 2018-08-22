package service

import (
	"context"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/GenesisKernel/go-genesis/packages/blockchain"
	"github.com/GenesisKernel/go-genesis/packages/conf"
	"github.com/GenesisKernel/go-genesis/packages/conf/syspar"
	"github.com/GenesisKernel/go-genesis/packages/consts"
	"github.com/GenesisKernel/go-genesis/packages/utils"
)

// DefaultBlockchainGap is default value for the number of lagging blocks
const DefaultBlockchainGap int64 = 10

type NodeActualizer struct {
	availableBlockchainGap int64
}

func NewNodeActualizer(availableBlockchainGap int64) NodeActualizer {
	return NodeActualizer{
		availableBlockchainGap: availableBlockchainGap,
	}
}

// Run is starting node monitoring
func (n *NodeActualizer) Run() {
	go func() {
		log.Info("Node Actualizer monitoring starting")
		for {
			actual, err := n.checkBlockchainActuality()
			if err != nil {
				log.WithFields(log.Fields{"type": consts.BCActualizationError, "err": err}).Error("checking blockchain actuality")
				return
			}

			if !actual && !IsNodePaused() {
				log.Info("Node Actualizer is pausing node activity")
				n.pauseNodeActivity()
			}

			if actual && IsNodePaused() {
				log.Info("Node Actualizer is resuming node activity")
				n.resumeNodeActivity()
			}

			time.Sleep(time.Minute * 5)
		}
	}()
}

func (n *NodeActualizer) checkBlockchainActuality() (bool, error) {
	block, found, err := blockchain.GetLastBlock()
	if err != nil {
		return false, err
	}

	remoteHosts := syspar.GetRemoteHosts()
	ctx, _ := context.WithCancel(context.Background())
	_, maxBlockID, err := utils.ChooseBestHost(ctx, remoteHosts, &log.Entry{Logger: &log.Logger{}})
	if err != nil {
		return false, errors.Wrapf(err, "choosing best host")
	}
	var curBlockID int64
	if found {
		curBlockID = block.Header.BlockID
	}

	// Currently this node is downloading blockchain
	if curBlockID == 0 || curBlockID+n.availableBlockchainGap < maxBlockID {
		return false, nil
	}

	foreignBlock, found, err := blockchain.GetMaxForeignBlock(conf.Config.KeyID)
	if err != nil {
		return false, errors.Wrapf(err, "retrieving last foreign block")
	}

	if !found {
		return false, nil
	}

	// Node did not accept any blocks for an hour
	t := time.Unix(foreignBlock.Header.Time, 0)
	if time.Since(t).Minutes() > 30 && len(remoteHosts) > 1 {
		return false, nil
	}

	return true, nil
}

func (n *NodeActualizer) pauseNodeActivity() {
	np.Set(PauseTypeUpdatingBlockchain)
}

func (n *NodeActualizer) resumeNodeActivity() {
	np.Unset()
}
