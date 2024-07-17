package test

import (
	"fmt"
	"testing"

	"github.com/sebastianmontero/bennyfi-go-client/ibc"
	"github.com/sebastianmontero/bennyfi-go-client/ibc/ibclocal"
	"github.com/sebastianmontero/bennyfi-go-client/ibc/test"
	"github.com/sebastianmontero/eos-go"
	"gotest.tools/assert"
)

type TestUtil struct {
	ibcLocalClient *ibclocal.IBCLocalContract
	*test.TestUtil
}

func NewTestUtil(t *testing.T, ibcLocalClient *ibclocal.IBCLocalContract) *TestUtil {
	return &TestUtil{
		ibcLocalClient: ibcLocalClient,
		TestUtil:       test.NewTestUtil(t, ibcLocalClient.IBCContract),
	}
}

func (m *TestUtil) AssertGlobal(expected *ibclocal.Global) {
	actual, err := m.ibcLocalClient.GetGlobal()
	assert.NilError(m.T, err)
	assert.Check(m.T, actual != nil)
	assert.Equal(m.T, actual.ChainId.String(), expected.ChainId.String())
	assert.Equal(m.T, actual.BridgeContract, expected.BridgeContract)
	assert.Equal(m.T, actual.PairedChainId.String(), expected.PairedChainId.String())
	assert.Equal(m.T, actual.PairedWraplockContract, expected.PairedWraplockContract)
	assert.Equal(m.T, actual.PairedTokenContract, expected.PairedTokenContract)
	assert.Equal(m.T, actual.StakeLocalContract, expected.StakeLocalContract)
	assert.Equal(m.T, actual.Enabled, expected.Enabled)
}

func (m *TestUtil) AssertGlobalNotExists() {
	actual, err := m.ibcLocalClient.GetGlobal()
	assert.NilError(m.T, err)
	assert.Assert(m.T, actual == nil)
}

func (m *TestUtil) AssertEmitStakeAction(emitStakeParams *ibc.EmitStakeParams) map[string]interface{} {

	actionData := map[string]interface{}{
		"stake": map[string]interface{}{
			"pool_id":            emitStakeParams.PoolId,
			"yield_source":       emitStakeParams.YieldSource,
			"quantity":           emitStakeParams.Quantity,
			"staking_period_hrs": emitStakeParams.StakingPeriodHrs,
		},
	}
	return m.AssertAction(eos.AccountName(m.ibcLocalClient.IBCContract.ContractName), "emitstake", actionData, 0)
}

func (m *TestUtil) AssertYieldReturnTransfer(stakeLocalContract eos.AccountName, roundId uint64, quantity eos.Asset) map[string]interface{} {
	actionData := map[string]interface{}{
		"from":     m.ibcLocalClient.IBCContract.ContractName,
		"to":       stakeLocalContract,
		"quantity": quantity,
		"memo":     fmt.Sprintf("pool id: %v", roundId),
	}
	return m.AssertAction(eos.AccountName(m.ibcLocalClient.IBCContract.ContractName), "transfer", actionData, 0)
}
