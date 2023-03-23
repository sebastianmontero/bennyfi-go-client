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
package dummyyield

import (
	"fmt"

	"github.com/eoscanada/eos-go"
	"github.com/sebastianmontero/eos-go-toolbox/contract"
	"github.com/sebastianmontero/eos-go-toolbox/service"
)

var (
	SettingBennyfiContract      = "BENNYFI_CONTRACT"
	SettingTokenContract        = "TOKEN_CONTRACT"
	SettingBatchSize            = "BATCH_SIZE"
	SettingMinStakingPeriod_hrs = "MIN_STAKING_PERIOD_HRS"
	SettingMaxStakingPeriod_hrs = "MAX_STAKING_PERIOD_HRS"
	SettingMinStakeAmount       = "MIN_STAKE_AMOUNT"
	SettingMaxStakeAmount       = "MAX_STAKE_AMOUNT"
	StateStaked                 = eos.Name("staked")
	StateWithdrawn              = eos.Name("withdrawn")
)

type DummyYieldContract struct {
	*contract.SettingsContract
	callCounter uint64
}

func NewDummyYieldContract(eos *service.EOS, contractName string) *DummyYieldContract {
	return &DummyYieldContract{
		contract.NewSettingsContract(eos, contractName),
		0,
	}
}

func (m *DummyYieldContract) NextCallCounter() uint64 {
	m.callCounter++
	return m.callCounter
}

func (m *DummyYieldContract) ExecAction(permissionLevel interface{}, action string, actionData interface{}) (string, error) {
	resp, err := m.Contract.ExecAction(permissionLevel, action, actionData)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Tx ID: %v", resp.TransactionID), nil
}

func (m *DummyYieldContract) Reset(limit uint64, toDelete []string) (string, error) {
	actionData := make(map[string]interface{})
	actionData["limit"] = limit
	actionData["to_delete"] = toDelete
	actionData["call_counter"] = m.NextCallCounter()
	return m.ExecAction(eos.AN(m.ContractName), "reset", actionData)
}
