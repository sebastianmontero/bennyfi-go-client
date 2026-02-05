package yield_source_registry

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// IYieldSourceRegistryABI is the ABI of the IYieldSourceRegistry contract.
const IYieldSourceRegistryABI = `[
    {
      "inputs": [
        {
          "internalType": "string",
          "name": "_name",
          "type": "string"
        }
      ],
      "name": "getYieldSource",
      "outputs": [
        {
          "components": [
            {
              "internalType": "string",
              "name": "name",
              "type": "string"
            },
            {
              "internalType": "address",
              "name": "adaptorContract",
              "type": "address"
            },
            {
              "internalType": "bool",
              "name": "active",
              "type": "bool"
            }
          ],
          "internalType": "struct IYieldSourceRegistry.YieldSource",
          "name": "",
          "type": "tuple"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "string",
          "name": "_name",
          "type": "string"
        }
      ],
      "name": "isYieldSourceRegistered",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    }
]`

// YieldSource represents a yield source in the registry.
type YieldSource struct {
	Name            string         `json:"name"`
	AdaptorContract common.Address `json:"adaptorContract"`
	Active          bool           `json:"active"`
}

// Registry interacts with the IYieldSourceRegistry contract.
type Registry struct {
	contract *bind.BoundContract
	address  common.Address
}

// NewRead creates a new Registry client for read-only operations.
func NewRead(rpcUrl string, contractAddr common.Address) (*Registry, error) {
	// 1. Connect to Eth Client
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to dial rpc: %w", err)
	}

	return NewReadWithClient(client, contractAddr)
}

// NewReadWithClient creates a new Registry client using an existing contract backend.
func NewReadWithClient(client *ethclient.Client, contractAddr common.Address) (*Registry, error) {
	// 1. Parse ABI
	parsedABI, err := abi.JSON(strings.NewReader(IYieldSourceRegistryABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse abi: %w", err)
	}

	// 2. Create Bound Contract
	contract := bind.NewBoundContract(contractAddr, parsedABI, client, client, client)
	return &Registry{
		contract: contract,
		address:  contractAddr,
	}, nil
}

// GetYieldSource retrieves a yield source by name.
func (r *Registry) GetYieldSource(name string) (*YieldSource, error) {
	// Use generic unmarshalling to avoid type mismatches with anonymous structs
	var results []interface{}
	err := r.contract.Call(nil, &results, "getYieldSource", name)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no result returned")
	}

	// Use JSON roundtrip to convert anonymous struct to YieldSource
	// This works because the ABI tuple unpacking produces an anonymous struct with json tags matching YieldSource
	data, err := json.Marshal(results[0])
	if err != nil {
		return nil, fmt.Errorf("failed to marshal result: %w", err)
	}

	result := new(YieldSource)
	if err := json.Unmarshal(data, result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal into YieldSource: %w. Data: %s", err, string(data))
	}

	return result, nil
}

// IsYieldSourceRegistered checks if a yield source is registered.
func (r *Registry) IsYieldSourceRegistered(name string) (bool, error) {
	var out []interface{}
	err := r.contract.Call(nil, &out, "isYieldSourceRegistered", name)
	if err != nil {
		return false, err
	}
	return out[0].(bool), nil
}
