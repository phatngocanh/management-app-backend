package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pna/management-app-backend/internal/domain/entity"
)

type ProductBomRepository interface {
	GetAllQuery(ctx context.Context, tx *sqlx.Tx) ([]entity.ProductBom, error)
	GetOneByIDQuery(ctx context.Context, id int, tx *sqlx.Tx) (*entity.ProductBom, error)
	GetByParentProductIDQuery(ctx context.Context, parentProductID int, tx *sqlx.Tx) ([]entity.ProductBom, error)
	GetByComponentProductIDQuery(ctx context.Context, componentProductID int, tx *sqlx.Tx) ([]entity.ProductBom, error)
	CreateCommand(ctx context.Context, bom *entity.ProductBom, tx *sqlx.Tx) error
	UpdateCommand(ctx context.Context, bom *entity.ProductBom, tx *sqlx.Tx) error
	DeleteCommand(ctx context.Context, id int, tx *sqlx.Tx) error
}
