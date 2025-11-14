package week7

import (
	"testing"
)

func TestNewPermissionedBlockchain(t *testing.T) {
	pb := NewPermissionedBlockchain()

	if pb == nil {
		t.Error("Failed to create permissioned blockchain")
	}

	if len(pb.Channels) != 0 {
		t.Error("Channels should be empty initially")
	}

	if len(pb.Members) != 0 {
		t.Error("Members should be empty initially")
	}

	if pb.CA == nil {
		t.Error("Certificate authority should be created")
	}

	if len(pb.Policies) != 0 {
		t.Error("Policies should be empty initially")
	}
}

func TestNewCertificateAuthority(t *testing.T) {
	ca := NewCertificateAuthority()

	if ca == nil {
		t.Error("Failed to create certificate authority")
	}

	if ca.Certificate == nil {
		t.Error("Certificate should be created")
	}

	if ca.PrivateKey == nil {
		t.Error("Private key should be created")
	}

	if ca.PublicKey == nil {
		t.Error("Public key should be created")
	}
}

func TestRegisterMember(t *testing.T) {
	pb := NewPermissionedBlockchain()

	// Register a member
	member, err := pb.RegisterMember("member1", "Test Member", "user")
	if err != nil {
		t.Errorf("Failed to register member: %s", err)
	}

	if member == nil {
		t.Error("Member should be created")
	}

	if member.ID != "member1" {
		t.Error("Member ID is incorrect")
	}

	if member.Name != "Test Member" {
		t.Error("Member name is incorrect")
	}

	if member.Role != "user" {
		t.Error("Member role is incorrect")
	}

	if member.Certificate == nil {
		t.Error("Member certificate should be created")
	}

	if member.PrivateKey == nil {
		t.Error("Member private key should be created")
	}

	if member.PublicKey == nil {
		t.Error("Member public key should be created")
	}

	// Check that member is added to the blockchain
	if len(pb.Members) != 1 {
		t.Error("Member should be added to the blockchain")
	}

	// Try to register the same member again
	_, err = pb.RegisterMember("member1", "Test Member", "user")
	if err == nil {
		t.Error("Should not be able to register the same member twice")
	}
}

func TestCreateChannel(t *testing.T) {
	pb := NewPermissionedBlockchain()

	// Register a member
	_, err := pb.RegisterMember("member1", "Test Member", "user")
	if err != nil {
		t.Errorf("Failed to register member: %s", err)
	}

	// Create a channel with the member
	channel, err := pb.CreateChannel("channel1", []string{"member1"})
	if err != nil {
		t.Errorf("Failed to create channel: %s", err)
	}

	if channel == nil {
		t.Error("Channel should be created")
	}

	if channel.Name != "channel1" {
		t.Error("Channel name is incorrect")
	}

	if len(channel.Members) != 1 {
		t.Error("Channel should have one member")
	}

	// Check that the member is in the channel
	if !channel.Members["member1"] {
		t.Error("Member should be in the channel")
	}

	// Check that channel is added to the blockchain
	if len(pb.Channels) != 1 {
		t.Error("Channel should be added to the blockchain")
	}

	// Try to create the same channel again
	_, err = pb.CreateChannel("channel1", []string{"member1"})
	if err == nil {
		t.Error("Should not be able to create the same channel twice")
	}

	// Try to create a channel with a non-existent member
	_, err = pb.CreateChannel("channel2", []string{"member2"})
	if err == nil {
		t.Error("Should not be able to create a channel with a non-existent member")
	}
}

func TestAddRemoveMemberFromChannel(t *testing.T) {
	pb := NewPermissionedBlockchain()

	// Register two members
	_, err := pb.RegisterMember("member1", "Test Member 1", "user")
	if err != nil {
		t.Errorf("Failed to register member1: %s", err)
	}

	_, err = pb.RegisterMember("member2", "Test Member 2", "user")
	if err != nil {
		t.Errorf("Failed to register member2: %s", err)
	}

	// Create a channel with the first member
	_, err = pb.CreateChannel("channel1", []string{"member1"})
	if err != nil {
		t.Errorf("Failed to create channel: %s", err)
	}

	// Add the second member to the channel
	err = pb.AddMemberToChannel("channel1", "member2")
	if err != nil {
		t.Errorf("Failed to add member to channel: %s", err)
	}

	// Check that both members are in the channel
	channel := pb.Channels["channel1"]
	if len(channel.Members) != 2 {
		t.Error("Channel should have two members")
	}

	if !channel.Members["member1"] {
		t.Error("Member1 should be in the channel")
	}

	if !channel.Members["member2"] {
		t.Error("Member2 should be in the channel")
	}

	// Remove the second member from the channel
	err = pb.RemoveMemberFromChannel("channel1", "member2")
	if err != nil {
		t.Errorf("Failed to remove member from channel: %s", err)
	}

	// Check that only the first member is in the channel
	if len(channel.Members) != 1 {
		t.Error("Channel should have one member")
	}

	if !channel.Members["member1"] {
		t.Error("Member1 should be in the channel")
	}

	if channel.Members["member2"] {
		t.Error("Member2 should not be in the channel")
	}

	// Try to remove a member from a non-existent channel
	err = pb.RemoveMemberFromChannel("channel2", "member1")
	if err == nil {
		t.Error("Should not be able to remove member from non-existent channel")
	}

	// Try to remove a non-existent member from a channel
	err = pb.RemoveMemberFromChannel("channel1", "member3")
	if err == nil {
		t.Error("Should not be able to remove non-existent member from channel")
	}
}

