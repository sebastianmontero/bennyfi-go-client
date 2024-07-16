package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/sebastianmontero/bennyfi-go-client/yield/source/adaptor/stakelocal"
	"github.com/sebastianmontero/eos-go"
	"github.com/sebastianmontero/eos-go-toolbox/test"
	"github.com/sebastianmontero/eos-go-toolbox/util"
	"gotest.tools/assert"
)

type TestUtil struct {
	stakeLocalClient *stakelocal.StakeLocalContract
	bennyfiContract  eos.AccountName
	*test.TestUtil
}

func NewTestUtil(t *testing.T, stakeLocalClient *stakelocal.StakeLocalContract, bennyfiContract eos.AccountName) *TestUtil {
	return &TestUtil{
		stakeLocalClient: stakeLocalClient,
		bennyfiContract:  bennyfiContract,
		TestUtil:         test.NewTestUtil(t, stakeLocalClient.EOS),
	}
}

func (m *TestUtil) AssertStakeByID(expected *stakelocal.Stake) *stakelocal.Stake {
	actual, err := m.stakeLocalClient.GetStake(expected.RoundID)
	assert.NilError(m.T, err)
	m.AssertStake(actual, expected)
	return actual
}

func (m *TestUtil) AssertStake(actual, expected *stakelocal.Stake) {
	// time.Sleep(time.Millisecond * 500)
	assert.Check(m.T, actual != nil)
	assert.Equal(m.T, actual.RoundID, expected.RoundID)
	assert.Equal(m.T, actual.TotalStake, expected.TotalStake)
	assert.Equal(m.T, actual.TotalReturn, expected.TotalReturn)
	assert.Equal(m.T, actual.State, expected.State)
	assert.DeepEqual(m.T, actual.StakingPeriod, expected.StakingPeriod)
	stakingPeriod := actual.StakingPeriod.AsTimeDuration()
	shift := stakingPeriod + time.Minute
	dateLimit := time.Now().Add(shift * -1)
	stakedTime := actual.StakedTime
	assert.Assert(m.T, !util.IsNullTimePoint(actual.StakedTime))
	assert.Assert(m.T, stakedTime.Time().After(dateLimit), "Expected staked time to be set")
	if actual.State == stakelocal.StateStaked {
		assert.Assert(m.T, util.IsNullTimePoint(actual.StakeEndTime))
	} else {
		assert.Assert(m.T, !util.IsNullTimePoint(actual.StakeEndTime))
	}
}

func (m *TestUtil) AssertYieldSourceByID(expected *stakelocal.YieldSource) *stakelocal.YieldSource {
	actual, err := m.stakeLocalClient.GetYieldSource(expected.YieldSource)
	assert.NilError(m.T, err)
	m.AssertYieldSource(actual, expected)
	return actual
}

func (m *TestUtil) AssertYieldSource(actual, expected *stakelocal.YieldSource) {
	// time.Sleep(time.Millisecond * 500)
	assert.Check(m.T, actual != nil)
	assert.Equal(m.T, actual.YieldSource, expected.YieldSource)
	assert.Equal(m.T, actual.TokenContract, expected.TokenContract)
	assert.Equal(m.T, actual.MinStakingPeriodHrs, expected.MinStakingPeriodHrs)
	assert.Equal(m.T, actual.MaxStakingPeriodHrs, expected.MaxStakingPeriodHrs)
	assert.DeepEqual(m.T, actual.MinStakeAmount, expected.MinStakeAmount)
	assert.DeepEqual(m.T, actual.MaxStakeAmount, expected.MaxStakeAmount)
}

func (m *TestUtil) AssertYieldNotExists(yieldSource eos.Name) {
	actual, err := m.stakeLocalClient.GetYieldSource(yieldSource)
	assert.NilError(m.T, err)
	assert.Assert(m.T, actual == nil)
}

func (m *TestUtil) AssertStakeAction(tokenContract eos.AccountName, roundId uint64, yieldSource eos.Name, quantity eos.Asset, stakingPeriodHrs int64) map[string]interface{} {
	actionData := map[string]interface{}{
		"pool_id":            roundId,
		"yield_source":       yieldSource,
		"quantity":           quantity,
		"staking_period_hrs": stakingPeriodHrs,
	}
	return m.AssertAction(tokenContract, "stake", actionData, 0)
}

func (m *TestUtil) AssertYieldReturnTransfer(tokenContract eos.AccountName, roundId uint64, quantity eos.Asset) map[string]interface{} {
	actionData := map[string]interface{}{
		"from":     m.stakeLocalClient.ContractName,
		"to":       m.bennyfiContract,
		"quantity": quantity,
		"memo":     fmt.Sprintf("pool id: %v", roundId),
	}
	return m.AssertAction(tokenContract, "transfer", actionData, 0)
}
