package eosrex

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/sebastianmontero/bennyfi-go-client/common/types"
	"github.com/sebastianmontero/eos-go"
	"github.com/sebastianmontero/eos-go-toolbox/dto"
)

var RexLockPeriodDays = 21

type Stake struct {
	RoundID              uint64            `json:"pool_id"`
	TotalStake           eos.Asset         `json:"total_stake"`
	RexBalance           eos.Asset         `json:"rex_balance"`
	TotalReturn          eos.Asset         `json:"total_return"`
	RexState             eos.Name          `json:"rex_state"`
	StakingPeriod        *dto.Microseconds `json:"staking_period"`
	StakedTime           eos.TimePoint     `json:"staked_time"`
	MovedFromSavingsTime eos.TimePoint     `json:"moved_from_savings_time"`
	MaturityTime         eos.TimePoint     `json:"maturity_time"`
	StakeEndTime         eos.TimePoint     `json:"stake_end_time"`
	UpdatedDate          eos.TimePoint     `json:"updated_date"`
	// NOT USED AT THE MOMENT
	AdditionalFields types.AdditionalFields `json:"additional_fields"`
}

func (m *Stake) String() string {
	result, err := json.Marshal(m)
	if err != nil {
		panic(fmt.Sprintf("Failed marshalling round: %v", err))
	}
	return string(result)
}

func (m *Stake) CalculateRexLockPeriodTime() eos.TimePoint {

	return eos.TimePoint(m.StakedTime.Time().Add(time.Hour * time.Duration(m.StakingPeriod.Hrs()-int64(24*RexLockPeriodDays))).UnixMicro())
}

func (m *Stake) CalculateSellRexTime() eos.TimePoint {

	return eos.TimePoint(m.MovedFromSavingsTime.Time().Add(time.Hour * time.Duration(24*RexLockPeriodDays)).UnixMicro())
}

func (m *Stake) CalculateStakeEndTime() eos.TimePoint {
	return eos.TimePoint(m.StakedTime.Time().Add(time.Hour * time.Duration(m.StakingPeriod.Hrs())).UnixMicro())
}

func (m *Stake) GetStateWhenStopped() eos.Name {
	if m.RexState != RexStateStopped {
		panic("Stake is not stopped")
	}
	return m.AdditionalFields.GetValue(FieldStateWhenStopped).Name()
}

func (m *Stake) SetStateWhenStopped(stateWhenStopped eos.Name) {
	m.AdditionalFields.Set(FieldStateWhenStopped, dto.FlexValueFromName(stateWhenStopped))
}

func (m *EosRexContract) CheckStakeParameters(authorizer, tokenContract eos.AccountName, minStakeAmount eos.Asset, maxStakeAmount eos.Asset, stakingPeriodHrs uint32) (string, error) {
	actionData := struct {
		TokenContract    eos.AccountName
		MinStakeAmount   eos.Asset
		MaxStakeAmount   eos.Asset
		StakingPeriodHrs uint32
	}{tokenContract, minStakeAmount, maxStakeAmount, stakingPeriodHrs}
	return m.ExecAction(authorizer, "chckstkparam", actionData)
}

func (m *EosRexContract) MoveRoundFromSavings(roundId uint64) (string, error) {
	return m.ExecAction(fmt.Sprintf("%v@open", m.ContractName), "mvfrmsvngsrn", roundId)
}

func (m *EosRexContract) CalculateProceedsRoundRex(roundId uint64) (string, error) {
	return m.ExecAction(fmt.Sprintf("%v@open", m.ContractName), "clcproceedrn", roundId)
}

func (m *EosRexContract) WithdrawRoundRex(roundId uint64) (string, error) {
	return m.ExecAction(fmt.Sprintf("%v@open", m.ContractName), "withdrwrexrn", roundId)
}

func (m *EosRexContract) MoveFromSavings(callCounter uint64) (string, error) {
	return m.ExecAction(fmt.Sprintf("%v@open", m.ContractName), "mvfrmsavings", callCounter)
}

func (m *EosRexContract) UpdateRex(callCounter uint64) (string, error) {
	return m.ExecAction(fmt.Sprintf("%v@open", m.ContractName), "updaterex", callCounter)
}

func (m *EosRexContract) CalculateProceeds(callCounter uint64) (string, error) {
	return m.ExecAction(fmt.Sprintf("%v@open", m.ContractName), "calcproceeds", callCounter)
}

func (m *EosRexContract) SellRex(callCounter uint64) (string, error) {
	return m.ExecAction(fmt.Sprintf("%v@open", m.ContractName), "sellrex", callCounter)
}

func (m *EosRexContract) WithdrawRex(callCounter uint64) (string, error) {
	return m.ExecAction(fmt.Sprintf("%v@open", m.ContractName), "withdrawrex", callCounter)
}

func (m *EosRexContract) StopStake(roundId uint64, authorizer interface{}) (string, error) {
	return m.ExecAction(m.GetValueOrContract(authorizer), "stopstake", roundId)
}

func (m *EosRexContract) StopStakes(callCounter uint64, authorizer interface{}) (string, error) {
	return m.ExecAction(m.GetValueOrContract(authorizer), "stopstakes", callCounter)
}

func (m *EosRexContract) SetRexBalance(roundId uint64, rexBalance eos.Asset, authorizer interface{}) (string, error) {
	actionData := struct {
		RoundId    uint64
		RexBalance eos.Asset
	}{roundId, rexBalance}
	return m.ExecAction(m.GetValueOrContract(authorizer), "setrexbal", actionData)
}

func (m *EosRexContract) SetTotalReturn(roundId uint64, totalReturn eos.Asset, updateState bool, authorizer interface{}) (string, error) {
	actionData := struct {
		RoundId     uint64
		TotalReturn eos.Asset
		UpdateState bool
	}{roundId, totalReturn, updateState}
	return m.ExecAction(m.GetValueOrContract(authorizer), "settotalrtrn", actionData)
}

func (m *EosRexContract) TstLapseTime(roundId uint64, inTests bool) (string, error) {
	actionData := struct {
		RoundId     uint64
		InTests     bool
		CallCounter uint64
	}{roundId, inTests, m.NextCallCounter()}
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
		MaturityTime:         m.MaturityTime,
		StakeEndTime:         m.StakeEndTime,
		UpdatedDate:          m.UpdatedDate,
	}
}

func (m *EosRexContract) GetStake(roundID uint64) (*Stake, error) {
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

func (m *EosRexContract) GetLastStake() (*Stake, error) {
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

func (m *EosRexContract) GetStakesReq(req *eos.GetTableRowsRequest) ([]Stake, error) {

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

func (m *EosRexContract) GetStakesByRexStateAndId(state eos.Name) ([]Stake, error) {
	request := &eos.GetTableRowsRequest{}
	err := m.FilterStakesByRexStateAndId(request, state)
	if err != nil {
		return nil, err
	}
	return m.GetStakesReq(request)
}

func (m *EosRexContract) FilterStakesByRexStateAndId(req *eos.GetTableRowsRequest, state eos.Name) error {

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
	// fmt.Println("LB: ", stateAndRndLB, "UB: ", stateAndRndUB)
	req.LowerBound = stateAndRndLB
	req.UpperBound = stateAndRndUB
	return err
}
