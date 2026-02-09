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
	"strconv"

	"github.com/sebastianmontero/eos-go"
	"github.com/sebastianmontero/eos-go-toolbox/contract"
	"github.com/sebastianmontero/eos-go-toolbox/service"
)

var (
	StateStopped = eos.Name("stopped")
)

type BaseContract struct {
	*contract.SettingsContract
	callCounter uint64
}

func NewBaseContract(eos *service.EOS, contractName string) *BaseContract {
	return &BaseContract{
		contract.NewSettingsContract(eos, contractName),
		0,
	}
}

func (m *BaseContract) NextCallCounter() uint64 {
	m.callCounter++
	return m.callCounter
}

func (m *BaseContract) ExecAction(permissionLevel interface{}, action string, actionData interface{}) (string, error) {
	resp, err := m.Contract.ExecAction(permissionLevel, action, actionData)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Tx ID: %v", resp.TransactionID), nil
}

func (m *BaseContract) getStoppedStakeIndexValue(keyValue string) (string, error) {
	if keyValue == "" {
		keyValue = "0"
	}
	// fmt.Printf("key value: %v \n", keyValue)
	return m.EOS.GetComposedIndexValue(StateStopped, keyValue)
}

func (m *BaseContract) GetAllStoppedStakesAsMap() ([]map[string]interface{}, error) {
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

func (m *BaseContract) GetAllStoppedBasicStakes(statePropertyName, roundIdPropertyName string) ([]StoppedStake, error) {
	stakes, err := m.GetAllStoppedStakesAsMap()
	if err != nil {
		return nil, err
	}
	var basicStakes []StoppedStake
	for _, stake := range stakes {
		roundId, err := strconv.ParseUint(fmt.Sprintf("%v", stake[roundIdPropertyName]), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing round id: %v, error: %v", stake[roundIdPropertyName], err)
		}
		basicStakes = append(basicStakes, StoppedStake{
			RoundId:   roundId,
			IsStopped: eos.Name(stake[statePropertyName].(string)) == StateStopped,
		})
	}
	return basicStakes, nil
}
