package exsatfrs

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	evmbase "github.com/sebastianmontero/bennyfi-go-client/evm/base"
	"github.com/sebastianmontero/bennyfi-go-client/evm/common/eth"
)

// IFixedRateStakingABI is the ABI of the FixedRateStaking contract.
// Note: We flatten the output tuple to simplify dynamic decoding into []interface{}.
const IFixedRateStakingABI = `[
	{
		"inputs": [
			{ "internalType": "address", "name": "agent", "type": "address" },
			{ "internalType": "bytes32", "name": "subId", "type": "bytes32" }
		],
		"name": "getPosition",
		"outputs": [
			{ "internalType": "address", "name": "agent", "type": "address" },
			{ "internalType": "bytes32", "name": "subId", "type": "bytes32" },
			{ "internalType": "uint256", "name": "principal", "type": "uint256" },
			{ "internalType": "uint256", "name": "rate", "type": "uint256" },
			{ "internalType": "uint256", "name": "startTime", "type": "uint256" },
			{ "internalType": "uint256", "name": "duration", "type": "uint256" },
			{ "internalType": "uint256", "name": "lastSettlementTime", "type": "uint256" },
			{ "internalType": "uint256", "name": "totalDistributed", "type": "uint256" },
			{ "internalType": "uint256", "name": "pendingReturn", "type": "uint256" },
			{ "internalType": "bool", "name": "isSettled", "type": "bool" },
			{ "internalType": "bool", "name": "isWithdrawn", "type": "bool" }
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{ "internalType": "address", "name": "agent", "type": "address" },
			{ "internalType": "bytes32", "name": "subId", "type": "bytes32" },
			{ "internalType": "uint256", "name": "settlementTimestamp", "type": "uint256" }
		],
		"name": "distribute",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{ "internalType": "bytes32", "name": "subId", "type": "bytes32" }
		],
		"name": "claimReturn",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{ "internalType": "bytes32", "name": "subId", "type": "bytes32" }
		],
		"name": "withdraw",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{ "internalType": "address", "name": "agent", "type": "address" },
			{ "internalType": "bytes32", "name": "subId", "type": "bytes32" }
		],
		"name": "unlockPosition",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{ "internalType": "address", "name": "agent", "type": "address" },
			{ "internalType": "bytes32", "name": "subId", "type": "bytes32" },
			{ "internalType": "uint256", "name": "intervalDays", "type": "uint256" }
		],
		"name": "shiftStartTime",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"anonymous": false,
		"inputs": [
			{ "indexed": true, "internalType": "address", "name": "agent", "type": "address" },
			{ "indexed": true, "internalType": "bytes32", "name": "subId", "type": "bytes32" },
			{ "indexed": false, "internalType": "uint256", "name": "amount", "type": "uint256" },
			{ "indexed": false, "internalType": "uint256", "name": "settlementTimestamp", "type": "uint256" },
			{ "indexed": false, "internalType": "bool", "name": "isFinal", "type": "bool" }
		],
		"name": "ReturnDistributed",
		"type": "event"
	},
	{
		"inputs": [],
		"name": "stakingToken",
		"outputs": [
			{
				"internalType": "contract IERC20",
				"name": "",
				"type": "address"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{ "internalType": "address", "name": "agent", "type": "address" },
			{ "internalType": "bytes32", "name": "subId", "type": "bytes32" },
			{ "internalType": "uint256", "name": "settlementTimestamp", "type": "uint256" }
		],
		"name": "calculateReturn",
		"outputs": [
			{ "internalType": "uint256", "name": "", "type": "uint256" }
		],
		"stateMutability": "view",
		"type": "function"
	}
]`

// ERC20ApproveABI is a minimal ABI for the ERC20 approve function.
const ERC20ApproveABI = `[
	{
		"constant": false,
		"inputs": [
			{
				"name": "_spender",
				"type": "address"
			},
			{
				"name": "_value",
				"type": "uint256"
			}
		],
		"name": "approve",
		"outputs": [
			{
				"name": "",
				"type": "bool"
			}
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	}
]`

// Position represents the fixed-rate position details.
type Position struct {
	Agent              common.Address
	SubId              [32]byte
	Principal          *big.Int
	Rate               *big.Int
	StartTime          *big.Int
	Duration           *big.Int
	LastSettlementTime *big.Int
	TotalDistributed   *big.Int
	PendingReturn      *big.Int
	IsSettled          bool
	IsWithdrawn        bool
}

// ReturnDistributedEvent represents a ReturnDistributed event emitted by the FixedRateStaking contract.
type ReturnDistributedEvent struct {
	Agent               common.Address
	SubId               [32]byte
	Amount              *big.Int
	SettlementTimestamp *big.Int
	IsFinal             bool
	Raw                 types.Log
}

