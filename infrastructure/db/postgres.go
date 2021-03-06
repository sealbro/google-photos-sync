package db

import (
	"google-photos-sync/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type PostgresConfig struct {
	Connection string
	Schema     string
}

type DB struct {
	*gorm.DB
}

func MakePostgresConnection(config *PostgresConfig) (*DB, error) {
	open, err := gorm.Open(postgres.Open(config.Connection), &gorm.Config{
		Logger: &logger.GormLogger{},
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: config.Schema + ".",
		},
	})

	db := &DB{
		DB: open,
	}

	return db, err
}
