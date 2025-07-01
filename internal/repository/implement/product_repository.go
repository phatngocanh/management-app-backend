package repositoryimplement

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pna/management-app-backend/internal/database"
	"github.com/pna/management-app-backend/internal/domain/entity"
	"github.com/pna/management-app-backend/internal/repository"
	"github.com/pna/management-app-backend/internal/utils/error_utils"
)

type ProductRepository struct {
	db *sqlx.DB
}

func NewProductRepository(db database.Db) repository.ProductRepository {
	return &ProductRepository{db: db}
}

func (repo *ProductRepository) GetAllQuery(ctx context.Context, categoryFilter string, operationTypeFilter string, tx *sqlx.Tx) ([]entity.Product, error) {
	var products []entity.Product
	var query string
	var args []interface{}
	var conditions []string

	// Build WHERE conditions dynamically
	if categoryFilter != "" {
		conditions = append(conditions, "category_id IN ("+categoryFilter+")")
	}

	if operationTypeFilter != "" {
		conditions = append(conditions, "operation_type = ?")
		args = append(args, operationTypeFilter)
	}

	// Construct the final query
	if len(conditions) > 0 {
		query = "SELECT * FROM products WHERE " + conditions[0]
		for i := 1; i < len(conditions); i++ {
			query += " AND " + conditions[i]
		}
	} else {
		query = "SELECT * FROM products"
	}

	query += " ORDER BY id"

	var err error
	if tx != nil {
		err = tx.SelectContext(ctx, &products, query, args...)
	} else {
		err = repo.db.SelectContext(ctx, &products, query, args...)
	}

	if err != nil {
		return nil, err
	}

	if products == nil {
		return []entity.Product{}, nil
	}

	return products, nil
}

func (repo *ProductRepository) GetOneByIDQuery(ctx context.Context, id int, tx *sqlx.Tx) (*entity.Product, error) {
	var product entity.Product
	query := "SELECT * FROM products WHERE id = ?"
	var err error

	if tx != nil {
		err = tx.GetContext(ctx, &product, query, id)
	} else {
		err = repo.db.GetContext(ctx, &product, query, id)
	}

	if err != nil {
		if err.Error() == error_utils.SystemErrorMessage.SqlxNoRow {
			return nil, nil
		}
		return nil, err
	}

	return &product, nil
}

func (repo *ProductRepository) CreateCommand(ctx context.Context, product *entity.Product, tx *sqlx.Tx) error {
	// First insert without code (code will be generated after getting ID)
	insertQuery := `INSERT INTO products(code, name, cost, category_id, unit_id, description, operation_type) 
					VALUES ('TEMP', :name, :cost, :category_id, :unit_id, :description, :operation_type)`

	var result sql.Result
	var err error

	if tx != nil {
		result, err = tx.NamedExecContext(ctx, insertQuery, product)
	} else {
		result, err = repo.db.NamedExecContext(ctx, insertQuery, product)
	}

	if err != nil {
		return err
	}

	// Get the inserted ID
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	product.ID = int(id)

	// Generate code based on ID (SP + 5-digit format)
	code := fmt.Sprintf("SP%05d", product.ID)
	product.Code = code

	// Update the record with the generated code
	updateCodeQuery := `UPDATE products SET code = ? WHERE id = ?`

	if tx != nil {
		_, err = tx.ExecContext(ctx, updateCodeQuery, code, product.ID)
	} else {
		_, err = repo.db.ExecContext(ctx, updateCodeQuery, code, product.ID)
	}

	return err
}

func (repo *ProductRepository) UpdateCommand(ctx context.Context, product *entity.Product, tx *sqlx.Tx) error {
	updateQuery := `UPDATE products SET name = :name, cost = :cost, category_id = :category_id, unit_id = :unit_id, description = :description, operation_type = :operation_type WHERE id = :id`

	if tx != nil {
		_, err := tx.NamedExecContext(ctx, updateQuery, product)
		return err
	}
	_, err := repo.db.NamedExecContext(ctx, updateQuery, product)
	return err
}
