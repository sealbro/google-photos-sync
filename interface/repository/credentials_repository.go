package repository

import (
	"context"
	"google-photos-sync/domain/model"
)

type CredentialsRepository interface {
	Set(ctx context.Context, accountType model.AccountType, token []byte) error
	Get(ctx context.Context) (*model.User, error)
}
