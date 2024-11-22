package database

import (
	"fmt"
	"github.com/rodrinoblega/stori/config"
	"github.com/rodrinoblega/stori/src/adapters/repositories"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"sync"
)

var (
	once     sync.Once
	instance *repositories.Database
)

func New(env *config.Config) *repositories.Database {
	once.Do(func() {
		instance = postgresDB(env)
	})

	return instance
}

func postgresDB(env *config.Config) *repositories.Database {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		env.PostgresHost,
		env.PgUser,
		env.PgPassword,
		env.PgDatabase,
		env.PostgresPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return &repositories.Database{DB: db}
}
