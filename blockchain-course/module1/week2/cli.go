package week2

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// CLI represents the command line interface
type CLI struct {
	bc *Blockchain
}

// NewCLI creates a new CLI
func NewCLI(bc *Blockchain) *CLI {
	return &CLI{bc: bc}
}

// Run starts the CLI
func (cli *CLI) Run() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("Blockchain CLI")
		fmt.Println("1. Add block")
		fmt.Println("2. Print blockchain")
		fmt.Println("3. Validate blockchain")
		fmt.Println("4. Exit")
		fmt.Print("Enter choice: ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			cli.addBlock()
		case "2":
			cli.printBlockchain()
		case "3":
			cli.validateBlockchain()
		case "4":
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}

func (cli *CLI) addBlock() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter block data: ")
	data, _ := reader.ReadString('\n')
	data = strings.TrimSpace(data)

	cli.bc.AddBlock(data)
	fmt.Println("Block added successfully!")
}

func (cli *CLI) printBlockchain() {
	for i, block := range cli.bc.Blocks {
		fmt.Printf("Block %d:\n", i)
		fmt.Printf("  Timestamp: %d\n", block.Timestamp)
		fmt.Printf("  Data: %s\n", block.Data)
		fmt.Printf("  Previous Hash: %x\n", block.PrevBlockHash)
		fmt.Printf("  Hash: %x\n", block.Hash)
		fmt.Printf("  Nonce: %d\n", block.Nonce)
		fmt.Println()
	}
}

func (cli *CLI) validateBlockchain() {
	if cli.bc.IsValid() {
		fmt.Println("Blockchain is valid!")
	} else {
		fmt.Println("Blockchain is invalid!")
	}
}
