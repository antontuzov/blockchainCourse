package week7

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"math/big"
	"time"
)

// PermissionedBlockchain represents a permissioned blockchain
type PermissionedBlockchain struct {
	Channels   map[string]*Channel
	Members    map[string]*Member
	CA         *CertificateAuthority
	Policies   []Policy
	Blockchain *Blockchain
}

// Channel represents a channel in the permissioned blockchain
type Channel struct {
	Name       string
	Members    map[string]bool
	Blockchain *Blockchain
	Policies   []Policy
}

// Member represents a member of the permissioned blockchain
type Member struct {
	ID          string
	Name        string
	Role        string
	Certificate *x509.Certificate
	PrivateKey  *ecdsa.PrivateKey
	PublicKey   *ecdsa.PublicKey
}

// CertificateAuthority represents a certificate authority
type CertificateAuthority struct {
	Certificate *x509.Certificate
	PrivateKey  *ecdsa.PrivateKey
	PublicKey   *ecdsa.PublicKey
}

// Policy represents an access control policy
type Policy struct {
	ID          string
	Name        string
	Description string
	Rules       []PolicyRule
}

// PolicyRule represents a rule in an access control policy
type PolicyRule struct {
	Role        string
	Permissions []string
}

// NewPermissionedBlockchain creates a new permissioned blockchain
func NewPermissionedBlockchain() *PermissionedBlockchain {
	// Create a certificate authority
	ca := NewCertificateAuthority()

	return &PermissionedBlockchain{
		Channels: make(map[string]*Channel),
		Members:  make(map[string]*Member),
		CA:       ca,
		Policies: make([]Policy, 0),
		Blockchain: &Blockchain{
			Blocks: make([]*Block, 0),
		},
	}
}

// NewCertificateAuthority creates a new certificate authority
func NewCertificateAuthority() *CertificateAuthority {
	// Generate a private key
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		fmt.Printf("Error generating private key: %s\n", err)
		return nil
	}

	// Create a certificate template
	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Blockchain CA"},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:  x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageClientAuth,
			x509.ExtKeyUsageServerAuth,
		},
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	// Create a self-signed certificate
	certBytes, err := x509.CreateCertificate(rand.Reader, template, template, &privateKey.PublicKey, privateKey)
	if err != nil {
		fmt.Printf("Error creating certificate: %s\n", err)
		return nil
	}

	// Parse the certificate
	cert, err := x509.ParseCertificate(certBytes)
	if err != nil {
		fmt.Printf("Error parsing certificate: %s\n", err)
		return nil
	}

	return &CertificateAuthority{
		Certificate: cert,
		PrivateKey:  privateKey,
		PublicKey:   &privateKey.PublicKey,
	}
}

// RegisterMember registers a new member in the permissioned blockchain
func (pb *PermissionedBlockchain) RegisterMember(id, name, role string) (*Member, error) {
	// Check if member already exists
	if _, exists := pb.Members[id]; exists {
		return nil, fmt.Errorf("member already exists")
	}

	// Generate a private key for the member
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("error generating private key: %s", err)
	}

	// Create a certificate template for the member
	template := &x509.Certificate{
		SerialNumber: big.NewInt(time.Now().Unix()),
		Subject: pkix.Name{
			Organization: []string{name},
			CommonName:   id,
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:  x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageClientAuth,
		},
		BasicConstraintsValid: true,
	}

	// Create a certificate for the member
	certBytes, err := x509.CreateCertificate(rand.Reader, template, pb.CA.Certificate, &privateKey.PublicKey, pb.CA.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("error creating certificate: %s", err)
	}

	// Parse the certificate
	cert, err := x509.ParseCertificate(certBytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing certificate: %s", err)
	}

	// Create the member
	member := &Member{
		ID:          id,
		Name:        name,
		Role:        role,
		Certificate: cert,
		PrivateKey:  privateKey,
		PublicKey:   &privateKey.PublicKey,
	}

	// Add the member to the blockchain
	pb.Members[id] = member

	fmt.Printf("Registered member: %s (%s)\n", name, id)
	return member, nil
}

