package machine

import (
	"context"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google-photos-sync/domain/model"
	"google-photos-sync/infrastructure/state"
	"google-photos-sync/infrastructure/web/auth"
	"google-photos-sync/usecases/machine/states"
	"testing"
)

type MockCredentialsRepository struct {
	mock.Mock
}

func (m *MockCredentialsRepository) Set(ctx context.Context, accountType model.AccountType, token []byte) error {
	args := m.Called(ctx, accountType, token)

	return args.Error(0)
}
func (m *MockCredentialsRepository) Get(ctx context.Context) (*model.User, error) {
	args := m.Called(ctx)

	return args.Get(0).(*model.User), args.Error(1)
}

type FakeStateStorage struct {
	mock.Mock
	SqlStorage
}

func (m *FakeStateStorage) GetLastState(ctx context.Context) (state.State, error) {
	args := m.Called(ctx)

	return args.Get(0).(state.State), args.Error(1)
}

func (m *FakeStateStorage) AddNewState(ctx context.Context, oldState, newState state.State) error {
	args := m.Called(ctx, oldState, newState)

	return args.Error(0)
}

func TestReadyToSyncStateMachineCurrent(t *testing.T) {
	// Arrange
	userId := "123"
	ctx := auth.WrapWithUserId(context.Background(), userId)
	repository := &MockCredentialsRepository{}
	repository.On("Exists", userId, model.From).Return(true, nil).Once()
	repository.On("Exists", userId, model.To).Return(true, nil).Once()
	syncState := states.NewReadyToSyncState(states.NewEmptyState(userId), repository)
	storage := &FakeStateStorage{}
	storage.On("GetLastState", ctx).Return(syncState, nil).Twice()
	machine := NewSyncMachine(userId, storage)

	// Act
	err1 := machine.Load(ctx)
	err2 := machine.Load(ctx)
	err3 := machine.Action(ctx)

	// Assert
	require.NoError(t, err1)
	require.NoError(t, err2)
	require.NoError(t, err3)
}
