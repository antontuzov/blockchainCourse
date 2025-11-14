package week6

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
)

// Consensus defines the interface for consensus algorithms
type Consensus interface {
	Start() error
	ValidateBlock(block *Block) bool
	ProposeBlock(block *Block) error
}

// Block represents a block in the blockchain
type Block struct {
	Index     int64
	Timestamp int64
	Data      []byte
	PrevHash  []byte
	Hash      []byte
	Nonce     int64
	Validator string // For PoS
}

// PoW represents Proof of Work consensus
type PoW struct {
	Blockchain *Blockchain
	Difficulty int
}

// PoS represents Proof of Stake consensus
type PoS struct {
	Blockchain *Blockchain
	Validators map[string]int64 // address -> stake
}

// PBFT represents Practical Byzantine Fault Tolerance consensus
type PBFT struct {
	Blockchain  *Blockchain
	Nodes       []string
	PrimaryNode string
	ViewID      int64
	SequenceID  int64
}

// NewPoW creates a new Proof of Work consensus
func NewPoW(blockchain *Blockchain, difficulty int) *PoW {
	return &PoW{
		Blockchain: blockchain,
		Difficulty: difficulty,
	}
}

// Start starts the PoW consensus
func (pow *PoW) Start() error {
	fmt.Println("Starting Proof of Work consensus")
	return nil
}

// ValidateBlock validates a block using PoW
func (pow *PoW) ValidateBlock(block *Block) bool {
	// Check if the hash meets the difficulty requirement
	// For difficulty D, target = max_target / D
	// max_target for 256-bit hash is 2^224 (most significant 28 bytes can be non-zero)
	maxTarget := new(big.Int)
	maxTarget.Exp(big.NewInt(2), big.NewInt(224), nil) // 2^224

	target := new(big.Int)
	target.Div(maxTarget, big.NewInt(int64(pow.Difficulty)))

	hashInt := new(big.Int)
	hashInt.SetBytes(block.Hash)

	return hashInt.Cmp(target) == -1
}

// ProposeBlock proposes a new block using PoW
func (pow *PoW) ProposeBlock(block *Block) error {
	fmt.Println("Proposing block with Proof of Work")

	// Perform mining
	nonce, hash := pow.MineBlock(block)
	block.Nonce = nonce
	block.Hash = hash

	return nil
}

// MineBlock performs the mining process
func (pow *PoW) MineBlock(block *Block) (int64, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := int64(0)

	target := big.NewInt(1)
	target.Lsh(target, uint(256-pow.Difficulty))

	fmt.Printf("Mining block with difficulty %d...\n", pow.Difficulty)

	for nonce < math.MaxInt64 {
		data := pow.prepareData(block, nonce)
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(target) == -1 {
			fmt.Printf("Block mined with nonce %d\n", nonce)
			break
		} else {
			nonce++
		}
	}

	return nonce, hash[:]
}

// prepareData prepares the data for hashing
func (pow *PoW) prepareData(block *Block, nonce int64) []byte {
	data := bytes.Join(
		[][]byte{
			block.PrevHash,
			block.Data,
			IntToHex(block.Timestamp),
			IntToHex(int64(pow.Difficulty)),
			IntToHex(nonce),
		},
		[]byte{},
	)

	return data
}

// NewPoS creates a new Proof of Stake consensus
func NewPoS(blockchain *Blockchain) *PoS {
	return &PoS{
		Blockchain: blockchain,
		Validators: make(map[string]int64),
	}
}

// Start starts the PoS consensus
func (pos *PoS) Start() error {
	fmt.Println("Starting Proof of Stake consensus")
	return nil
}

// ValidateBlock validates a block using PoS
func (pos *PoS) ValidateBlock(block *Block) bool {
	// Check if the validator has enough stake
	stake, exists := pos.Validators[block.Validator]
	if !exists || stake <= 0 {
		return false
	}

	// In a real implementation, this would check:
	// 1. Validator selection based on stake
	// 2. Signature validation
	// 3. Timestamp validation

	return true
}

// ProposeBlock proposes a new block using PoS
func (pos *PoS) ProposeBlock(block *Block) error {
	fmt.Println("Proposing block with Proof of Stake")

	// Select validator based on stake
	validator := pos.selectValidator()
	block.Validator = validator

	// In a real implementation, this would:
	// 1. Sign the block with the validator's private key
	// 2. Broadcast the block to other nodes

	return nil
}

// selectValidator selects a validator based on stake
func (pos *PoS) selectValidator() string {
	// In a real implementation, this would use a more sophisticated
	// selection algorithm based on stake and possibly other factors

	// For this example, we'll just return the first validator
	for validator := range pos.Validators {
		return validator
	}

	return ""
}

// AddValidator adds a validator to the PoS consensus
func (pos *PoS) AddValidator(address string, stake int64) {
	pos.Validators[address] = stake
}

// RemoveValidator removes a validator from the PoS consensus
func (pos *PoS) RemoveValidator(address string) {
	delete(pos.Validators, address)
}

// NewPBFT creates a new PBFT consensus
func NewPBFT(blockchain *Blockchain, nodes []string) *PBFT {
	return &PBFT{
		Blockchain: blockchain,
		Nodes:      nodes,
		ViewID:     0,
		SequenceID: 0,
	}
}

// Start starts the PBFT consensus
func (pbft *PBFT) Start() error {
	fmt.Println("Starting PBFT consensus")

	// Select primary node
	if len(pbft.Nodes) > 0 {
		pbft.PrimaryNode = pbft.Nodes[0]
	}

	return nil
}

// ValidateBlock validates a block using PBFT
func (pbft *PBFT) ValidateBlock(block *Block) bool {
	// In a real PBFT implementation, this would:
	// 1. Check if the block has sufficient pre-prepare and prepare messages
	// 2. Verify signatures
	// 3. Check sequence numbers

	// For this example, we'll just return true
	return true
}

// ProposeBlock proposes a new block using PBFT
func (pbft *PBFT) ProposeBlock(block *Block) error {
	fmt.Println("Proposing block with PBFT")

	// In a real PBFT implementation, this would:
	// 1. If primary, broadcast pre-prepare message
	// 2. If replica, send prepare messages
	// 3. Collect sufficient prepare messages
	// 4. Send commit messages
	// 5. Wait for sufficient commit messages
	// 6. Add block to blockchain

	return nil
}

// HandlePrePrepare handles a pre-prepare message
func (pbft *PBFT) HandlePrePrepare(viewID, sequenceID int64, block *Block) {
	// In a real implementation, this would:
	// 1. Verify the primary signature
	// 2. Check view and sequence numbers
	// 3. Broadcast prepare messages
}

// HandlePrepare handles a prepare message
func (pbft *PBFT) HandlePrepare(viewID, sequenceID int64, blockHash []byte, nodeID string) {
	// In a real implementation, this would:
	// 1. Verify the node signature
	// 2. Check view and sequence numbers
	// 3. Store the prepare message
	// 4. If sufficient prepare messages, broadcast commit messages
}

// HandleCommit handles a commit message
func (pbft *PBFT) HandleCommit(viewID, sequenceID int64, blockHash []byte, nodeID string) {
	// In a real implementation, this would:
	// 1. Verify the node signature
	// 2. Check view and sequence numbers
	// 3. Store the commit message
	// 4. If sufficient commit messages, add block to blockchain
}

// IntToHex converts an int64 to a byte array
func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	binary.Write(buff, binary.BigEndian, num)
	return buff.Bytes()
}

// Blockchain represents a simple blockchain structure for consensus
type Blockchain struct {
	Blocks []*Block
}
