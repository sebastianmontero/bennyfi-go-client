package base

import (
	"context"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// IYieldSourceAdaptorABI is the ABI of the IYieldSourceAdaptor contract.
const IYieldSourceAdaptorABI = `[
    {
      "inputs": [
        {
          "internalType": "uint64",
          "name": "_poolId",
          "type": "uint64"
        },
        {
          "internalType": "string",
          "name": "_yieldSource",
          "type": "string"
        },
        {
          "internalType": "address",
          "name": "_token",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "_amount",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "_stakingPeriodHrs",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "_stakingPeriodCycleHrs",
          "type": "uint256"
        }
      ],
      "name": "stake",
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
      "name": "stopStake",
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
      "name": "isStoppable",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "uint64",
          "name": "poolId",
          "type": "uint64"
        }
      ],
      "name": "StakeStopped",
      "type": "event"
    }
]`

// BaseRead represents a base yield source adaptor client for read-only operations.
type BaseRead struct {
	Contract *bind.BoundContract
	Address  common.Address
}

// BaseWrite represents a base yield source adaptor client for write operations.
type BaseWrite struct {
	*BaseRead
	Auth *bind.TransactOpts
}

// New creates a new BaseWrite client.
func New(rpcUrl string, privateKeyHex string, contractAddr common.Address, additionalABIs ...string) (*BaseWrite, error) {
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to dial rpc: %w", err)
	}
	return NewWithClient(client, privateKeyHex, contractAddr, additionalABIs...)
}

// NewWithClient creates a new BaseWrite client with an existing ethclient.
func NewWithClient(client *ethclient.Client, privateKeyHex string, contractAddr common.Address, additionalABIs ...string) (*BaseWrite, error) {
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get chain id: %w", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create auth: %w", err)
	}

	readClient, err := NewReadWithClient(client, contractAddr, additionalABIs...)
	if err != nil {
		return nil, err
	}

	return &BaseWrite{
		BaseRead: readClient,
		Auth:     auth,
	}, nil
}

// NewClientFactory creates a factory function that produces BaseWrite clients for a given address.
// It captures the client and credentials to reuse them for creating multiple clients.
func NewClientFactory(client *ethclient.Client, privateKeyHex string, additionalABIs ...string) (func(common.Address) (*BaseWrite, error), error) {
	// Verify private key validity early
	_, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %w", err)
	}

	return func(contractAddr common.Address) (*BaseWrite, error) {
		return NewWithClient(client, privateKeyHex, contractAddr, additionalABIs...)
	}, nil
}

// NewRead creates a new BaseRead client.
func NewRead(rpcUrl string, contractAddr common.Address, additionalABIs ...string) (*BaseRead, error) {
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to dial rpc: %w", err)
	}
	return NewReadWithClient(client, contractAddr, additionalABIs...)
}

// NewReadWithClient creates a new BaseRead client with an existing ethclient.
func NewReadWithClient(client *ethclient.Client, contractAddr common.Address, additionalABIs ...string) (*BaseRead, error) {
	abis := append([]string{IYieldSourceAdaptorABI}, additionalABIs...)
	mergedABI, err := MergeABIs(abis...)
	if err != nil {
		return nil, fmt.Errorf("failed to merge ABIs: %w", err)
	}

	parsedABI, err := abi.JSON(strings.NewReader(mergedABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse abi: %w", err)
	}

	contract := bind.NewBoundContract(contractAddr, parsedABI, client, client, client)

	return &BaseRead{
		Contract: contract,
		Address:  contractAddr,
	}, nil
}

// MergeABIs merges multiple ABI JSON strings into a single ABI string.
// It naively concatenates the arrays, verifying they are valid JSON arrays.
func MergeABIs(abis ...string) (string, error) {
	var builder strings.Builder
	builder.WriteString("[")
	first := true
	for _, abiStr := range abis {
		trimmed := strings.TrimSpace(abiStr)
		if len(trimmed) < 2 {
			continue
		} // "[]" or empty

		if !strings.HasPrefix(trimmed, "[") || !strings.HasSuffix(trimmed, "]") {
			return "", fmt.Errorf("invalid ABI JSON: must be an array")
		}

		inner := trimmed[1 : len(trimmed)-1]
		if strings.TrimSpace(inner) == "" {
			continue
		}

		if !first {
			builder.WriteString(",")
		}
		builder.WriteString(inner)
		first = false
	}
	builder.WriteString("]")
	return builder.String(), nil
}

// NewReadFromContract creates a new BaseRead client from an existing bound contract.
func NewReadFromContract(contract *bind.BoundContract, contractAddr common.Address) *BaseRead {
	return &BaseRead{
		Contract: contract,
		Address:  contractAddr,
	}
}

// StopStake calls the stopStake function on the contract.
// Returns the transaction hash.
func (b *BaseWrite) StopStake(poolId uint64) (string, error) {
	tx, err := b.Contract.Transact(b.Auth, "stopStake", poolId)
	if err != nil {
		return "", fmt.Errorf("failed to clean stake: %w", err)
	}
	return tx.Hash().Hex(), nil
}

// IsStoppable checks if the stake can be stopped.
func (b *BaseRead) IsStoppable(poolId uint64) (bool, error) {
	var out []interface{}
	err := b.Contract.Call(nil, &out, "isStoppable", poolId)
	if err != nil {
		return false, fmt.Errorf("failed to check if stoppable: %w", err)
	}
	return out[0].(bool), nil
}
