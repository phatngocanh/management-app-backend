package repositoryimplement

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/pna/management-app-backend/internal/database"
	"github.com/pna/management-app-backend/internal/domain/entity"
	"github.com/pna/management-app-backend/internal/repository"
	"github.com/pna/management-app-backend/internal/utils/error_utils"
)

type UnitOfMeasureRepository struct {
	db *sqlx.DB
}

func NewUnitOfMeasureRepository(db database.Db) repository.UnitOfMeasureRepository {
	return &UnitOfMeasureRepository{db: db}
}

func (repo *UnitOfMeasureRepository) GetAllQuery(ctx context.Context, tx *sqlx.Tx) ([]entity.UnitOfMeasure, error) {
	var units []entity.UnitOfMeasure
	query := "SELECT * FROM units_of_measure ORDER BY id"
	var err error

	if tx != nil {
		err = tx.SelectContext(ctx, &units, query)
	} else {
		err = repo.db.SelectContext(ctx, &units, query)
	}

	if err != nil {
		return nil, err
	}

	if units == nil {
		return []entity.UnitOfMeasure{}, nil
	}

	return units, nil
}

func (repo *UnitOfMeasureRepository) GetOneByIDQuery(ctx context.Context, id int, tx *sqlx.Tx) (*entity.UnitOfMeasure, error) {
	var unit entity.UnitOfMeasure
	query := "SELECT * FROM units_of_measure WHERE id = ?"
	var err error

	if tx != nil {
		err = tx.GetContext(ctx, &unit, query, id)
	} else {
		err = repo.db.GetContext(ctx, &unit, query, id)
	}

	if err != nil {
		if err.Error() == error_utils.SystemErrorMessage.SqlxNoRow {
			return nil, nil
		}
		return nil, err
	}

	return &unit, nil
}

func (repo *UnitOfMeasureRepository) GetOneByCodeQuery(ctx context.Context, code string, tx *sqlx.Tx) (*entity.UnitOfMeasure, error) {
	var unit entity.UnitOfMeasure
	query := "SELECT * FROM units_of_measure WHERE code = ?"
	var err error

	if tx != nil {
		err = tx.GetContext(ctx, &unit, query, code)
	} else {
		err = repo.db.GetContext(ctx, &unit, query, code)
	}

	if err != nil {
		if err.Error() == error_utils.SystemErrorMessage.SqlxNoRow {
			return nil, nil
		}
		return nil, err
	}

	return &unit, nil
}

func (repo *UnitOfMeasureRepository) CreateCommand(ctx context.Context, unit *entity.UnitOfMeasure, tx *sqlx.Tx) error {
	insertQuery := `INSERT INTO units_of_measure(name, code, description) VALUES (:name, :code, :description)`

	var result sql.Result
	var err error

	if tx != nil {
		result, err = tx.NamedExecContext(ctx, insertQuery, unit)
	} else {
		result, err = repo.db.NamedExecContext(ctx, insertQuery, unit)
	}

	if err != nil {
		return err
	}

	// Get the last inserted ID
	lastID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	// Set the ID to the unit entity
	unit.ID = int(lastID)
	return nil
}

func (repo *UnitOfMeasureRepository) UpdateCommand(ctx context.Context, unit *entity.UnitOfMeasure, tx *sqlx.Tx) error {
	updateQuery := `UPDATE units_of_measure SET name = :name, code = :code, description = :description WHERE id = :id`

	if tx != nil {
		_, err := tx.NamedExecContext(ctx, updateQuery, unit)
		return err
	}
	_, err := repo.db.NamedExecContext(ctx, updateQuery, unit)
	return err
}
