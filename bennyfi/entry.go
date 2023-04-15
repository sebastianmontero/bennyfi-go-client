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

	eos "github.com/sebastianmontero/eos-go"
)

var (
	EntryStaked     = eos.Name("entrystaked")
	EntryReturnPaid = eos.Name("returnpaid")
	EntryUnstaked   = eos.Name("unstaked")
	EntryEarlyExit  = eos.Name("earlyexit")
)

type EnterRoundArgs struct {
	RoundID     uint64          `json:"round_id"`
	Participant eos.AccountName `json:"participant"`
}

type Entry struct {
	EntryID      uint64          `json:"entry_id"`
	RoundID      uint64          `json:"round_id"`
	Position     uint64          `json:"position"`
	Participant  eos.AccountName `json:"participant"`
	EntryStake   eos.Asset       `json:"entry_stake"`
	Returns      ReturnEntries   `json:"returns"`
	EntryStatus  eos.Name        `json:"entry_status"`
	VestingState eos.Name        `json:"vesting_state"`
	EnteredDate  eos.TimePoint   `json:"entered_date"`
}

type EntryCustomJSON struct {
	Returns map[eos.Name]interface{} `json:"returns"`
	Entry
}

func (m Entry) ToCustomJSON() EntryCustomJSON {
	return EntryCustomJSON{
		Returns: m.Returns.ToMap(),
		Entry:   m,
	}
}

func (m *Entry) UpsertReturn(name eos.Name, ret interface{}) {
	if m.Returns == nil {
		m.Returns = make(ReturnEntries, 0, 1)
	}
	m.Returns.Upsert(name, ret)
}

func (m *Entry) RemoveReturn(name eos.Name) {
	m.Returns.Remove(name)
}

func (m *BennyfiContract) EnterRound(roundId uint64, participant eos.AccountName) (string, error) {
	actionData := &EnterRoundArgs{
		RoundID:     roundId,
		Participant: participant,
	}

	return m.ExecAction(participant, "enterround", actionData)
}

func (m *BennyfiContract) ClaimReturn(entryId uint64, claimer eos.AccountName) (string, error) {
	return m.ExecAction(claimer, "claimreturn", entryId)
}

func (m *BennyfiContract) Unstake(entryId uint64, permissionLevel interface{}) (string, error) {
	return m.ExecAction(permissionLevel, "unstake", entryId)
}

func (m *BennyfiContract) UnstakeOpen(entryId uint64) (string, error) {
	return m.ExecAction(fmt.Sprintf("%v@open", m.ContractName), "unstakeopen", entryId)
}

func (m *BennyfiContract) Vesting(entryId uint64, permissionLevel interface{}) (string, error) {
	return m.ExecAction(permissionLevel, "vesting", entryId)
}

func (m *BennyfiContract) GetEntries() ([]Entry, error) {

	return m.GetEntriesReq(&eos.GetTableRowsRequest{})
}

func (m *BennyfiContract) GetAllEntriesAsMap() ([]map[string]interface{}, error) {
	req := eos.GetTableRowsRequest{
		Table: "entries",
	}
	return m.GetAllTableRowsAsMap(req, "entry_id")
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
