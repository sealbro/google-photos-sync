package states

import (
	"context"
	"google-photos-sync/interface/repository"
)

type ReadyToSyncState struct {
	EmptyState
	CredentialsRepository repository.CredentialsRepository
}

func NewReadyToSyncState(state EmptyState, credentialsRepository repository.CredentialsRepository) *ReadyToSyncState {
	state.Type = ReadyToSync

	return &ReadyToSyncState{
		EmptyState:            state,
		CredentialsRepository: credentialsRepository,
	}
}

func (state *ReadyToSyncState) Action(ctx context.Context) error {
	return nil
}