// CreateChannel creates a new channel in the permissioned blockchain
func (pb *PermissionedBlockchain) CreateChannel(name string, memberIDs []string) (*Channel, error) {
	// Check if channel already exists
	if _, exists := pb.Channels[name]; exists {
		return nil, fmt.Errorf("channel already exists")
	}

	// Create the channel
	channel := &Channel{
		Name:       name,
		Members:    make(map[string]bool),
		Blockchain: &Blockchain{Blocks: make([]*Block, 0)},
		Policies:   make([]Policy, 0),
	}

	// Add members to the channel
	for _, memberID := range memberIDs {
		// Check if member exists
		if _, exists := pb.Members[memberID]; !exists {
			return nil, fmt.Errorf("member %s does not exist", memberID)
		}

		// Add member to channel
		channel.Members[memberID] = true
	}

	// Add the channel to the blockchain
	pb.Channels[name] = channel

	fmt.Printf("Created channel: %s\n", name)
	return channel, nil
}

// AddMemberToChannel adds a member to a channel
func (pb *PermissionedBlockchain) AddMemberToChannel(channelName, memberID string) error {
	// Check if channel exists
	channel, exists := pb.Channels[channelName]
	if !exists {
		return fmt.Errorf("channel %s does not exist", channelName)
	}

	// Check if member exists
	if _, exists := pb.Members[memberID]; !exists {
		return fmt.Errorf("member %s does not exist", memberID)
	}

	// Add member to channel
	channel.Members[memberID] = true

	fmt.Printf("Added member %s to channel %s\n", memberID, channelName)
	return nil
}

// RemoveMemberFromChannel removes a member from a channel
func (pb *PermissionedBlockchain) RemoveMemberFromChannel(channelName, memberID string) error {
	// Check if channel exists
	channel, exists := pb.Channels[channelName]
	if !exists {
		return fmt.Errorf("channel %s does not exist", channelName)
	}

	// Check if member exists in channel
	if _, exists := channel.Members[memberID]; !exists {
		return fmt.Errorf("member %s is not in channel %s", memberID, channelName)
	}

	// Remove member from channel
	delete(channel.Members, memberID)

	fmt.Printf("Removed member %s from channel %s\n", memberID, channelName)
	return nil
}

// CreatePolicy creates a new access control policy
func (pb *PermissionedBlockchain) CreatePolicy(id, name, description string, rules []PolicyRule) *Policy {
	policy := &Policy{
		ID:          id,
		Name:        name,
		Description: description,
		Rules:       rules,
	}

	// Add the policy to the blockchain
	pb.Policies = append(pb.Policies, *policy)

	fmt.Printf("Created policy: %s\n", name)
	return policy
}

// ApplyPolicyToChannel applies a policy to a channel
func (pb *PermissionedBlockchain) ApplyPolicyToChannel(channelName, policyID string) error {
	// Check if channel exists
	channel, exists := pb.Channels[channelName]
	if !exists {
		return fmt.Errorf("channel %s does not exist", channelName)
	}

	// Find the policy
	var policy *Policy
	for _, p := range pb.Policies {
		if p.ID == policyID {
			p := p
			policy = &p
			break
		}
	}

	if policy == nil {
		return fmt.Errorf("policy %s does not exist", policyID)
	}

	// Apply the policy to the channel
	channel.Policies = append(channel.Policies, *policy)

	fmt.Printf("Applied policy %s to channel %s\n", policy.Name, channelName)
	return nil
}

// VerifyMember verifies a member's certificate
func (pb *PermissionedBlockchain) VerifyMember(memberID string) bool {
	// Check if member exists
	member, exists := pb.Members[memberID]
	if !exists {
		return false
	}

	// Verify the certificate
	roots := x509.NewCertPool()
	roots.AddCert(pb.CA.Certificate)

	_, err := member.Certificate.Verify(x509.VerifyOptions{
		Roots:         roots,
		Intermediates: nil,
		KeyUsages:     []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	})

	return err == nil
}

// GetChannelMembers returns the members of a channel
func (pb *PermissionedBlockchain) GetChannelMembers(channelName string) ([]*Member, error) {
	// Check if channel exists
	channel, exists := pb.Channels[channelName]
	if !exists {
		return nil, fmt.Errorf("channel %s does not exist", channelName)
	}

	// Get the members
	members := make([]*Member, 0)
	for memberID := range channel.Members {
		if member, exists := pb.Members[memberID]; exists {
			members = append(members, member)
		}
	}

	return members, nil
}

// Block represents a block in the blockchain
type Block struct {
	Index      int64
	Timestamp  int64
	Data       []byte
	PrevHash   []byte
	Hash       []byte
	Nonce      int64
	Channel    string
	Creator    string
	Signatures map[string][]byte
}

// Blockchain represents a simple blockchain structure
type Blockchain struct {
	Blocks []*Block
}
