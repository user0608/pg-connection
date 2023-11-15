package connection

import (
	"context"

	"gorm.io/gorm"
)

type StorageManager interface {
	Conn(ctx context.Context) *gorm.DB
	WithTx(ctx context.Context, fc func(ctx context.Context) error) error
}