// String returns a string representation of the ReturnDistributedEvent.
func (e *ReturnDistributedEvent) String() string {
	subId := new(big.Int).SetBytes(e.SubId[:]).Uint64()
	return fmt.Sprintf("Agent: %s, SubId: %d, Amount: %s, SettlementTimestamp: %s, IsFinal: %v",
		e.Agent.Hex(), subId, e.Amount.String(), e.SettlementTimestamp.String(), e.IsFinal)
}

// ExSatFRSRead interacts with the FixedRateStaking contract for read-only operations.
type ExSatFRSRead struct {
	*evmbase.ReadClient
}

// ExSatFRSWrite interacts with the FixedRateStaking contract for write operations.
type ExSatFRSWrite struct {
	*ExSatFRSRead
	*evmbase.WriteClient
}

// NewRead creates a new ExSatFRSRead client for read-only operations.
func NewRead(rpcUrl string, contractAddr common.Address) (*ExSatFRSRead, error) {
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to dial rpc: %w", err)
	}
	return NewReadWithClient(client, contractAddr)
}

// NewReadWithClient creates a new ExSatFRSRead client using an existing contract backend.
func NewReadWithClient(client eth.EthClient, contractAddr common.Address) (*ExSatFRSRead, error) {
	baseClient, err := evmbase.NewReadWithClient(client, contractAddr, IFixedRateStakingABI)
	if err != nil {
		return nil, err
	}
	return &ExSatFRSRead{ReadClient: baseClient}, nil
}

// NewWriteField creates a new ExSatFRSWrite client for write operations.
func NewWrite(rpcUrl string, privateKeyHex string, contractAddr common.Address) (*ExSatFRSWrite, error) {
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to dial rpc: %w", err)
	}
	return NewWriteWithClient(client, privateKeyHex, contractAddr)
}

// NewWriteWithClient creates a new ExSatFRSWrite client with an existing ethclient.
func NewWriteWithClient(client eth.EthClient, privateKeyHex string, contractAddr common.Address) (*ExSatFRSWrite, error) {
	baseClient, err := evmbase.NewWithClient(client, privateKeyHex, contractAddr, IFixedRateStakingABI)
	if err != nil {
		return nil, err
	}

	readClient := &ExSatFRSRead{ReadClient: baseClient.ReadClient}

	return &ExSatFRSWrite{
		ExSatFRSRead: readClient,
		WriteClient:  baseClient,
	}, nil
}

// GetPosition retrieves the position details for a given agent and subId.
func (c *ExSatFRSRead) GetPosition(agent common.Address, subId uint64) (*Position, error) {
	var out []interface{}

	// Convert uint64 subId to bytes32 (Big Endian)
	var subIdBytes [32]byte
	bigSubId := new(big.Int).SetUint64(subId)
	bigSubIdBytes := bigSubId.Bytes()
	copy(subIdBytes[32-len(bigSubIdBytes):], bigSubIdBytes)

	err := c.Contract.Call(nil, &out, "getPosition", agent, subIdBytes)
	if err != nil {
		return nil, err
	}

	principal := out[2].(*big.Int)
	if principal.Sign() == 0 {
		return nil, fmt.Errorf("position not found for agent %s and subId %d", agent.Hex(), subId)
	}

	return &Position{
		Agent:              out[0].(common.Address),
		SubId:              out[1].([32]byte),
		Principal:          principal,
		Rate:               out[3].(*big.Int),
		StartTime:          out[4].(*big.Int),
		Duration:           out[5].(*big.Int),
		LastSettlementTime: out[6].(*big.Int),
		TotalDistributed:   out[7].(*big.Int),
		PendingReturn:      out[8].(*big.Int),
		IsSettled:          out[9].(bool),
		IsWithdrawn:        out[10].(bool),
	}, nil
}

// StakingToken retrieves the address of the staking token.
func (c *ExSatFRSRead) StakingToken() (common.Address, error) {
	var out []interface{}
	err := c.Contract.Call(nil, &out, "stakingToken")
	if err != nil {
		return common.Address{}, err
	}
	return out[0].(common.Address), nil
}

// CalculateReturn calculates the return for a given agent and subId.
func (c *ExSatFRSRead) CalculateReturn(agent common.Address, subId uint64, settlementTimestamp *big.Int) (*big.Int, error) {
	var out []interface{}

	// Convert uint64 subId to bytes32 (Big Endian)
	var subIdBytes [32]byte
	bigSubId := new(big.Int).SetUint64(subId)
	bigSubIdBytes := bigSubId.Bytes()
	copy(subIdBytes[32-len(bigSubIdBytes):], bigSubIdBytes)

	err := c.Contract.Call(nil, &out, "calculateReturn", agent, subIdBytes, settlementTimestamp)
	if err != nil {
		return nil, err
	}

	return out[0].(*big.Int), nil
}

