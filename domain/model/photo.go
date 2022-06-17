package model

import "time"

type PhotoUpload struct {
	Id       int64 `gorm:"primaryKey"`
	Created  time.Time
	Modified time.Time
	FromId   string
	ToId     string
}

type AlbumUpload struct {
	Id          int64 `gorm:"primaryKey"`
	Created     time.Time
	Modified    time.Time
	FromAlbumId string
	ToAlbumId   string
	CountPhotos int
}

type AlbumPhotosUpload struct {
	AlbumId string
	PhotoId string
}
