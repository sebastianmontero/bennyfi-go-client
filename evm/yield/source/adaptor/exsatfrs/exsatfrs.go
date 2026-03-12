package exsatfrs

import (
	"context"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sebastianmontero/bennyfi-go-client/evm/common/eth"
	"github.com/sebastianmontero/bennyfi-go-client/evm/yield/source/adaptor/base"
	"golang.org/x/sync/errgroup"
)

// ExSatBankFixedRateStakingYieldSourceAdaptorABI is the ABI of the ExSatBankFixedRateStakingYieldSourceAdaptor contract.
const ExSatBankFixedRateStakingYieldSourceAdaptorABI = `[
    {
      "inputs": [
        {
          "internalType": "uint64",
          "name": "",
          "type": "uint64"
        }
      ],
      "name": "stakes",
      "outputs": [
        {
          "internalType": "uint64",
          "name": "poolId",
          "type": "uint64"
        },
        {
          "internalType": "uint256",
          "name": "totalStake",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "stakeTime",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "unlockTime",
          "type": "uint256"
        },
        {
          "internalType": "uint8",
          "name": "state",
          "type": "uint8"
        },
        {
          "internalType": "uint256",
          "name": "totalReturn",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint64",
          "name": "_poolId",
          "type": "uint64"
        }
      ],
      "name": "unstake",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint64",
          "name": "_poolId",
          "type": "uint64"
        }
      ],
      "name": "unlockStake",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "fixedStakingAddress",
      "outputs": [
        {
          "internalType": "address",
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
          "internalType": "uint64",
          "name": "_poolId",
          "type": "uint64"
        }
      ],
      "name": "claimPartialReturn",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
		{
      "inputs": [],
      "name": "minPartialReturnThreshold",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    }
]`

// StakeState represents the state of a stake.
type StakeState uint8

const (
	StakeStateUninitialized StakeState = 0
	StakeStateActive        StakeState = 1
	StakeStateUnstaked      StakeState = 2
)

// Stake represents the stake information returned by the smart contract.
type Stake struct {
	PoolId      uint64
	TotalStake  *big.Int
	StakeTime   *big.Int
	UnlockTime  *big.Int
	State       StakeState
	TotalReturn *big.Int
}

// ExSatFRSRead interacts with the ExSatBankFixedRateStakingYieldSourceAdaptor contract for read-only operations.
type ExSatFRSRead struct {
	*base.BaseRead
}

// ExSatFRSWrite interacts with the ExSatBankFixedRateStakingYieldSourceAdaptor contract for write operations.
type ExSatFRSWrite struct {
	*ExSatFRSRead
	*base.BaseWrite
}

// New creates a new ExSatFRSWrite client.
// It dials the RPC URL, parses the private key, and sets up the transaction authorizer.
func New(rpcUrl string, privateKeyHex string, contractAddr common.Address) (*ExSatFRSWrite, error) {
	// 1. Connect to Eth Client
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		return nil, err
	}

	return NewWithClient(client, privateKeyHex, contractAddr)
}

// NewWithClient creates a new ExSatFRSWrite client using an existing contract backend and private key.
func NewWithClient(client eth.EthClient, privateKeyHex string, contractAddr common.Address) (*ExSatFRSWrite, error) {
	// 1. Create BaseWrite Client (merged ABI handled in base)
	baseClient, err := base.NewWithClient(client, privateKeyHex, contractAddr, ExSatBankFixedRateStakingYieldSourceAdaptorABI)
	if err != nil {
		return nil, err
	}

	// 2. Create Read Client wrapper (reuse BaseRead from BaseWrite)
	readClient := &ExSatFRSRead{
		BaseRead: baseClient.BaseRead,
	}

	return &ExSatFRSWrite{
		ExSatFRSRead: readClient,
		BaseWrite:    baseClient,
	}, nil
}

// NewRead creates a new ExSatFRSRead client for read-only operations.
// It dials the RPC URL but does not set up a transaction authorizer.
func NewRead(rpcUrl string, contractAddr common.Address) (*ExSatFRSRead, error) {
	// 1. Connect to Eth Client
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		return nil, err
	}

	return NewReadWithClient(client, contractAddr)
}

// NewReadWithClient creates a new ExSatFRSRead client using an existing contract backend.
func NewReadWithClient(client eth.EthClient, contractAddr common.Address) (*ExSatFRSRead, error) {
	// 1. Create BaseRead Client with merged ABI
	baseClient, err := base.NewReadWithClient(client, contractAddr, ExSatBankFixedRateStakingYieldSourceAdaptorABI)
	if err != nil {
		return nil, err
	}

	return &ExSatFRSRead{
		BaseRead: baseClient,
	}, nil
}

// GetStake retrieves the stake information for a given pool ID.
func (c *ExSatFRSRead) GetStake(poolId uint64) (*Stake, error) {
	var out []interface{}
	err := c.BaseRead.Contract.Call(nil, &out, "stakes", poolId)
	if err != nil {
		return nil, err
	}

	return &Stake{
		PoolId:      out[0].(uint64),
		TotalStake:  out[1].(*big.Int),
		StakeTime:   out[2].(*big.Int),
		UnlockTime:  out[3].(*big.Int),
		State:       StakeState(out[4].(uint8)),
		TotalReturn: out[5].(*big.Int),
	}, nil
}

// GetStakes retrieves multiple stake information for a given list of pool IDs in parallel.
func (c *ExSatFRSRead) GetStakes(poolIds []uint64) (map[uint64]*Stake, error) {
	results := make(map[uint64]*Stake, len(poolIds))
	var mu sync.Mutex
	g, _ := errgroup.WithContext(context.Background())
	g.SetLimit(20)

	for _, pid := range poolIds {
		pid := pid // capture variable
		g.Go(func() error {
			stake, err := c.GetStake(pid)
			if err != nil {
				return err
			}

			mu.Lock()
			results[pid] = stake
			mu.Unlock()
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}
	return results, nil
}

// FixedStakingAddress retrieves the fixed staking address.
func (c *ExSatFRSRead) FixedStakingAddress() (common.Address, error) {
	var out []interface{}
	err := c.BaseRead.Contract.Call(nil, &out, "fixedStakingAddress")
	if err != nil {
		return common.Address{}, err
	}

	return out[0].(common.Address), nil
}

// MinPartialReturnThreshold retrieves the minimum partial return threshold.
func (c *ExSatFRSRead) MinPartialReturnThreshold() (*big.Int, error) {
	var out []interface{}
	err := c.Contract.Call(nil, &out, "minPartialReturnThreshold")
	if err != nil {
		return nil, err
	}
	return out[0].(*big.Int), nil
}

// Unstake calls the unstake function on the contract.
func (c *ExSatFRSWrite) Unstake(poolId uint64) (*types.Transaction, error) {
	return c.BaseWrite.Contract.Transact(c.Auth, "unstake", poolId)
}

// UnlockStake calls the unlockStake function on the contract.
func (c *ExSatFRSWrite) UnlockStake(poolId uint64) (*types.Transaction, error) {
	return c.BaseWrite.Contract.Transact(c.Auth, "unlockStake", poolId)
}

// ClaimPartialReturn calls the claimPartialReturn function on the contract.
func (c *ExSatFRSWrite) ClaimPartialReturn(poolId uint64) (*types.Transaction, error) {
	return c.BaseWrite.Contract.Transact(c.Auth, "claimPartialReturn", poolId)
}
