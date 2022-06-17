package machine

import (
	"context"
	"fmt"
	"google-photos-sync/domain/model"
	"google-photos-sync/infrastructure/state"
	"google-photos-sync/infrastructure/web/auth"
	"google-photos-sync/interface/repository"
	"google-photos-sync/usecases/machine/states"
	"time"
)

type StateStorage interface {
	Convert(userId string, stateType states.StateType) state.State
	GetLastState(ctx context.Context) (state.State, error)
	AddNewState(ctx context.Context, oldState, newState state.State) error
}

type SqlStorage struct {
	stateRepository       repository.StateRepository
	credentialsRepository repository.CredentialsRepository
}

func MakeSqlStorage(stateRepository repository.StateRepository, credentialsRepository repository.CredentialsRepository) StateStorage {
	return &SqlStorage{
		stateRepository:       stateRepository,
		credentialsRepository: credentialsRepository,
	}
}

func (s *SqlStorage) GetLastState(ctx context.Context) (state.State, error) {
	userId := auth.GetUserId(ctx)
	lastState, err := s.stateRepository.GetLastState(ctx)
	if err != nil {
		return nil, fmt.Errorf("SqlStorage - GetLastState: %w", err)
	}

	stateType := states.Empty
	if lastState != nil {
		stateType = lastState.To
	}

	return s.Convert(userId, stateType), nil
}

func (s *SqlStorage) AddNewState(ctx context.Context, oldState, newState state.State) error {
	userId := auth.GetUserId(ctx)
	err := s.stateRepository.AddNewState(ctx, &model.StateHistory{
		Created: time.Now(),
		UserId:  userId,
		From:    oldState.GetType(),
		To:      newState.GetType(),
	})

	if err != nil {
		return fmt.Errorf("SqlStorage - AddNewState: %w", err)
	}
	return nil
}

func (s *SqlStorage) Convert(userId string, stateType states.StateType) state.State {
	state := states.NewEmptyState(userId)

	switch stateType {
	case states.ReadyToSync:
		return states.NewReadyToSyncState(state, s.credentialsRepository)
	case states.Sync:
		return states.NewSyncState(state)
	case states.Complete:
		return states.NewCompleteState(state)
	}

	return &state
}
