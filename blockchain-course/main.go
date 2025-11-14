package main

import (
	"fmt"
	"os"

	"blockchain-course/module1/week2"
)

func main() {
	// Create a new blockchain
	bc := week2.NewBlockchain()

	// Add some sample blocks
	bc.AddBlock("First transaction")
	bc.AddBlock("Second transaction")
	bc.AddBlock("Third transaction")

	// Print the blockchain
	for i, block := range bc.Blocks {
		fmt.Printf("Block %d:\n", i)
		fmt.Printf("  Timestamp: %d\n", block.Timestamp)
		fmt.Printf("  Data: %s\n", block.Data)
		fmt.Printf("  Previous Hash: %x\n", block.PrevBlockHash)
		fmt.Printf("  Hash: %x\n", block.Hash)
		fmt.Printf("  Nonce: %d\n", block.Nonce)
		fmt.Println()
	}

	// Validate the blockchain
	if bc.IsValid() {
		fmt.Println("Blockchain is valid!")
	} else {
		fmt.Println("Blockchain is invalid!")
		os.Exit(1)
	}

	// Start the CLI
	cli := week2.NewCLI(bc)
	cli.Run()
}
