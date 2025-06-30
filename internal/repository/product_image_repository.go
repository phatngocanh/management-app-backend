package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pna/management-app-backend/internal/domain/entity"
)

type ProductImageRepository interface {
	GetAllQuery(ctx context.Context, tx *sqlx.Tx) ([]entity.ProductImage, error)
	GetOneByIDQuery(ctx context.Context, id int, tx *sqlx.Tx) (*entity.ProductImage, error)
	GetByProductIDQuery(ctx context.Context, productID int, tx *sqlx.Tx) ([]entity.ProductImage, error)
	CreateCommand(ctx context.Context, image *entity.ProductImage, tx *sqlx.Tx) error
	UpdateCommand(ctx context.Context, image *entity.ProductImage, tx *sqlx.Tx) error
	DeleteCommand(ctx context.Context, id int, tx *sqlx.Tx) error
}
