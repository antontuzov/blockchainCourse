package week1

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"time"

	"blockchain-course/module1/transaction"
)

// Block represents a block in the blockchain
type Block struct {
	Index         int64
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
	Transactions  []*transaction.Transaction
}

// NewBlock creates and returns a new Block
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{
		Index:         0,
		Timestamp:     time.Now().Unix(),
		Data:          []byte(data),
		PrevBlockHash: prevBlockHash,
		Hash:          []byte{},
		Nonce:         0,
		Transactions:  make([]*transaction.Transaction, 0),
	}
	block.SetHash()
	return block
}

// SetHash calculates and sets the hash of the block
func (b *Block) SetHash() {
	headers := [][]byte{
		b.PrevBlockHash,
		b.Data,
		IntToHex(b.Timestamp),
		IntToHex(int64(b.Nonce)),
	}
	header := bytes.Join(headers, []byte{})
	hash := sha256.Sum256(header)
	b.Hash = hash[:]
}

// IntToHex converts an int64 to a byte array
func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	binary.Write(buff, binary.BigEndian, num)
	return buff.Bytes()
}

// GenesisBlock creates and returns the genesis block
func GenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}
