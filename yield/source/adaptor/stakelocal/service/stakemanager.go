package service

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	registry "github.com/sebastianmontero/bennyfi-go-client/evm/bridge/yield-source-registry"
	"github.com/sebastianmontero/bennyfi-go-client/yield/source/adaptor/stakelocal"
	"github.com/sebastianmontero/eos-go"
)

// StakeContract defines the interface for interacting with the local stake contract.
type StakeContract interface {
	GetStake(roundID uint64) (*stakelocal.Stake, error)
	GetAllStakesByStateAndYieldSource(state, yieldSource eos.Name) ([]stakelocal.Stake, error)
	StopStake(roundId uint64) (string, error)
}

// BaseClient defines the interface for interacting with the EVM base client.
type BaseClient interface {
	IsStoppable(poolId uint64) (bool, error)
	StopStake(poolId uint64) (string, error)
}

// RegistryClient defines the interface for interacting with the yield source registry.
type RegistryClient interface {
	GetYieldSource(name string) (*registry.YieldSource, error)
}

// BaseClientFactory is a function that creates a BaseClient for a given address.
type BaseClientFactory func(common.Address) (BaseClient, error)

// StakeManager manages the interaction between the local stake contract and the EVM adaptor.
type StakeManager struct {
	LocalContract StakeContract
	Registry      RegistryClient
	ClientFactory BaseClientFactory
}

// New creates a new StakeManager instance.
func New(localContract StakeContract, registry RegistryClient, clientFactory BaseClientFactory) *StakeManager {
	return &StakeManager{
		LocalContract: localContract,
		Registry:      registry,
		ClientFactory: clientFactory,
	}
}

// StopStake stops a stake if it is in the active state.
// It checks if the stake is stoppable on the EVM side before proceeding.
func (m *StakeManager) StopStake(roundId uint64) error {
	// 1. Get the stake from the local contract
	stake, err := m.LocalContract.GetStake(roundId)
	if err != nil {
		return fmt.Errorf("failed to get stake %d: %w", roundId, err)
	}
	if stake == nil {
		return fmt.Errorf("stake %d not found", roundId)
	}

	// 2. Check if the stake is in the 'staked' state
	if stake.State != stakelocal.StakeStateStaked {
		return fmt.Errorf("stake %d is not in 'staked' state, current state: %s", roundId, stake.State)
	}

	// 3. Resolve EVM Client
	baseClient, err := m.getBaseClient(string(stake.YieldSource))
	if err != nil {
		return err
	}

	return m.stopStake(stake, baseClient)
}

// StopStakesByYieldSource stops all stakes in 'staked' state for a given yield source.
func (m *StakeManager) StopStakesByYieldSource(yieldSource eos.Name) []error {
	stakes, err := m.LocalContract.GetAllStakesByStateAndYieldSource(stakelocal.StakeStateStaked, yieldSource)
	if err != nil {
		return []error{fmt.Errorf("failed to get stakes for yield source %s: %w", yieldSource, err)}
	}

	fmt.Printf("Found %d stakes for yield source %s\n", len(stakes), yieldSource)

	if len(stakes) == 0 {
		return nil
	}

	// Resolve BaseClient once for all stakes
	baseClient, err := m.getBaseClient(string(yieldSource))
	if err != nil {
		return []error{fmt.Errorf("failed to resolve base client for yield source %s: %w", yieldSource, err)}
	}

	var errs []error
	for _, stake := range stakes {
		stake := stake // capture loop variable
		if stake.YieldSource != yieldSource {
			continue
		}
		if err := m.stopStake(&stake, baseClient); err != nil {
			errs = append(errs, fmt.Errorf("failed to stop stake %d: %w", stake.RoundID, err))
			continue
		}
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

// getBaseClient resolves the EVM base client for a given yield source.
func (m *StakeManager) getBaseClient(yieldSourceStr string) (BaseClient, error) {
	yieldSourceInfo, err := m.Registry.GetYieldSource(yieldSourceStr)
	if err != nil {
		return nil, fmt.Errorf("failed to get yield source info for %s: %w", yieldSourceStr, err)
	}

	baseClient, err := m.ClientFactory(yieldSourceInfo.AdaptorContract)
	if err != nil {
		return nil, fmt.Errorf("failed to create base client for address %s: %w", yieldSourceInfo.AdaptorContract.Hex(), err)
	}
	return baseClient, nil
}

// stopStake is a helper function that orchestrates the stop stake process for a single stake.
func (m *StakeManager) stopStake(stake *stakelocal.Stake, baseClient BaseClient) error {
	roundId := stake.RoundID

	// 1. Check if the stake is in the 'staked' state
	if stake.State != stakelocal.StakeStateStaked {
		return fmt.Errorf("stake %d is not in 'staked' state, current state: %s", roundId, stake.State)
	}

	// 2. Check if the stake is stoppable on the EVM side
	isStoppable, err := baseClient.IsStoppable(roundId)
	if err != nil {
		return fmt.Errorf("failed to check if stake %d is stoppable: %w", roundId, err)
	}

	// 3. If stoppable, call StopStake on the EVM base client
	if isStoppable {
		txHash, err := baseClient.StopStake(roundId)
		if err != nil {
			return fmt.Errorf("failed to stop stake %d on EVM: %w", roundId, err)
		}
		fmt.Printf("Stopped stake %d on EVM. Tx Hash: %s\n", roundId, txHash)
	}

	// 4. Call StopStake on the local contract
	localTxId, err := m.LocalContract.StopStake(roundId)
	if err != nil {
		return fmt.Errorf("failed to stop stake %d on local contract: %w", roundId, err)
	}
	fmt.Printf("Stopped stake %d on local contract. Tx ID: %s\n", roundId, localTxId)

	return nil
}
