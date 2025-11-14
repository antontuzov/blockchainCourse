package week1

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"math/big"
)

// HashData computes the SHA-256 hash of the input data
func HashData(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

// GenerateKeyPair generates a new ECDSA key pair
func GenerateKeyPair() (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	return privateKey, &privateKey.PublicKey, nil
}

// SignData signs the data with the private key
func SignData(privateKey *ecdsa.PrivateKey, data []byte) (r, s *big.Int, err error) {
	hash := sha256.Sum256(data)
	return ecdsa.Sign(rand.Reader, privateKey, hash[:])
}

// VerifySignature verifies the signature with the public key
func VerifySignature(publicKey *ecdsa.PublicKey, data []byte, r, s *big.Int) bool {
	hash := sha256.Sum256(data)
	return ecdsa.Verify(publicKey, hash[:], r, s)
}

// MerkleTree represents a Merkle tree structure
type MerkleTree struct {
	Root *MerkleNode
}

// MerkleNode represents a node in the Merkle tree
type MerkleNode struct {
	Left  *MerkleNode
	Right *MerkleNode
	Data  []byte
}

// NewMerkleNode creates a new Merkle node
func NewMerkleNode(left, right *MerkleNode, data []byte) *MerkleNode {
	node := &MerkleNode{
		Left:  left,
		Right: right,
		Data:  data,
	}

	if left == nil && right == nil {
		// Leaf node
		hash := sha256.Sum256(data)
		node.Data = hash[:]
	} else {
		// Internal node
		prevHashes := append(left.Data, right.Data...)
		hash := sha256.Sum256(prevHashes)
		node.Data = hash[:]
	}

	return node
}

// NewMerkleTree creates a new Merkle tree from data
func NewMerkleTree(data [][]byte) *MerkleTree {
	var nodes []*MerkleNode

	// Create leaf nodes
	for _, datum := range data {
		node := NewMerkleNode(nil, nil, datum)
		nodes = append(nodes, node)
	}

	// Create internal nodes
	for len(nodes) > 1 {
		var newLevel []*MerkleNode

		for i := 0; i < len(nodes); i += 2 {
			left := nodes[i]
			var right *MerkleNode

			if i+1 < len(nodes) {
				right = nodes[i+1]
			} else {
				// Duplicate the node if odd number of nodes
				right = left
			}

			newNode := NewMerkleNode(left, right, nil)
			newLevel = append(newLevel, newNode)
		}

		nodes = newLevel
	}

	tree := &MerkleTree{
		Root: nodes[0],
	}

	return tree
}
