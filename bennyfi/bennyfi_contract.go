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

	eos "github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/ecc"
	"github.com/sebastianmontero/eos-go-toolbox/contract"
	"github.com/sebastianmontero/eos-go-toolbox/service"
	"github.com/sebastianmontero/eos-go-toolbox/util"
)

var PercentageAdjustment = float64(10000000)

type BennyfiContract struct {
	*contract.Contract
}

func NewBennyfiContract(eos *service.EOS, contractName string) *BennyfiContract {
	return &BennyfiContract{
		&contract.Contract{
			EOS:          eos,
			ContractName: contractName,
		},
	}
}

func (m *BennyfiContract) ExecAction(permissionLevel interface{}, action string, actionData interface{}) (string, error) {
	resp, err := m.Contract.ExecAction(permissionLevel, action, actionData)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Tx ID: %v", resp.TransactionID), nil
}

func (m *BennyfiContract) ConfigureOpenPermission(publicKey *ecc.PublicKey) error {
	err := m.EOS.CreateSimplePermission(m.ContractName, "open", publicKey)
	if err != nil {
		return fmt.Errorf("failed to create open permission, error: %v", err)
	}
	err = m.EOS.LinkPermission(m.ContractName, "timedevents", "open")
	if err != nil {
		return fmt.Errorf("failed to link open permission to the timedevents action, error: %v", err)
	}

	err = m.EOS.LinkPermission(m.ContractName, "unstakeopen", "open")
	if err != nil {
		return fmt.Errorf("failed to link open permission to the unstakeopen action, error: %v", err)
	}
	return nil
}

func CalculatePercentage(amount interface{}, percentage int64) (eos.Asset, error) {
	amnt, err := util.ToAsset(amount)
	if err != nil {
		return eos.Asset{}, err
	}
	perctAmnt := float64(amnt.Amount) * float64((float64(percentage) / PercentageAdjustment))
	return eos.Asset{Amount: eos.Int64(perctAmnt), Symbol: amnt.Symbol}, nil
}
