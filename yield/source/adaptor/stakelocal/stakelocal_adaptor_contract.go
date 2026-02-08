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
package stakelocal

import (
	"github.com/sebastianmontero/bennyfi-go-client/yield/source/adaptor/common"
	"github.com/sebastianmontero/eos-go"
	"github.com/sebastianmontero/eos-go-toolbox/service"
)

var (
	SettingBennyfiContract       = "BENNYFI_CONTRACT"
	SettingBatchSize             = "BATCH_SIZE"
	SettingBennyEvmBridgeAddress = "BENNY_EVM_BRIDGE_ADDRESS"
	SettingExsatBridgeContract   = "EXSAT_BRIDGE_CONTRACT"
	SettingEvmGasLimit           = "EVM_GAS_LIMIT"
	SettingEvmAccount            = "EVM_ACCOUNT"
)

type StakeLocalContract struct {
	*common.BaseContract
}

func NewStakeLocalContract(eos *service.EOS, contractName string) *StakeLocalContract {
	return &StakeLocalContract{
		common.NewBaseContract(eos, contractName),
	}
}

func (m *StakeLocalContract) Reset(limit uint64, toDelete []string) (string, error) {
	actionData := struct {
		Limit       uint64
		ToDelete    []string
		CallCounter uint64
	}{limit, toDelete, m.NextCallCounter()}
	return m.ExecAction(eos.AN(m.ContractName), "reset", actionData)
}
