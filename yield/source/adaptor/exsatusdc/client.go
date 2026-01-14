package exsatusdc

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/sync/errgroup"
)

// ExSatBankYieldSourceAdaptorABI is the ABI of the ExSatBankYieldSourceAdaptor contract.
const ExSatBankYieldSourceAdaptorABI = `[
    {
      "inputs": [
        {
          "internalType": "uint64",
          "name": "_poolId",
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
          "name": "totalShares",
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
          "internalType": "enum ExSatBankYieldSourceAdaptor.StakeState",
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
      "name": "triggerUnstake",
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
	TotalShares *big.Int
	StakeTime   *big.Int
	UnlockTime  *big.Int
	State       uint8
	TotalReturn *big.Int
}

// Client interacts with the ExSatBankYieldSourceAdaptor contract.
type Client struct {
	contract *bind.BoundContract
	address  common.Address
	auth     *bind.TransactOpts
}

// NewClient creates a new Client.
// It dials the RPC URL, parses the private key, and sets up the transaction authorizer.
func NewClient(rpcUrl string, privateKeyHex string, contractAddr common.Address) (*Client, error) {
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
	parsedABI, err := abi.JSON(strings.NewReader(ExSatBankYieldSourceAdaptorABI))
	if err != nil {
		return nil, err
	}

	// 6. Create Bound Contract
	contract := bind.NewBoundContract(contractAddr, parsedABI, client, client, client)
	return &Client{
		contract: contract,
		address:  contractAddr,
		auth:     auth,
	}, nil
}

// NewReadClient creates a new Client for read-only operations.
// It dials the RPC URL but does not set up a transaction authorizer.
func NewReadClient(rpcUrl string, contractAddr common.Address) (*Client, error) {
	// 1. Connect to Eth Client
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		return nil, err
	}

	// 2. Parse ABI
	parsedABI, err := abi.JSON(strings.NewReader(ExSatBankYieldSourceAdaptorABI))
	if err != nil {
		return nil, err
	}

	// 3. Create Bound Contract
	contract := bind.NewBoundContract(contractAddr, parsedABI, client, client, client)
	return &Client{
		contract: contract,
		address:  contractAddr,
		auth:     nil,
	}, nil
}

// GetStake retrieves the stake information for a given pool ID.
func (c *Client) GetStake(poolId uint64) (*Stake, error) {
	var out []interface{}
	err := c.contract.Call(nil, &out, "stakes", poolId)
	if err != nil {
		return nil, err
	}

	return &Stake{
		PoolId:      out[0].(uint64),
		TotalStake:  out[1].(*big.Int),
		TotalShares: out[2].(*big.Int),
		StakeTime:   out[3].(*big.Int),
		UnlockTime:  out[4].(*big.Int),
		State:       out[5].(uint8),
		TotalReturn: out[6].(*big.Int),
	}, nil
}

// GetStakes retrieves multiple stake information for a given list of pool IDs in parallel.
func (c *Client) GetStakes(poolIds []uint64) ([]*Stake, error) {
	results := make([]*Stake, len(poolIds))
	g, _ := errgroup.WithContext(context.Background())

	for i, pid := range poolIds {
		i, pid := i, pid // capture variables
		g.Go(func() error {
			stake, err := c.GetStake(pid)
			if err != nil {
				return err
			}
			results[i] = stake
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}
	return results, nil
}

// TriggerUnstake calls the triggerUnstake function on the contract.
func (c *Client) TriggerUnstake(poolId uint64) (*types.Transaction, error) {
	if c.auth == nil {
		return nil, fmt.Errorf("client is read-only")
	}
	return c.contract.Transact(c.auth, "triggerUnstake", poolId)
}
