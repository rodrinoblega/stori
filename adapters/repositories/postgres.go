package repositories

import (
	"fmt"
	"github.com/rodrinoblega/stori/config"
	"github.com/rodrinoblega/stori/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"sync"
)

type Database struct {
	*gorm.DB
}

var (
	once     sync.Once
	instance *Database
)

func (d *Database) StoreTransactions(transactions entities.Transactions) error {
	result := d.Create(transactions)
	if result.Error != nil {
		return fmt.Errorf("failed to store transactions: %w", result.Error)
	}
	return nil
}

func (d *Database) GetAccountById(accountID int) (*entities.Account, error) {
	var account entities.Account
	if err := d.First(&account, accountID).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func New(env *config.Config) *Database {
	once.Do(func() {
		instance = postgresDB(env)
	})

	return instance
}

func postgresDB(env *config.Config) *Database {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		env.PostgresHost,
		env.PgUser,
		env.PgPassword,
		env.PgDatabase,
		env.PostgresPort,
	)
	fmt.Println(dsn)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return &Database{db}
}
