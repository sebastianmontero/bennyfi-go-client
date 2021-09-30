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
	"strconv"

	eos "github.com/eoscanada/eos-go"
)

var (
	EntryStaked     = eos.Name("entrystaked")
	EntryReturnPaid = eos.Name("returnpaid")
	EntryUnstaked   = eos.Name("unstaked")
	EntryEarlyExit  = eos.Name("earlyexit")
)

type Entry struct {
	EntryID       uint64          `json:"entry_id"`
	RoundID       uint64          `json:"round_id"`
	Position      uint64          `json:"position"`
	Participant   eos.AccountName `json:"participant"`
	EntryStake    string          `json:"entry_stake"`
	Prize         string          `json:"prize"`
	MinimumPayout string          `json:"minimum_payout"`
	ReturnAmount  string          `json:"return_amount"`
	EntryStatus   eos.Name        `json:"entry_status"`
	EnteredDate   string          `json:"entered_date"`
}

func (m *BennyfiContract) EnterRound(roundId uint64, participant eos.AccountName) (string, error) {
	actionData := make(map[string]interface{})
	actionData["round_id"] = roundId
	actionData["participant"] = participant

	return m.ExecAction(participant, "enterround", actionData)
}

func (m *BennyfiContract) ClaimReturn(entryId uint64, claimer eos.AccountName) (string, error) {
	actionData := make(map[string]interface{})
	actionData["entry_id"] = entryId
	return m.ExecAction(claimer, "claimreturn", actionData)
}

func (m *BennyfiContract) Unstake(entryId uint64, permissionLevel interface{}) (string, error) {
	actionData := make(map[string]interface{})
	actionData["entry_id"] = entryId
	return m.ExecAction(permissionLevel, "unstake", actionData)
}

func (m *BennyfiContract) UnstakeOpen(entryId uint64) (string, error) {
	actionData := make(map[string]interface{})
	actionData["entry_id"] = entryId
	return m.ExecAction(fmt.Sprintf("%v@open", m.ContractName), "unstakeopen", actionData)
}

func (m *BennyfiContract) GetEntries() ([]Entry, error) {

	return m.GetEntriesReq(&eos.GetTableRowsRequest{})
}

func (m *BennyfiContract) GetEntriesbyParticipant(participant eos.AccountName) ([]Entry, error) {
	request := &eos.GetTableRowsRequest{}
	m.FilterEntriesbyParticipant(request, participant)
	return m.GetEntriesReq(request)
}

func (m *BennyfiContract) FilterEntriesbyParticipant(req *eos.GetTableRowsRequest, participant eos.AccountName) {
	req.Index = "4"
	req.KeyType = "name"
	req.LowerBound = string(participant)
	req.UpperBound = string(participant)
}

func (m *BennyfiContract) GetEntriesbyRound(roundID uint64) ([]Entry, error) {
	request := &eos.GetTableRowsRequest{}
	m.FilterEntriesbyRound(request, roundID)
	return m.GetEntriesReq(request)
}

func (m *BennyfiContract) FilterEntriesbyRound(req *eos.GetTableRowsRequest, roundID uint64) {

	req.Index = "5"
	req.KeyType = "i64"
	req.LowerBound = strconv.FormatUint(roundID, 10)
	req.UpperBound = strconv.FormatUint(roundID, 10)

}

func (m *BennyfiContract) GetEntriesbyStatus(status eos.Name) ([]Entry, error) {
	request := &eos.GetTableRowsRequest{}
	m.FilterEntriesbyStatus(request, status)
	return m.GetEntriesReq(request)
}

func (m *BennyfiContract) FilterEntriesbyStatus(req *eos.GetTableRowsRequest, status eos.Name) {

	req.Index = "2"
	req.KeyType = "name"
	req.LowerBound = string(status)
	req.UpperBound = string(status)
}

func (m *BennyfiContract) GetEntriesbyRoundAndPos(roundID uint64, pos uint64) ([]Entry, error) {
	request := &eos.GetTableRowsRequest{}
	err := m.FilterEntriesbyRoundAndPos(request, roundID, pos)
	if err != nil {
		return nil, err
	}
	return m.GetEntriesReq(request)
}

func (m *BennyfiContract) FilterEntriesbyRoundAndPos(req *eos.GetTableRowsRequest, roundID uint64, pos uint64) error {

	req.Index = "7"
	req.KeyType = "i128"
	rndAndPos, err := m.EOS.GetComposedIndexValue(roundID, pos)
	if err != nil {
		return fmt.Errorf("failed to generate composed index, err: %v", err)
	}
	req.LowerBound = rndAndPos
	req.UpperBound = rndAndPos
	return err
}

func (m *BennyfiContract) GetEntriesbyRoundAndStatus(roundID uint64, status eos.Name) ([]Entry, error) {
	request := &eos.GetTableRowsRequest{}
	err := m.FilterEntriesbyRoundAndStatus(request, roundID, status)
	if err != nil {
		return nil, err
	}
	return m.GetEntriesReq(request)
}

func (m *BennyfiContract) FilterEntriesbyRoundAndStatus(req *eos.GetTableRowsRequest, roundID uint64, status eos.Name) error {

	req.Index = "8"
	req.KeyType = "i128"
	rndAndStatus, err := m.EOS.GetComposedIndexValue(roundID, status)
	if err != nil {
		return fmt.Errorf("failed to generate composed index, err: %v", err)
	}
	req.LowerBound = rndAndStatus
	req.UpperBound = rndAndStatus
	return err
}

func (m *BennyfiContract) GetEntryByParticipantAndRound(participant eos.AccountName, roundID uint64) (*Entry, error) {
	entries, err := m.GetEntriesbyParticipant(participant)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if entry.RoundID == roundID {
			return &entry, nil
		}
	}
	return nil, nil
}

func (m *BennyfiContract) GetLastEntry() (*Entry, error) {
	entries, err := m.GetEntriesReq(&eos.GetTableRowsRequest{
		Reverse: true,
		Limit:   1,
	})
	if err != nil {
		return nil, err
	}
	if len(entries) > 0 {
		return &entries[0], nil
	}
	return nil, nil
}

func (m *BennyfiContract) GetEntryById(entryID uint64) (*Entry, error) {
	entries, err := m.GetEntriesReq(&eos.GetTableRowsRequest{
		LowerBound: strconv.FormatUint(entryID, 10),
		UpperBound: strconv.FormatUint(entryID, 10),
		Limit:      1,
	})
	if err != nil {
		return nil, err
	}
	if len(entries) > 0 {
		return &entries[0], nil
	}
	return nil, nil
}

func (m *BennyfiContract) GetEntriesReq(req *eos.GetTableRowsRequest) ([]Entry, error) {

	var entries []Entry
	if req == nil {
		req = &eos.GetTableRowsRequest{}
	}
	req.Table = "entries"
	err := m.GetTableRows(*req, &entries)
	if err != nil {
		return nil, fmt.Errorf("get table rows %v", err)
	}
	return entries, nil
}
