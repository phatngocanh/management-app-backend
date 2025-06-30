package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pna/management-app-backend/internal/domain/entity"
)

type ProductCategoryRepository interface {
	GetAllQuery(ctx context.Context, tx *sqlx.Tx) ([]entity.ProductCategory, error)
	GetOneByIDQuery(ctx context.Context, id int, tx *sqlx.Tx) (*entity.ProductCategory, error)
	GetOneByCodeQuery(ctx context.Context, code string, tx *sqlx.Tx) (*entity.ProductCategory, error)
	CreateCommand(ctx context.Context, category *entity.ProductCategory, tx *sqlx.Tx) error
	UpdateCommand(ctx context.Context, category *entity.ProductCategory, tx *sqlx.Tx) error
}
