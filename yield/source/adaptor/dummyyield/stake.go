package dummyyield

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/sebastianmontero/eos-go"
)

type Stake struct {
	RoundID      uint64   `json:"pool_id"`
	TotalStake   string   `json:"total_stake"`
	TotalReturn  string   `json:"total_return"`
	CurrentState eos.Name `json:"current_state"`
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

func (m *Stake) GetTotalReturn() eos.Asset {
	totalReturn, err := eos.NewAssetFromString(m.TotalReturn)
	if err != nil {
		panic(fmt.Sprintf("Unable to parse total return: %v to asset", m.TotalReturn))
	}
	return totalReturn
}

func (m *DummyYieldContract) CheckStakeParameters(authorizer, tokenContract eos.AccountName, stakeAmount eos.Asset, stakingPeriodHrs uint32) (string, error) {
	actionData := struct {
		TokenContract    eos.AccountName
		StakeAmount      eos.Asset
		StakingPeriodHrs uint32
	}{tokenContract, stakeAmount, stakingPeriodHrs}
	return m.ExecAction(authorizer, "chckstkparam", actionData)
}

func (m *DummyYieldContract) Withdraw(roundID uint64, yieldPercx100000 uint32) (string, error) {
	actionData := struct {
		RoundID          uint64
		YieldPercx100000 uint32
	}{roundID, yieldPercx100000}
	return m.ExecAction(m.ContractName, "withdraw", actionData)
}

func (m *Stake) Clone() *Stake {
	return &Stake{
		RoundID:     m.RoundID,
		TotalStake:  m.TotalStake,
		TotalReturn: m.TotalReturn,
	}
}

func (m *DummyYieldContract) GetStake(roundID uint64) (*Stake, error) {
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

func (m *DummyYieldContract) GetLastStake() (*Stake, error) {
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

func (m *DummyYieldContract) GetStakesReq(req *eos.GetTableRowsRequest) ([]Stake, error) {

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

func (m *DummyYieldContract) GetStakesByStateAndId(state eos.Name) ([]Stake, error) {
	request := &eos.GetTableRowsRequest{}
	err := m.FilterStakesByStateAndId(request, state)
	if err != nil {
		return nil, err
	}
	return m.GetStakesReq(request)
}

func (m *DummyYieldContract) FilterStakesByStateAndId(req *eos.GetTableRowsRequest, state eos.Name) error {

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
