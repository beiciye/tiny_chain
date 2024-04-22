package block

import (
	"log"
	"strings"
)

var initialDifficulty int = 4

const BLOCK_GENERATION_INTERVAL = 10

const DIFFICULTY_ADJUSTMENT_INTERVAL = 10

func getDifficulty() int {
	lastBlock := GetLatestBlock()
	if lastBlock.Index%DIFFICULTY_ADJUSTMENT_INTERVAL == 0 {
		return getAdjustedDifficulty()
	} else {
		return lastBlock.Difficulty
	}
}

func getAdjustedDifficulty() int {
	if len(blockchain) <= DIFFICULTY_ADJUSTMENT_INTERVAL {
		return initialDifficulty
	}

	lastBlock := GetLatestBlock()
	prevAdjustmentBlock := blockchain[len(blockchain)-DIFFICULTY_ADJUSTMENT_INTERVAL]
	timeExpected := int64(BLOCK_GENERATION_INTERVAL * DIFFICULTY_ADJUSTMENT_INTERVAL)
	timeTaken := lastBlock.Timestamp - prevAdjustmentBlock.Timestamp

	if timeTaken < timeExpected/2 {
		return prevAdjustmentBlock.Difficulty + 1
	} else if timeTaken > timeExpected*2 {
		return prevAdjustmentBlock.Difficulty - 1
	} else {
		return prevAdjustmentBlock.Difficulty
	}
}

func MineBlock() *Block {
	latestBlock := GetLatestBlock()

	var nextBlock *Block = newBlock(latestBlock.Index+1, []byte("New block"), latestBlock.Hash, getDifficulty(), int64(0))
	for {
		hash := nextBlock.CalculateHash()
		if matchProofOfWork(hash) {
			log.Default().Printf("Found hash: %s Nonce: %d", hash, nextBlock.Nonce)
			blockchain = append(blockchain, nextBlock)
			break
		}
		nextBlock.Nonce = nextBlock.Nonce + 1
	}
	return nextBlock
}

func matchProofOfWork(hash string) bool {
	return hash[:getDifficulty()] == strings.Repeat("0", getDifficulty())
}
