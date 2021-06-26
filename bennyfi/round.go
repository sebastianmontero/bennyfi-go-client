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
	"encoding/json"
	"fmt"
	"strconv"

	eos "github.com/eoscanada/eos-go"
)

var (
	RoundAcceptingEntries = eos.Name("acceptentrys")
	RoundDrawing          = eos.Name("rounddrawing")
	RoundOpen             = eos.Name("roundopen")
	RoundClosed           = eos.Name("roundclosed")
	RoundUnlocked         = eos.Name("rndunlocked")
	RoundUnlockedUnstaked = eos.Name("rndunlckdutk")
	RoundTimedOut         = eos.Name("rndtimedout")
	RoundTimedOutUnstaked = eos.Name("rndtmdoututk")
	EntryStaked           = eos.Name("entrystaked")
	EntryReturnPaid       = eos.Name("returnpaid")
	EntryUnstaked         = eos.Name("unstaked")
	EntryEarlyExit        = eos.Name("earlyexit")
)

var microsecondsPerHr int64 = 60 * 60 * 1000000

type Microseconds struct {
	Microseconds string `json:"_count"`
}

func NewMicroseconds(hrs int64) *Microseconds {
	return &Microseconds{
		Microseconds: strconv.FormatInt(hrs*microsecondsPerHr, 10),
	}
}

func (m *Microseconds) Hrs() int64 {
	ms, _ := strconv.ParseInt(m.Microseconds, 10, 64)
	return ms / microsecondsPerHr
}

func (m *Microseconds) UnmarshalJSON(b []byte) error {
	ms := make(map[string]interface{})
	if err := json.Unmarshal(b, &ms); err != nil {
		return err
	}
	if countI, ok := ms["_count"]; ok {
		var microseconds string
		switch count := countI.(type) {
		case float64:
			microseconds = strconv.FormatFloat(count, 'f', 0, 64)
		case string:
			microseconds = count
		default:
			return fmt.Errorf("Microseconds count of unknown type: %T", count)
		}
		*m = Microseconds{
			Microseconds: microseconds,
		}
	} else {
		return fmt.Errorf("Error unmarshalling microseconds no '_count' property found: %v", ms)
	}
	return nil
}

type Term struct {
	TermID              uint64          `json:"term_id"`
	TermName            string          `json:"term_name"`
	AllParticipantsPerc uint32          `json:"all_participants_perc_x100000"`
	RoundManager        eos.AccountName `json:"round_manager"`
	Beneficiary         eos.AccountName `json:"beneficiary"`
	BeneficiaryPerc     uint32          `json:"beneficiary_perc_x100000"`
	CreatedDate         string          `json:"created_date"`
	UpdatedDate         string          `json:"updated_date"`
}

type NewTermArgs struct {
	TermName            string          `json:"term_name"`
	AllParticipantsPerc uint32          `json:"all_participants_perc_x100000"`
	RoundManager        eos.AccountName `json:"round_manager"`
	Beneficiary         eos.AccountName `json:"beneficiary"`
	BeneficiaryPerc     uint32          `json:"beneficiary_perc_x100000"`
}

func TermToNewTermArgs(terms *Term) *NewTermArgs {
	return &NewTermArgs{
		TermName:            terms.TermName,
		AllParticipantsPerc: terms.AllParticipantsPerc,
		RoundManager:        terms.RoundManager,
		Beneficiary:         terms.Beneficiary,
		BeneficiaryPerc:     terms.BeneficiaryPerc,
	}
}

type Winner struct {
	Participant   eos.AccountName `json:"participant"`
	Prize         string          `json:"prize"`
	EntryPosition uint64          `json:"entry_position"`
}

func NewWinner(participant eos.AccountName, prize string, entryPosition uint64) *Winner {
	return &Winner{
		Participant:   participant,
		Prize:         prize,
		EntryPosition: entryPosition,
	}
}

type Winners []*Winner

type Round struct {
	RoundID                uint64          `json:"round_id"`
	TermID                 uint64          `json:"term_id"`
	RoundName              string          `json:"round_name"`
	StakingPeriod          *Microseconds   `json:"staking_period"`
	EnrollmentTimeOut      *Microseconds   `json:"enrollment_time_out"`
	NumParticipants        uint32          `json:"num_participants"`
	EntryStake             string          `json:"entry_stake"`
	TotalReward            string          `json:"total_reward"`
	RewardTokenContract    eos.AccountName `json:"reward_token_contract"`
	NumParticipantsEntered uint32          `json:"num_participants_entered"`
	NumClaimedReturns      uint32          `json:"num_claimed_returns"`
	NumUnstaked            uint32          `json:"num_unstaked"`
	NumEarlyExits          uint32          `json:"num_early_exits"`
	CurrentState           eos.Name        `json:"current_state"`
	TotalDeposits          string          `json:"total_deposits"`
	Winners                Winners         `json:"winners"`
	BeneficiaryReward      string          `json:"beneficiary_reward"`
	MinParticipantReward   string          `json:"min_participant_reward"`
	TotalEarlyExitStake    string          `json:"total_early_exit_stake"`
	TotalEarlyExitReward   string          `json:"total_early_exit_reward"`
	RoundManager           eos.AccountName `json:"round_manager"`
	StakedTime             string          `json:"staked_time"`
	CreatedDate            string          `json:"created_date"`
	UpdatedDate            string          `json:"updated_date"`
}

