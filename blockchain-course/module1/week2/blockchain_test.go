package week2

import (
	"testing"

	week1 "blockchain-course/module1/week1"
)

func TestNewBlockchain(t *testing.T) {
	bc := NewBlockchain()

	if len(bc.Blocks) != 1 {
		t.Errorf("Expected blockchain to have 1 block, got %d", len(bc.Blocks))
	}

	if string(bc.Blocks[0].Data) != "Genesis Block" {
		t.Error("Genesis block data is incorrect")
	}
}

func TestAddBlock(t *testing.T) {
	bc := NewBlockchain()

	// Add a block
	bc.AddBlock("Test data")

	if len(bc.Blocks) != 2 {
		t.Errorf("Expected blockchain to have 2 blocks, got %d", len(bc.Blocks))
	}

	if string(bc.Blocks[1].Data) != "Test data" {
		t.Error("Block data is incorrect")
	}
}

func TestIsValid(t *testing.T) {
	bc := NewBlockchain()

	// Add some blocks
	bc.AddBlock("Block 1")
	bc.AddBlock("Block 2")

	// Check if blockchain is valid
	if !bc.IsValid() {
		t.Error("Blockchain should be valid")
	}

	// Tamper with a block
	bc.Blocks[1].Data = []byte("Tampered data")
	bc.Blocks[1].SetHash()

	// Check if blockchain is invalid
	if bc.IsValid() {
		t.Error("Blockchain should be invalid after tampering")
	}
}

func TestProofOfWork(t *testing.T) {
	block := week1.NewBlock("Test data", []byte{})
	pow := NewProofOfWork(block)

	nonce, hash := pow.Run()

	if nonce < 0 {
		t.Error("Nonce should be non-negative")
	}

	if len(hash) != 32 {
		t.Errorf("Expected hash length of 32, got %d", len(hash))
	}

	// Validate the proof of work
	block.Nonce = nonce
	block.Hash = hash
	valid := pow.Validate()

	if !valid {
		t.Error("Proof of work should be valid")
	}
}
