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
package ibcremote

import (
	"fmt"
	"strconv"

	"github.com/sebastianmontero/bennyfi-go-client/ibc"
	"github.com/sebastianmontero/bennyfi-go-client/ibc/bridge"
	eos "github.com/sebastianmontero/eos-go"
	"github.com/sebastianmontero/eos-go-toolbox/service"
)

type InitParams struct {
	ChainId        eos.Checksum256 `json:"chain_id"`
	BridgeContract eos.AccountName `json:"bridge_contract"`
	PairedChainId  eos.Checksum256 `json:"paired_chain_id"`
}

func (m *InitParams) Clone() *InitParams {
	return &InitParams{
		ChainId:        m.ChainId,
		BridgeContract: m.BridgeContract,
		PairedChainId:  m.PairedChainId,
	}
}
func (m *InitParams) ToGlobal() *Global {
	return &Global{m.Clone(), 0}
}

type Global struct {
	*InitParams
	Enabled uint8 `json:"enabled"`
}

type ContractMapping struct {
	NativeTokenContract     eos.AccountName `json:"native_token_contract"`
	PairedWraptokenContract eos.AccountName `json:"paired_wraptoken_contract"`
}

type YieldSource struct {
	YieldSource     eos.Name        `json:"yield_source"`
	AdaptorContract eos.AccountName `json:"adaptor_contract"`
}

type Reserve struct {
	Balance eos.Asset `json:"balance"`
}

type IBCRemoteContract struct {
	*ibc.IBCContract
}

func NewIBCRemoteContract(eos *service.EOS, contractName string) *IBCRemoteContract {
	return &IBCRemoteContract{
		ibc.NewIBCContract(eos, contractName),
	}
}

func (m *IBCRemoteContract) Init(initParams *InitParams, authorizer interface{}) (string, error) {
	return m.IBCContract.ExecActionStr(m.GetValueOrContract(authorizer), "init", initParams)
}

func (m *IBCRemoteContract) AddContract(contractMapping *ContractMapping, authorizer interface{}) (string, error) {
	return m.IBCContract.ExecActionStr(m.GetValueOrContract(authorizer), "addcontract", contractMapping)
}

func (m *IBCRemoteContract) DeleteContract(nativeTokenContract eos.AccountName, authorizer interface{}) (string, error) {
	actionData := struct {
		NativeTokenContract eos.AccountName
	}{nativeTokenContract}
	return m.IBCContract.ExecActionStr(m.GetValueOrContract(authorizer), "delcontract", actionData)
}

func (m *IBCRemoteContract) AddYieldSource(yieldSource *YieldSource, authorizer interface{}) (string, error) {
	return m.IBCContract.ExecActionStr(m.GetValueOrContract(authorizer), "addyieldsrc", yieldSource)
}

func (m *IBCRemoteContract) DeleteYieldSource(yieldSource eos.Name, authorizer interface{}) (string, error) {
	actionData := struct {
		YieldSource eos.Name
	}{yieldSource}
	return m.IBCContract.ExecActionStr(m.GetValueOrContract(authorizer), "delyieldsrc", actionData)
}

func (m *IBCRemoteContract) StakeA(prover eos.AccountName, heavyProof *bridge.HeavyProof, actionProof *bridge.ActionProof, authorizer interface{}) (string, error) {
	return m.ProofA(prover, "stakea", heavyProof, actionProof, authorizer)
}

func (m *IBCRemoteContract) StakeB(prover eos.AccountName, lightProof *bridge.LightProof, actionProof *bridge.ActionProof, authorizer interface{}) (string, error) {
	return m.ProofB(prover, "stakeb", lightProof, actionProof, authorizer)
}

func (m *IBCRemoteContract) Stake(stakeParams *ibc.StakeParams, authorizer interface{}) (string, error) {
	return m.IBCContract.ExecActionStr(m.GetValueOrContract(authorizer), "stake", stakeParams)
}

func (m *IBCRemoteContract) GetGlobal() (*Global, error) {
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

func (m *IBCRemoteContract) GetYieldSource(yieldSource eos.Name) (*YieldSource, error) {
	yieldSources, err := m.GetYieldSourcesReq(&eos.GetTableRowsRequest{
		LowerBound: string(yieldSource),
		UpperBound: string(yieldSource),
		Limit:      1,
	})
	if err != nil {
		return nil, err
	}
	if len(yieldSources) > 0 {
		return yieldSources[0], nil
	}
	return nil, nil
}

func (m *IBCRemoteContract) GetYieldSourcesReq(req *eos.GetTableRowsRequest) ([]*YieldSource, error) {

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

func (m *IBCRemoteContract) GetContractMapping(nativeTokenContract eos.AccountName) (*ContractMapping, error) {
	contractMappings, err := m.GetContractMappingsReq(&eos.GetTableRowsRequest{
		LowerBound: string(nativeTokenContract),
		UpperBound: string(nativeTokenContract),
		Limit:      1,
	})
	if err != nil {
		return nil, err
	}
	if len(contractMappings) > 0 {
		return contractMappings[0], nil
	}
	return nil, nil
}

func (m *IBCRemoteContract) GetContractMappingsReq(req *eos.GetTableRowsRequest) ([]*ContractMapping, error) {

	var contractMappings []*ContractMapping
	if req == nil {
		req = &eos.GetTableRowsRequest{}
	}
	req.Table = "contractmap"
	err := m.GetTableRows(*req, &contractMappings)
	if err != nil {
		return nil, fmt.Errorf("get table rows %v", err)
	}
	return contractMappings, nil
}

func (m *IBCRemoteContract) GetReserve(tokenContract eos.AccountName, symbol eos.Symbol) (*Reserve, error) {
	symbolCode, err := symbol.SymbolCode()
	if err != nil {
		return nil, fmt.Errorf("failed getting symbol code, error: %v", err)
	}
	reserves, err := m.GetReservesReq(&eos.GetTableRowsRequest{
		Scope:      tokenContract.String(),
		LowerBound: strconv.FormatUint(uint64(symbolCode), 10),
		UpperBound: strconv.FormatUint(uint64(symbolCode), 10),
		Limit:      1,
	})
	if err != nil {
		return nil, err
	}
	if len(reserves) > 0 {
		return reserves[0], nil
	}
	return nil, nil
}

func (m *IBCRemoteContract) GetReservesReq(req *eos.GetTableRowsRequest) ([]*Reserve, error) {

	var reserves []*Reserve
	if req == nil {
		req = &eos.GetTableRowsRequest{}
	}
	req.Table = "reserves"
	err := m.GetTableRows(*req, &reserves)
	if err != nil {
		return nil, fmt.Errorf("get table rows %v", err)
	}
	return reserves, nil
}
