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
	RoundPending               = eos.Name("pending")
	RoundAcceptingEntries      = eos.Name("acceptentrys")
	RoundDrawing               = eos.Name("rounddrawing")
	RoundOpen                  = eos.Name("roundopen")
	RoundClosed                = eos.Name("roundclosed")
	RoundUnlocked              = eos.Name("rndunlocked")
	RoundUnlockedUnstaked      = eos.Name("rndunlckdutk")
	RoundTimedOut              = eos.Name("rndtimedout")
	RoundTimedOutUnstaked      = eos.Name("rndtmdoututk")
	RoundTypeManagerFunded     = eos.Name("mgrfunded")
	RoundTypeRexPool           = eos.Name("rexpool")
	RoundAccessPrivate         = eos.Name("private")
	RoundAccessPublic          = eos.Name("public")
	RexStateNotApplicable      = eos.Name("notaplicable")
	RexStatePreRex             = eos.Name("prerex")
	RexStateInSavings          = eos.Name("insavings")
	RexStateInLockPeriod       = eos.Name("lockperiod")
	RexStateSold               = eos.Name("sold")
	RexStateWithdrawn          = eos.Name("withdrawn")
	RewardFundingStatePending  = eos.Name("pending")
	RewardFundingStateFunded   = eos.Name("funded")
	RewardFundingStateRefunded = eos.Name("refunded")
	RewardFundingStateCommited = eos.Name("commited")
	RewardFundingStateRex      = eos.Name("rex")
	DistributionMain           = "main"
	DistributionProjectToken   = "project-token"
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

type Round struct {
	RoundID                  uint64                   `json:"round_id"`
	TermID                   uint64                   `json:"term_id"`
	ProjectID                uint64                   `json:"project_id"`
	RoundName                string                   `json:"round_name"`
	RoundType                eos.Name                 `json:"round_type"`
	RoundAccess              eos.Name                 `json:"round_access"`
	StakingPeriod            *Microseconds            `json:"staking_period"`
	EnrollmentTimeOut        *Microseconds            `json:"enrollment_time_out"`
	NumParticipants          uint32                   `json:"num_participants"`
	RoundFee                 string                   `json:"round_fee"`
	EntryStake               string                   `json:"entry_stake"`
	Rewards                  Rewards                  `json:"rewards"`
	RexBalance               string                   `json:"rex_balance"`
	NumParticipantsEntered   uint32                   `json:"num_participants_entered"`
	NumClaimedReturns        uint32                   `json:"num_claimed_returns"`
	NumUnstaked              uint32                   `json:"num_unstaked"`
	NumEarlyExits            uint32                   `json:"num_early_exits"`
	CurrentState             eos.Name                 `json:"current_state"`
	RexState                 eos.Name                 `json:"rex_state"`
	TotalDeposits            string                   `json:"total_deposits"`
	Winners                  Winners                  `json:"winners"`
	Beneficiary              eos.AccountName          `json:"beneficiary"`
	Distributions            Distributions            `json:"distributions"`
	TotalEarlyExitStake      string                   `json:"total_early_exit_stake"`
	TotalEarlyExitRewardFees TotalEarlyExitRewardFees `json:"total_early_exit_reward_fees"`
	RoundManager             eos.AccountName          `json:"round_manager"`
	ClosedTime               string                   `json:"closed_time"`
	StakedTime               string                   `json:"staked_time"`
	MovedFromSavingsTime     string                   `json:"moved_from_savings_time"`
	StakeEndTime             string                   `json:"stake_end_time"`
	EnrollmentTimeEnd        string                   `json:"enrollment_time_end"`
	CreatedDate              string                   `json:"created_date"`
	UpdatedDate              string                   `json:"updated_date"`
}

func (m *Round) GetTotalDeposits() eos.Asset {
	totalDeposits, err := eos.NewAssetFromString(m.TotalDeposits)
	if err != nil {
		panic(fmt.Sprintf("Unable to parse total deposits: %v to asset", m.TotalDeposits))
	}
	return totalDeposits
}

