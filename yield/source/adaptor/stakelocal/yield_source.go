package stakelocal

import (
	"encoding/json"
	"fmt"

	"github.com/sebastianmontero/eos-go"
)

type YieldSource struct {
	YieldSource            eos.Name        `json:"yield_source"`
	StakeTokenContract     eos.AccountName `json:"stake_token_contract"`
	MinStakingPeriodHrs    uint32          `json:"min_staking_period_hrs"`
	MaxStakingPeriodHrs    uint32          `json:"max_staking_period_hrs"`
	MinStakeAmount         eos.Asset       `json:"min_stake_amount"`
	MaxStakeAmount         eos.Asset       `json:"max_stake_amount"`
	RewardTokenContract    eos.AccountName `json:"reward_token_contract"`
	RewardTokenSymbol      eos.Symbol      `json:"reward_token_symbol"`
	SupportsPartialReturns bool            `json:"supports_partial_returns"`
}

func (m *YieldSource) String() string {
	result, err := json.Marshal(m)
	if err != nil {
		panic(fmt.Sprintf("Failed marshalling round: %v", err))
	}
	return string(result)
}

func (m *StakeLocalContract) SetYieldSource(yieldSource YieldSource, authorizer interface{}) (string, error) {
	return m.ExecAction(m.GetValueOrContract(authorizer), "setyieldsrc", yieldSource)
}

func (m *StakeLocalContract) EraseYieldSource(yieldSource eos.Name, authorizer interface{}) (string, error) {
	return m.ExecAction(m.GetValueOrContract(authorizer), "eraseyldsrc", yieldSource)
}

func (m *YieldSource) Clone() *YieldSource {
	return &YieldSource{
		YieldSource:            m.YieldSource,
		StakeTokenContract:     m.StakeTokenContract,
		MinStakingPeriodHrs:    m.MinStakingPeriodHrs,
		MaxStakingPeriodHrs:    m.MaxStakingPeriodHrs,
		MinStakeAmount:         m.MinStakeAmount,
		MaxStakeAmount:         m.MaxStakeAmount,
		RewardTokenContract:    m.RewardTokenContract,
		RewardTokenSymbol:      m.RewardTokenSymbol,
		SupportsPartialReturns: m.SupportsPartialReturns,
	}
}

func (m *StakeLocalContract) GetYieldSourcesReq(req *eos.GetTableRowsRequest) ([]*YieldSource, error) {

	var yieldSources []*YieldSource
	if req == nil {
		req = &eos.GetTableRowsRequest{}
	}
	req.Table = "yieldsources"
	err := m.GetTableRows(*req, &yieldSources)
	if err != nil {
		return nil, fmt.Errorf("get table rows %v", err)
	}
	return yieldSources, nil
}

func (m *StakeLocalContract) GetYieldSourceById(yieldSource eos.Name) (*YieldSource, error) {
	request := &eos.GetTableRowsRequest{}
	m.FilterYieldSourcesById(request, yieldSource)
	yieldSources, err := m.GetYieldSourcesReq(request)
	if err != nil {
		return nil, err
	}
	if len(yieldSources) > 0 {
		return yieldSources[0], nil
	}
	return nil, nil
}

func (m *StakeLocalContract) FilterYieldSourcesById(req *eos.GetTableRowsRequest, yieldSource eos.Name) {
	req.LowerBound = string(yieldSource)
	req.UpperBound = string(yieldSource)
}
