package states

import "context"

type SyncState struct {
	EmptyState
}

func NewSyncState(state EmptyState) *SyncState {
	state.Type = Sync

	return &SyncState{
		EmptyState: state,
	}
}

func (state *SyncState) Action(ctx context.Context) error {
	return nil
}
