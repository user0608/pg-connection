package connection

import (
	"context"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func NewConnection(config DBConfigParams) (StorageManager, error) {
	var conn = connection{config: config}
	if err := conn.openConnection(); err != nil {
		return nil, err
	}
	return &conn, nil
}

type key int

var context_connextion_key key

type connection struct {
	conn   *gorm.DB
	config DBConfigParams
}

func (*connection) level(s string) logger.LogLevel {
	switch s {
	case "info":
		return logger.Info
	case "error":
		return logger.Error
	case "warn":
		return logger.Warn
	}
	return logger.Silent
}

func (c *connection) openConnection() error {
	var config = c.config
	var level = c.level(c.config.DBLogLevel)
	const layer = "host=%s user=%s password=%s dbname=%s port=%s sslmode=disable"
	var dsn = fmt.Sprintf(layer, config.DBHost, config.DBUsername, config.DBPassword, config.DBName, config.DBPort)
	dialector := postgres.Open(dsn)
	var err error
	c.conn, err = gorm.Open(dialector, &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 logger.Default.LogMode(level),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err == nil && c.conn != nil {
		log.Println("Database connection established successfully.")
	}
	return err
}

func (c *connection) Conn(ctx context.Context) *gorm.DB {
	value := ctx.Value(context_connextion_key)
	if value == nil {
		return c.conn.WithContext(ctx)
	}
	connection, ok := value.(*gorm.DB)
	if !ok {
		return connection.WithContext(ctx)
	}
	return connection
}

func (c *connection) WithTx(ctx context.Context, txFunc func(ctx context.Context) error) error {
	if ctx.Value(context_connextion_key) != nil {
		return txFunc(ctx) // returns the current transaction
	}
	return c.Conn(ctx).Transaction(func(tx *gorm.DB) error {
		return txFunc(context.WithValue(ctx, context_connextion_key, tx))
	})
}
