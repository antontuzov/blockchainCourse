package week1

import (
	"testing"
)

func TestHashData(t *testing.T) {
	data := []byte("Hello, Blockchain!")
	hash := HashData(data)

	if len(hash) != 32 {
		t.Errorf("Expected hash length of 32, got %d", len(hash))
	}

	// Test that the same input produces the same output
	hash2 := HashData(data)
	for i := range hash {
		if hash[i] != hash2[i] {
			t.Error("Hash function is not deterministic")
		}
	}
}

func TestGenerateKeyPair(t *testing.T) {
	privateKey, publicKey, err := GenerateKeyPair()

	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	if privateKey == nil {
		t.Error("Private key is nil")
	}

	if publicKey == nil {
		t.Error("Public key is nil")
	}
}

func TestSignAndVerify(t *testing.T) {
	privateKey, publicKey, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	data := []byte("Data to sign")
	r, s, err := SignData(privateKey, data)
	if err != nil {
		t.Fatalf("Failed to sign data: %v", err)
	}

	valid := VerifySignature(publicKey, data, r, s)
	if !valid {
		t.Error("Signature verification failed")
	}

	// Test with wrong data
	wrongData := []byte("Wrong data")
	valid = VerifySignature(publicKey, wrongData, r, s)
	if valid {
		t.Error("Signature verification should fail with wrong data")
	}
}

func TestMerkleTree(t *testing.T) {
	data := [][]byte{
		[]byte("data1"),
		[]byte("data2"),
		[]byte("data3"),
		[]byte("data4"),
	}

	tree := NewMerkleTree(data)

	if tree.Root == nil {
		t.Error("Merkle tree root is nil")
	}

	if len(tree.Root.Data) != 32 {
		t.Errorf("Expected root hash length of 32, got %d", len(tree.Root.Data))
	}
}
