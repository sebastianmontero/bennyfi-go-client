package tlosrex

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/eoscanada/eos-go"
	"github.com/sebastianmontero/eos-go-toolbox/dto"
	"github.com/sebastianmontero/eos-go-toolbox/util"
)

var RexLockPeriodDays = 5

type Stake struct {
	RoundID              uint64            `json:"round_id"`
	TotalStake           string            `json:"total_stake"`
	RexBalance           string            `json:"rex_balance"`
	TotalReturn          string            `json:"total_return"`
	RexState             eos.Name          `json:"rex_state"`
	StakingPeriod        *dto.Microseconds `json:"staking_period"`
	StakedTime           string            `json:"staked_time"`
	MovedFromSavingsTime string            `json:"moved_from_savings_time"`
	StakeEndTime         string            `json:"stake_end_time"`
	UpdatedDate          string            `json:"updated_date"`
}

func (m *Stake) String() string {
	result, err := json.Marshal(m)
	if err != nil {
		panic(fmt.Sprintf("Failed marshalling round: %v", err))
	}
	return string(result)
}

func (m *Stake) GetTotalStake() eos.Asset {
	totalStake, err := eos.NewAssetFromString(m.TotalStake)
	if err != nil {
		panic(fmt.Sprintf("Unable to parse total stake: %v to asset", m.TotalStake))
	}
	return totalStake
}

func (m *Stake) GetRexBalance() eos.Asset {
	rexBalance, err := eos.NewAssetFromString(m.RexBalance)
	if err != nil {
		panic(fmt.Sprintf("Unable to parse rex balance: %v to asset", m.RexBalance))
	}
	return rexBalance
}

func (m *Stake) GetTotalReturn() eos.Asset {
	totalReturn, err := eos.NewAssetFromString(m.TotalReturn)
	if err != nil {
		panic(fmt.Sprintf("Unable to parse total return: %v to asset", m.TotalReturn))
	}
	return totalReturn
}

func (m *Stake) GetStakedTime() time.Time {
	stakedTime, err := util.ToTime(m.StakedTime)
	if err != nil {
		panic(fmt.Sprintf("Unable to parse stakedTime: %v to asset", m.StakedTime))
	}
	return stakedTime
}

func (m *Stake) GetMovedFromSavingsTime() time.Time {
	movedFromSavingsTime, err := util.ToTime(m.MovedFromSavingsTime)
	if err != nil {
		panic(fmt.Sprintf("Unable to parse movedFromSavingsTime: %v to asset", m.MovedFromSavingsTime))
	}
	return movedFromSavingsTime
}

func (m *Stake) GetStakeEndTime() time.Time {
	stakeEndTime, err := util.ToTime(m.StakeEndTime)
	if err != nil {
		panic(fmt.Sprintf("Unable to parse stakeEndTime: %v to asset", m.StakeEndTime))
	}
	return stakeEndTime
}

func (m *Stake) CalculateRexLockPeriodTime() time.Time {
	stakedTime, err := util.ToTime(m.StakedTime)
	if err != nil {
		panic(fmt.Sprintf("failed to calculate Rex Lock Period Time, could not parse staked time: %v, error: %v", m.StakedTime, err))
	}
	return stakedTime.Add(time.Hour * time.Duration(m.StakingPeriod.Hrs()-int64(24*RexLockPeriodDays)))
}

func (m *Stake) CalculateSellRexTime() time.Time {
	movedFromSavingsTime, err := util.ToTime(m.MovedFromSavingsTime)
	if err != nil {
		panic(fmt.Sprintf("failed to calculate Sell Rex Time, could not parse moved from savings time: %v, error: %v", m.MovedFromSavingsTime, err))
	}
	return movedFromSavingsTime.Add(time.Hour * time.Duration(24*RexLockPeriodDays))
}

func (m *Stake) CalculateStakeEndTime() time.Time {
	stakedTime, err := util.ToTime(m.StakedTime)
	if err != nil {
		panic(fmt.Sprintf("failed to calculate Unlock Time, could not parse staked time: %v, error: %v", m.StakedTime, err))
	}
	return stakedTime.Add(time.Hour * time.Duration(m.StakingPeriod.Hrs()))
}

func (m *TlosRexContract) CheckStakeParameters(authorizer, tokenContract, stakeAmount interface{}, stakingPeriodHrs uint32) (string, error) {
	actionData := make(map[string]interface{})
	actionData["token_contract"] = tokenContract
	actionData["stake_amount"] = stakeAmount
	actionData["staking_period_hrs"] = stakingPeriodHrs
	return m.ExecAction(authorizer, "chckstkparam", actionData)
}

func (m *TlosRexContract) MoveRoundFromSavings(roundId uint64) (string, error) {
	actionData := make(map[string]interface{})
	actionData["round_id"] = roundId
	return m.ExecAction(fmt.Sprintf("%v@open", m.ContractName), "mvfrmsvngsrn", actionData)
}

