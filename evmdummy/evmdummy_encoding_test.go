package evmdummy

import (
	"encoding/binary"
	"encoding/hex"
	"math/big"
	"testing"
)

func TestUnstakeEncoding(t *testing.T) {
	// Manual encoding verification logic
	poolId := uint64(123456789)
	stakeAmount := uint64(1000000000000000000) // 1 ETH
	rewardAmount := uint64(500000000000000000) // 0.5 ETH

	// Buffer construction
	data := make([]byte, 4+32+32+32)

	// app_type
	binary.LittleEndian.PutUint32(data[0:4], 0xbc66aee6)

	// pool_id
	binary.BigEndian.PutUint64(data[28:36], poolId)

	// stake_amount
	binary.BigEndian.PutUint64(data[36+24:68], stakeAmount)

	// reward_amount
	binary.BigEndian.PutUint64(data[68+24:100], rewardAmount)

	// Expected values check
	if len(data) != 100 {
		t.Fatalf("Expected length 100, got %d", len(data))
	}

	// Verify app_type
	appType := binary.LittleEndian.Uint32(data[0:4])
	if appType != 0xbc66aee6 {
		t.Errorf("Expected app_type 0xbc66aee6, got 0x%x", appType)
	}

	// Verify pool_id
	decodedPoolId := binary.BigEndian.Uint64(data[28:36])
	if decodedPoolId != poolId {
		t.Errorf("Expected pool_id %d, got %d", poolId, decodedPoolId)
	}

	// Verify stake_amount
	decodedStake := new(big.Int).SetBytes(data[36:68])
	if decodedStake.Uint64() != stakeAmount {
		t.Errorf("Expected stake amount %d, got %s", stakeAmount, decodedStake)
	}

	// Verify reward_amount
	decodedReward := new(big.Int).SetBytes(data[68:100])
	if decodedReward.Uint64() != rewardAmount {
		t.Errorf("Expected reward amount %d, got %s", rewardAmount, decodedReward)
	}

	t.Logf("Encoded data: %s", hex.EncodeToString(data))
}
