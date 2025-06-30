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

type ProductCategoryRepository struct {
	db *sqlx.DB
}

func NewProductCategoryRepository(db database.Db) repository.ProductCategoryRepository {
	return &ProductCategoryRepository{db: db}
}

func (repo *ProductCategoryRepository) GetAllQuery(ctx context.Context, tx *sqlx.Tx) ([]entity.ProductCategory, error) {
	var categories []entity.ProductCategory
	query := "SELECT * FROM product_categories ORDER BY id"
	var err error

	if tx != nil {
		err = tx.SelectContext(ctx, &categories, query)
	} else {
		err = repo.db.SelectContext(ctx, &categories, query)
	}

	if err != nil {
		return nil, err
	}

	if categories == nil {
		return []entity.ProductCategory{}, nil
	}

	return categories, nil
}

func (repo *ProductCategoryRepository) GetOneByIDQuery(ctx context.Context, id int, tx *sqlx.Tx) (*entity.ProductCategory, error) {
	var category entity.ProductCategory
	query := "SELECT * FROM product_categories WHERE id = ?"
	var err error

	if tx != nil {
		err = tx.GetContext(ctx, &category, query, id)
	} else {
		err = repo.db.GetContext(ctx, &category, query, id)
	}

	if err != nil {
		if err.Error() == error_utils.SystemErrorMessage.SqlxNoRow {
			return nil, nil
		}
		return nil, err
	}

	return &category, nil
}

func (repo *ProductCategoryRepository) GetOneByCodeQuery(ctx context.Context, code string, tx *sqlx.Tx) (*entity.ProductCategory, error) {
	var category entity.ProductCategory
	query := "SELECT * FROM product_categories WHERE code = ?"
	var err error

	if tx != nil {
		err = tx.GetContext(ctx, &category, query, code)
	} else {
		err = repo.db.GetContext(ctx, &category, query, code)
	}

	if err != nil {
		if err.Error() == error_utils.SystemErrorMessage.SqlxNoRow {
			return nil, nil
		}
		return nil, err
	}

	return &category, nil
}

func (repo *ProductCategoryRepository) CreateCommand(ctx context.Context, category *entity.ProductCategory, tx *sqlx.Tx) error {
	insertQuery := `INSERT INTO product_categories(name, code, description) VALUES (:name, :code, :description)`

	var result sql.Result
	var err error

	if tx != nil {
		result, err = tx.NamedExecContext(ctx, insertQuery, category)
	} else {
		result, err = repo.db.NamedExecContext(ctx, insertQuery, category)
	}

	if err != nil {
		return err
	}

	// Get the last inserted ID
	lastID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	// Set the ID to the category entity
	category.ID = int(lastID)
	return nil
}

func (repo *ProductCategoryRepository) UpdateCommand(ctx context.Context, category *entity.ProductCategory, tx *sqlx.Tx) error {
	updateQuery := `UPDATE product_categories SET name = :name, code = :code, description = :description WHERE id = :id`

	if tx != nil {
		_, err := tx.NamedExecContext(ctx, updateQuery, category)
		return err
	}
	_, err := repo.db.NamedExecContext(ctx, updateQuery, category)
	return err
}
