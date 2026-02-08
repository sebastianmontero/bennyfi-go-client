package erc20dummy

import (
	"fmt"
	"strconv"

	eos "github.com/sebastianmontero/eos-go"
	"github.com/sebastianmontero/eos-go-toolbox/contract"
	"github.com/sebastianmontero/eos-go-toolbox/service"
)

type Token struct {
	ID                        uint64          `json:"id"`
	TokenContract             eos.AccountName `json:"token_contract"`
	Address                   eos.HexBytes    `json:"address"`
	IngressFee                eos.Asset       `json:"ingress_fee"`
	Balance                   eos.Asset       `json:"balance"`
	FeeBalance                eos.Asset       `json:"fee_balance"`
	Erc20Precision            uint8           `json:"erc20_precision"`
	FromEvmToNative           uint8           `json:"from_evm_to_native,omitempty"`
	OriginalErc20TokenAddress eos.HexBytes    `json:"original_erc20_token_address,omitempty"`
	MinIngress                eos.Asset       `json:"min_ingress,omitempty"`
}

type ERC20DummyContract struct {
	*contract.Contract
	callCounter uint64
}

func NewERC20DummyContract(eos *service.EOS, contractName string) *ERC20DummyContract {
	return &ERC20DummyContract{
		contract.NewContract(eos, contractName),
		0,
	}
}

func (m *ERC20DummyContract) ExecAction(permissionLevel interface{}, action string, actionData interface{}) (string, error) {
	resp, err := m.Contract.ExecAction(permissionLevel, action, actionData)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Tx ID: %v", resp.TransactionID), nil
}

func (m *ERC20DummyContract) SetToken(token *Token) (string, error) {
	return m.ExecAction(m.ContractName, "settoken", token)
}

func (m *ERC20DummyContract) GetToken(id uint64) (*Token, error) {
	tokens, err := m.GetTokensReq(&eos.GetTableRowsRequest{
		LowerBound: strconv.FormatUint(id, 10),
		UpperBound: strconv.FormatUint(id, 10),
		Limit:      1,
	})
	if err != nil {
		return nil, err
	}
	if len(tokens) > 0 {
		return &tokens[0], nil
	}
	return nil, nil
}

func (m *ERC20DummyContract) GetTokenByContractSymbol(contract eos.AccountName, symbol eos.Symbol) (*Token, error) {
	key, err := m.EOS.GetComposedIndexValue(contract, symbol)
	if err != nil {
		return nil, err
	}
	tokens, err := m.GetTokensReq(&eos.GetTableRowsRequest{
		Index:      "2",
		KeyType:    "i128",
		LowerBound: key,
		UpperBound: key,
		Limit:      1,
	})
	if err != nil {
		return nil, err
	}
	if len(tokens) > 0 {
		return &tokens[0], nil
	}
	return nil, nil
}

func (m *ERC20DummyContract) GetTokensReq(req *eos.GetTableRowsRequest) ([]Token, error) {
	var tokens []Token
	if req == nil {
		req = &eos.GetTableRowsRequest{}
	}
	req.Table = "tokens"
	err := m.GetTableRows(*req, &tokens)
	if err != nil {
		return nil, fmt.Errorf("get table rows %v", err)
	}
	return tokens, nil
}
