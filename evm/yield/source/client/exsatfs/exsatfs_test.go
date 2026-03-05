package exsatfs

import (
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"gotest.tools/assert"
)

func TestABI(t *testing.T) {
	parsedABI, err := abi.JSON(strings.NewReader(IFixedStakingABI))
	assert.NilError(t, err)

	// Verify 'getPosition' method
	method, ok := parsedABI.Methods["getPosition"]
	assert.Assert(t, ok, "Method 'getPosition' not found in ABI")
	assert.Equal(t, len(method.Inputs), 2) // agent, subId
	assert.Equal(t, len(method.Outputs), 8) // Flattened outputs

	// Verify Input types
	assert.Equal(t, method.Inputs[0].Type.String(), "address")
	assert.Equal(t, method.Inputs[1].Type.String(), "bytes32")

	// Verify 'stakingToken' method
	method, ok = parsedABI.Methods["stakingToken"]
	assert.Assert(t, ok, "Method 'stakingToken' not found in ABI")
	assert.Equal(t, len(method.Inputs), 0)
	assert.Equal(t, len(method.Outputs), 1)
	assert.Equal(t, method.Outputs[0].Type.String(), "address")
}

func TestSubIdConversion(t *testing.T) {
	// Manual verification of the logic used in GetPosition
	subId := uint64(1)
	var subIdBytes [32]byte
	bigSubId := new(big.Int).SetUint64(subId)
	bigSubIdBytes := bigSubId.Bytes()
	copy(subIdBytes[32-len(bigSubIdBytes):], bigSubIdBytes)

	// Expect 000...0001 (Big Endian)
	assert.Equal(t, subIdBytes[31], byte(1))
	assert.Equal(t, subIdBytes[0], byte(0))

	// Test case for mapping back to hex (common.Hash)
	hash := common.BytesToHash(subIdBytes[:])
	assert.Equal(t, hash.Big().Uint64(), subId)
}
