package exsatfs

import (
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"gotest.tools/assert"
)

func TestABI(t *testing.T) {
	parsedABI, err := abi.JSON(strings.NewReader(ExSatBankFixedStakingYieldSourceAdaptorABI))
	assert.NilError(t, err)

	// Verify 'stakes' method
	method, ok := parsedABI.Methods["stakes"]
	assert.Assert(t, ok, "Method 'stakes' not found in ABI")
	assert.Equal(t, len(method.Inputs), 1)
	assert.Equal(t, len(method.Outputs), 7) // poolId, yieldSourceName, totalStake, stakeTime, unlockTime, state, totalReturn

	// Verify 'unstake' method
	method, ok = parsedABI.Methods["unstake"]
	assert.Assert(t, ok, "Method 'unstake' not found in ABI")
	assert.Equal(t, len(method.Inputs), 1)
	assert.Equal(t, len(method.Outputs), 0)

	// Verify 'unlockStake' method
	method, ok = parsedABI.Methods["unlockStake"]
	assert.Assert(t, ok, "Method 'unlockStake' not found in ABI")
	assert.Equal(t, len(method.Inputs), 1)
	assert.Equal(t, len(method.Outputs), 0)
}