func (m *Round) NumEntriesToClose() uint32 {
	return m.NumParticipants - m.NumParticipantsEntered
}

type NewRoundArgs struct {
	TermID               uint64          `json:"term_id"`
	RoundName            string          `json:"round_name"`
	StakingPeriodHrs     uint32          `json:"staking_period_hrs"`
	EnrollmentTimeOutHrs uint32          `json:"enrollment_time_out_hrs"`
	NumParticipants      uint32          `json:"num_participants"`
	EntryStake           string          `json:"entry_stake"`
	TotalReward          string          `json:"total_reward"`
	RoundManager         eos.AccountName `json:"round_manager"`
}

func RoundToNewRoundArgs(round *Round) *NewRoundArgs {
	return &NewRoundArgs{
		TermID:               round.TermID,
		RoundName:            round.RoundName,
		StakingPeriodHrs:     uint32(round.StakingPeriod.Hrs()),
		EnrollmentTimeOutHrs: uint32(round.EnrollmentTimeOut.Hrs()),
		NumParticipants:      round.NumParticipants,
		EntryStake:           round.EntryStake,
		TotalReward:          round.TotalReward,
		RoundManager:         round.RoundManager,
	}
}

func (m *Round) Clone() *Round {
	return &Round{
		RoundID:                m.RoundID,
		TermID:                 m.TermID,
		RoundName:              m.RoundName,
		StakingPeriod:          m.StakingPeriod,
		EnrollmentTimeOut:      m.EnrollmentTimeOut,
		NumParticipants:        m.NumParticipants,
		EntryStake:             m.EntryStake,
		TotalReward:            m.TotalReward,
		RewardTokenContract:    m.RewardTokenContract,
		NumParticipantsEntered: m.NumParticipantsEntered,
		NumClaimedReturns:      m.NumClaimedReturns,
		NumUnstaked:            m.NumUnstaked,
		NumEarlyExits:          m.NumEarlyExits,
		CurrentState:           m.CurrentState,
		TotalDeposits:          m.TotalDeposits,
		Winners:                m.Winners,
		BeneficiaryReward:      m.BeneficiaryReward,
		MinParticipantReward:   m.MinParticipantReward,
		TotalEarlyExitStake:    m.TotalEarlyExitStake,
		TotalEarlyExitReward:   m.TotalEarlyExitReward,
		RoundManager:           m.RoundManager,
		StakedTime:             m.StakedTime,
		CreatedDate:            m.CreatedDate,
		UpdatedDate:            m.UpdatedDate,
	}
}

type Entry struct {
	EntryID      uint64          `json:"entry_id"`
	RoundID      uint64          `json:"round_id"`
	Position     uint64          `json:"position"`
	Participant  eos.AccountName `json:"participant"`
	EntryStake   string          `json:"entry_stake"`
	ReturnAmount string          `json:"return_amount"`
	EntryStatus  eos.Name        `json:"entry_status"`
	EnteredDate  string          `json:"entered_date"`
}

func (m *BennyfiContract) NewTerm(term *Term) (string, error) {
	return m.NewTermFromTermArgs(TermToNewTermArgs(term))
}

func (m *BennyfiContract) NewTermFromTermArgs(termArgs *NewTermArgs) (string, error) {
	actionData := make(map[string]interface{})
	actionData["round_manager"] = termArgs.RoundManager
	actionData["term_name"] = termArgs.TermName
	actionData["all_participants_perc_x100000"] = termArgs.AllParticipantsPerc
	actionData["beneficiary"] = termArgs.Beneficiary
	actionData["beneficiary_perc_x100000"] = termArgs.BeneficiaryPerc

	return m.ExecAction(termArgs.RoundManager, "newterm", actionData)
}

func (m *BennyfiContract) NewRound(round *Round) (string, error) {
	return m.NewRoundFromRoundArgs(RoundToNewRoundArgs(round))
}

