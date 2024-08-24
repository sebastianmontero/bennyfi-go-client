package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/sebastianmontero/bennyfi-go-client/yield/source/adaptor/eosrex"
	"github.com/sebastianmontero/eos-go"
	"github.com/sebastianmontero/eos-go-toolbox/test"
	"github.com/sebastianmontero/eos-go-toolbox/util"
	"gotest.tools/assert"
)

type TestUtil struct {
	eosRexClient *eosrex.EosRexContract
	rexContract  eos.AccountName
	*test.TestUtil
}

func NewTestUtil(t *testing.T, eosRexClient *eosrex.EosRexContract, rexContract eos.AccountName) *TestUtil {
	return &TestUtil{
		eosRexClient: eosRexClient,
		rexContract:  rexContract,
		TestUtil:     test.NewTestUtil(t, eosRexClient.EOS),
	}
}

func (m *TestUtil) AssertStakeByID(expected *eosrex.Stake, sellDelay bool) *eosrex.Stake {
	actual, err := m.eosRexClient.GetStake(expected.RoundID)
	assert.NilError(m.T, err)
	m.AssertStake(actual, expected, sellDelay)
	return actual

}

func (m *TestUtil) AssertStake(actual, expected *eosrex.Stake, sellDelay bool) {
	// time.Sleep(time.Millisecond * 500)
	assert.Check(m.T, actual != nil)
	assert.Equal(m.T, actual.RoundID, expected.RoundID)
	assert.Equal(m.T, actual.TotalStake, expected.TotalStake)
	assert.Equal(m.T, actual.RexBalance, expected.RexBalance)
	assert.Equal(m.T, actual.TotalReturn, expected.TotalReturn)
	assert.Equal(m.T, actual.RexState, expected.RexState)
	assert.DeepEqual(m.T, actual.StakingPeriod, expected.StakingPeriod)
	stakingPeriod := actual.StakingPeriod.AsTimeDuration()
	shift := stakingPeriod + time.Minute
	dateLimit := time.Now().Add(shift * -1)
	stakedTime := actual.StakedTime
	stakeEndTime := actual.StakeEndTime
	assert.Assert(m.T, stakedTime.Time().Add(-1*time.Second).Before(time.Now()))
	assert.Assert(m.T, stakedTime.Time().After(dateLimit), "Expected staked time to be set")
	if actual.RexState == eosrex.RexStateInSavings {
		assert.Assert(m.T, util.IsNullTimePoint(actual.MovedFromSavingsTime), "Expected moved from savings time not to be set")
		assert.Assert(m.T, util.IsNullTimePoint(actual.MaturityTime), "Expected maturity time not to be set")
		assert.Assert(m.T, stakeEndTime.Time().After(time.Now()))
		assert.Assert(m.T, stakeEndTime.Time().Equal(stakedTime.Time().Add(stakingPeriod)), "Expected stake end time to be set and be staking period apart from staked time")
	} else {
		movedFromSavingsTime := actual.MovedFromSavingsTime
		maturityTime := actual.MaturityTime
		rexLockPeriod := time.Duration(int64(eosrex.RexLockPeriodDays) * int64(time.Hour) * 24)
		expectedMoveFromSavingsTime := stakeEndTime.Time().Add(rexLockPeriod * -1)
		assert.Equal(m.T, movedFromSavingsTime.Time(), expectedMoveFromSavingsTime)
		assert.Equal(m.T, movedFromSavingsTime.Time().Add(rexLockPeriod), stakeEndTime.Time())
		assert.Assert(m.T, time.Since(movedFromSavingsTime.Time().Add(-1*time.Second)) >= 0)
		if actual.RexState == eosrex.RexStateInLockPeriod {
			stakeEndTimeUTC := stakeEndTime.Time().UTC()
			assert.Equal(m.T, maturityTime.Time().UTC(), time.Date(stakeEndTimeUTC.Year(), stakeEndTimeUTC.Month(), stakeEndTimeUTC.Day(), 0, 0, 0, 0, stakeEndTimeUTC.Location()), fmt.Sprintf("Expected maturity time to be start of day UTC of stake end time, maturity time: %v, stake end time: %v", maturityTime.Time().UTC(), stakeEndTime.Time().UTC()))
			assert.Assert(m.T, time.Now().Before(maturityTime.Time()))
		} else if actual.RexState == eosrex.RexStateProceedsCalculated {
			bufferPeriod, err := m.eosRexClient.SettingAsUint32(eosrex.SettingProceedsCalculationBufferPeriodMins)
			assert.NilError(m.T, err)
			assert.Assert(m.T, time.Since(maturityTime.Time().Add(-1*time.Second).Add(-1*time.Minute*time.Duration(bufferPeriod))) >= 0)
		} else if actual.RexState == eosrex.RexStateWithdrawn {
			assert.Assert(m.T, time.Since(stakeEndTime.Time().Add(-1*time.Second)) >= 0)
		} else {
			if !sellDelay {
				assert.Assert(m.T, stakeEndTime.Time().After(time.Now()))
			}
		}
	}
}

func (m *TestUtil) AssertTREXNotifications(account eos.AccountName, notifications [][]string, verifyExactQuantity bool) {
	m.BaseAssertNotifications(account, "notifytrex", notifications, verifyExactQuantity)
}

func (m *TestUtil) AssertInvestedInRex(stakeAmount, rexAmount eos.Asset) {
	m.AssertAction(m.rexContract, "deposit", map[string]interface{}{
		"owner":  m.eosRexClient.ContractName,
		"amount": stakeAmount.String(),
	}, 0)

	m.AssertAction(m.rexContract, "buyrex", map[string]interface{}{
		"from":   m.eosRexClient.ContractName,
		"amount": stakeAmount.String(),
	}, 0)

	m.AssertNumActions(m.rexContract, "mvtosavings", 0)
}

func (m *TestUtil) AssertRexMovedFromSavings(rexBalance eos.Asset) {
	m.AssertAction(m.rexContract, "mvfrsavings", map[string]interface{}{
		"owner": m.eosRexClient.ContractName,
		"rex":   rexBalance.String(),
	}, 0)
}

func (m *TestUtil) AssertRexSold(rexBalance eos.Asset) {
	m.AssertAction(m.rexContract, "sellrex", map[string]interface{}{
		"from": m.eosRexClient.ContractName,
		"rex":  rexBalance.String(),
	}, 0)
}

func (m *TestUtil) AssertRexWithdrawn(amount eos.Asset) {
	m.AssertAction(m.rexContract, "withdraw", map[string]interface{}{
		"owner":  m.eosRexClient.ContractName,
		"amount": amount.String(),
	}, 0)
}
