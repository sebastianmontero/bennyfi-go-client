package exsatfs

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

// IFixedStakingABI is the ABI of the IFixedStaking contract.
// Note: We flatten the output tuple to simplify dynamic decoding into []interface{}.
const IFixedStakingABI = `[
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "agent",
          "type": "address"
        },
        {
          "internalType": "bytes32",
          "name": "subId",
          "type": "bytes32"
        }
      ],
      "name": "getPosition",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "principal",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "startTime",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "duration",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "returnAmount",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "totalCoupons",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "pendingRewards",
          "type": "uint256"
        },
        {
          "internalType": "bool",
          "name": "isSettled",
          "type": "bool"
        },
        {
          "internalType": "bool",
          "name": "isWithdrawn",
          "type": "bool"
        }
      ],
      "stateMutability": "view",
      "type": "function"
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
        {
          "internalType": "address",
          "name": "agent",
          "type": "address"
        },
        {
          "internalType": "bytes32",
          "name": "subId",
          "type": "bytes32"
        },
        {
          "internalType": "uint256",
          "name": "returnAmount",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "topupAmount",
          "type": "uint256"
        },
        {
          "internalType": "bool",
          "name": "forceEarly",
          "type": "bool"
        }
      ],
      "name": "settle",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "agent",
          "type": "address"
        },
        {
          "internalType": "bytes32",
          "name": "subId",
          "type": "bytes32"
        }
      ],
      "name": "unlockPosition",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "agent",
          "type": "address"
        },
        {
          "internalType": "bytes32",
          "name": "subId",
          "type": "bytes32"
        },
        {
          "internalType": "uint256",
          "name": "amount",
          "type": "uint256"
        }
      ],
      "name": "distributeReturn",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "agent",
          "type": "address"
        },
        {
          "indexed": true,
          "internalType": "bytes32",
          "name": "subId",
          "type": "bytes32"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "returnAmount",
          "type": "uint256"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "topupAmount",
          "type": "uint256"
        },
        {
          "indexed": false,
          "internalType": "bool",
          "name": "forcedEarly",
          "type": "bool"
        }
      ],
      "name": "Settled",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "agent",
          "type": "address"
        },
        {
          "indexed": true,
          "internalType": "bytes32",
          "name": "subId",
          "type": "bytes32"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "amount",
          "type": "uint256"
        }
      ],
      "name": "CouponDistributed",
      "type": "event"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "agent",
          "type": "address"
        },
        {
          "internalType": "bytes32",
          "name": "subId",
          "type": "bytes32"
        },
        {
          "internalType": "uint256",
          "name": "intervalDays",
          "type": "uint256"
        }
      ],
      "name": "shiftStartTime",
      "outputs": [],
      "stateMutability": "nonpayable",
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

// Position represents the fixed-term position details.
type Position struct {
	Principal      *big.Int
	StartTime      *big.Int
	Duration       *big.Int
	ReturnAmount   *big.Int
	TotalCoupons   *big.Int
	PendingRewards *big.Int
	IsSettled      bool
	IsWithdrawn    bool
}

// SettledEvent represents a Settled event emitted by the IFixedStaking contract.
type SettledEvent struct {
	Agent        common.Address
	SubId        [32]byte
	ReturnAmount *big.Int
	TopupAmount  *big.Int
	ForcedEarly  bool
	Raw          types.Log
}

// String returns a string representation of the SettledEvent.
func (e *SettledEvent) String() string {
	subId := new(big.Int).SetBytes(e.SubId[:]).Uint64()
	return fmt.Sprintf("Agent: %s, SubId: %d, ReturnAmount: %s, TopupAmount: %s, ForcedEarly: %v",
		e.Agent.Hex(), subId, e.ReturnAmount.String(), e.TopupAmount.String(), e.ForcedEarly)
}

// CouponDistributedEvent represents a CouponDistributed event emitted by the IFixedStaking contract.
type CouponDistributedEvent struct {
	Agent  common.Address
	SubId  [32]byte
	Amount *big.Int
	Raw    types.Log
}

// String returns a string representation of the CouponDistributedEvent.
func (e *CouponDistributedEvent) String() string {
	subId := new(big.Int).SetBytes(e.SubId[:]).Uint64()
	return fmt.Sprintf("Agent: %s, SubId: %d, Amount: %s",
		e.Agent.Hex(), subId, e.Amount.String())
}

// ExSatFSRead interacts with the IFixedStaking contract for read-only operations.
type ExSatFSRead struct {
	*evmbase.ReadClient
}

// ExSatFSWrite interacts with the IFixedStaking contract for write operations.
type ExSatFSWrite struct {
	*ExSatFSRead
	*evmbase.WriteClient
}

// NewRead creates a new ExSatFSRead client for read-only operations.
func NewRead(rpcUrl string, contractAddr common.Address) (*ExSatFSRead, error) {
	// 1. Connect to Eth Client
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to dial rpc: %w", err)
	}

	return NewReadWithClient(client, contractAddr)
}

// NewReadWithClient creates a new ExSatFSRead client using an existing contract backend.
func NewReadWithClient(client eth.EthClient, contractAddr common.Address) (*ExSatFSRead, error) {
	// 1. Create ReadClient with merged ABI
	baseClient, err := evmbase.NewReadWithClient(client, contractAddr, IFixedStakingABI)
	if err != nil {
		return nil, err
	}

	return &ExSatFSRead{
		ReadClient: baseClient,
	}, nil
}

// NewWriteField creates a new ExSatFSWrite client for write operations.
func NewWrite(rpcUrl string, privateKeyHex string, contractAddr common.Address) (*ExSatFSWrite, error) {
	// 1. Connect to Eth Client
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to dial rpc: %w", err)
	}

	return NewWriteWithClient(client, privateKeyHex, contractAddr)
}

// NewWriteWithClient creates a new ExSatFSWrite client with an existing ethclient.
func NewWriteWithClient(client eth.EthClient, privateKeyHex string, contractAddr common.Address) (*ExSatFSWrite, error) {
	// 1. Create WriteClient with merged ABI
	baseClient, err := evmbase.NewWithClient(client, privateKeyHex, contractAddr, IFixedStakingABI)
	if err != nil {
		return nil, err
	}

	readClient := &ExSatFSRead{
		ReadClient: baseClient.ReadClient,
	}

	return &ExSatFSWrite{
		ExSatFSRead: readClient,
		WriteClient: baseClient,
	}, nil
}

// StakingToken retrieves the address of the staking token.
func (c *ExSatFSRead) StakingToken() (common.Address, error) {
	var out []interface{}
	err := c.Contract.Call(nil, &out, "stakingToken")
	if err != nil {
		return common.Address{}, err
	}
	return out[0].(common.Address), nil
}

// GetPosition retrieves the position details for a given agent and subId.
func (c *ExSatFSRead) GetPosition(agent common.Address, subId uint64) (*Position, error) {
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

	principal := out[0].(*big.Int)
	if principal.Sign() == 0 {
		return nil, fmt.Errorf("position not found for agent %s and subId %d", agent.Hex(), subId)
	}

	return &Position{
		Principal:      principal,
		StartTime:      out[1].(*big.Int),
		Duration:       out[2].(*big.Int),
		ReturnAmount:   out[3].(*big.Int),
		TotalCoupons:   out[4].(*big.Int),
		PendingRewards: out[5].(*big.Int),
		IsSettled:      out[6].(bool),
		IsWithdrawn:    out[7].(bool),
	}, nil
}

// Settle settles a position, defining the return amount.
func (c *ExSatFSWrite) Settle(agent common.Address, subId uint64, returnAmount *big.Int, topupAmount *big.Int, forceEarly bool) (string, error) {
	// Convert uint64 subId to bytes32 (Big Endian)
	var subIdBytes [32]byte
	bigSubId := new(big.Int).SetUint64(subId)
	bigSubIdBytes := bigSubId.Bytes()
	copy(subIdBytes[32-len(bigSubIdBytes):], bigSubIdBytes)

	tx, err := c.WriteClient.Contract.Transact(c.Auth, "settle", agent, subIdBytes, returnAmount, topupAmount, forceEarly)
	if err != nil {
		err = c.WriteClient.TryGetRevertReason(err, "settle", agent, subIdBytes, returnAmount, topupAmount, forceEarly)
		return "", fmt.Errorf("failed to settle position for agent %s, subId %d, returnAmount %s, topupAmount %s, forceEarly %t: %w", agent.Hex(), subId, returnAmount.String(), topupAmount.String(), forceEarly, err)
	}
	return tx.Hash().Hex(), nil
}

// UnlockPosition unlocks a position by updating its start time.
func (c *ExSatFSWrite) UnlockPosition(agent common.Address, subId uint64) (string, error) {
	// Convert uint64 subId to bytes32 (Big Endian)
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

// DistributeReturn distributes an interim return (coupon) to a specific position.
func (c *ExSatFSWrite) DistributeReturn(agent common.Address, subId uint64, amount *big.Int) (string, error) {
	// 1. Approve the staking token to be spent by the ExSatFS contract
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
		From:     c.Auth.From,
		Nonce:    new(big.Int).SetUint64(nonce),
		Signer:   c.Auth.Signer,
		Value:    c.Auth.Value,
		GasPrice: c.Auth.GasPrice,
		GasFeeCap: c.Auth.GasFeeCap,
		GasTipCap: c.Auth.GasTipCap,
		GasLimit: c.Auth.GasLimit,
		Context:  c.Auth.Context,
		NoSend:   c.Auth.NoSend,
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

	// 2. Execute distributeReturn using the updated pending nonce
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

	// Convert uint64 subId to bytes32 (Big Endian)
	var subIdBytes [32]byte
	bigSubId := new(big.Int).SetUint64(subId)
	bigSubIdBytes := bigSubId.Bytes()
	copy(subIdBytes[32-len(bigSubIdBytes):], bigSubIdBytes)

	tx, err := c.WriteClient.Contract.Transact(distributeAuth, "distributeReturn", agent, subIdBytes, amount)
	if err != nil {
		err = c.WriteClient.TryGetRevertReason(err, "distributeReturn", agent, subIdBytes, amount)
		return "", fmt.Errorf("failed to distribute return for agent %s, subId %d, amount %s: %w", agent.Hex(), subId, amount.String(), err)
	}
	return tx.Hash().Hex(), nil
}

// ShiftStartTime shifts the startTime of a position to the start of the next interval.
func (c *ExSatFSWrite) ShiftStartTime(agent common.Address, subId uint64, intervalDays *big.Int) (string, error) {
	// Convert uint64 subId to bytes32 (Big Endian)
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

func (c *ExSatFSRead) filterEvents(ctx context.Context, eventName string, startBlock uint64, endBlock *uint64, agent *common.Address) ([]types.Log, *abi.ABI, error) {
	parsedABI, err := abi.JSON(strings.NewReader(IFixedStakingABI))
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

// FilterSettledEvents retrieves Settled events for a specific agent from a starting block.
// endBlock is optional (can be nil).
func (c *ExSatFSRead) FilterSettledEvents(ctx context.Context, startBlock uint64, endBlock *uint64, agent common.Address) ([]*SettledEvent, error) {
	logs, parsedABI, err := c.filterEvents(ctx, "Settled", startBlock, endBlock, &agent)
	if err != nil {
		return nil, err
	}

	var events []*SettledEvent
	for _, vLog := range logs {
		var event SettledEvent

		err := parsedABI.UnpackIntoInterface(&event, "Settled", vLog.Data)
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

// FilterCouponDistributedEvents retrieves CouponDistributed events for a specific agent from a starting block.
// endBlock is optional (can be nil).
func (c *ExSatFSRead) FilterCouponDistributedEvents(ctx context.Context, startBlock uint64, endBlock *uint64, agent common.Address) ([]*CouponDistributedEvent, error) {
	logs, parsedABI, err := c.filterEvents(ctx, "CouponDistributed", startBlock, endBlock, &agent)
	if err != nil {
		return nil, err
	}

	var events []*CouponDistributedEvent
	for _, vLog := range logs {
		var event CouponDistributedEvent

		err := parsedABI.UnpackIntoInterface(&event, "CouponDistributed", vLog.Data)
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
