// The MIT License (MIT)

// Copyright (c) 2020, Digital Scarcity

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.
package ibclocal

import (
	"fmt"

	"github.com/sebastianmontero/bennyfi-go-client/ibc"
	"github.com/sebastianmontero/bennyfi-go-client/ibc/bridge"
	eos "github.com/sebastianmontero/eos-go"
	"github.com/sebastianmontero/eos-go-toolbox/contract"
	"github.com/sebastianmontero/eos-go-toolbox/service"
)

type InitParams struct {
	ChainId                eos.Checksum256 `json:"chain_id"`
	BridgeContract         eos.AccountName `json:"bridge_contract"`
	PairedChainId          eos.Checksum256 `json:"paired_chain_id"`
	PairedWraplockContract eos.AccountName `json:"paired_wraplock_contract"`
	PairedTokenContract    eos.AccountName `json:"paired_token_contract"`
	StakeLocalContract     eos.AccountName `json:"stake_local_contract"`
}

func (m *InitParams) Clone() *InitParams {
	return &InitParams{
		ChainId:                m.ChainId,
		BridgeContract:         m.BridgeContract,
		PairedChainId:          m.PairedChainId,
		PairedWraplockContract: m.PairedWraplockContract,
		PairedTokenContract:    m.PairedTokenContract,
		StakeLocalContract:     m.StakeLocalContract,
	}
}
func (m *InitParams) ToGlobal() *Global {
	return &Global{m.Clone(), 0}
}

type Global struct {
	*InitParams
	Enabled uint8 `json:"enabled"`
}

type IBCLocalContract struct {
	*contract.TokenContract
	*ibc.IBCContract
}

func NewIBCLocalContract(eos *service.EOS, contractName string) *IBCLocalContract {
	return &IBCLocalContract{
		contract.NewTokenContractWithDefaultContract(eos, contractName),
		ibc.NewIBCContract(eos, contractName),
	}
}

func (m *IBCLocalContract) Init(initParams *InitParams, authorizer interface{}) (string, error) {
	return m.IBCContract.ExecActionStr(m.getAuthorizer(authorizer), "init", initParams)
}

func (m *IBCLocalContract) IssueA(prover eos.AccountName, heavyProof *bridge.HeavyProof, actionProof *bridge.ActionProof, authorizer interface{}) (string, error) {
	return m.ProofA(prover, "issuea", heavyProof, actionProof, authorizer)
}

func (m *IBCLocalContract) IssueB(prover eos.AccountName, lightProof *bridge.LightProof, actionProof *bridge.ActionProof, authorizer interface{}) (string, error) {
	return m.ProofB(prover, "issueb", lightProof, actionProof, authorizer)
}

func (m *IBCLocalContract) UnstakeA(prover eos.AccountName, heavyProof *bridge.HeavyProof, actionProof *bridge.ActionProof, authorizer interface{}) (string, error) {
	return m.ProofA(prover, "unstakea", heavyProof, actionProof, authorizer)
}

func (m *IBCLocalContract) UnstakeB(prover eos.AccountName, lightProof *bridge.LightProof, actionProof *bridge.ActionProof, authorizer interface{}) (string, error) {
	return m.ProofB(prover, "unstakeb", lightProof, actionProof, authorizer)
}

func (m *IBCLocalContract) Stake(stakeParams *ibc.StakeParams, authorizer interface{}) (string, error) {
	return m.IBCContract.ExecActionStr(m.getAuthorizer(authorizer), "stake", stakeParams)
}

func (m *IBCLocalContract) Enable(authorizer interface{}, enable bool) (string, error) {
	actionData := struct {
		Enable bool
	}{enable}
	return m.IBCContract.ExecActionStr(m.getAuthorizer(authorizer), "enable", actionData)
}

func (m *IBCLocalContract) GetGlobal() (*Global, error) {
	var global []*Global

	req := &eos.GetTableRowsRequest{
		Table: "global",
		Limit: 1,
	}
	err := m.IBCContract.GetTableRows(*req, &global)
	if err != nil {
		return nil, fmt.Errorf("get table rows %v", err)
	}
	if len(global) > 0 {
		return global[0], nil
	}
	return nil, nil
}

func (m *IBCLocalContract) getAuthorizer(authorizer interface{}) interface{} {
	if authorizer == nil {
		authorizer = m.IBCContract.ContractName
	}
	return authorizer
}
