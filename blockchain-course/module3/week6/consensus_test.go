package week6

import (
	"testing"
)

func TestNewPoW(t *testing.T) {
	blockchain := &Blockchain{}
	pow := NewPoW(blockchain, 4)

	if pow == nil {
		t.Error("Failed to create PoW consensus")
	}

	if pow.Blockchain != blockchain {
		t.Error("PoW blockchain reference is incorrect")
	}

	if pow.Difficulty != 4 {
		t.Error("PoW difficulty is incorrect")
	}
}

func TestPoWValidateBlock(t *testing.T) {
	blockchain := &Blockchain{}
	pow := NewPoW(blockchain, 4)

	// Create a block with an invalid hash for difficulty 4
	block := &Block{
		Index:     0,
		Timestamp: 1234567890,
		Data:      []byte("test data"),
		PrevHash:  []byte("previous hash"),
		Hash:      []byte{0x01, 0x00, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d},
		Nonce:     12345,
	}

	// This should fail because the hash doesn't meet the difficulty requirement
	if pow.ValidateBlock(block) {
		t.Error("Block should not be valid with this hash")
	}

	// Create a block with a valid hash for difficulty 4
	block2 := &Block{
		Index:     0,
		Timestamp: 1234567890,
		Data:      []byte("test data"),
		PrevHash:  []byte("previous hash"),
		Hash:      []byte{0x00, 0x00, 0x00, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c},
		Nonce:     12345,
	}

	// This should pass because the hash meets the difficulty requirement
	if !pow.ValidateBlock(block2) {
		t.Error("Block should be valid with this hash")
	}
}

func TestNewPoS(t *testing.T) {
	blockchain := &Blockchain{}
	pos := NewPoS(blockchain)

	if pos == nil {
		t.Error("Failed to create PoS consensus")
	}

	if pos.Blockchain != blockchain {
		t.Error("PoS blockchain reference is incorrect")
	}

	if len(pos.Validators) != 0 {
		t.Error("PoS validators should be empty initially")
	}
}

func TestPoSValidatorManagement(t *testing.T) {
	blockchain := &Blockchain{}
	pos := NewPoS(blockchain)

	// Add a validator
	pos.AddValidator("validator1", 100)

	if len(pos.Validators) != 1 {
		t.Error("Validator should be added")
	}

	stake, exists := pos.Validators["validator1"]
	if !exists {
		t.Error("Validator should exist")
	}

	if stake != 100 {
		t.Error("Validator stake is incorrect")
	}

	// Remove a validator
	pos.RemoveValidator("validator1")

	if len(pos.Validators) != 0 {
		t.Error("Validator should be removed")
	}
}

func TestPoSValidateBlock(t *testing.T) {
	blockchain := &Blockchain{}
	pos := NewPoS(blockchain)

	// Add a validator
	pos.AddValidator("validator1", 100)

	// Create a block with a valid validator
	block := &Block{
		Index:     0,
		Timestamp: 1234567890,
		Data:      []byte("test data"),
		PrevHash:  []byte("previous hash"),
		Hash:      []byte("block hash"),
		Nonce:     12345,
		Validator: "validator1",
	}

	// This should pass because the validator exists and has stake
	if !pos.ValidateBlock(block) {
		t.Error("Block should be valid with this validator")
	}

	// Create a block with an invalid validator
	block2 := &Block{
		Index:     0,
		Timestamp: 1234567890,
		Data:      []byte("test data"),
		PrevHash:  []byte("previous hash"),
		Hash:      []byte("block hash"),
		Nonce:     12345,
		Validator: "validator2",
	}

	// This should fail because the validator doesn't exist
	if pos.ValidateBlock(block2) {
		t.Error("Block should not be valid with this validator")
	}
}

func TestNewPBFT(t *testing.T) {
	blockchain := &Blockchain{}
	nodes := []string{"node1", "node2", "node3"}
	pbft := NewPBFT(blockchain, nodes)

	if pbft == nil {
		t.Error("Failed to create PBFT consensus")
	}

	if pbft.Blockchain != blockchain {
		t.Error("PBFT blockchain reference is incorrect")
	}

	if len(pbft.Nodes) != 3 {
		t.Error("PBFT nodes count is incorrect")
	}

	if pbft.ViewID != 0 {
		t.Error("PBFT view ID should be 0 initially")
	}

	if pbft.SequenceID != 0 {
		t.Error("PBFT sequence ID should be 0 initially")
	}
}

func TestPBFTStart(t *testing.T) {
	blockchain := &Blockchain{}
	nodes := []string{"node1", "node2", "node3"}
	pbft := NewPBFT(blockchain, nodes)

	err := pbft.Start()
	if err != nil {
		t.Errorf("Failed to start PBFT: %s", err)
	}

	// Check that primary node is set
	if pbft.PrimaryNode != "node1" {
		t.Error("Primary node should be set to first node")
	}
}

func TestIntToHex(t *testing.T) {
	// Test conversion of various numbers
	testCases := []int64{0, 1, 10, 100, 1000, 1234567890}

	for _, num := range testCases {
		result := IntToHex(num)
		if len(result) != 8 {
			t.Errorf("Expected 8 bytes for %d, got %d", num, len(result))
		}
	}
}
