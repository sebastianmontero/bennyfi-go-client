package stakelocal

import (
	"encoding/json"
	"fmt"

	"github.com/sebastianmontero/eos-go"
)

type YieldSource struct {
	YieldSource         eos.Name  `json:"yield_source"`
	TokenContract       eos.Name  `json:"token_contract"`
	MinStakingPeriodHrs uint32    `json:"min_staking_period_hrs"`
	MaxStakingPeriodHrs uint32    `json:"max_staking_period_hrs"`
	MinStakeAmount      eos.Asset `json:"min_stake_amount"`
	MaxStakeAmount      eos.Asset `json:"max_stake_amount"`
	// NOT USED AT THE MOMENT
	// AdditionalFields types.AdditionalFields `json:"additional_fields"`
}

func (m *YieldSource) String() string {
	result, err := json.Marshal(m)
	if err != nil {
		panic(fmt.Sprintf("Failed marshalling round: %v", err))
	}
	return string(result)
}

func (m *YieldSource) Clone() *YieldSource {
	return &YieldSource{
		YieldSource:         m.YieldSource,
		TokenContract:       m.TokenContract,
		MinStakingPeriodHrs: m.MinStakingPeriodHrs,
		MaxStakingPeriodHrs: m.MaxStakingPeriodHrs,
		MinStakeAmount:      m.MinStakeAmount,
		MaxStakeAmount:      m.MaxStakeAmount,
	}
}

func (m *StakeLocalContract) SetYieldSource(yieldSource *YieldSource) (string, error) {
	return m.ExecAction(m.ContractName, "setyieldsrc", yieldSource)
}

func (m *StakeLocalContract) EraseYieldSource(yieldSource eos.Name) (string, error) {
	actionData := struct {
		YieldSource eos.Name
	}{yieldSource}
	return m.ExecAction(m.ContractName, "setyieldsrc", actionData)
}

func (m *StakeLocalContract) GetYieldSource(yieldSource eos.Name) (*YieldSource, error) {
	yieldSources, err := m.GetYieldSourcesReq(&eos.GetTableRowsRequest{
		LowerBound: string(yieldSource),
		UpperBound: string(yieldSource),
		Limit:      1,
	})
	if err != nil {
		return nil, err
	}
	if len(yieldSources) > 0 {
		return &yieldSources[0], nil
	}
	return nil, nil
}

func (m *StakeLocalContract) GetLastYieldSource() (*YieldSource, error) {
	yieldSources, err := m.GetYieldSourcesReq(&eos.GetTableRowsRequest{
		Reverse: true,
		Limit:   1,
	})
	if err != nil {
		return nil, err
	}
	if len(yieldSources) > 0 {
		return &yieldSources[0], nil
	}
	return nil, nil
}

func (m *StakeLocalContract) GetYieldSourcesReq(req *eos.GetTableRowsRequest) ([]YieldSource, error) {

	var yieldSources []YieldSource
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
