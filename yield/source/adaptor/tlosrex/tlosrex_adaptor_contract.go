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
package tlosrex

import (
	"fmt"

	"github.com/sebastianmontero/eos-go"
	"github.com/sebastianmontero/eos-go-toolbox/contract"
	"github.com/sebastianmontero/eos-go-toolbox/service"
	"github.com/sebastianmontero/eos-go/ecc"
)

var (
	SettingBennyfiContract      = "BENNYFI_CONTRACT"
	SettingTokenContract        = "TOKEN_CONTRACT"
	SettingRexContract          = "REX_CONTRACT"
	SettingRexDeposit_account   = "REX_DEPOSIT_ACCOUNT"
	SettingBatchSize            = "BATCH_SIZE"
	SettingMinStakingPeriod_hrs = "MIN_STAKING_PERIOD_HRS"
	SettingMaxStakingPeriod_hrs = "MAX_STAKING_PERIOD_HRS"
	SettingMinStakeAmount       = "MIN_STAKE_AMOUNT"
	SettingMaxStakeAmount       = "MAX_STAKE_AMOUNT"
	RexStatePreRex              = eos.Name("prerex")
	RexStateInSavings           = eos.Name("insavings")
	RexStateInLockPeriod        = eos.Name("lockperiod")
	RexStateSold                = eos.Name("sold")
	RexStateWithdrawn           = eos.Name("withdrawn")
)

type TlosRexContract struct {
	*contract.SettingsContract
	callCounter uint64
}

func NewTlosRexContract(eos *service.EOS, contractName string) *TlosRexContract {
	return &TlosRexContract{
		contract.NewSettingsContract(eos, contractName),
		0,
	}
}

func (m *TlosRexContract) NextCallCounter() uint64 {
	m.callCounter++
	return m.callCounter
}

func (m *TlosRexContract) ExecAction(permissionLevel interface{}, action string, actionData interface{}) (string, error) {
	resp, err := m.Contract.ExecAction(permissionLevel, action, actionData)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Tx ID: %v", resp.TransactionID), nil
}

func (m *TlosRexContract) ConfigureOpenPermission(publicKey *ecc.PublicKey) error {
	openActions := []string{
		"mvfrmsvngsrn",
		"sellrexrn",
		"withdrwrexrn",
		"mvfrmsavings",
		"sellrex",
		"withdrawrex",
	}
	err := m.EOS.CreateSimplePermission(m.ContractName, "open", publicKey)
	if err != nil {
		return fmt.Errorf("failed to create open permission, error: %v", err)
	}
	for _, action := range openActions {
		err = m.EOS.LinkPermission(m.ContractName, action, "open", false)
		if err != nil {
			return fmt.Errorf("failed to link open permission to the %v action, error: %v", action, err)
		}
	}
	return nil
}

func (m *TlosRexContract) Reset(limit uint64, toDelete []string) (string, error) {
	actionData := struct {
		Limit       uint64
		ToDelete    []string
		CallCounter uint64
	}{limit, toDelete, m.NextCallCounter()}
	return m.ExecAction(eos.AN(m.ContractName), "reset", actionData)
}