func (m *BennyfiContract) NewRoundFromRoundArgs(roundArgs *NewRoundArgs) (string, error) {
	actionData := make(map[string]interface{})
	actionData["round_manager"] = roundArgs.RoundManager
	actionData["round_name"] = roundArgs.RoundName
	actionData["term_id"] = roundArgs.TermID
	actionData["entry_stake"] = roundArgs.EntryStake
	actionData["total_reward"] = roundArgs.TotalReward
	actionData["num_participants"] = roundArgs.NumParticipants
	actionData["staking_period_hrs"] = roundArgs.StakingPeriodHrs
	actionData["enrollment_time_out_hrs"] = roundArgs.EnrollmentTimeOutHrs

	return m.ExecAction(roundArgs.RoundManager, "newround", actionData)
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

func (m *BennyfiContract) TimedEvents() (string, error) {
	return m.ExecAction(fmt.Sprintf("%v@open", m.ContractName), "timedevents", nil)
}

func (m *BennyfiContract) TstLapseTime(roundId uint64) (string, error) {
	actionData := make(map[string]interface{})
	actionData["round_id"] = roundId
	return m.ExecAction(eos.AN(m.ContractName), "tstlapsetime", actionData)
}

func (m *BennyfiContract) ReceiveRand(actor eos.AccountName, roundId uint64, randomNumbers []uint64) (string, error) {
	actionData := make(map[string]interface{})
	actionData["assoc_id"] = roundId

	proofs := make([]map[string]interface{}, 0, len(randomNumbers))
	for _, randomNumber := range randomNumbers {
		proofs = append(proofs, createEOSProof(randomNumber))
	}
	actionData["proofs"] = proofs

	return m.ExecAction(actor, "receiverand", actionData)
}

func (m *BennyfiContract) GetTerms() ([]Term, error) {
	return m.GetTermsReq(nil)
}

func (m *BennyfiContract) GetTermsbyManager(termManager eos.AccountName) ([]Term, error) {
	request := &eos.GetTableRowsRequest{
		Index:      "2",
		KeyType:    "name",
		LowerBound: string(termManager),
	}
	return m.GetTermsReq(request)
}

func (m *BennyfiContract) GetLastTerm() (*Term, error) {
	terms, err := m.GetTermsReq(&eos.GetTableRowsRequest{
		Reverse: true,
		Limit:   1,
	})
	if err != nil {
		return nil, err
	}
	if len(terms) > 0 {
		return &terms[0], nil
	}
	return nil, nil
}

func (m *BennyfiContract) GetTermsReq(req *eos.GetTableRowsRequest) ([]Term, error) {

	var terms []Term
	if req == nil {
		req = &eos.GetTableRowsRequest{}
	}
	req.Table = "terms"
	err := m.GetTableRows(*req, &terms)
	if err != nil {
		return nil, fmt.Errorf("get table rows %v", err)
	}
	return terms, nil
}

func (m *BennyfiContract) GetRounds() ([]Round, error) {

	return m.GetRoundsReq(nil)
}

func (m *BennyfiContract) GetRoundsbyManager(roundManager eos.AccountName) ([]Round, error) {
	request := &eos.GetTableRowsRequest{
		Index:      "2",
		KeyType:    "name",
		LowerBound: string(roundManager),
		UpperBound: string(roundManager),
	}
	return m.GetRoundsReq(request)
}

func (m *BennyfiContract) GetRoundsbyTerm(termID uint64) ([]Round, error) {
	request := &eos.GetTableRowsRequest{}
	m.FilterRoundsbyTerm(request, termID)
	return m.GetRoundsReq(request)
}

func (m *BennyfiContract) FilterRoundsbyTerm(req *eos.GetTableRowsRequest, termID uint64) {
	req.Index = "3"
	req.KeyType = "i64"
	req.LowerBound = strconv.FormatUint(termID, 10)
	req.UpperBound = strconv.FormatUint(termID, 10)
}

func (m *BennyfiContract) GetRoundsbyState(state eos.Name) ([]Round, error) {
	request := &eos.GetTableRowsRequest{}
	m.FilterRoundsbyState(request, state)
	return m.GetRoundsReq(request)
}

func (m *BennyfiContract) FilterRoundsbyState(req *eos.GetTableRowsRequest, state eos.Name) {
	req.Index = "4"
	req.KeyType = "name"
	req.LowerBound = string(state)
	req.UpperBound = req.LowerBound
}

func (m *BennyfiContract) GetRound(roundID uint64) (*Round, error) {
	rounds, err := m.GetRoundsReq(&eos.GetTableRowsRequest{
		LowerBound: strconv.FormatUint(roundID, 10),
		UpperBound: strconv.FormatUint(roundID, 10),
		Limit:      1,
	})
	if err != nil {
		return nil, err
	}
	if len(rounds) > 0 {
		return &rounds[0], nil
	}
	return nil, nil
}

func (m *BennyfiContract) GetLastRound() (*Round, error) {
	rounds, err := m.GetRoundsReq(&eos.GetTableRowsRequest{
		Reverse: true,
		Limit:   1,
	})
	if err != nil {
		return nil, err
	}
	if len(rounds) > 0 {
		return &rounds[0], nil
	}
	return nil, nil
}

func (m *BennyfiContract) GetRoundsReq(req *eos.GetTableRowsRequest) ([]Round, error) {

	var rounds []Round
	if req == nil {
		req = &eos.GetTableRowsRequest{}
	}
	req.Table = "rounds"
	err := m.GetTableRows(*req, &rounds)
	if err != nil {
		return nil, fmt.Errorf("get table rows %v", err)
	}
	return rounds, nil
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

func createEOSProof(randomNumber uint64) map[string]interface{} {
	proof := make(map[string]interface{})
	proof["block_num"] = 1
	proof["block_id"] = "blockid"
	proof["seed"] = 1
	proof["final_seed"] = "finalseed"
	proof["public_key"] = "publickey"
	proof["gamma"] = "gamma"
	proof["c"] = "c"
	proof["s"] = "s"
	proof["output_u256"] = "u256"
	proof["output_u64"] = randomNumber
	return proof
}
