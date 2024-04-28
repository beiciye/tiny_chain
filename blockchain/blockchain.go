package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

type Block struct {
	Index      int
	Timestamp  int64
	Data       []byte
	Hash       string
	PrevHash   string
	Difficulty int
	Nonce      int64
}

var firstBlock = &Block{
	Index:      0,
	Data:       []byte("First block"),
	PrevHash:   "",
	Timestamp:  1714154017,
	Difficulty: initialDifficulty,
	Hash:       "4653dbe9183c6ed83761da0dac13410b154c4146217d7a2a6708e02868b4fe2d",
	Nonce:      0,
}

var chains = []*Block{firstBlock}

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

func ReplaceChain(newBlocks []*Block) {
	if (len(newBlocks) > len(chains)) && ValidateBlockChain(newBlocks) {
		fmt.Println("Replacing blockchain")
		chains = newBlocks
	}
}

func blockchainAddBlock(nextBlock *Block) {
	if nextBlock.IsValid(GetLatestBlock()) {
		chains = append(chains, nextBlock)
	}

}

func GetLatestBlock() *Block {
	return chains[len(chains)-1]
}

func newBlock(index int, data []byte, prevHash string, difficulty int, nonce int64) *Block {
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
	return chains
}
