package repository

import (
	"context"
	"errors"
	"google-photos-sync/domain/model"
	"google-photos-sync/infrastructure/db"
	"google-photos-sync/infrastructure/web/auth"
	"google-photos-sync/interface/repository"
	"gorm.io/gorm"
)

type SqlStateRepository struct {
	db *db.DB
}

func MakeSqlStateRepository(db *db.DB) repository.StateRepository {
	db.AutoMigrate(&model.StateHistory{})
	return &SqlStateRepository{
		db: db,
	}
}

func (r *SqlStateRepository) GetLastState(ctx context.Context) (*model.StateHistory, error) {
	state := &model.StateHistory{}
	last := r.db.WithContext(ctx).Last(state, "user_id = ?", auth.GetUserId(ctx))
	if errors.Is(last.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return state, last.Error
}

func (r *SqlStateRepository) AddNewState(ctx context.Context, newState *model.StateHistory) error {
	create := r.db.WithContext(ctx).Create(newState)
	return create.Error
}
