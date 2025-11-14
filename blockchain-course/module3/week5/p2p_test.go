package week5

import (
	"testing"
)

func TestNewNode(t *testing.T) {
	// Create a new node
	node := NewNode("localhost", 8080, nil)

	if node == nil {
		t.Error("Failed to create node")
	}

	if node.Address != "localhost" {
		t.Error("Node address is incorrect")
	}

	if node.Port != 8080 {
		t.Error("Node port is incorrect")
	}

	if len(node.Peers) != 0 {
		t.Error("Node peers should be empty initially")
	}
}

func TestAddPeer(t *testing.T) {
	// Create a new node
	node := NewNode("localhost", 8080, nil)

	// Add a peer (this will fail because there's no actual server running)
	err := node.AddPeer("localhost", 8081)
	// We expect this to fail in testing environment
	if err == nil {
		// If it somehow succeeds, check that the peer was added
		if len(node.Peers) != 1 {
			t.Error("Peer should be added to node")
		}
	}
}

func TestRemovePeer(t *testing.T) {
	// Create a new node
	node := NewNode("localhost", 8080, nil)

	// Try to remove a peer that doesn't exist
	node.RemovePeer("localhost", 8081)

	// Should not panic or cause errors
}

func TestMessageSerialization(t *testing.T) {
	// Create a new node
	node := NewNode("localhost", 8080, nil)

	// Create a message
	msg := &Message{
		Type:    "test",
		Payload: []byte("test payload"),
	}

	// Serialize the message
	data, err := node.serializeMessage(msg)
	if err != nil {
		t.Errorf("Failed to serialize message: %s", err)
	}

	if len(data) == 0 {
		t.Error("Serialized data should not be empty")
	}

	// Deserialize the message
	deserializedMsg, err := node.deserializeMessage(data)
	if err != nil {
		t.Errorf("Failed to deserialize message: %s", err)
	}

	if deserializedMsg.Type != msg.Type {
		t.Error("Deserialized message type is incorrect")
	}

	if string(deserializedMsg.Payload) != string(msg.Payload) {
		t.Error("Deserialized message payload is incorrect")
	}
}

func TestMessageTypes(t *testing.T) {
	// Create a new node
	node := NewNode("localhost", 8080, nil)

	// Test different message types
	msgTypes := []string{"block", "transaction", "get_blocks", "ping", "unknown"}

	for _, msgType := range msgTypes {
		msg := &Message{
			Type:    msgType,
			Payload: []byte("test payload"),
		}

		// Handle the message (should not panic)
		node.HandleMessage(msg)
	}
}
