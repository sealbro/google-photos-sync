package repository

import (
	"context"
	"google-photos-sync/domain/model"
)

type StateRepository interface {
	GetLastState(ctx context.Context) (*model.StateHistory, error)
	AddNewState(ctx context.Context, newState *model.StateHistory) error
}
