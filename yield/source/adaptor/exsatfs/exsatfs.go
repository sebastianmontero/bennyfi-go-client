package exsatfs

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/sync/errgroup"
)

// ExSatBankFixedStakingYieldSourceAdaptorABI is the ABI of the ExSatBankFixedStakingYieldSourceAdaptor contract.
const ExSatBankFixedStakingYieldSourceAdaptorABI = `[
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
	State       uint8
	TotalReturn *big.Int
}

// ExSatFS interacts with the ExSatBankFixedStakingYieldSourceAdaptor contract.
type ExSatFS struct {
	contract *bind.BoundContract
	address  common.Address
	auth     *bind.TransactOpts
}

// New creates a new ExSatFS client.
// It dials the RPC URL, parses the private key, and sets up the transaction authorizer.
func New(rpcUrl string, privateKeyHex string, contractAddr common.Address) (*ExSatFS, error) {
	// 1. Parse Private Key
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, err
	}

	// 2. Connect to Eth Client
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		return nil, err
	}

	// 3. Get Chain ID
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return nil, err
	}

	// 4. Create TransactOpts
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return nil, err
	}

	// 5. Parse ABI
	parsedABI, err := abi.JSON(strings.NewReader(ExSatBankFixedStakingYieldSourceAdaptorABI))
	if err != nil {
		return nil, err
	}

	// 6. Create Bound Contract
	contract := bind.NewBoundContract(contractAddr, parsedABI, client, client, client)
	return &ExSatFS{
		contract: contract,
		address:  contractAddr,
		auth:     auth,
	}, nil
}

// NewRead creates a new ExSatFS client for read-only operations.
// It dials the RPC URL but does not set up a transaction authorizer.
func NewRead(rpcUrl string, contractAddr common.Address) (*ExSatFS, error) {
	// 1. Connect to Eth Client
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		return nil, err
	}

	// 2. Parse ABI
	parsedABI, err := abi.JSON(strings.NewReader(ExSatBankFixedStakingYieldSourceAdaptorABI))
	if err != nil {
		return nil, err
	}

	// 3. Create Bound Contract
	contract := bind.NewBoundContract(contractAddr, parsedABI, client, client, client)
	return &ExSatFS{
		contract: contract,
		address:  contractAddr,
		auth:     nil,
	}, nil
}

// GetStake retrieves the stake information for a given pool ID.
func (c *ExSatFS) GetStake(poolId uint64) (*Stake, error) {
	var out []interface{}
	err := c.contract.Call(nil, &out, "stakes", poolId)
	if err != nil {
		return nil, err
	}

	return &Stake{
		PoolId:      out[0].(uint64),
		TotalStake:  out[1].(*big.Int),
		StakeTime:   out[2].(*big.Int),
		UnlockTime:  out[3].(*big.Int),
		State:       out[4].(uint8),
		TotalReturn: out[5].(*big.Int),
	}, nil
}

// GetStakes retrieves multiple stake information for a given list of pool IDs in parallel.
func (c *ExSatFS) GetStakes(poolIds []uint64) (map[uint64]*Stake, error) {
	results := make(map[uint64]*Stake, len(poolIds))
	var mu sync.Mutex
	g, _ := errgroup.WithContext(context.Background())

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

// Unstake calls the unstake function on the contract.
func (c *ExSatFS) Unstake(poolId uint64) (*types.Transaction, error) {
	if c.auth == nil {
		return nil, fmt.Errorf("client is read-only")
	}
	return c.contract.Transact(c.auth, "unstake", poolId)
}
