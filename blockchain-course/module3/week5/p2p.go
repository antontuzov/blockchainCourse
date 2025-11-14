package week5

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	"net/http"
	"sync"

	"blockchain-course/module1/week2"
)

// Node represents a P2P node in the network
type Node struct {
	Address    string
	Port       int
	Peers      map[string]*Peer
	peersMutex sync.RWMutex
	Server     *http.Server
	Blockchain *week2.Blockchain
}

// Peer represents a peer in the network
type Peer struct {
	Address string
	Port    int
	Conn    net.Conn
}

// Message represents a network message
type Message struct {
	Type    string
	Payload []byte
}

// NewNode creates a new P2P node
func NewNode(address string, port int, blockchain *week2.Blockchain) *Node {
	node := &Node{
		Address:    address,
		Port:       port,
		Peers:      make(map[string]*Peer),
		Blockchain: blockchain,
	}

	// Set up HTTP server for handling requests
	node.Server = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", address, port),
		Handler: node.createRouter(),
	}

	return node
}

// Start starts the P2P node
func (n *Node) Start() error {
	fmt.Printf("Starting node at %s:%d\n", n.Address, n.Port)

	// Start listening for connections
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", n.Address, n.Port))
	if err != nil {
		return err
	}
	defer listener.Close()

	// Start HTTP server in a goroutine
	go func() {
		if err := n.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("HTTP server error: %s\n", err)
		}
	}()

	// Accept incoming connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %s\n", err)
			continue
		}

		// Handle connection in a goroutine
		go n.handleConnection(conn)
	}
}

// Stop stops the P2P node
func (n *Node) Stop() error {
	fmt.Println("Stopping node")
	return n.Server.Close()
}

// AddPeer adds a peer to the node
func (n *Node) AddPeer(address string, port int) error {
	n.peersMutex.Lock()
	defer n.peersMutex.Unlock()

	peerAddress := fmt.Sprintf("%s:%d", address, port)

	// Check if peer already exists
	if _, exists := n.Peers[peerAddress]; exists {
		return fmt.Errorf("peer already exists")
	}

	// Connect to peer
	conn, err := net.Dial("tcp", peerAddress)
	if err != nil {
		return err
	}

	// Create peer
	peer := &Peer{
		Address: address,
		Port:    port,
		Conn:    conn,
	}

	// Add peer to map
	n.Peers[peerAddress] = peer

	fmt.Printf("Added peer: %s\n", peerAddress)
	return nil
}

// RemovePeer removes a peer from the node
func (n *Node) RemovePeer(address string, port int) {
	n.peersMutex.Lock()
	defer n.peersMutex.Unlock()

	peerAddress := fmt.Sprintf("%s:%d", address, port)

	// Check if peer exists
	peer, exists := n.Peers[peerAddress]
	if !exists {
		return
	}

	// Close connection
	if peer.Conn != nil {
		peer.Conn.Close()
	}

	// Remove peer from map
	delete(n.Peers, peerAddress)

	fmt.Printf("Removed peer: %s\n", peerAddress)
}

// BroadcastMessage broadcasts a message to all peers
func (n *Node) BroadcastMessage(msg *Message) error {
	n.peersMutex.RLock()
	defer n.peersMutex.RUnlock()

	// Serialize message
	data, err := n.serializeMessage(msg)
	if err != nil {
		return err
	}

	// Send message to all peers
	for _, peer := range n.Peers {
		if peer.Conn != nil {
			_, err := peer.Conn.Write(data)
			if err != nil {
				fmt.Printf("Error sending message to peer %s:%d: %s\n", peer.Address, peer.Port, err)
			}
		}
	}

	return nil
}

// HandleMessage handles an incoming message
func (n *Node) HandleMessage(msg *Message) {
	switch msg.Type {
	case "block":
		// Handle block message
		fmt.Println("Received block message")
		// In a real implementation, this would deserialize and add the block to the blockchain
	case "transaction":
		// Handle transaction message
		fmt.Println("Received transaction message")
		// In a real implementation, this would deserialize and add the transaction to the mempool
	case "get_blocks":
		// Handle get blocks request
		fmt.Println("Received get blocks request")
		// In a real implementation, this would send blocks to the requesting peer
	case "ping":
		// Handle ping message
		fmt.Println("Received ping message")
	default:
		fmt.Printf("Unknown message type: %s\n", msg.Type)
	}
}

// DiscoverPeers discovers new peers in the network
func (n *Node) DiscoverPeers(bootstrapNodes []string) {
	fmt.Println("Discovering peers...")

	// In a real implementation, this would:
	// 1. Connect to bootstrap nodes
	// 2. Request peer lists
	// 3. Add discovered peers to the peer list
	// 4. Recursively discover more peers

	for _, bootstrapNode := range bootstrapNodes {
		fmt.Printf("Connecting to bootstrap node: %s\n", bootstrapNode)
		// Add logic to connect and discover peers
	}
}

// SyncBlockchain synchronizes the blockchain with peers
func (n *Node) SyncBlockchain() error {
	fmt.Println("Synchronizing blockchain...")

	// In a real implementation, this would:
	// 1. Request the latest block height from peers
	// 2. Identify missing blocks
	// 3. Request missing blocks
	// 4. Validate and add blocks to the local blockchain

	return nil
}

// serializeMessage serializes a message to bytes
func (n *Node) serializeMessage(msg *Message) ([]byte, error) {
	var encoded bytes.Buffer
	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(msg)
	if err != nil {
		return nil, err
	}
	return encoded.Bytes(), nil
}

// deserializeMessage deserializes bytes to a message
func (n *Node) deserializeMessage(data []byte) (*Message, error) {
	var msg Message
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&msg)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

// handleConnection handles an incoming connection
func (n *Node) handleConnection(conn net.Conn) {
	defer conn.Close()

	// Read data from connection
	buffer := make([]byte, 4096)
	for {
		nBytes, err := conn.Read(buffer)
		if err != nil {
			fmt.Printf("Error reading from connection: %s\n", err)
			return
		}

		// Deserialize message
		msg, err := n.deserializeMessage(buffer[:nBytes])
		if err != nil {
			fmt.Printf("Error deserializing message: %s\n", err)
			continue
		}

		// Handle message
		n.HandleMessage(msg)
	}
}

// createRouter creates an HTTP router for the node
func (n *Node) createRouter() *http.ServeMux {
	router := http.NewServeMux()

	// Add routes
	router.HandleFunc("/peers", n.handlePeers)
	router.HandleFunc("/blocks", n.handleBlocks)
	router.HandleFunc("/transactions", n.handleTransactions)

	return router
}

// handlePeers handles peer-related HTTP requests
func (n *Node) handlePeers(w http.ResponseWriter, r *http.Request) {
	// Implementation would handle peer management via HTTP
	fmt.Fprintf(w, "Peer management endpoint")
}

// handleBlocks handles block-related HTTP requests
func (n *Node) handleBlocks(w http.ResponseWriter, r *http.Request) {
	// Implementation would handle block requests via HTTP
	fmt.Fprintf(w, "Block management endpoint")
}

// handleTransactions handles transaction-related HTTP requests
func (n *Node) handleTransactions(w http.ResponseWriter, r *http.Request) {
	// Implementation would handle transaction requests via HTTP
	fmt.Fprintf(w, "Transaction management endpoint")
}
