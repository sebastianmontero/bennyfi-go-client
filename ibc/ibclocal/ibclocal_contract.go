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
	"time"

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

type StakeParams struct {
	PoolId           uint64    `json:"pool_id"`
	YieldSource      eos.Name  `json:"yield_source"`
	Quantity         eos.Asset `json:"quantity"`
	StakingPeriodHrs uint32    `json:"staking_period_hrs"`
}

type IBCLocalContract struct {
	*contract.TokenContract
	callCounter uint64
}

func NewIBCLocalContract(eos *service.EOS, contractName string) *IBCLocalContract {
	return &IBCLocalContract{
		contract.NewTokenContractWithDefaultContract(eos, contractName),
		0,
	}
}

func (m *IBCLocalContract) NextCallCounter() uint64 {
	m.callCounter++
	return m.callCounter
}

func (m *IBCLocalContract) ExecAction(permissionLevel interface{}, action string, actionData interface{}) (string, error) {
	resp, err := m.Contract.ExecAction(permissionLevel, action, actionData)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Tx ID: %v", resp.TransactionID), nil
}

func (m *IBCLocalContract) ExecActions(actions ...*eos.Action) (string, error) {
	resp, err := m.Contract.ExecActions(actions...)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Tx ID: %v", resp.TransactionID), nil
}

func (m *IBCLocalContract) ProposeAction(proposerName interface{}, requested []eos.PermissionLevel, expireIn time.Duration, permissionLevel, actionName, data interface{}) (string, error) {
	resp, err := m.Contract.ProposeAction(proposerName, requested, expireIn, permissionLevel, actionName, data)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Proposal Name: %v, Tx ID: %v", resp.ProposalName, resp.PushTransactionFullResp.TransactionID), nil
}

func (m *IBCLocalContract) Init(initParams *InitParams, authorizer interface{}) (string, error) {
	return m.ExecAction(m.getAuthorizer(authorizer), "init", initParams)
}

func (m *IBCLocalContract) Stake(stakeParams *StakeParams, authorizer interface{}) (string, error) {
	return m.ExecAction(m.getAuthorizer(authorizer), "stake", stakeParams)
}

func (m *IBCLocalContract) Enable(authorizer interface{}) (string, error) {
	return m.ExecAction(m.getAuthorizer(authorizer), "enable", nil)
}

func (m *IBCLocalContract) Disable(authorizer interface{}) (string, error) {
	return m.ExecAction(m.getAuthorizer(authorizer), "disable", nil)
}

func (m *IBCLocalContract) GetGlobal() (*Global, error) {
	var global []*Global

	req := &eos.GetTableRowsRequest{
		Table: "global",
		Limit: 1,
	}
	err := m.GetTableRows(*req, &global)
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
		authorizer = m.ContractName
	}
	return authorizer
}
