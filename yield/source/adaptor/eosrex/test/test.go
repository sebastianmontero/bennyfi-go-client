package test

import (
	"fmt"
	"testing"
	"time"

	ltest "github.com/sebastianmontero/bennyfi-go-client/util/test"
	"github.com/sebastianmontero/bennyfi-go-client/yield/source/adaptor/eosrex"
	"github.com/sebastianmontero/eos-go"
	"github.com/sebastianmontero/eos-go-toolbox/dto"
	"github.com/sebastianmontero/eos-go-toolbox/test"
	"github.com/sebastianmontero/eos-go-toolbox/util"
	"gotest.tools/assert"
)

type TestUtil struct {
	eosRexClient     *eosrex.EosRexContract
	rexContract      eos.AccountName
	rexProxyContract eos.AccountName
	*test.TestUtil
}

func NewTestUtil(t *testing.T, eosRexClient *eosrex.EosRexContract, rexContract eos.AccountName, rexProxyContract eos.AccountName) *TestUtil {
	return &TestUtil{
		eosRexClient:     eosRexClient,
		rexContract:      rexContract,
		rexProxyContract: rexProxyContract,
		TestUtil:         test.NewTestUtil(t, eosRexClient.EOS),
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
	shift := stakingPeriod + time.Minute + (time.Second * 30)
	dateLimit := time.Now().Add(shift * -1)
	stakedTime := actual.StakedTime
	stakeEndTime := actual.StakeEndTime
	assert.Assert(m.T, stakedTime.Time().Add(-1*time.Second).Before(time.Now()))
	assert.Assert(m.T, stakedTime.Time().After(dateLimit), "Expected staked time to be set staked time: %v date limit: %v", stakedTime.Time(), dateLimit)
	rexState := actual.RexState
	if actual.RexState == eosrex.RexStateStopped {
		rexState = actual.GetStateWhenStopped()
	}
	if rexState == eosrex.RexStateInSavings {
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
		if rexState == eosrex.RexStateInLockPeriod {
			stakeEndTimeUTC := stakeEndTime.Time().UTC()
			assert.Equal(m.T, maturityTime.Time().UTC(), time.Date(stakeEndTimeUTC.Year(), stakeEndTimeUTC.Month(), stakeEndTimeUTC.Day(), 0, 0, 0, 0, stakeEndTimeUTC.Location()), fmt.Sprintf("Expected maturity time to be start of day UTC of stake end time, maturity time: %v, stake end time: %v", maturityTime.Time().UTC(), stakeEndTime.Time().UTC()))
			assert.Assert(m.T, time.Now().Before(maturityTime.Time()))
		} else if rexState == eosrex.RexStateProceedsCalculated {
			bufferPeriod, err := m.eosRexClient.SettingAsUint32(eosrex.SettingProceedsCalculationBufferPeriodMins)
			assert.NilError(m.T, err)
			assert.Assert(m.T, time.Since(maturityTime.Time().Add(-1*time.Second).Add(-1*time.Minute*time.Duration(bufferPeriod))) >= 0)
		} else if rexState == eosrex.RexStateWithdrawn {
			assert.Assert(m.T, time.Since(stakeEndTime.Time().Add(-1*time.Second)) >= 0)
		} else {
			if !sellDelay {
				assert.Assert(m.T, stakeEndTime.Time().After(time.Now()))
			}
		}
	}
	ltest.AssertAdditionalFields(m.T, actual.AdditionalFields, expected.AdditionalFields)
}

func (m *TestUtil) AssertTREXNotifications(account eos.AccountName, notifications [][]string, verifyExactQuantity bool) {
	m.BaseAssertNotifications(account, "notifytrex", notifications, verifyExactQuantity)
}

func (m *TestUtil) AssertInvestedInRex(stakeAmount, rexAmount eos.Asset) {
	m.assertInvestedInRex(m.rexProxyContract, stakeAmount, rexAmount)
	m.assertInvestedInRex(m.rexContract, stakeAmount, rexAmount)
}

func (m *TestUtil) assertInvestedInRex(contract eos.AccountName, stakeAmount, rexAmount eos.Asset) {
	m.AssertAction(contract, "deposit", map[string]interface{}{
		"owner":  m.eosRexClient.ContractName,
		"amount": stakeAmount.String(),
	}, 0)

	m.AssertAction(contract, "buyrex", map[string]interface{}{
		"from":   m.eosRexClient.ContractName,
		"amount": stakeAmount.String(),
	}, 0)

	m.AssertNumActions(contract, "mvtosavings", 0)
}

func (m *TestUtil) AssertRexMovedFromSavings(rexBalance eos.Asset) {
	m.assertRexMovedFromSavings(m.rexProxyContract, rexBalance)
	m.assertRexMovedFromSavings(m.rexContract, rexBalance)
}

func (m *TestUtil) assertRexMovedFromSavings(contract eos.AccountName, rexBalance eos.Asset) {
	m.AssertAction(contract, "mvfrsavings", map[string]interface{}{
		"owner": m.eosRexClient.ContractName,
		"rex":   rexBalance.String(),
	}, 0)
}

func (m *TestUtil) AssertUpdateRexCalled() {
	// The updaterex action was not added to the vaulta rex contract so it has to be called directly on the rex contract(eosio)
	m.assertUpdateRexCalled(m.rexContract)
}

func (m *TestUtil) assertUpdateRexCalled(contract eos.AccountName) {
	m.AssertAction(contract, "updaterex", map[string]interface{}{
		"owner": m.eosRexClient.ContractName,
	}, 0)
}

func (m *TestUtil) AssertUpdateRexCalledNTimes(times int) {
	m.assertUpdateRexCalledNTimes(m.rexProxyContract, times)
	m.assertUpdateRexCalledNTimes(m.rexContract, times)
}

func (m *TestUtil) assertUpdateRexCalledNTimes(contract eos.AccountName, times int) {
	m.AssertNumActions(contract, "updaterex", times)
}

func (m *TestUtil) AssertRexSold(rexBalance eos.Asset) {
	m.assertRexSold(m.rexProxyContract, rexBalance)
	m.assertRexSold(m.rexContract, rexBalance)
}

func (m *TestUtil) assertRexSold(contract eos.AccountName, rexBalance eos.Asset) {
	m.AssertAction(contract, "sellrex", map[string]interface{}{
		"from": m.eosRexClient.ContractName,
		"rex":  rexBalance.String(),
	}, 0)
}

func (m *TestUtil) AssertRexSellCalledNTimes(times int) {
	m.assertRexSellCalledNTimes(m.rexProxyContract, times)
	m.assertRexSellCalledNTimes(m.rexContract, times)
}

func (m *TestUtil) assertRexSellCalledNTimes(contract eos.AccountName, times int) {
	m.AssertNumActions(contract, "sellrex", times)
}

func (m *TestUtil) AssertRexWithdrawn(amount eos.Asset) {
	m.assertRexWithdrawn(m.rexProxyContract, amount)
	m.assertRexWithdrawn(m.rexContract, amount)
}

func (m *TestUtil) assertRexWithdrawn(contract eos.AccountName, amount eos.Asset) {
	m.AssertAction(contract, "withdraw", map[string]interface{}{
		"owner":  m.eosRexClient.ContractName,
		"amount": amount.String(),
	}, 0)
}

func (m *TestUtil) AssertRexWithdrawCalledNTimes(times int) {
	m.assertRexWithdrawCalledNTimes(m.rexProxyContract, times)
	m.assertRexWithdrawCalledNTimes(m.rexContract, times)
}

func (m *TestUtil) assertRexWithdrawCalledNTimes(contract eos.AccountName, times int) {
	m.AssertNumActions(contract, "withdraw", times)
}

func (m *TestUtil) LapseLastNotifiedTime(lastNotifiedSetting string, shift time.Duration) {
	notificationPeriodMins, err := m.eosRexClient.SettingAsUint32(eosrex.SettingNotificationPeriodMins)
	assert.NilError(m.T, err)
	lastNotified := time.Now().Add(time.Duration(notificationPeriodMins)*time.Minute*-1 + shift)
	_, err = m.eosRexClient.SetSetting(eos.AccountName(m.eosRexClient.ContractName), lastNotifiedSetting, dto.FlexValueFromTime(lastNotified))
	assert.NilError(m.T, err)
}
