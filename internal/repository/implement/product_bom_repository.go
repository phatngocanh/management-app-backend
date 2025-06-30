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

type ProductBomRepository struct {
	db *sqlx.DB
}

func NewProductBomRepository(db database.Db) repository.ProductBomRepository {
	return &ProductBomRepository{db: db}
}

func (repo *ProductBomRepository) GetAllQuery(ctx context.Context, tx *sqlx.Tx) ([]entity.ProductBom, error) {
	var boms []entity.ProductBom
	query := "SELECT * FROM product_boms ORDER BY id"
	var err error

	if tx != nil {
		err = tx.SelectContext(ctx, &boms, query)
	} else {
		err = repo.db.SelectContext(ctx, &boms, query)
	}

	if err != nil {
		return nil, err
	}

	if boms == nil {
		return []entity.ProductBom{}, nil
	}

	return boms, nil
}

func (repo *ProductBomRepository) GetOneByIDQuery(ctx context.Context, id int, tx *sqlx.Tx) (*entity.ProductBom, error) {
	var bom entity.ProductBom
	query := "SELECT * FROM product_boms WHERE id = ?"
	var err error

	if tx != nil {
		err = tx.GetContext(ctx, &bom, query, id)
	} else {
		err = repo.db.GetContext(ctx, &bom, query, id)
	}

	if err != nil {
		if err.Error() == error_utils.SystemErrorMessage.SqlxNoRow {
			return nil, nil
		}
		return nil, err
	}

	return &bom, nil
}

func (repo *ProductBomRepository) GetByParentProductIDQuery(ctx context.Context, parentProductID int, tx *sqlx.Tx) ([]entity.ProductBom, error) {
	var boms []entity.ProductBom
	query := "SELECT * FROM product_boms WHERE parent_product_id = ? ORDER BY id"
	var err error

	if tx != nil {
		err = tx.SelectContext(ctx, &boms, query, parentProductID)
	} else {
		err = repo.db.SelectContext(ctx, &boms, query, parentProductID)
	}

	if err != nil {
		return nil, err
	}

	if boms == nil {
		return []entity.ProductBom{}, nil
	}

	return boms, nil
}

func (repo *ProductBomRepository) GetByComponentProductIDQuery(ctx context.Context, componentProductID int, tx *sqlx.Tx) ([]entity.ProductBom, error) {
	var boms []entity.ProductBom
	query := "SELECT * FROM product_boms WHERE component_product_id = ? ORDER BY id"
	var err error

	if tx != nil {
		err = tx.SelectContext(ctx, &boms, query, componentProductID)
	} else {
		err = repo.db.SelectContext(ctx, &boms, query, componentProductID)
	}

	if err != nil {
		return nil, err
	}

	if boms == nil {
		return []entity.ProductBom{}, nil
	}

	return boms, nil
}

func (repo *ProductBomRepository) CreateCommand(ctx context.Context, bom *entity.ProductBom, tx *sqlx.Tx) error {
	insertQuery := `INSERT INTO product_boms(parent_product_id, component_product_id, quantity) VALUES (:parent_product_id, :component_product_id, :quantity)`

	var result sql.Result
	var err error

	if tx != nil {
		result, err = tx.NamedExecContext(ctx, insertQuery, bom)
	} else {
		result, err = repo.db.NamedExecContext(ctx, insertQuery, bom)
	}

	if err != nil {
		return err
	}

	// Get the last inserted ID
	lastID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	// Set the ID to the bom entity
	bom.ID = int(lastID)
	return nil
}

func (repo *ProductBomRepository) UpdateCommand(ctx context.Context, bom *entity.ProductBom, tx *sqlx.Tx) error {
	updateQuery := `UPDATE product_boms SET parent_product_id = :parent_product_id, component_product_id = :component_product_id, quantity = :quantity WHERE id = :id`

	if tx != nil {
		_, err := tx.NamedExecContext(ctx, updateQuery, bom)
		return err
	}
	_, err := repo.db.NamedExecContext(ctx, updateQuery, bom)
	return err
}

func (repo *ProductBomRepository) DeleteCommand(ctx context.Context, id int, tx *sqlx.Tx) error {
	deleteQuery := `DELETE FROM product_boms WHERE id = ?`

	if tx != nil {
		_, err := tx.ExecContext(ctx, deleteQuery, id)
		return err
	}
	_, err := repo.db.ExecContext(ctx, deleteQuery, id)
	return err
}
