package week3

import (
	"testing"
)

func TestNewTransaction(t *testing.T) {
	// Create a new wallet
	wallet := NewWallet()
	address := wallet.GetAddress()

	// Create a coinbase transaction
	cbTx := NewCoinbaseTX(string(address), "Test Coinbase")

	if cbTx == nil {
		t.Error("Failed to create coinbase transaction")
	}

	if !cbTx.IsCoinbase() {
		t.Error("Transaction should be coinbase")
	}
}

func TestTXOutputLock(t *testing.T) {
	// Create a new wallet
	wallet := NewWallet()
	address := wallet.GetAddress()

	// Create a new TXOutput
	output := TXOutput{100, nil}

	// Lock the output with the address
	output.Lock(address)

	if output.PubKeyHash == nil {
		t.Error("PubKeyHash should not be nil after locking")
	}
}

func TestTXOutputIsLockedWithKey(t *testing.T) {
	// Create a new wallet
	wallet := NewWallet()
	pubKeyHash := HashPubKey(wallet.PublicKey)

	// Create a new TXOutput and lock it
	output := TXOutput{100, nil}
	output.Lock(wallet.GetAddress())

	// Check if the output is locked with the correct key
	if !output.IsLockedWithKey(pubKeyHash) {
		t.Error("Output should be locked with the correct key")
	}

	// Create another wallet
	anotherWallet := NewWallet()
	anotherPubKeyHash := HashPubKey(anotherWallet.PublicKey)

	// Check if the output is not locked with another key
	if output.IsLockedWithKey(anotherPubKeyHash) {
		t.Error("Output should not be locked with another key")
	}
}

func TestTXInputUsesKey(t *testing.T) {
	// Create a new wallet
	wallet := NewWallet()
	pubKeyHash := HashPubKey(wallet.PublicKey)

	// Create a new TXInput
	input := TXInput{[]byte("txid"), 0, nil, wallet.PublicKey}

	// Check if the input uses the correct key
	if !input.UsesKey(pubKeyHash) {
		t.Error("Input should use the correct key")
	}

	// Create another wallet
	anotherWallet := NewWallet()
	anotherPubKeyHash := HashPubKey(anotherWallet.PublicKey)

	// Check if the input does not use another key
	if input.UsesKey(anotherPubKeyHash) {
		t.Error("Input should not use another key")
	}
}

func TestTransactionHash(t *testing.T) {
	// Create a new transaction
	tx := Transaction{nil, nil, nil}
	hash := tx.Hash()

	if len(hash) != 32 {
		t.Errorf("Expected hash length of 32, got %d", len(hash))
	}
}

func TestTransactionIsCoinbase(t *testing.T) {
	// Create a coinbase transaction
	cbTx := NewCoinbaseTX("address", "Test Coinbase")

	if !cbTx.IsCoinbase() {
		t.Error("Transaction should be coinbase")
	}

	// Create a regular transaction
	tx := Transaction{nil, []TXInput{{[]byte("txid"), 0, nil, nil}}, nil}

	if tx.IsCoinbase() {
		t.Error("Transaction should not be coinbase")
	}
}