func (m *TlosRexContract) SellRoundRex(roundId uint64) (string, error) {
	actionData := make(map[string]interface{})
	actionData["round_id"] = roundId
	return m.ExecAction(fmt.Sprintf("%v@open", m.ContractName), "sellrexrn", actionData)
}

func (m *TlosRexContract) WithdrawRoundRex(roundId uint64) (string, error) {
	actionData := make(map[string]interface{})
	actionData["round_id"] = roundId
	return m.ExecAction(fmt.Sprintf("%v@open", m.ContractName), "withdrwrexrn", actionData)
}

func (m *TlosRexContract) MoveFromSavings(callCounter uint64) (string, error) {
	actionData := make(map[string]interface{})
	actionData["call_counter"] = callCounter
	return m.ExecAction(fmt.Sprintf("%v@open", m.ContractName), "mvfrmsavings", actionData)
}

func (m *TlosRexContract) SellRex(callCounter uint64) (string, error) {
	actionData := make(map[string]interface{})
	actionData["call_counter"] = callCounter
	return m.ExecAction(fmt.Sprintf("%v@open", m.ContractName), "sellrex", actionData)
}

func (m *TlosRexContract) WithdrawRex(callCounter uint64) (string, error) {
	actionData := make(map[string]interface{})
	actionData["call_counter"] = callCounter
	return m.ExecAction(fmt.Sprintf("%v@open", m.ContractName), "withdrawrex", actionData)
}

func (m *TlosRexContract) TstLapseTime(roundId uint64) (string, error) {
	actionData := make(map[string]interface{})
	actionData["round_id"] = roundId
	actionData["call_counter"] = m.NextCallCounter()
	return m.ExecAction(eos.AN(m.ContractName), "tstlapsetime", actionData)
}

func (m *Stake) Clone() *Stake {
	return &Stake{
		RoundID:              m.RoundID,
		TotalStake:           m.TotalStake,
		RexBalance:           m.RexBalance,
		TotalReturn:          m.TotalReturn,
		RexState:             m.RexState,
		StakingPeriod:        m.StakingPeriod,
		StakedTime:           m.StakedTime,
		MovedFromSavingsTime: m.MovedFromSavingsTime,
		StakeEndTime:         m.StakeEndTime,
		UpdatedDate:          m.UpdatedDate,
	}
}

func (m *TlosRexContract) GetStake(roundID uint64) (*Stake, error) {
	stakes, err := m.GetStakesReq(&eos.GetTableRowsRequest{
		LowerBound: strconv.FormatUint(roundID, 10),
		UpperBound: strconv.FormatUint(roundID, 10),
		Limit:      1,
	})
	if err != nil {
		return nil, err
	}
	if len(stakes) > 0 {
		return &stakes[0], nil
	}
	return nil, nil
}

func (m *TlosRexContract) GetLastStake() (*Stake, error) {
	stakes, err := m.GetStakesReq(&eos.GetTableRowsRequest{
		Reverse: true,
		Limit:   1,
	})
	if err != nil {
		return nil, err
	}
	if len(stakes) > 0 {
		return &stakes[0], nil
	}
	return nil, nil
}

func (m *TlosRexContract) GetStakesReq(req *eos.GetTableRowsRequest) ([]Stake, error) {

	var stakes []Stake
	if req == nil {
		req = &eos.GetTableRowsRequest{}
	}
	req.Table = "stakes"
	err := m.GetTableRows(*req, &stakes)
	if err != nil {
		return nil, fmt.Errorf("get table rows %v", err)
	}
	return stakes, nil
}

func (m *TlosRexContract) GetStakesByRexStateAndId(state eos.Name) ([]Stake, error) {
	request := &eos.GetTableRowsRequest{}
	err := m.FilterStakesByRexStateAndId(request, state)
	if err != nil {
		return nil, err
	}
	return m.GetStakesReq(request)
}

func (m *TlosRexContract) FilterStakesByRexStateAndId(req *eos.GetTableRowsRequest, state eos.Name) error {

	req.Index = "2"
	req.KeyType = "i128"
	req.Reverse = true
	stateAndRndLB, err := m.EOS.GetComposedIndexValue(state, 0)
	if err != nil {
		return fmt.Errorf("failed to generate lower bound composed index, err: %v", err)
	}
	stateAndRndUB, err := m.EOS.GetComposedIndexValue(state, uint64(18446744073709551615))
	if err != nil {
		return fmt.Errorf("failed to generate upper bound composed index, err: %v", err)
	}
	fmt.Println("LB: ", stateAndRndLB, "UB: ", stateAndRndUB)
	req.LowerBound = stateAndRndLB
	req.UpperBound = stateAndRndUB
	return err
}
