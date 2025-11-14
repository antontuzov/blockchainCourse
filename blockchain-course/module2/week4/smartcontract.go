package week4

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// SmartContract represents a smart contract
type SmartContract struct {
	Code    []byte
	Address string
	Storage map[string]string
	Balance int
}

// VM represents a virtual machine for executing smart contracts
type VM struct {
	stack   []interface{}
	memory  map[string]interface{}
	storage map[string]string
	gas     int
}

// NewSmartContract creates a new smart contract
func NewSmartContract(code []byte) *SmartContract {
	contract := &SmartContract{
		Code:    code,
		Address: generateContractAddress(code),
		Storage: make(map[string]string),
		Balance: 0,
	}
	return contract
}

// NewVM creates a new virtual machine
func NewVM() *VM {
	vm := &VM{
		stack:   make([]interface{}, 0),
		memory:  make(map[string]interface{}),
		storage: make(map[string]string),
		gas:     10000, // Initial gas limit
	}
	return vm
}

// Execute executes a smart contract
func (vm *VM) Execute(contract *SmartContract) error {
	// Deduct gas for execution
	vm.gas -= 100

	// Simple execution logic based on contract code
	code := string(contract.Code)

	switch code {
	case "increment":
		// Get current counter value
		counterStr := contract.Storage["counter"]
		if counterStr == "" {
			counterStr = "0"
		}

		// Increment counter
		var counter int
		fmt.Sscanf(counterStr, "%d", &counter)
		counter++

		// Store updated counter
		contract.Storage["counter"] = fmt.Sprintf("%d", counter)
	default:
		// For other codes, just store in storage
		contract.Storage["executed"] = "true"
	}

	return nil
}

// DeployContract deploys a smart contract
func (vm *VM) DeployContract(contract *SmartContract) string {
	// Deduct gas for deployment
	vm.gas -= 500

	// In a real implementation, this would involve more complex logic
	// For now, we just return the contract address
	return contract.Address
}

// CallContract calls a smart contract method
func (vm *VM) CallContract(contract *SmartContract, method string) error {
	// Deduct gas for calling
	vm.gas -= 50

	// Simple method calling logic
	switch method {
	case "increment":
		return vm.Execute(contract)
	default:
		// Store method call in storage
		contract.Storage["last_call"] = method
		return nil
	}
}

// GetBalance returns the contract balance
func (contract *SmartContract) GetBalance() int {
	return contract.Balance
}

// AddBalance adds to the contract balance
func (contract *SmartContract) AddBalance(amount int) {
	contract.Balance += amount
}

// SetStorageValue sets a value in contract storage
func (contract *SmartContract) SetStorageValue(key, value string) {
	contract.Storage[key] = value
}

// GetStorageValue gets a value from contract storage
func (contract *SmartContract) GetStorageValue(key string) string {
	return contract.Storage[key]
}

// GetGas returns the current gas
func (vm *VM) GetGas() int {
	return vm.gas
}

// SetGas sets the gas
func (vm *VM) SetGas(gas int) {
	vm.gas = gas
}

// generateContractAddress generates a contract address from code
func generateContractAddress(code []byte) string {
	hash := sha256.Sum256(code)
	return hex.EncodeToString(hash[:])
}