func (m *Round) NumEntriesToClose() uint32 {
	return m.NumParticipants - m.NumParticipantsEntered
}

func (m *Round) UpsertDistribution(name string, distribution *Distribution) {
	if m.Distributions == nil {
		m.Distributions = make(Distributions, 0, 1)
	}
	m.Distributions.Upsert(name, distribution)
}

func (m *Round) AssignWinnerPrizes(distName string, dist *Distribution) error {
	winnersEntry := m.Winners.Find(distName)
	if winnersEntry == nil {
		return fmt.Errorf("failed assigning winner prizes, there is no winners array for distribution name: %v", distName)
	}
	err := winnersEntry.Value.AssignPrizes(dist.WinnerPrizes)
	if err != nil {
		return fmt.Errorf("failed assigning winner prizes for dist: %v, error: %v", distName, err)
	}
	return nil
}

func (m *Round) RemoveDistribution(name string) {
	m.Distributions.Remove(name)
}

func (m *Round) UpsertReward(name string, reward *Reward) {
	if m.Rewards == nil {
		m.Rewards = make(Rewards, 0, 1)
	}
	m.Rewards.Upsert(name, reward)
}

func (m *Round) RemoveReward(name string) {
	m.Rewards.Remove(name)
}

func (m *Round) UpsertEarlyExitRewardFee(name string, earlyExitRewardFee string) {
	if m.TotalEarlyExitRewardFees == nil {
		m.TotalEarlyExitRewardFees = make(TotalEarlyExitRewardFees, 0, 1)
	}
	m.TotalEarlyExitRewardFees.Upsert(name, earlyExitRewardFee)
}

func (m *Round) RemoveEarlyExitRewardFee(name string) {
	m.TotalEarlyExitRewardFees.Remove(name)
}

func (m *Round) UpsertWinner(name string, winner *Winner) {
	if m.Winners == nil {
		m.Winners = make(Winners, 0, 1)
	}
	m.Winners.Upsert(name, winner)
}

func (m *Round) RemoveWinner(name string) {
	m.Winners.Remove(name)
}

func (m *Round) UpdateFundingStateAll(state eos.Name) {
	m.Rewards.UpdateFundingStateAll(state)
}

func (m *Round) UpdateFundingState(dist string, state eos.Name) {
	m.Rewards.UpdateFundingState(dist, state)
}

type NewRoundArgs struct {
	TermID               uint64          `json:"term_id"`
	ProjectID            uint64          `json:"project_id"`
	RoundName            string          `json:"round_name"`
	StakingPeriodHrs     uint32          `json:"staking_period_hrs"`
	EnrollmentTimeOutHrs uint32          `json:"enrollment_time_out_hrs"`
	NumParticipants      uint32          `json:"num_participants"`
	RoundFee             string          `json:"round_fee"`
	EntryStake           string          `json:"entry_stake"`
	Rewards              RewardsArg      `json:"rewards"`
	RoundManager         eos.AccountName `json:"round_manager"`
}

func RoundToNewRoundArgs(round *Round) *NewRoundArgs {
	return &NewRoundArgs{
		TermID:               round.TermID,
		ProjectID:            round.ProjectID,
		RoundName:            round.RoundName,
		StakingPeriodHrs:     uint32(round.StakingPeriod.Hrs()),
		EnrollmentTimeOutHrs: uint32(round.EnrollmentTimeOut.Hrs()),
		NumParticipants:      round.NumParticipants,
		RoundFee:             round.RoundFee,
		EntryStake:           round.EntryStake,
		Rewards:              round.Rewards.ToRewardsArg(),
		RoundManager:         round.RoundManager,
	}
}

