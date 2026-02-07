package bennydummy

import (
	"fmt"
	"strconv"

	"github.com/sebastianmontero/bennyfi-go-client/bennyfi"
	eos "github.com/sebastianmontero/eos-go"
	"github.com/sebastianmontero/eos-go-toolbox/contract"
	"github.com/sebastianmontero/eos-go-toolbox/service"
	"github.com/sebastianmontero/eos-go-toolbox/util"
	"github.com/sebastianmontero/eos-go/token"
)

type BennyDummyContract struct {
	*contract.Contract
	bennyfiContract *bennyfi.BennyfiContract
	callCounter     uint64
}

func NewBennyDummyContract(eos *service.EOS, contractName string) *BennyDummyContract {
	return &BennyDummyContract{
		contract.NewContract(eos, contractName),
		bennyfi.NewBennyfiContract(eos, contractName),
		0,
	}
}

type Config struct {
	StakeLocalContract eos.Name `json:"stakelocal_contract"`
}
type Transfer struct {
	TransferID uint64          `json:"transfer_id"`
	From       eos.AccountName `json:"from"`
	To         eos.AccountName `json:"to"`
	Amount     eos.Asset       `json:"amount"`
	Memo       string          `json:"memo"`
}

func (m *BennyDummyContract) ExecAction(permissionLevel interface{}, action string, actionData interface{}) (string, error) {
	resp, err := m.Contract.ExecAction(permissionLevel, action, actionData)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Tx ID: %v", resp.TransactionID), nil
}

func (m *BennyDummyContract) SetConfig(stakeLocalContract eos.Name) (string, error) {
	actionData := struct {
		StakeLocalContract eos.Name
	}{stakeLocalContract}
	return m.ExecAction(m.ContractName, "setconfig", actionData)
}

func (m *BennyDummyContract) SetTerm(termId uint64, yieldSource eos.Name) (string, error) {
	actionData := struct {
		TermId      uint64
		YieldSource eos.Name
	}{termId, yieldSource}
	return m.ExecAction(m.ContractName, "setterm", actionData)
}

func (m *BennyDummyContract) SetPool(poolId uint64, termId uint64, currentState eos.Name, stakingPeriodHrs uint32) (string, error) {
	actionData := struct {
		PoolId           uint64
		TermId           uint64
		CurrentState     eos.Name
		PoolType         eos.Name
		StakingPeriodHrs uint32
	}{poolId, termId, currentState, bennyfi.RoundTypeYield, stakingPeriodHrs}
	return m.ExecAction(m.ContractName, "setpool", actionData)
}

func (m *BennyDummyContract) GetStakeTransfer(to, quantity interface{}, poolId uint64) (*token.Transfer, error) {

	toAN, err := util.ToAccountName(to)
	if err != nil {
		return nil, err
	}

	qty, err := util.ToAsset(quantity)
	if err != nil {
		return nil, err
	}

	return &token.Transfer{
		From:     eos.AccountName(m.ContractName),
		To:       toAN,
		Quantity: qty,
		Memo:     fmt.Sprintf("pool id: %v", poolId),
	}, nil
}

func (m *BennyDummyContract) GetConfig(req *eos.GetTableRowsRequest) (*Config, error) {

	var config []Config
	if req == nil {
		req = &eos.GetTableRowsRequest{}
	}
	req.Table = "config"
	err := m.GetTableRows(*req, &config)
	if err != nil {
		return nil, fmt.Errorf("get table rows %v", err)
	}
	if len(config) == 0 {
		return nil, nil
	}
	return &config[0], nil
}

func (m *BennyDummyContract) GetTransfer(transferID uint64) (*Transfer, error) {
	transfers, err := m.GetTransfersReq(&eos.GetTableRowsRequest{
		LowerBound: strconv.FormatUint(transferID, 10),
		UpperBound: strconv.FormatUint(transferID, 10),
		Limit:      1,
	})
	if err != nil {
		return nil, err
	}
	if len(transfers) > 0 {
		return &transfers[0], nil
	}
	return nil, nil
}

func (m *BennyDummyContract) GetLastTransfer() (*Transfer, error) {
	transfers, err := m.GetTransfersReq(&eos.GetTableRowsRequest{
		Reverse: true,
		Limit:   1,
	})
	if err != nil {
		return nil, err
	}
	if len(transfers) > 0 {
		return &transfers[0], nil
	}
	return nil, nil
}

func (m *BennyDummyContract) GetTransfersReq(req *eos.GetTableRowsRequest) ([]Transfer, error) {

	var transfers []Transfer
	if req == nil {
		req = &eos.GetTableRowsRequest{}
	}
	req.Table = "transfers"
	err := m.GetTableRows(*req, &transfers)
	if err != nil {
		return nil, fmt.Errorf("get table rows %v", err)
	}
	return transfers, nil
}

func (m *BennyDummyContract) GetRound(roundID uint64) (*bennyfi.Round, error) {
	return m.bennyfiContract.GetRound(roundID)
}

func (m *BennyDummyContract) GetRoundsReq(req *eos.GetTableRowsRequest) ([]bennyfi.Round, error) {
	return m.bennyfiContract.GetRoundsReq(req)
}

func (m *BennyDummyContract) GetTerms(termID uint64) (*bennyfi.Term, error) {
	return m.bennyfiContract.GetTermsById(termID)
}

func (m *BennyDummyContract) GetTermsReq(req *eos.GetTableRowsRequest) ([]bennyfi.Term, error) {
	return m.bennyfiContract.GetTermsReq(req)
}
