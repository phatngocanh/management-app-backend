package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pna/management-app-backend/internal/domain/entity"
)

type UnitOfMeasureRepository interface {
	GetAllQuery(ctx context.Context, tx *sqlx.Tx) ([]entity.UnitOfMeasure, error)
	GetOneByIDQuery(ctx context.Context, id int, tx *sqlx.Tx) (*entity.UnitOfMeasure, error)
	GetOneByCodeQuery(ctx context.Context, code string, tx *sqlx.Tx) (*entity.UnitOfMeasure, error)
	CreateCommand(ctx context.Context, unit *entity.UnitOfMeasure, tx *sqlx.Tx) error
	UpdateCommand(ctx context.Context, unit *entity.UnitOfMeasure, tx *sqlx.Tx) error
}
