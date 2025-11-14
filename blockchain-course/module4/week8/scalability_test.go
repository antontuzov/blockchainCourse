package week8

import (
	"testing"
)

func TestNewShardManager(t *testing.T) {
	sm := NewShardManager(4)

	if sm == nil {
		t.Error("Failed to create shard manager")
	}

	if len(sm.Shards) != 4 {
		t.Error("Shard manager should have 4 shards")
	}

	if sm.Beacon == nil {
		t.Error("Beacon chain should be created")
	}

	// Check that each shard is properly initialized
	for i, shard := range sm.Shards {
		if shard.ID != i {
			t.Errorf("Shard %d has incorrect ID", i)
		}

		if len(shard.Nodes) != 0 {
			t.Errorf("Shard %d should have no nodes initially", i)
		}

		if shard.Blockchain == nil {
			t.Errorf("Shard %d should have a blockchain", i)
		}

		if len(shard.Blockchain.Blocks) != 0 {
			t.Errorf("Shard %d blockchain should be empty initially", i)
		}
	}
}

func TestAddRemoveNodeFromShard(t *testing.T) {
	sm := NewShardManager(2)

	// Create a node
	node := &Node{
		ID:      "node1",
		Address: "127.0.0.1:8080",
	}

	// Add node to shard 0
	err := sm.AddNodeToShard(0, node)
	if err != nil {
		t.Errorf("Failed to add node to shard: %s", err)
	}

	// Check that node is added
	if len(sm.Shards[0].Nodes) != 1 {
		t.Error("Node should be added to shard")
	}

	if sm.Shards[0].Nodes[0].ID != "node1" {
		t.Error("Node ID is incorrect")
	}

	// Try to add node to invalid shard
	err = sm.AddNodeToShard(5, node)
	if err == nil {
		t.Error("Should not be able to add node to invalid shard")
	}

	// Remove node from shard
	err = sm.RemoveNodeFromShard(0, "node1")
	if err != nil {
		t.Errorf("Failed to remove node from shard: %s", err)
	}

	// Check that node is removed
	if len(sm.Shards[0].Nodes) != 0 {
		t.Error("Node should be removed from shard")
	}

	// Try to remove node from invalid shard
	err = sm.RemoveNodeFromShard(5, "node1")
	if err == nil {
		t.Error("Should not be able to remove node from invalid shard")
	}

	// Try to remove non-existent node
	err = sm.RemoveNodeFromShard(0, "node2")
	if err == nil {
		t.Error("Should not be able to remove non-existent node")
	}
}

func TestCreatePaymentChannel(t *testing.T) {
	sm := NewShardManager(1)

	// Create a payment channel
	participants := []string{"alice", "bob"}
	pc := sm.CreatePaymentChannel("channel1", participants)

	if pc == nil {
		t.Error("Payment channel should be created")
	}

	if pc.ID != "channel1" {
		t.Error("Payment channel ID is incorrect")
	}

	if len(pc.Participants) != 2 {
		t.Error("Payment channel should have 2 participants")
	}

	if pc.Participants[0] != "alice" || pc.Participants[1] != "bob" {
		t.Error("Payment channel participants are incorrect")
	}

	if len(pc.Balances) != 2 {
		t.Error("Payment channel should have 2 balances")
	}

	if pc.Balances["alice"] != 0 || pc.Balances["bob"] != 0 {
		t.Error("Payment channel balances should be zero initially")
	}

	if len(pc.State) != 0 {
		t.Error("Payment channel state should be empty initially")
	}
}

func TestUpdatePaymentChannel(t *testing.T) {
	sm := NewShardManager(1)

	// Create a payment channel
	participants := []string{"alice", "bob"}
	pc := sm.CreatePaymentChannel("channel1", participants)

	// Update Alice's balance
	err := pc.UpdatePaymentChannel("alice", 100)
	if err != nil {
		t.Errorf("Failed to update payment channel: %s", err)
	}

	// Check that Alice's balance is updated
	if pc.Balances["alice"] != 100 {
		t.Error("Alice's balance should be 100")
	}

	// Check that Bob's balance is still 0
	if pc.Balances["bob"] != 0 {
		t.Error("Bob's balance should still be 0")
	}

	// Check that state is updated
	if len(pc.State) != 1 {
		t.Error("Payment channel should have 1 state")
	}

	// Try to update balance for non-existent participant
	err = pc.UpdatePaymentChannel("charlie", 50)
	if err == nil {
		t.Error("Should not be able to update balance for non-existent participant")
	}
}

