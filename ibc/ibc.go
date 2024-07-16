package ibc

import (
	"fmt"

	"github.com/sebastianmontero/bennyfi-go-client/ibc/bridge"
	"github.com/sebastianmontero/eos-go"
	"github.com/sebastianmontero/eos-go-toolbox/contract"
	"github.com/sebastianmontero/eos-go-toolbox/service"
)

type StakeParams struct {
	PoolId           uint64    `json:"pool_id"`
	YieldSource      eos.Name  `json:"yield_source"`
	Quantity         eos.Asset `json:"quantity"`
	StakingPeriodHrs uint32    `json:"staking_period_hrs"`
}

func (m *StakeParams) ToEmitStakeParams(tokenContract eos.AccountName) *EmitStakeParams {
	return &EmitStakeParams{
		PoolId:           m.PoolId,
		YieldSource:      m.YieldSource,
		Quantity:         eos.ExtendedAsset{Asset: m.Quantity, Contract: tokenContract},
		StakingPeriodHrs: m.StakingPeriodHrs,
	}
}

type EmitXferParams struct {
	Owner       eos.AccountName   `json:"owner"`
	Quantity    eos.ExtendedAsset `json:"quantity"`
	Beneficiary eos.AccountName   `json:"beneficiary"`
}

type EmitStakeParams struct {
	PoolId           uint64            `json:"pool_id"`
	YieldSource      eos.Name          `json:"yield_source"`
	Quantity         eos.ExtendedAsset `json:"quantity"`
	StakingPeriodHrs uint32            `json:"staking_period_hrs"`
}

type EmitUnstakeParams struct {
	PoolId      uint64            `json:"pool_id"`
	YieldSource eos.Name          `json:"yield_source"`
	TotalReturn eos.ExtendedAsset `json:"total_return"`
}

type LightProofRecord struct {
	ID         uint64             `json:"id"`
	LightProof *bridge.LightProof `json:"lp"`
}

type HeavyProofRecord struct {
	ID         uint64             `json:"id"`
	HeavyProof *bridge.HeavyProof `json:"hp"`
}
type IBCContract struct {
	*contract.Contract
}

func NewIBCContract(eos *service.EOS, contractName string) *IBCContract {
	return &IBCContract{
		contract.NewContract(eos, contractName),
	}
}

func (m *IBCContract) ProofA(prover eos.AccountName, action string, heavyProof *bridge.HeavyProof, actionProof *bridge.ActionProof, authorizer interface{}) (string, error) {
	actionData := struct {
		Prover      eos.AccountName
		HeavyProof  *bridge.HeavyProof
		ActionProof *bridge.ActionProof
	}{prover, heavyProof, actionProof}
	if authorizer == nil {
		authorizer = prover
	}
	return m.ExecActionStr(authorizer, action, actionData)
}

func (m *IBCContract) ProofB(prover eos.AccountName, action string, lightProof *bridge.LightProof, actionProof *bridge.ActionProof, authorizer interface{}) (string, error) {
	actionData := struct {
		Prover      eos.AccountName
		LightProof  *bridge.LightProof
		ActionProof *bridge.ActionProof
	}{prover, lightProof, actionProof}
	if authorizer == nil {
		authorizer = prover
	}
	return m.ExecActionStr(authorizer, action, actionData)
}

func (m *IBCContract) GetLightProof() (*LightProofRecord, error) {
	var lpr []*LightProofRecord

	req := &eos.GetTableRowsRequest{
		Table: "lightproof",
		Limit: 1,
	}
	err := m.GetTableRows(*req, &lpr)
	if err != nil {
		return nil, fmt.Errorf("get table rows %v", err)
	}
	if len(lpr) > 0 {
		return lpr[0], nil
	}
	return nil, nil
}

func (m *IBCContract) GetHeavyProof() (*HeavyProofRecord, error) {
	var hpr []*HeavyProofRecord

	req := &eos.GetTableRowsRequest{
		Table: "heavyproof",
		Limit: 1,
	}
	err := m.GetTableRows(*req, &hpr)
	if err != nil {
		return nil, fmt.Errorf("get table rows %v", err)
	}
	if len(hpr) > 0 {
		return hpr[0], nil
	}
	return nil, nil
}