func TestCreatePolicy(t *testing.T) {
	pb := NewPermissionedBlockchain()

	// Create a policy
	rules := []PolicyRule{
		{
			Role:        "admin",
			Permissions: []string{"read", "write", "delete"},
		},
		{
			Role:        "user",
			Permissions: []string{"read"},
		},
	}

	policy := pb.CreatePolicy("policy1", "Test Policy", "A test policy", rules)

	if policy == nil {
		t.Error("Policy should be created")
	}

	if policy.ID != "policy1" {
		t.Error("Policy ID is incorrect")
	}

	if policy.Name != "Test Policy" {
		t.Error("Policy name is incorrect")
	}

	if policy.Description != "A test policy" {
		t.Error("Policy description is incorrect")
	}

	if len(policy.Rules) != 2 {
		t.Error("Policy should have two rules")
	}

	// Check that policy is added to the blockchain
	if len(pb.Policies) != 1 {
		t.Error("Policy should be added to the blockchain")
	}
}

func TestApplyPolicyToChannel(t *testing.T) {
	pb := NewPermissionedBlockchain()

	// Register a member
	_, err := pb.RegisterMember("member1", "Test Member", "user")
	if err != nil {
		t.Errorf("Failed to register member: %s", err)
	}

	// Create a channel
	_, err = pb.CreateChannel("channel1", []string{"member1"})
	if err != nil {
		t.Errorf("Failed to create channel: %s", err)
	}

	// Create a policy
	rules := []PolicyRule{
		{
			Role:        "admin",
			Permissions: []string{"read", "write", "delete"},
		},
	}

	pb.CreatePolicy("policy1", "Test Policy", "A test policy", rules)

	// Apply the policy to the channel
	err = pb.ApplyPolicyToChannel("channel1", "policy1")
	if err != nil {
		t.Errorf("Failed to apply policy to channel: %s", err)
	}

	// Check that the policy is applied to the channel
	channel := pb.Channels["channel1"]
	if len(channel.Policies) != 1 {
		t.Error("Channel should have one policy")
	}

	if channel.Policies[0].ID != "policy1" {
		t.Error("Policy ID is incorrect")
	}

	// Try to apply a non-existent policy to a channel
	err = pb.ApplyPolicyToChannel("channel1", "policy2")
	if err == nil {
		t.Error("Should not be able to apply non-existent policy to channel")
	}

	// Try to apply a policy to a non-existent channel
	err = pb.ApplyPolicyToChannel("channel2", "policy1")
	if err == nil {
		t.Error("Should not be able to apply policy to non-existent channel")
	}
}

func TestVerifyMember(t *testing.T) {
	pb := NewPermissionedBlockchain()

	// Register a member
	_, err := pb.RegisterMember("member1", "Test Member", "user")
	if err != nil {
		t.Errorf("Failed to register member: %s", err)
	}

	// Verify the member
	valid := pb.VerifyMember("member1")
	if !valid {
		t.Error("Member should be valid")
	}

	// Try to verify a non-existent member
	valid = pb.VerifyMember("member2")
	if valid {
		t.Error("Non-existent member should not be valid")
	}
}

func TestGetChannelMembers(t *testing.T) {
	pb := NewPermissionedBlockchain()

	// Register two members
	_, err := pb.RegisterMember("member1", "Test Member 1", "user")
	if err != nil {
		t.Errorf("Failed to register member1: %s", err)
	}

	_, err = pb.RegisterMember("member2", "Test Member 2", "user")
	if err != nil {
		t.Errorf("Failed to register member2: %s", err)
	}

	// Create a channel with both members
	_, err = pb.CreateChannel("channel1", []string{"member1", "member2"})
	if err != nil {
		t.Errorf("Failed to create channel: %s", err)
	}

	// Get the channel members
	members, err := pb.GetChannelMembers("channel1")
	if err != nil {
		t.Errorf("Failed to get channel members: %s", err)
	}

	if len(members) != 2 {
		t.Error("Should have two members in the channel")
	}

	// Try to get members from a non-existent channel
	_, err = pb.GetChannelMembers("channel2")
	if err == nil {
		t.Error("Should not be able to get members from non-existent channel")
	}
}
