package test

import (
	"testing"

	"github.com/sebastianmontero/bennyfi-go-client/ibc/ibclocal"
	"github.com/sebastianmontero/eos-go"
	"github.com/sebastianmontero/eos-go-toolbox/test"
	"gotest.tools/assert"
)

type TestUtil struct {
	t              *testing.T
	ibcLocalClient *ibclocal.IBCLocalContract
	eosTestUtil    *test.TestUtil
}

func NewTestUtil(t *testing.T, ibcLocalClient *ibclocal.IBCLocalContract) *TestUtil {
	return &TestUtil{
		t:              t,
		ibcLocalClient: ibcLocalClient,
		eosTestUtil:    test.NewTestUtil(t, ibcLocalClient.EOS),
	}
}

func (m *TestUtil) AssertGlobal(expected *ibclocal.Global) {
	actual, err := m.ibcLocalClient.GetGlobal()
	assert.NilError(m.t, err)
	assert.Check(m.t, actual != nil)
	assert.Equal(m.t, actual.ChainId.String(), expected.ChainId.String())
	assert.Equal(m.t, actual.BridgeContract, expected.BridgeContract)
	assert.Equal(m.t, actual.PairedChainId.String(), expected.PairedChainId.String())
	assert.Equal(m.t, actual.PairedWraplockContract, expected.PairedWraplockContract)
	assert.Equal(m.t, actual.PairedTokenContract, expected.PairedTokenContract)
	assert.Equal(m.t, actual.StakeLocalContract, expected.StakeLocalContract)
	assert.Equal(m.t, actual.Enabled, expected.Enabled)
}

func (m *TestUtil) AssertGlobalNotExists() {
	actual, err := m.ibcLocalClient.GetGlobal()
	assert.NilError(m.t, err)
	assert.Assert(m.t, actual == nil)
}

func (m *TestUtil) AssertEmitStakeAction(tokenContract eos.AccountName, emitStakeParams *ibclocal.EmitStakeParams) map[string]interface{} {
	actionData := map[string]interface{}{
		"pool_id":            float64(emitStakeParams.PoolId),
		"yield_source":       emitStakeParams.YieldSource.String(),
		"quantity":           emitStakeParams.Quantity,
		"staking_period_hrs": float64(emitStakeParams.StakingPeriodHrs),
	}
	return m.eosTestUtil.AssertAction(tokenContract, "emitstake", actionData, 0)
}
