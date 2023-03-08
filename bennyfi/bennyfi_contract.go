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
package bennyfi

import (
	"fmt"
	"time"

	eos "github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/ecc"
	"github.com/sebastianmontero/eos-go-toolbox/contract"
	"github.com/sebastianmontero/eos-go-toolbox/service"
)

var (
	PAUSED   uint32 = 1
	UNPAUSED uint32 = 0
)

type BennyfiContract struct {
	*contract.SettingsContract
	callCounter uint64
}

func NewBennyfiContract(eos *service.EOS, contractName string) *BennyfiContract {
	return &BennyfiContract{
		contract.NewSettingsContract(eos, contractName),
		0,
	}
}

func (m *BennyfiContract) NextCallCounter() uint64 {
	m.callCounter++
	return m.callCounter
}

func (m *BennyfiContract) ExecAction(permissionLevel interface{}, action string, actionData interface{}) (string, error) {
	resp, err := m.Contract.ExecAction(permissionLevel, action, actionData)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Tx ID: %v", resp.TransactionID), nil
}

func (m *BennyfiContract) ProposeAction(proposerName interface{}, requested []eos.PermissionLevel, expireIn time.Duration, permissionLevel, actionName, data interface{}) (string, error) {
	resp, err := m.Contract.ProposeAction(proposerName, requested, expireIn, permissionLevel, actionName, data)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Proposal Name: %v, Tx ID: %v", resp.ProposalName, resp.PushTransactionFullResp.TransactionID), nil
}

func (m *BennyfiContract) ConfigureOpenPermission(publicKey *ecc.PublicKey) error {
	openActions := []string{
		"timedevents",
		"startround",
		"startrounds",
		"timeoutrnds",
		"mvfrmsavings",
		"sellrex",
		"withdrawrex",
		"unlockrnd",
		"unlockrnds",
		"redraw",
		"unstakeopen",
		"ustkulckrnds",
		"ustktmdrnds",
		"vestingrnds",
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

func (m *BennyfiContract) Pause(pause uint32) (string, error) {
	actionData := make(map[string]interface{})
	actionData["pause"] = pause
	return m.ExecAction(eos.AN(m.ContractName), "pause", actionData)
}

func (m *BennyfiContract) Reset(limit uint64, toDelete []string) (string, error) {
	actionData := make(map[string]interface{})
	actionData["limit"] = limit
	actionData["to_delete"] = toDelete
	actionData["call_counter"] = m.NextCallCounter()
	return m.ExecAction(eos.AN(m.ContractName), "reset", actionData)
}
