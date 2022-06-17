package repository

import (
	"context"
	"google-photos-sync/domain/model"
	"google-photos-sync/infrastructure/db"
	"google-photos-sync/infrastructure/web/auth"
	"google-photos-sync/interface/repository"
	"gorm.io/gorm/clause"
	"time"
)

type SqlCredentialsRepository struct {
	db *db.DB
}

func MakeSqlCredentialsRepository(db *db.DB) repository.CredentialsRepository {
	db.AutoMigrate(&model.User{})
	return &SqlCredentialsRepository{
		db: db,
	}
}

func (r *SqlCredentialsRepository) Set(ctx context.Context, accountType model.AccountType, token []byte) error {
	user := &model.User{
		Id:       auth.GetUserId(ctx),
		Created:  time.Now(),
		Modified: time.Now(),
	}

	columns := []string{"modified"}

	if accountType == model.From {
		user.FromToken = token
		columns = append(columns, "from_token")
	} else {
		user.ToToken = token
		columns = append(columns, "to_token")
	}

	tx := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns(columns),
	}).Create(&user)

	return tx.Error
}

func (r *SqlCredentialsRepository) Get(ctx context.Context) (*model.User, error) {
	userId := auth.GetUserId(ctx)

	user := &model.User{}
	find := r.db.WithContext(ctx).Find(user, "id = ?", userId)

	return user, find.Error
}
