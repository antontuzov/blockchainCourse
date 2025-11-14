package week2

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
	"os"
	"strings"

	"blockchain-course/module1/week1"
)

const targetBits = 8 // Reduced difficulty for testing - controls mining difficulty

// Blockchain represents the blockchain structure
type Blockchain struct {
	Blocks []*week1.Block
}

// ProofOfWork represents the proof of work structure
type ProofOfWork struct {
	block  *week1.Block
	target *big.Int
}

// NewProofOfWork creates a new ProofOfWork
func NewProofOfWork(block *week1.Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	pow := &ProofOfWork{block: block, target: target}
	return pow
}

// NewBlockchain creates a new Blockchain with genesis block
func NewBlockchain() *Blockchain {
	return &Blockchain{Blocks: []*week1.Block{week1.GenesisBlock()}}
}

// AddBlock adds a new block to the blockchain
func (bc *Blockchain) AddBlock(data string) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := week1.NewBlock(data, prevBlock.Hash)

	// Mine the block
	pow := NewProofOfWork(newBlock)
	nonce, hash := pow.Run()

	newBlock.Hash = hash[:]
	newBlock.Nonce = nonce

	bc.Blocks = append(bc.Blocks, newBlock)
}

// IsValid checks if the blockchain is valid
func (bc *Blockchain) IsValid() bool {
	for i := 1; i < len(bc.Blocks); i++ {
		currentBlock := bc.Blocks[i]
		prevBlock := bc.Blocks[i-1]

		// Check if the current block's hash is valid with proof of work
		pow := NewProofOfWork(currentBlock)
		if !pow.Validate() {
			return false
		}

		// Check if the previous block hash matches
		if !bytes.Equal(currentBlock.PrevBlockHash, prevBlock.Hash) {
			return false
		}
	}
	return true
}

// Run performs the proof of work algorithm
func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	var nonce int = 0

	// Only print mining progress if not in test mode
	if !isTesting() {
		fmt.Printf("Mining a new block...")
	}
	for nonce < math.MaxInt64 {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		// Only print progress if not in test mode
		if !isTesting() {
			fmt.Printf("\r%x", hash)
		}
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}
	// Only print newline if not in test mode
	if !isTesting() {
		fmt.Print("\n\n")
	}

	return nonce, hash[:]
}

// prepareData prepares the data for hashing in POW
func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.Data,
			week1.IntToHex(pow.block.Timestamp),
			week1.IntToHex(int64(targetBits)),
			week1.IntToHex(int64(nonce)),
		},
		[]byte{},
	)

	return data
}

// Validate validates the proof of work
func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(pow.target) == -1

	return isValid
}

// BlockchainIterator represents an iterator for the blockchain
type BlockchainIterator struct {
	current int
	bc      *Blockchain
}

// Iterator returns a blockchain iterator
func (bc *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{
		current: len(bc.Blocks) - 1,
		bc:      bc,
	}
}

// Next returns the next block in the blockchain
func (i *BlockchainIterator) Next() *week1.Block {
	if i.current < 0 {
		return nil
	}

	block := i.bc.Blocks[i.current]
	i.current--

	return block
}

// isTesting checks if we're running in test mode
func isTesting() bool {
	// Check if we're running tests by looking at the binary name
	// This is a simple way to detect test mode
	return strings.HasSuffix(os.Args[0], ".test") || strings.Contains(os.Args[0], "/_test/")
}