// Distribute calls the distribute action.
func (c *ExSatFRSWrite) Distribute(agent common.Address, subId uint64, settlementTimestamp *big.Int) (string, error) {
	// 1. Calculate the return amount to distribute
	amount, err := c.CalculateReturn(agent, subId, settlementTimestamp)
	if err != nil {
		return "", fmt.Errorf("failed to calculate return for agent %s, subId %d, settlementTimestamp %s: %w", agent.Hex(), subId, settlementTimestamp.String(), err)
	}

	// 2. Approve the staking token to be spent by the ExSatFRS contract
	stakingTokenAddr, err := c.StakingToken()
	if err != nil {
		return "", fmt.Errorf("failed to get staking token address: %w", err)
	}

	parsedABI, err := abi.JSON(strings.NewReader(ERC20ApproveABI))
	if err != nil {
		return "", fmt.Errorf("failed to parse ERC20 Approve ABI: %w", err)
	}

	tokenContract := bind.NewBoundContract(stakingTokenAddr, parsedABI, c.WriteClient.Client, c.WriteClient.Client, c.WriteClient.Client)

	// Get pending nonce to sequence transactions correctly
	nonce, err := c.WriteClient.Client.PendingNonceAt(context.Background(), c.Auth.From)
	if err != nil {
		return "", fmt.Errorf("failed to get pending nonce for %s: %w", c.Auth.From.Hex(), err)
	}

	approveAuth := &bind.TransactOpts{
		From:      c.Auth.From,
		Nonce:     new(big.Int).SetUint64(nonce),
		Signer:    c.Auth.Signer,
		Value:     c.Auth.Value,
		GasPrice:  c.Auth.GasPrice,
		GasFeeCap: c.Auth.GasFeeCap,
		GasTipCap: c.Auth.GasTipCap,
		GasLimit:  c.Auth.GasLimit,
		Context:   c.Auth.Context,
		NoSend:    c.Auth.NoSend,
	}

	txApprove, err := tokenContract.Transact(approveAuth, "approve", c.WriteClient.Address, amount)
	if err != nil {
		return "", fmt.Errorf("failed to approve token amount %s for contract %s: %w", amount.String(), c.WriteClient.Address.Hex(), err)
	}

	// Wait for the approve transaction to be mined
	receipt, err := bind.WaitMined(context.Background(), c.WriteClient.Client, txApprove)
	if err != nil {
		return "", fmt.Errorf("failed to wait for approve transaction to be mined: %w", err)
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		return "", fmt.Errorf("approve transaction reverted")
	}

	// 2. Execute distribute using the updated pending nonce
	distributeNonce, err := c.WriteClient.Client.PendingNonceAt(context.Background(), c.Auth.From)
	if err != nil {
		return "", fmt.Errorf("failed to get pending nonce for %s: %w", c.Auth.From.Hex(), err)
	}

	distributeAuth := &bind.TransactOpts{
		From:      c.Auth.From,
		Nonce:     new(big.Int).SetUint64(distributeNonce),
		Signer:    c.Auth.Signer,
		Value:     c.Auth.Value,
		GasPrice:  c.Auth.GasPrice,
		GasFeeCap: c.Auth.GasFeeCap,
		GasTipCap: c.Auth.GasTipCap,
		GasLimit:  c.Auth.GasLimit,
		Context:   c.Auth.Context,
		NoSend:    c.Auth.NoSend,
	}

	var subIdBytes [32]byte
	bigSubId := new(big.Int).SetUint64(subId)
	bigSubIdBytes := bigSubId.Bytes()
	copy(subIdBytes[32-len(bigSubIdBytes):], bigSubIdBytes)

	tx, err := c.WriteClient.Contract.Transact(distributeAuth, "distribute", agent, subIdBytes, settlementTimestamp)
	if err != nil {
		err = c.WriteClient.TryGetRevertReason(err, "distribute", agent, subIdBytes, settlementTimestamp)
		return "", fmt.Errorf("failed to distribute for agent %s, subId %d: %w", agent.Hex(), subId, err)
	}
	return tx.Hash().Hex(), nil
}

// ClaimReturn calls claimReturn.
func (c *ExSatFRSWrite) ClaimReturn(subId uint64) (string, error) {
	var subIdBytes [32]byte
	bigSubId := new(big.Int).SetUint64(subId)
	bigSubIdBytes := bigSubId.Bytes()
	copy(subIdBytes[32-len(bigSubIdBytes):], bigSubIdBytes)

	tx, err := c.WriteClient.Contract.Transact(c.Auth, "claimReturn", subIdBytes)
	if err != nil {
		err = c.WriteClient.TryGetRevertReason(err, "claimReturn", subIdBytes)
		return "", fmt.Errorf("failed to claimReturn for subId %d: %w", subId, err)
	}
	return tx.Hash().Hex(), nil
}

