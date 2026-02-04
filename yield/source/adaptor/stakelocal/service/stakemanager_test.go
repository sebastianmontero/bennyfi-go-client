package service

import (
	"errors"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	registry "github.com/sebastianmontero/bennyfi-go-client/evm/bridge/yield-source-registry"
	"github.com/sebastianmontero/bennyfi-go-client/yield/source/adaptor/stakelocal"
	"github.com/sebastianmontero/eos-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockStakeContract is a mock implementation of StakeContract
type MockStakeContract struct {
	mock.Mock
}

func (m *MockStakeContract) GetStake(roundID uint64) (*stakelocal.Stake, error) {
	args := m.Called(roundID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*stakelocal.Stake), args.Error(1)
}

func (m *MockStakeContract) GetAllStakesByStateAndYieldSource(state, yieldSource eos.Name) ([]stakelocal.Stake, error) {
	args := m.Called(state, yieldSource)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]stakelocal.Stake), args.Error(1)
}

func (m *MockStakeContract) StopStake(roundId uint64) (string, error) {
	args := m.Called(roundId)
	return args.String(0), args.Error(1)
}

// MockRegistryClient is a mock implementation of RegistryClient
type MockRegistryClient struct {
	mock.Mock
}

func (m *MockRegistryClient) GetYieldSource(name string) (*registry.YieldSource, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*registry.YieldSource), args.Error(1)
}

// MockBaseClient is a mock implementation of BaseClient
type MockBaseClient struct {
	mock.Mock
}

func (m *MockBaseClient) IsStoppable(poolId uint64) (bool, error) {
	args := m.Called(poolId)
	return args.Bool(0), args.Error(1)
}

func (m *MockBaseClient) StopStake(poolId uint64) (string, error) {
	args := m.Called(poolId)
	return args.String(0), args.Error(1)
}

func setupMocks(t *testing.T) (*MockStakeContract, *MockRegistryClient, *MockBaseClient, BaseClientFactory) {
	localMock := new(MockStakeContract)
	registryMock := new(MockRegistryClient)
	baseMock := new(MockBaseClient)

	factory := func(addr common.Address) (BaseClient, error) {
		return baseMock, nil
	}

	return localMock, registryMock, baseMock, factory
}

func TestStopStake_Success_StoppableOnEVM(t *testing.T) {
	localMock, registryMock, baseMock, factory := setupMocks(t)
	manager := New(localMock, registryMock, factory)

	roundID := uint64(123)
	yieldSource := "yield1"
	contractAddr := common.HexToAddress("0x123")

	stake := &stakelocal.Stake{
		RoundID:     roundID,
		State:       stakelocal.StakeStateStaked,
		YieldSource: eos.Name(yieldSource),
	}
	ys := &registry.YieldSource{
		Name:            yieldSource,
		AdaptorContract: contractAddr,
	}

	localMock.On("GetStake", roundID).Return(stake, nil)
	registryMock.On("GetYieldSource", yieldSource).Return(ys, nil)
	// Factor invoked internally, returns baseMock

	baseMock.On("IsStoppable", roundID).Return(true, nil)
	baseMock.On("StopStake", roundID).Return("evm_tx_hash", nil)
	localMock.On("StopStake", roundID).Return("local_tx_id", nil)

	err := manager.StopStake(roundID)
	assert.NoError(t, err)

	localMock.AssertExpectations(t)
	registryMock.AssertExpectations(t)
	baseMock.AssertExpectations(t)
}

func TestStopStake_Success_NotStoppableOnEVM(t *testing.T) {
	localMock, registryMock, baseMock, factory := setupMocks(t)
	manager := New(localMock, registryMock, factory)

	roundID := uint64(123)
	yieldSource := "yield1"
	contractAddr := common.HexToAddress("0x123")

	stake := &stakelocal.Stake{
		RoundID:     roundID,
		State:       stakelocal.StakeStateStaked,
		YieldSource: eos.Name(yieldSource),
	}
	ys := &registry.YieldSource{
		Name:            yieldSource,
		AdaptorContract: contractAddr,
	}

	localMock.On("GetStake", roundID).Return(stake, nil)
	registryMock.On("GetYieldSource", yieldSource).Return(ys, nil)

	baseMock.On("IsStoppable", roundID).Return(false, nil)
	localMock.On("StopStake", roundID).Return("local_tx_id", nil)

	err := manager.StopStake(roundID)
	assert.NoError(t, err)

	localMock.AssertExpectations(t)
	registryMock.AssertExpectations(t)
	baseMock.AssertExpectations(t)
}

func TestStopStake_Fail_StakeNotFound(t *testing.T) {
	localMock, registryMock, baseMock, factory := setupMocks(t)
	manager := New(localMock, registryMock, factory)

	roundID := uint64(999)
	localMock.On("GetStake", roundID).Return(nil, nil)

	err := manager.StopStake(roundID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "stake 999 not found")

	localMock.AssertExpectations(t)
	baseMock.AssertNotCalled(t, "IsStoppable")
}

