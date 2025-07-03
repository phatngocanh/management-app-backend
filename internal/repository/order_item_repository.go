package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pna/management-app-backend/internal/domain/entity"
)

type OrderItemRepository interface {
	CreateCommand(ctx context.Context, item *entity.OrderItem, tx *sqlx.Tx) error
	GetAllQuery(ctx context.Context, tx *sqlx.Tx) ([]entity.OrderItem, error)
	GetOneByIDQuery(ctx context.Context, id int, tx *sqlx.Tx) (*entity.OrderItem, error)
	GetAllByOrderIDQuery(ctx context.Context, orderID int, tx *sqlx.Tx) ([]entity.OrderItem, error)
	UpdateCommand(ctx context.Context, item *entity.OrderItem, tx *sqlx.Tx) error
	DeleteCommand(ctx context.Context, id int, tx *sqlx.Tx) error
	DeleteByOrderIDCommand(ctx context.Context, orderID int, tx *sqlx.Tx) error
}
