package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/sebastianmontero/bennyfi-go-client/yield/source/adaptor/stakelocal"
	"github.com/sebastianmontero/eos-go"
	"github.com/sebastianmontero/eos-go-toolbox/dto"
	"github.com/sebastianmontero/eos-go-toolbox/test"
	"github.com/sebastianmontero/eos-go-toolbox/util"
	"gotest.tools/assert"
)

type TestUtil struct {
	t                *testing.T
	stakeLocalClient *stakelocal.StakeLocalContract
	bennyfiContract  eos.AccountName
	eosTestUtil      *test.TestUtil
}

func NewTestUtil(t *testing.T, stakeLocalClient *stakelocal.StakeLocalContract, bennyfiContract eos.AccountName) *TestUtil {
	return &TestUtil{
		t:                t,
		stakeLocalClient: stakeLocalClient,
		bennyfiContract:  bennyfiContract,
		eosTestUtil:      test.NewTestUtil(t, stakeLocalClient.EOS),
	}
}

func (m *TestUtil) AssertStakeByID(expected *stakelocal.Stake) *stakelocal.Stake {
	actual, err := m.stakeLocalClient.GetStake(expected.RoundID)
	assert.NilError(m.t, err)
	m.AssertStake(actual, expected)
	return actual
}

func (m *TestUtil) AssertStake(actual, expected *stakelocal.Stake) {
	// time.Sleep(time.Millisecond * 500)
	assert.Check(m.t, actual != nil)
	assert.Equal(m.t, actual.RoundID, expected.RoundID)
	assert.Equal(m.t, actual.TotalStake, expected.TotalStake)
	assert.Equal(m.t, actual.TotalReturn, expected.TotalReturn)
	assert.Equal(m.t, actual.State, expected.State)
	assert.DeepEqual(m.t, actual.StakingPeriod, expected.StakingPeriod)
	stakingPeriod := actual.StakingPeriod.AsTimeDuration()
	shift := stakingPeriod + time.Minute
	dateLimit := time.Now().Add(shift * -1)
	stakedTime := actual.StakedTime
	assert.Assert(m.t, !util.IsNullTimePoint(actual.StakedTime))
	assert.Assert(m.t, stakedTime.Time().After(dateLimit), "Expected staked time to be set")
	if actual.State == stakelocal.StateStaked {
		assert.Assert(m.t, util.IsNullTimePoint(actual.StakeEndTime))
	} else {
		assert.Assert(m.t, !util.IsNullTimePoint(actual.StakeEndTime))
	}
}

func (m *TestUtil) AssertYieldSourceByID(expected *stakelocal.YieldSource) *stakelocal.YieldSource {
	actual, err := m.stakeLocalClient.GetYieldSource(expected.YieldSource)
	assert.NilError(m.t, err)
	m.AssertYieldSource(actual, expected)
	return actual
}

func (m *TestUtil) AssertYieldSource(actual, expected *stakelocal.YieldSource) {
	// time.Sleep(time.Millisecond * 500)
	assert.Check(m.t, actual != nil)
	assert.Equal(m.t, actual.YieldSource, expected.YieldSource)
	assert.Equal(m.t, actual.TokenContract, expected.TokenContract)
	assert.Equal(m.t, actual.MinStakingPeriodHrs, expected.MinStakingPeriodHrs)
	assert.Equal(m.t, actual.MaxStakingPeriodHrs, expected.MaxStakingPeriodHrs)
	assert.DeepEqual(m.t, actual.MinStakeAmount, expected.MinStakeAmount)
	assert.DeepEqual(m.t, actual.MaxStakeAmount, expected.MaxStakeAmount)
}

func (m *TestUtil) AssertYieldNotExists(yieldSource eos.Name) {
	actual, err := m.stakeLocalClient.GetYieldSource(yieldSource)
	assert.NilError(m.t, err)
	assert.Assert(m.t, actual == nil)
}

func (m *TestUtil) AssertStakeAction(tokenContract eos.AccountName, roundId uint64, yieldSource eos.Name, quantity eos.Asset, stakingPeriod *dto.Microseconds) map[string]interface{} {
	actionData := map[string]interface{}{
		"pool_id":        float64(roundId),
		"yield_source":   yieldSource.String(),
		"quantity":       quantity.String(),
		"staking_period": stakingPeriod.ToMap(),
	}
	return m.eosTestUtil.AssertAction(tokenContract, "stake", actionData, 0)
}

func (m *TestUtil) AssertYieldReturnTransfer(tokenContract eos.AccountName, roundId uint64, quantity eos.Asset) map[string]interface{} {
	actionData := map[string]interface{}{
		"from":     m.stakeLocalClient.ContractName,
		"to":       m.bennyfiContract.String(),
		"quantity": quantity.String(),
		"memo":     fmt.Sprintf("pool id: %v", roundId),
	}
	return m.eosTestUtil.AssertAction(tokenContract, "transfer", actionData, 0)
}
