package test

import (
	"fmt"
	"testing"

	"github.com/sebastianmontero/bennyfi-go-client/ibc"
	"github.com/sebastianmontero/bennyfi-go-client/ibc/ibcremote"
	"github.com/sebastianmontero/bennyfi-go-client/ibc/test"
	"github.com/sebastianmontero/eos-go"
	"gotest.tools/assert"
)

type TestUtil struct {
	ibcRemoteClient *ibcremote.IBCRemoteContract
	*test.TestUtil
}

func NewTestUtil(t *testing.T, ibcRemoteClient *ibcremote.IBCRemoteContract) *TestUtil {
	return &TestUtil{
		ibcRemoteClient: ibcRemoteClient,
		TestUtil:        test.NewTestUtil(t, ibcRemoteClient.IBCContract),
	}
}

func (m *TestUtil) AssertGlobal(expected *ibcremote.Global) {
	actual, err := m.ibcRemoteClient.GetGlobal()
	assert.NilError(m.T, err)
	assert.Assert(m.T, actual != nil)
	assert.Equal(m.T, actual.ChainId.String(), expected.ChainId.String())
	assert.Equal(m.T, actual.BridgeContract, expected.BridgeContract)
	assert.Equal(m.T, actual.PairedChainId.String(), expected.PairedChainId.String())
	assert.Equal(m.T, actual.Enabled, expected.Enabled)
}

func (m *TestUtil) AssertGlobalNotExists() {
	actual, err := m.ibcRemoteClient.GetGlobal()
	assert.NilError(m.T, err)
	assert.Assert(m.T, actual == nil)
}

func (m *TestUtil) AssertYieldSource(expected *ibcremote.YieldSource) {
	actual, err := m.ibcRemoteClient.GetYieldSource(expected.YieldSource)
	assert.NilError(m.T, err)
	assert.Assert(m.T, actual != nil)
	assert.Equal(m.T, actual.AdaptorContract, expected.AdaptorContract)
}

func (m *TestUtil) AssertContractMapping(expected *ibcremote.ContractMapping) {
	actual, err := m.ibcRemoteClient.GetContractMapping(expected.NativeTokenContract)
	assert.NilError(m.T, err)
	assert.Assert(m.T, actual != nil)
	assert.Equal(m.T, actual.PairedWraptokenContract, expected.PairedWraptokenContract)
}

func (m *TestUtil) AssertReserve(tokenContract eos.AccountName, expected eos.Asset) {
	actual, err := m.ibcRemoteClient.GetReserve(tokenContract, expected.Symbol)
	assert.NilError(m.T, err)
	// fmt.Println("Actual: ", actual, actual != nil)
	assert.Assert(m.T, actual != nil)
	assert.Equal(m.T, actual.Balance, expected)
}

func (m *TestUtil) AssertEmitUnstakeAction(emitUnstakeParams *ibc.EmitUnstakeParams) map[string]interface{} {

	actionData := map[string]interface{}{
		"stake": map[string]interface{}{
			"pool_id":      emitUnstakeParams.PoolId,
			"yield_source": emitUnstakeParams.YieldSource,
			"total_return": emitUnstakeParams.TotalReturn,
		},
	}
	return m.AssertAction(eos.AccountName(m.ibcRemoteClient.IBCContract.ContractName), "emitunstake", actionData, 0)
}

func (m *TestUtil) AssertStakeTransfer(yieldAdaptorContract eos.AccountName, emitStakeParams *ibc.EmitStakeParams) map[string]interface{} {
	actionData := map[string]interface{}{
		"from":     m.ibcRemoteClient.IBCContract.ContractName,
		"to":       yieldAdaptorContract,
		"quantity": emitStakeParams.Quantity.Asset,
		"memo":     fmt.Sprintf("pool id: %v staking period hrs: %v", emitStakeParams.PoolId, emitStakeParams.StakingPeriodHrs),
	}
	return m.AssertAction(eos.AccountName(m.ibcRemoteClient.IBCContract.ContractName), "transfer", actionData, 0)
}