func TestStopStake_Fail_WrongState(t *testing.T) {
	localMock, _, _, factory := setupMocks(t)
	manager := New(localMock, nil, factory)

	roundID := uint64(123)
	stake := &stakelocal.Stake{
		RoundID: roundID,
		State:   stakelocal.StakeStateUnstaked,
	}

	localMock.On("GetStake", roundID).Return(stake, nil)

	err := manager.StopStake(roundID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not in 'staked' state")

	localMock.AssertExpectations(t)
}

func TestStopStake_Fail_EVMError(t *testing.T) {
	localMock, registryMock, baseMock, factory := setupMocks(t)
	manager := New(localMock, registryMock, factory)

	roundID := uint64(123)
	yieldSource := "yield1"
	contractAddr := common.HexToAddress("0x123")

	stake := &stakelocal.Stake{
		RoundID:     roundID,
		State:       stakelocal.StakeStateStaked,
		YieldSource: eos.Name(yieldSource),
	}
	ys := &registry.YieldSource{
		Name:            yieldSource,
		AdaptorContract: contractAddr,
	}

	localMock.On("GetStake", roundID).Return(stake, nil)
	registryMock.On("GetYieldSource", yieldSource).Return(ys, nil)

	baseMock.On("IsStoppable", roundID).Return(true, nil)
	baseMock.On("StopStake", roundID).Return("", errors.New("evm failed"))

	err := manager.StopStake(roundID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to stop stake 123 on EVM")

	localMock.AssertExpectations(t)
	registryMock.AssertExpectations(t)
	baseMock.AssertExpectations(t)
}

func TestStopStakesByYieldSource_Success(t *testing.T) {
	localMock, registryMock, baseMock, factory := setupMocks(t)
	manager := New(localMock, registryMock, factory)

	yieldSource := eos.Name("yield1")
	yieldSourceStr := "yield1"
	contractAddr := common.HexToAddress("0x123")

	stakes := []stakelocal.Stake{
		{RoundID: 1, State: stakelocal.StakeStateStaked, YieldSource: yieldSource},
		{RoundID: 2, State: stakelocal.StakeStateStaked, YieldSource: yieldSource},
	}
	ys := &registry.YieldSource{
		Name:            yieldSourceStr,
		AdaptorContract: contractAddr,
	}

	localMock.On("GetAllStakesByStateAndYieldSource", stakelocal.StakeStateStaked, yieldSource).Return(stakes, nil)

	// Stake 1
	registryMock.On("GetYieldSource", yieldSourceStr).Return(ys, nil).Once() // Called only ONCE for the batch
	baseMock.On("IsStoppable", uint64(1)).Return(false, nil)
	localMock.On("StopStake", uint64(1)).Return("tx1", nil)

	// Stake 2
	baseMock.On("IsStoppable", uint64(2)).Return(true, nil)
	baseMock.On("StopStake", uint64(2)).Return("tx_evm_2", nil)
	localMock.On("StopStake", uint64(2)).Return("tx2", nil)

	errs := manager.StopStakesByYieldSource(yieldSource)
	assert.Nil(t, errs)

	localMock.AssertExpectations(t)
	registryMock.AssertExpectations(t)
	baseMock.AssertExpectations(t)
}

func TestStopStakesByYieldSource_PartialFailure(t *testing.T) {
	localMock, registryMock, baseMock, factory := setupMocks(t)
	manager := New(localMock, registryMock, factory)

	yieldSource := eos.Name("yield1")
	yieldSourceStr := "yield1"
	contractAddr := common.HexToAddress("0x123")

	stakes := []stakelocal.Stake{
		{RoundID: 1, State: stakelocal.StakeStateStaked, YieldSource: yieldSource},
		{RoundID: 2, State: stakelocal.StakeStateStaked, YieldSource: yieldSource},
	}
	ys := &registry.YieldSource{
		Name:            yieldSourceStr,
		AdaptorContract: contractAddr,
	}

	localMock.On("GetAllStakesByStateAndYieldSource", stakelocal.StakeStateStaked, yieldSource).Return(stakes, nil)
	registryMock.On("GetYieldSource", yieldSourceStr).Return(ys, nil).Once() // Called only ONCE

	// Stake 1 fails on EVM check
	baseMock.On("IsStoppable", uint64(1)).Return(false, errors.New("check failed"))

	// Stake 2 succeeds
	baseMock.On("IsStoppable", uint64(2)).Return(false, nil)
	localMock.On("StopStake", uint64(2)).Return("tx2", nil)

	errs := manager.StopStakesByYieldSource(yieldSource)
	assert.NotNil(t, errs)
	assert.Len(t, errs, 1)
	assert.Contains(t, errs[0].Error(), "failed to stop stake 1")

	localMock.AssertExpectations(t)
	registryMock.AssertExpectations(t)
	baseMock.AssertExpectations(t)
}

func TestStopStakesByYieldSource_Filtering(t *testing.T) {
	localMock, registryMock, baseMock, factory := setupMocks(t)
	manager := New(localMock, registryMock, factory)

	yieldSource := eos.Name("yield1")
	otherYieldSource := eos.Name("yield2")
	yieldSourceStr := "yield1"
	contractAddr := common.HexToAddress("0x123")

	stakes := []stakelocal.Stake{
		{RoundID: 1, State: stakelocal.StakeStateStaked, YieldSource: yieldSource},
		{RoundID: 2, State: stakelocal.StakeStateStaked, YieldSource: otherYieldSource}, // Should be skipped
	}
	ys := &registry.YieldSource{
		Name:            yieldSourceStr,
		AdaptorContract: contractAddr,
	}

	localMock.On("GetAllStakesByStateAndYieldSource", stakelocal.StakeStateStaked, yieldSource).Return(stakes, nil)

	// Stake 1 should be processed
	registryMock.On("GetYieldSource", yieldSourceStr).Return(ys, nil).Once()
	baseMock.On("IsStoppable", uint64(1)).Return(false, nil)
	localMock.On("StopStake", uint64(1)).Return("tx1", nil)

	// Stake 2 should check yield source and skip, so no mock calls for ID 2

	errs := manager.StopStakesByYieldSource(yieldSource)
	assert.Nil(t, errs)

	localMock.AssertExpectations(t)
	registryMock.AssertExpectations(t)
	baseMock.AssertExpectations(t)
}
