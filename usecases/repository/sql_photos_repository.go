package repository

import (
	"context"
	"google-photos-sync/domain/model"
	"google-photos-sync/infrastructure/db"
	"google-photos-sync/interface/repository"
	"gorm.io/gorm/clause"
	"time"
)

type SqlPhotosRepository struct {
	db *db.DB
}

func MakeSqlPhotosRepository(db *db.DB) repository.PhotosRepository {
	db.AutoMigrate(&model.PhotoUpload{})
	return &SqlPhotosRepository{
		db: db,
	}
}

func (r *SqlPhotosRepository) AddPhotos(ctx context.Context, photoIds []string) error {
	length := len(photoIds)
	uploads := make([]model.PhotoUpload, length)
	for i, photoId := range photoIds {
		uploads[i].Created = time.Now()
		uploads[i].Modified = time.Now()
		uploads[i].FromId = photoId
	}

	batches := r.db.WithContext(ctx).CreateInBatches(uploads, length)
	return batches.Error
}

func (r *SqlPhotosRepository) UpdatePhotos(ctx context.Context, photos []*model.PhotoUpload) error {
	length := len(photos)
	for _, photo := range photos {
		photo.Modified = time.Now()
	}

	columns := []string{"modified", "to_id"}

	batches := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "from_id"}},
		DoUpdates: clause.AssignmentColumns(columns),
	}).CreateInBatches(photos, length)
	return batches.Error
}
