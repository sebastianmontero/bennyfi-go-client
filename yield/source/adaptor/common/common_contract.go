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
package common

import (
	"fmt"

	"github.com/sebastianmontero/eos-go"
	"github.com/sebastianmontero/eos-go-toolbox/contract"
	"github.com/sebastianmontero/eos-go-toolbox/service"
)

var (
	StateStopped = eos.Name("stopped")
)

type CommonContract struct {
	*contract.SettingsContract
	callCounter uint64
}

func NewCommonContract(eos *service.EOS, contractName string) *CommonContract {
	return &CommonContract{
		contract.NewSettingsContract(eos, contractName),
		0,
	}
}

func (m *CommonContract) NextCallCounter() uint64 {
	m.callCounter++
	return m.callCounter
}

func (m *CommonContract) ExecAction(permissionLevel interface{}, action string, actionData interface{}) (string, error) {
	resp, err := m.Contract.ExecAction(permissionLevel, action, actionData)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Tx ID: %v", resp.TransactionID), nil
}

func (m *CommonContract) getStoppedStakeIndexValue(keyValue string) (string, error) {
	if keyValue == "" {
		keyValue = "0"
	}
	// fmt.Printf("key value: %v \n", keyValue)
	return m.EOS.GetComposedIndexValue(StateStopped, keyValue)
}

func (m *CommonContract) GetAllStoppedStakesAsMap() ([]map[string]interface{}, error) {
	req := eos.GetTableRowsRequest{
		Table:   "stakes",
		Index:   "2",
		KeyType: "i128",
	}
	stateAndRndUB, err := m.EOS.GetComposedIndexValue(StateStopped, uint64(18446744073709551615))
	if err != nil {
		return nil, fmt.Errorf("failed to generate upper bound composed index, err: %v", err)
	}
	return m.GetAllTableRowsFromTillAsMap(req, "pool_id", "0", m.getStoppedStakeIndexValue, stateAndRndUB)
}
