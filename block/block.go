package block

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

type Block struct {
	Index      int64
	Timestamp  int64
	Data       []byte
	Hash       string
	PrevHash   string
	Difficulty int
	Nonce      int64
}

var firstBlock = newBlock(0, []byte("First block"), "", initialDifficulty, int64(0))

var blockchain = []*Block{firstBlock}

func (b *Block) CalculateHash() string {
	str := fmt.Sprintf("%d%d%s%s%d%d", b.Index, b.Timestamp, string(b.Data), b.PrevHash, b.Difficulty, b.Nonce)
	sha := sha256.Sum256([]byte(str))
	return hex.EncodeToString(sha[:])
}

func (b *Block) IsValid(preBlock *Block) bool {
	if b.Index != preBlock.Index+1 {
		return false
	}
	if b.PrevHash != preBlock.Hash {
		return false
	}

	return true
}

func ValidateBlockChain(blockchain []*Block) bool {
	isValidGenesisBlock := func(b *Block) bool {
		if b.Index != firstBlock.Index {
			return false
		}
		if !bytes.Equal(b.Data, firstBlock.Data) {
			return false
		}

		if (b.PrevHash != firstBlock.PrevHash) || (b.Hash != firstBlock.Hash) {
			return false
		}

		return true
	}
	if !isValidGenesisBlock(blockchain[0]) {
		return false
	}

	for i := 1; i < len(blockchain); i++ {
		if !blockchain[i].IsValid(blockchain[i-1]) {
			return false
		}
	}
	return true
}

func ReplaceChain(newChain []*Block) {
	if (len(newChain) > len(blockchain)) && ValidateBlockChain(newChain) {
		blockchain = newChain
		// Todo: publish to other nodes
		// broadcastLatest()
	}
}

func GetLatestBlock() *Block {
	return blockchain[len(blockchain)-1]
}

func newBlock(index int64, data []byte, prevHash string, difficulty int, nonce int64) *Block {
	b := &Block{
		Index:      index,
		Data:       data,
		PrevHash:   prevHash,
		Timestamp:  time.Now().Unix(),
		Difficulty: difficulty,
		Nonce:      nonce,
	}
	b.Hash = b.CalculateHash()
	return b
}

func GetBlockChain() []*Block {
	return blockchain
}
