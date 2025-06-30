package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pna/management-app-backend/internal/domain/entity"
)

type InventoryReceiptItemRepository interface {
	CreateCommand(ctx context.Context, item *entity.InventoryReceiptItem, tx *sqlx.Tx) error
	GetAllQuery(ctx context.Context, tx *sqlx.Tx) ([]entity.InventoryReceiptItem, error)
	GetOneByIDQuery(ctx context.Context, id int, tx *sqlx.Tx) (*entity.InventoryReceiptItem, error)
	GetByInventoryReceiptIDQuery(ctx context.Context, inventoryReceiptID int, tx *sqlx.Tx) ([]entity.InventoryReceiptItem, error)
	GetByProductIDQuery(ctx context.Context, productID int, tx *sqlx.Tx) ([]entity.InventoryReceiptItem, error)
	UpdateCommand(ctx context.Context, item *entity.InventoryReceiptItem, tx *sqlx.Tx) error
	DeleteCommand(ctx context.Context, id int, tx *sqlx.Tx) error
	DeleteByInventoryReceiptIDCommand(ctx context.Context, inventoryReceiptID int, tx *sqlx.Tx) error
}