func TestClosePaymentChannel(t *testing.T) {
	sm := NewShardManager(1)

	// Create a payment channel
	participants := []string{"alice", "bob"}
	pc := sm.CreatePaymentChannel("channel1", participants)

	// Update balances
	pc.UpdatePaymentChannel("alice", 100)
	pc.UpdatePaymentChannel("bob", 50)

	// Close the payment channel
	balances := pc.ClosePaymentChannel()

	if len(balances) != 2 {
		t.Error("Closed channel should return 2 balances")
	}

	if balances["alice"] != 100 {
		t.Error("Alice's final balance should be 100")
	}

	if balances["bob"] != 50 {
		t.Error("Bob's final balance should be 50")
	}
}

func TestCreateSidechain(t *testing.T) {
	sm := NewShardManager(1)

	// Create a sidechain
	sc := sm.CreateSidechain("sidechain1", "mainchain", "sidechain1")

	if sc == nil {
		t.Error("Sidechain should be created")
	}

	if sc.ID != "sidechain1" {
		t.Error("Sidechain ID is incorrect")
	}

	if sc.Blockchain == nil {
		t.Error("Sidechain should have a blockchain")
	}

	if len(sc.Blockchain.Blocks) != 0 {
		t.Error("Sidechain blockchain should be empty initially")
	}

	if sc.Bridge == nil {
		t.Error("Sidechain should have a bridge")
	}

	if sc.Bridge.SourceChain != "mainchain" {
		t.Error("Bridge source chain is incorrect")
	}

	if sc.Bridge.DestinationChain != "sidechain1" {
		t.Error("Bridge destination chain is incorrect")
	}
}

func TestTransferToSidechain(t *testing.T) {
	sm := NewShardManager(1)

	// Create a sidechain
	sc := sm.CreateSidechain("sidechain1", "mainchain", "sidechain1")

	// Transfer assets to sidechain
	err := sc.TransferToSidechain("ETH", 10, "recipient1")
	if err != nil {
		t.Errorf("Failed to transfer to sidechain: %s", err)
	}

	// Check that a block was added to the sidechain
	if len(sc.Blockchain.Blocks) != 1 {
		t.Error("Sidechain should have 1 block after transfer")
	}

	block := sc.Blockchain.Blocks[0]
	if string(block.Data) != "Transfer 10 ETH to recipient1" {
		t.Error("Block data is incorrect")
	}
}

func TestAddBlockToShard(t *testing.T) {
	sm := NewShardManager(2)

	// Create a block
	block := &Block{
		Index:     0,
		Timestamp: 1234567890,
		Data:      []byte("test block"),
		PrevHash:  []byte("previous_hash"),
		Hash:      []byte("block_hash"),
		Nonce:     12345,
	}

	// Add block to shard 0
	err := sm.AddBlockToShard(0, block)
	if err != nil {
		t.Errorf("Failed to add block to shard: %s", err)
	}

	// Check that block is added
	if len(sm.Shards[0].Blockchain.Blocks) != 1 {
		t.Error("Block should be added to shard")
	}

	// Try to add block to invalid shard
	err = sm.AddBlockToShard(5, block)
	if err == nil {
		t.Error("Should not be able to add block to invalid shard")
	}
}

func TestGetShardInfo(t *testing.T) {
	sm := NewShardManager(2)

	// Create a node and add it to shard 0
	node := &Node{
		ID:      "node1",
		Address: "127.0.0.1:8080",
	}
	sm.AddNodeToShard(0, node)

	// Create a block and add it to shard 0
	block := &Block{
		Index:     0,
		Timestamp: 1234567890,
		Data:      []byte("test block"),
		PrevHash:  []byte("previous_hash"),
		Hash:      []byte("block_hash"),
		Nonce:     12345,
	}
	sm.AddBlockToShard(0, block)

	// Get shard info
	info, err := sm.GetShardInfo(0)
	if err != nil {
		t.Errorf("Failed to get shard info: %s", err)
	}

	if info == nil {
		t.Error("Shard info should be returned")
	}

	if info.ID != 0 {
		t.Error("Shard ID is incorrect")
	}

	if info.NodeCount != 1 {
		t.Error("Shard should have 1 node")
	}

	if info.BlockCount != 1 {
		t.Error("Shard should have 1 block")
	}

	if info.LastBlockID != 0 {
		t.Error("Last block ID is incorrect")
	}

	// Try to get info for invalid shard
	_, err = sm.GetShardInfo(5)
	if err == nil {
		t.Error("Should not be able to get info for invalid shard")
	}
}
