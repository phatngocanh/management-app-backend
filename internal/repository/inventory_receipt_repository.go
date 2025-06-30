package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pna/management-app-backend/internal/domain/entity"
)

type InventoryReceiptRepository interface {
	CreateCommand(ctx context.Context, receipt *entity.InventoryReceipt, tx *sqlx.Tx) error
	GetAllQuery(ctx context.Context, tx *sqlx.Tx) ([]entity.InventoryReceipt, error)
	GetOneByIDQuery(ctx context.Context, id int, tx *sqlx.Tx) (*entity.InventoryReceipt, error)
	GetOneByCodeQuery(ctx context.Context, code string, tx *sqlx.Tx) (*entity.InventoryReceipt, error)
	GetByUserIDQuery(ctx context.Context, userID int, tx *sqlx.Tx) ([]entity.InventoryReceipt, error)
	UpdateCommand(ctx context.Context, receipt *entity.InventoryReceipt, tx *sqlx.Tx) error
	DeleteCommand(ctx context.Context, id int, tx *sqlx.Tx) error
}
