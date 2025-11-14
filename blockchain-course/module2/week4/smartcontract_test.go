package week4

import (
	"testing"
)

func TestNewSmartContract(t *testing.T) {
	code := []byte("test code")
	contract := NewSmartContract(code)

	if contract == nil {
		t.Error("Failed to create smart contract")
	}

	if string(contract.Code) != "test code" {
		t.Error("Contract code is incorrect")
	}

	if contract.Address == "" {
		t.Error("Contract address should not be empty")
	}

	if len(contract.Storage) != 0 {
		t.Error("Contract storage should be empty initially")
	}

	if contract.Balance != 0 {
		t.Error("Contract balance should be zero initially")
	}
}

func TestNewVM(t *testing.T) {
	vm := NewVM()

	if vm == nil {
		t.Error("Failed to create VM")
	}

	if len(vm.stack) != 0 {
		t.Error("VM stack should be empty initially")
	}

	if len(vm.memory) != 0 {
		t.Error("VM memory should be empty initially")
	}

	if len(vm.storage) != 0 {
		t.Error("VM storage should be empty initially")
	}

	if vm.gas != 10000 {
		t.Error("VM gas should be 10000 initially")
	}
}

func TestExecute(t *testing.T) {
	vm := NewVM()
	contract := NewSmartContract([]byte("increment"))

	// Execute the contract
	err := vm.Execute(contract)
	if err != nil {
		t.Errorf("Failed to execute contract: %s", err)
	}

	// Check that the counter was incremented
	if contract.Storage["counter"] != "1" {
		t.Error("Counter should be 1 after increment")
	}

	// Execute again
	err = vm.Execute(contract)
	if err != nil {
		t.Errorf("Failed to execute contract: %s", err)
	}

	// Check that the counter was incremented again
	if contract.Storage["counter"] != "2" {
		t.Error("Counter should be 2 after second increment")
	}
}

func TestDeployContract(t *testing.T) {
	vm := NewVM()
	contract := NewSmartContract([]byte("test code"))

	address := vm.DeployContract(contract)

	if address == "" {
		t.Error("Contract address should not be empty")
	}

	if address != contract.Address {
		t.Error("Deployed address should match contract address")
	}
}

func TestCallContract(t *testing.T) {
	vm := NewVM()
	contract := NewSmartContract([]byte("increment"))

	// Call the contract
	err := vm.CallContract(contract, "increment")
	if err != nil {
		t.Errorf("Failed to call contract: %s", err)
	}

	// Check that the counter was incremented
	if contract.Storage["counter"] != "1" {
		t.Error("Counter should be 1 after increment")
	}
}

func TestContractBalance(t *testing.T) {
	contract := NewSmartContract([]byte("test code"))

	// Check initial balance
	if contract.GetBalance() != 0 {
		t.Error("Initial balance should be zero")
	}

	// Add balance
	contract.AddBalance(100)

	// Check updated balance
	if contract.GetBalance() != 100 {
		t.Error("Balance should be 100 after adding 100")
	}

	// Add more balance
	contract.AddBalance(50)

	// Check final balance
	if contract.GetBalance() != 150 {
		t.Error("Balance should be 150 after adding 50 more")
	}
}

func TestStorageOperations(t *testing.T) {
	contract := NewSmartContract([]byte("test code"))

	// Set a value in storage
	contract.SetStorageValue("key1", "value1")

	// Get the value from storage
	value := contract.GetStorageValue("key1")
	if value != "value1" {
		t.Error("Storage value should be 'value1'")
	}

	// Set another value
	contract.SetStorageValue("key2", "value2")

	// Get the second value
	value = contract.GetStorageValue("key2")
	if value != "value2" {
		t.Error("Storage value should be 'value2'")
	}

	// Update first value
	contract.SetStorageValue("key1", "updated_value1")

	// Get the updated value
	value = contract.GetStorageValue("key1")
	if value != "updated_value1" {
		t.Error("Storage value should be 'updated_value1'")
	}
}

func TestGasOperations(t *testing.T) {
	vm := NewVM()

	// Check initial gas
	if vm.GetGas() != 10000 {
		t.Error("Initial gas should be 10000")
	}

	// Set gas
	vm.SetGas(5000)

	// Check updated gas
	if vm.GetGas() != 5000 {
		t.Error("Gas should be 5000 after setting")
	}
}

func TestGenerateContractAddress(t *testing.T) {
	code := []byte("test code")
	address := generateContractAddress(code)

	if address == "" {
		t.Error("Contract address should not be empty")
	}

	if len(address) < 10 {
		t.Error("Contract address should be reasonably long")
	}
}
