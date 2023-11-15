package connection

import (
	"context"
	"os"
	"sync"

	"gorm.io/gorm"
)

var once sync.Once
var manager StorageManager

func readConfigs() DBConfigParams {
	var host = os.Getenv("PG_HOST")
	var port = os.Getenv("PG_PORT")
	var user = os.Getenv("PG_USER")
	var password = os.Getenv("PG_PASSWORD")
	var database = os.Getenv("PG_DATABASE")
	var level = os.Getenv("PG_DATABASE_LOG_LEVEL")
	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "5432"
	}
	return DBConfigParams{
		DBHost:     host,
		DBPort:     port,
		DBName:     database,
		DBUsername: user,
		DBPassword: password,
		DBLogLevel: level,
	}
}

// this function will use the following environment variables
// to open the PostgreSQL connection: PG_HOST, PG_PORT, PG_USER, PG_PASSWORD, PG_DATABASE.
func Open() (StorageManager, error) {
	var err error
	once.Do(func() {
		manager, err = NewConnection(readConfigs())
	})
	return manager, err
}

func OpenWithConfigs(params DBConfigParams) (StorageManager, error) {
	var err error
	once.Do(func() {
		manager, err = NewConnection(params)
	})
	return manager, err
}

func Conn(ctx context.Context) *gorm.DB {
	return manager.Conn(ctx)
}

func WithTx(ctx context.Context, fc func(ctx context.Context) error) error {
	return manager.WithTx(ctx, fc)
}