func (m *Round) Clone() *Round {
	return &Round{
		RoundID:                  m.RoundID,
		TermID:                   m.TermID,
		ProjectID:                m.ProjectID,
		RoundName:                m.RoundName,
		RoundType:                m.RoundType,
		RoundAccess:              m.RoundAccess,
		StakingPeriod:            m.StakingPeriod,
		EnrollmentTimeOut:        m.EnrollmentTimeOut,
		NumParticipants:          m.NumParticipants,
		RoundFee:                 m.RoundFee,
		EntryStake:               m.EntryStake,
		Rewards:                  m.Rewards.Clone(),
		RexBalance:               m.RexBalance,
		NumParticipantsEntered:   m.NumParticipantsEntered,
		NumClaimedReturns:        m.NumClaimedReturns,
		NumUnstaked:              m.NumUnstaked,
		NumEarlyExits:            m.NumEarlyExits,
		CurrentState:             m.CurrentState,
		RexState:                 m.RexState,
		TotalDeposits:            m.TotalDeposits,
		Winners:                  m.Winners,
		Beneficiary:              m.Beneficiary,
		Distributions:            m.Distributions,
		TotalEarlyExitStake:      m.TotalEarlyExitStake,
		TotalEarlyExitRewardFees: m.TotalEarlyExitRewardFees,
		RoundManager:             m.RoundManager,
		ClosedTime:               m.ClosedTime,
		StakedTime:               m.StakedTime,
		MovedFromSavingsTime:     m.MovedFromSavingsTime,
		StakeEndTime:             m.StakeEndTime,
		EnrollmentTimeEnd:        m.EnrollmentTimeEnd,
		CreatedDate:              m.CreatedDate,
		UpdatedDate:              m.UpdatedDate,
	}
}

func (m *BennyfiContract) NewRound(round *Round) (string, error) {
	return m.NewRoundFromRoundArgs(RoundToNewRoundArgs(round))
}

func (m *BennyfiContract) NewRoundFromRoundArgs(roundArgs *NewRoundArgs) (string, error) {
	actionData := make(map[string]interface{})
	actionData["round_manager"] = roundArgs.RoundManager
	actionData["round_name"] = roundArgs.RoundName
	actionData["term_id"] = roundArgs.TermID
	actionData["project_id"] = roundArgs.ProjectID
	actionData["round_fee"] = roundArgs.RoundFee
	actionData["entry_stake"] = roundArgs.EntryStake
	actionData["rewards"] = roundArgs.Rewards
	actionData["num_participants"] = roundArgs.NumParticipants
	actionData["staking_period_hrs"] = roundArgs.StakingPeriodHrs
	actionData["enrollment_time_out_hrs"] = roundArgs.EnrollmentTimeOutHrs

	return m.ExecAction(roundArgs.RoundManager, "newround", actionData)
}

func (m *BennyfiContract) FundRound(roundID uint64, beneficiary interface{}) (string, error) {
	actionData := make(map[string]interface{})
	actionData["round_id"] = roundID
	actionData["beneficiary"] = beneficiary

	return m.ExecAction(beneficiary, "fundround", actionData)
}

func (m *BennyfiContract) TimedEvents() (string, error) {
	return m.ExecAction(fmt.Sprintf("%v@open", m.ContractName), "timedevents", nil)
}

func (m *BennyfiContract) TimeoutRounds(callCounter uint64) (string, error) {
	actionData := make(map[string]interface{})
	actionData["call_counter"] = callCounter
	return m.ExecAction(fmt.Sprintf("%v@open", m.ContractName), "timeoutrnds", actionData)
}

func (m *BennyfiContract) MoveFromSavings(callCounter uint64) (string, error) {
	actionData := make(map[string]interface{})
	actionData["call_counter"] = callCounter
	return m.ExecAction(fmt.Sprintf("%v@open", m.ContractName), "mvfrmsavings", actionData)
}

func (m *BennyfiContract) SellRex(callCounter uint64) (string, error) {
	actionData := make(map[string]interface{})
	actionData["call_counter"] = callCounter
	return m.ExecAction(fmt.Sprintf("%v@open", m.ContractName), "sellrex", actionData)
}

func (m *BennyfiContract) WithdrawRex(callCounter uint64) (string, error) {
	actionData := make(map[string]interface{})
	actionData["call_counter"] = callCounter
	return m.ExecAction(fmt.Sprintf("%v@open", m.ContractName), "withdrawrex", actionData)
}

