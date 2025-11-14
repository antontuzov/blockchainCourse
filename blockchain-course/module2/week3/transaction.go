package week3

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"math/big"

	"blockchain-course/module1/transaction"
)

// Transaction represents a transaction in the blockchain
type Transaction struct {
	ID   []byte
	Vin  []TXInput
	Vout []TXOutput
}

// TXOutput represents a transaction output
type TXOutput struct {
	Value      int
	PubKeyHash []byte
}

// Lock signs the output
func (out *TXOutput) Lock(address []byte) {
	pubKeyHash := Base58Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	out.PubKeyHash = pubKeyHash
}

// IsLockedWithKey checks if the output can be used by the owner of the pubkey
func (out *TXOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}

// Hash returns the hash of the Transaction
func (tx *Transaction) Hash() []byte {
	var hash [32]byte

	txCopy := *tx
	txCopy.ID = []byte{}

	hash = sha256.Sum256(txCopy.Serialize())
	return hash[:]
}

// Serialize returns a serialized Transaction
func (tx Transaction) Serialize() []byte {
	var encoded bytes.Buffer

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}

	return encoded.Bytes()
}

// IsCoinbase checks whether the transaction is coinbase
func (tx Transaction) IsCoinbase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1
}

// TXInput represents a transaction input
type TXInput struct {
	Txid      []byte
	Vout      int
	Signature []byte
	PubKey    []byte
}

// UsesKey checks whether the address initiated the transaction
func (in *TXInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := HashPubKey(in.PubKey)
	return bytes.Compare(lockingHash, pubKeyHash) == 0
}

// NewTransaction creates a new transaction
func NewTransaction(from, to string, amount int, bc *Blockchain) *transaction.Transaction {
	var inputs []transaction.TXInput
	var outputs []transaction.TXOutput

	// Get wallet for the sender
	wallets, err := NewWallets()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return nil
	}
	wallet := wallets.GetWallet(from)
	pubKeyHash := HashPubKey(wallet.PublicKey)

	// Find spendable outputs
	acc, validOutputs := bc.FindSpendableOutputs(pubKeyHash, amount)

	if acc < amount {
		fmt.Printf("Error: Not enough funds\n")
		return nil
	}

	// Build a list of inputs
	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return nil
		}

		for _, out := range outs {
			input := transaction.TXInput{txID, out, nil, wallet.PublicKey}
			inputs = append(inputs, input)
		}
	}

	// Build a list of outputs
	outputs = append(outputs, *NewTXOutput(amount, to))
	if acc > amount {
		outputs = append(outputs, *NewTXOutput(acc-amount, from)) // a change
	}

	tx := transaction.Transaction{nil, inputs, outputs}
	tx.ID = tx.Hash()
	bc.SignTransaction(&tx, wallet.PrivateKey)

	return &tx
}

// NewCoinbaseTX creates a new coinbase transaction
func NewCoinbaseTX(to, data string) *transaction.Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to \"%s\"", to)
	}

	txin := transaction.TXInput{[]byte{}, -1, nil, []byte(data)}
	txout := NewTXOutput(10, to)
	tx := transaction.Transaction{nil, []transaction.TXInput{txin}, []transaction.TXOutput{*txout}}
	tx.ID = tx.Hash()

	return &tx
}

// NewTXOutput creates a new TXOutput
func NewTXOutput(value int, address string) *transaction.TXOutput {
	txo := &transaction.TXOutput{value, nil}
	txo.Lock([]byte(address))

	return txo
}

// SerializeTransaction serializes a Transaction
func SerializeTransaction(tx transaction.Transaction) []byte {
	var encoded bytes.Buffer

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}

	return encoded.Bytes()
}

// IsCoinbaseTransaction checks whether the transaction is coinbase
func IsCoinbaseTransaction(tx transaction.Transaction) bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1
}

// SignTransaction signs each input of a Transaction
func SignTransaction(tx *transaction.Transaction, privKey ecdsa.PrivateKey, prevTXs map[string]transaction.Transaction) {
	if IsCoinbaseTransaction(*tx) {
		return
	}

	for _, vin := range tx.Vin {
		if prevTXs[hex.EncodeToString(vin.Txid)].ID == nil {
			fmt.Printf("Error: Previous transaction is not correct\n")
			return
		}
	}

	txCopy := TrimmedCopy(*tx)

	for inID, vin := range txCopy.Vin {
		prevTx := prevTXs[hex.EncodeToString(vin.Txid)]
		txCopy.Vin[inID].Signature = nil
		txCopy.Vin[inID].PubKey = prevTx.Vout[vin.Vout].PubKeyHash

		dataToSign := fmt.Sprintf("%x\n", txCopy)

		r, s, err := ecdsa.Sign(rand.Reader, &privKey, []byte(dataToSign))
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}
		signature := append(r.Bytes(), s.Bytes()...)

		tx.Vin[inID].Signature = signature
		txCopy.Vin[inID].PubKey = nil
	}
}

// TrimmedCopy creates a trimmed copy of Transaction to be used in signing
func TrimmedCopy(tx transaction.Transaction) transaction.Transaction {
	var inputs []transaction.TXInput
	var outputs []transaction.TXOutput

	for _, vin := range tx.Vin {
		inputs = append(inputs, transaction.TXInput{vin.Txid, vin.Vout, nil, nil})
	}

	for _, vout := range tx.Vout {
		outputs = append(outputs, transaction.TXOutput{vout.Value, vout.PubKeyHash})
	}

	txCopy := transaction.Transaction{tx.ID, inputs, outputs}

	return txCopy
}

// VerifyTransaction verifies signatures of Transaction inputs
func VerifyTransaction(tx transaction.Transaction, prevTXs map[string]transaction.Transaction) bool {
	if IsCoinbaseTransaction(tx) {
		return true
	}

	for _, vin := range tx.Vin {
		if prevTXs[hex.EncodeToString(vin.Txid)].ID == nil {
			fmt.Printf("Error: Previous transaction is not correct\n")
			return false
		}
	}

	txCopy := TrimmedCopy(tx)
	curve := elliptic.P256()

	for inID, vin := range tx.Vin {
		prevTx := prevTXs[hex.EncodeToString(vin.Txid)]
		txCopy.Vin[inID].Signature = nil
		txCopy.Vin[inID].PubKey = prevTx.Vout[vin.Vout].PubKeyHash

		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.Signature)
		r.SetBytes(vin.Signature[:(sigLen / 2)])
		s.SetBytes(vin.Signature[(sigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(vin.PubKey)
		x.SetBytes(vin.PubKey[:(keyLen / 2)])
		y.SetBytes(vin.PubKey[(keyLen / 2):])

		rawPubKey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}
		dataToSign := fmt.Sprintf("%x\n", txCopy)
		if ecdsa.Verify(&rawPubKey, []byte(dataToSign), &r, &s) == false {
			return false
		}
		txCopy.Vin[inID].PubKey = nil
	}

	return true
}
