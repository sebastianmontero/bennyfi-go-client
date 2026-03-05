package exsatfs

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
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

// CouponDistributedEvent represents a CouponDistributed event emitted by the IFixedStaking contract.
type CouponDistributedEvent struct {
	Agent  common.Address
	SubId  [32]byte
	Amount *big.Int
	Raw    types.Log
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
		return "", fmt.Errorf("failed to settle position: %w", err)
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
		return "", fmt.Errorf("failed to unlock position: %w", err)
	}
	return tx.Hash().Hex(), nil
}

// DistributeReturn distributes an interim return (coupon) to a specific position.
func (c *ExSatFSWrite) DistributeReturn(agent common.Address, subId uint64, amount *big.Int) (string, error) {
	// Convert uint64 subId to bytes32 (Big Endian)
	var subIdBytes [32]byte
	bigSubId := new(big.Int).SetUint64(subId)
	bigSubIdBytes := bigSubId.Bytes()
	copy(subIdBytes[32-len(bigSubIdBytes):], bigSubIdBytes)

	tx, err := c.WriteClient.Contract.Transact(c.Auth, "distributeReturn", agent, subIdBytes, amount)
	if err != nil {
		return "", fmt.Errorf("failed to distribute return: %w", err)
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