func (m *BennyfiContract) UnlockRounds(callCounter uint64) (string, error) {
	actionData := make(map[string]interface{})
	actionData["call_counter"] = callCounter
	return m.ExecAction(fmt.Sprintf("%v@open", m.ContractName), "unlockrnds", actionData)
}

func (m *BennyfiContract) UnstakeUnlockedRounds(callCounter uint64) (string, error) {
	actionData := make(map[string]interface{})
	actionData["call_counter"] = callCounter
	return m.ExecAction(fmt.Sprintf("%v@open", m.ContractName), "ustkulckrnds", actionData)
}

func (m *BennyfiContract) UnstakeTimedoutRounds(callCounter uint64) (string, error) {
	actionData := make(map[string]interface{})
	actionData["call_counter"] = callCounter
	return m.ExecAction(fmt.Sprintf("%v@open", m.ContractName), "ustktmdrnds", actionData)
}

func (m *BennyfiContract) Redraw() (string, error) {
	return m.ExecAction(fmt.Sprintf("%v@open", m.ContractName), "redraw", nil)
}

func (m *BennyfiContract) TstLapseTime(roundId uint64) (string, error) {
	actionData := make(map[string]interface{})
	actionData["round_id"] = roundId
	return m.ExecAction(eos.AN(m.ContractName), "tstlapsetime", actionData)
}

func (m *BennyfiContract) ReceiveRand(actor eos.AccountName, roundId uint64, randomNumber string) (string, error) {
	actionData := make(map[string]interface{})
	actionData["assoc_id"] = roundId
	actionData["random"] = randomNumber
	return m.ExecAction(actor, "receiverand", actionData)
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

func (m *BennyfiContract) GetRoundsbyStateAndId(state eos.Name) ([]Round, error) {
	request := &eos.GetTableRowsRequest{}
	err := m.FilterRoundsbyStateAndId(request, state)
	if err != nil {
		return nil, err
	}
	return m.GetRoundsReq(request)
}

func (m *BennyfiContract) FilterRoundsbyStateAndId(req *eos.GetTableRowsRequest, state eos.Name) error {

	req.Index = "7"
	req.KeyType = "i128"
	req.Reverse = true
	stateAndRndLB, err := m.EOS.GetComposedIndexValue(state, 0)
	if err != nil {
		return fmt.Errorf("failed to generate lower bound composed index, err: %v", err)
	}
	stateAndRndUB, err := m.EOS.GetComposedIndexValue(state, uint64(18446744073709551615))
	if err != nil {
		return fmt.Errorf("failed to generate upper bound composed index, err: %v", err)
	}
	fmt.Println("LB: ", stateAndRndLB, "UB: ", stateAndRndUB)
	req.LowerBound = stateAndRndLB
	req.UpperBound = stateAndRndUB
	return err
}

func (m *BennyfiContract) GetRoundsbyManagerAndId(manager interface{}) ([]Round, error) {
	request := &eos.GetTableRowsRequest{}
	err := m.FilterRoundsbyManagerAndId(request, manager)
	if err != nil {
		return nil, err
	}
	return m.GetRoundsReq(request)
}

func (m *BennyfiContract) FilterRoundsbyManagerAndId(req *eos.GetTableRowsRequest, manager interface{}) error {

	req.Index = "8"
	req.KeyType = "i128"
	req.Reverse = true
	mgrAndRndLB, err := m.EOS.GetComposedIndexValue(manager, 0)
	if err != nil {
		return fmt.Errorf("failed to generate lower bound composed index, err: %v", err)
	}
	mgrAndRndUB, err := m.EOS.GetComposedIndexValue(manager, uint64(18446744073709551615))
	if err != nil {
		return fmt.Errorf("failed to generate upper bound composed index, err: %v", err)
	}
	fmt.Println("LB: ", mgrAndRndLB, "UB: ", mgrAndRndUB)
	req.LowerBound = mgrAndRndLB
	req.UpperBound = mgrAndRndUB
	return err
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