// Withdraw calls withdraw.
func (c *ExSatFRSWrite) Withdraw(subId uint64) (string, error) {
	var subIdBytes [32]byte
	bigSubId := new(big.Int).SetUint64(subId)
	bigSubIdBytes := bigSubId.Bytes()
	copy(subIdBytes[32-len(bigSubIdBytes):], bigSubIdBytes)

	tx, err := c.WriteClient.Contract.Transact(c.Auth, "withdraw", subIdBytes)
	if err != nil {
		err = c.WriteClient.TryGetRevertReason(err, "withdraw", subIdBytes)
		return "", fmt.Errorf("failed to withdraw for subId %d: %w", subId, err)
	}
	return tx.Hash().Hex(), nil
}

// UnlockPosition unlocks a position by updating its start time.
func (c *ExSatFRSWrite) UnlockPosition(agent common.Address, subId uint64) (string, error) {
	var subIdBytes [32]byte
	bigSubId := new(big.Int).SetUint64(subId)
	bigSubIdBytes := bigSubId.Bytes()
	copy(subIdBytes[32-len(bigSubIdBytes):], bigSubIdBytes)

	tx, err := c.WriteClient.Contract.Transact(c.Auth, "unlockPosition", agent, subIdBytes)
	if err != nil {
		err = c.WriteClient.TryGetRevertReason(err, "unlockPosition", agent, subIdBytes)
		return "", fmt.Errorf("failed to unlock position for agent %s, subId %d: %w", agent.Hex(), subId, err)
	}
	return tx.Hash().Hex(), nil
}

// ShiftStartTime shifts the startTime of a position to the start of the next interval.
func (c *ExSatFRSWrite) ShiftStartTime(agent common.Address, subId uint64, intervalDays *big.Int) (string, error) {
	var subIdBytes [32]byte
	bigSubId := new(big.Int).SetUint64(subId)
	bigSubIdBytes := bigSubId.Bytes()
	copy(subIdBytes[32-len(bigSubIdBytes):], bigSubIdBytes)

	tx, err := c.WriteClient.Contract.Transact(c.Auth, "shiftStartTime", agent, subIdBytes, intervalDays)
	if err != nil {
		err = c.WriteClient.TryGetRevertReason(err, "shiftStartTime", agent, subIdBytes, intervalDays)
		return "", fmt.Errorf("failed to shift start time for agent %s, subId %d, intervalDays %s: %w", agent.Hex(), subId, intervalDays.String(), err)
	}
	return tx.Hash().Hex(), nil
}

func (c *ExSatFRSRead) filterEvents(ctx context.Context, eventName string, startBlock uint64, endBlock *uint64, agent *common.Address) ([]types.Log, *abi.ABI, error) {
	parsedABI, err := abi.JSON(strings.NewReader(IFixedRateStakingABI))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse ABI: %w", err)
	}

	eventAbi, ok := parsedABI.Events[eventName]
	if !ok {
		return nil, nil, fmt.Errorf("event %s not found in ABI", eventName)
	}

	var topics [][]common.Hash
	topics = append(topics, []common.Hash{eventAbi.ID})

	if agent != nil {
		topics = append(topics, []common.Hash{common.BytesToHash(agent.Bytes())})
	}

	query := ethereum.FilterQuery{
		FromBlock: new(big.Int).SetUint64(startBlock),
		Addresses: []common.Address{c.Address},
		Topics:    topics,
	}
	if endBlock != nil {
		query.ToBlock = new(big.Int).SetUint64(*endBlock)
	}

	logs, err := c.Client.FilterLogs(ctx, query)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to filter logs: %w", err)
	}

	return logs, &parsedABI, nil
}

// FilterReturnDistributedEvents retrieves ReturnDistributed events for a specific agent from a starting block.
// endBlock is optional (can be nil).
func (c *ExSatFRSRead) FilterReturnDistributedEvents(ctx context.Context, startBlock uint64, endBlock *uint64, agent common.Address) ([]*ReturnDistributedEvent, error) {
	logs, parsedABI, err := c.filterEvents(ctx, "ReturnDistributed", startBlock, endBlock, &agent)
	if err != nil {
		return nil, err
	}

	var events []*ReturnDistributedEvent
	for _, vLog := range logs {
		var event ReturnDistributedEvent

		err := parsedABI.UnpackIntoInterface(&event, "ReturnDistributed", vLog.Data)
		if err != nil {
			return nil, fmt.Errorf("failed to unpack event data: %w", err)
		}

		if len(vLog.Topics) > 1 {
			event.Agent = common.BytesToAddress(vLog.Topics[1].Bytes())
		}
		if len(vLog.Topics) > 2 {
			event.SubId = vLog.Topics[2]
		}

		event.Raw = vLog
		events = append(events, &event)
	}

	return events, nil
}
