package query

import (
	"fmt"
	"sync"
)

const maxBlockIDEndpoint = "/api/v2/maxblockid"
const blockInfoEndpoint = "/api/v2/block/%s"

type MaxBlockID struct {
	MaxBlockID int64  `json:"max_block_id"`
	Hash       string `json:"hash"`
}

type blockInfoResult struct {
	Hash          []byte `json:"hash"`
	EcosystemID   int64  `json:"ecosystem_id"`
	KeyID         int64  `json:"key_id"`
	Time          int64  `json:"time"`
	Tx            int32  `json:"tx_count"`
	RollbacksHash []byte `json:"rollbacks_hash"`
}

func MaxBlockIDs(nodesList []string) ([]*MaxBlockID, error) {
	wg := sync.WaitGroup{}
	workResults := ConcurrentMap{m: map[string]interface{}{}}
	for _, nodeUrl := range nodesList {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			maxBlockID := &MaxBlockID{}
			if err := sendGetRequest(url+maxBlockIDEndpoint, maxBlockID); err != nil {
				workResults.Set(url, err)
				return
			}
			workResults.Set(url, maxBlockID)
		}(nodeUrl)
	}
	wg.Wait()
	maxBlockIds := []*MaxBlockID{}
	for _, result := range workResults.m {
		switch res := result.(type) {
		case *MaxBlockID:
			maxBlockIds = append(maxBlockIds, res)
		case error:
			return nil, res
		}
	}
	return maxBlockIds, nil
}

func BlockInfo(nodesList []string, blockHash string) (map[string]*blockInfoResult, error) {
	wg := sync.WaitGroup{}
	workResults := ConcurrentMap{m: map[string]interface{}{}}
	for _, nodeUrl := range nodesList {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			blockInfo := &blockInfoResult{}
			if err := sendGetRequest(url+fmt.Sprintf(blockInfoEndpoint, blockHash), blockInfo); err != nil {
				workResults.Set(url, err)
				return
			}
			workResults.Set(url, blockInfo)
		}(nodeUrl)
	}
	wg.Wait()
	result := map[string]*blockInfoResult{}
	for nodeUrl, blockInfoOrError := range workResults.m {
		switch res := blockInfoOrError.(type) {
		case error:
			return nil, res
		case *blockInfoResult:
			result[nodeUrl] = res
		}
	}
	return result, nil
}
