package exsatfs

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sebastianmontero/bennyfi-go-client/common/eth"
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
    }
]`

// Position represents the fixed-term position details.
type Position struct {
	Principal    *big.Int
	StartTime    *big.Int
	Duration     *big.Int
	ReturnAmount *big.Int
	IsSettled    bool
	IsWithdrawn  bool
}

// ExSatFS interacts with the IFixedStaking contract.
type ExSatFS struct {
	contract *bind.BoundContract
	address  common.Address
}

// NewRead creates a new ExSatFS client for read-only operations.
func NewRead(rpcUrl string, contractAddr common.Address) (*ExSatFS, error) {
	// 1. Connect to Eth Client
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to dial rpc: %w", err)
	}

	return NewReadWithClient(client, contractAddr)
}

// NewReadWithClient creates a new ExSatFS client using an existing contract backend.
func NewReadWithClient(client eth.EthClient, contractAddr common.Address) (*ExSatFS, error) {
	// 1. Parse ABI
	parsedABI, err := abi.JSON(strings.NewReader(IFixedStakingABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse abi: %w", err)
	}

	// 2. Create Bound Contract
	contract := bind.NewBoundContract(contractAddr, parsedABI, client, client, client)
	return &ExSatFS{
		contract: contract,
		address:  contractAddr,
	}, nil
}

// StakingToken retrieves the address of the staking token.
func (c *ExSatFS) StakingToken() (common.Address, error) {
	var out []interface{}
	err := c.contract.Call(nil, &out, "stakingToken")
	if err != nil {
		return common.Address{}, err
	}
	return out[0].(common.Address), nil
}

// GetPosition retrieves the position details for a given agent and subId.
func (c *ExSatFS) GetPosition(agent common.Address, subId uint64) (*Position, error) {
	var out []interface{}
	
	// Convert uint64 subId to bytes32 (Big Endian)
	var subIdBytes [32]byte
	bigSubId := new(big.Int).SetUint64(subId)
	bigSubIdBytes := bigSubId.Bytes()
	copy(subIdBytes[32-len(bigSubIdBytes):], bigSubIdBytes)

	err := c.contract.Call(nil, &out, "getPosition", agent, subIdBytes)
	if err != nil {
		return nil, err
	}

	return &Position{
		Principal:    out[0].(*big.Int),
		StartTime:    out[1].(*big.Int),
		Duration:     out[2].(*big.Int),
		ReturnAmount: out[3].(*big.Int),
		IsSettled:    out[4].(bool),
		IsWithdrawn:  out[5].(bool),
	}, nil
}
