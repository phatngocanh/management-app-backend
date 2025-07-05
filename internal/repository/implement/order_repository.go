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

type OrderRepository struct {
	db *sqlx.DB
}

func NewOrderRepository(db database.Db) repository.OrderRepository {
	return &OrderRepository{db: db}
}

func (repo *OrderRepository) GetAllQuery(ctx context.Context, tx *sqlx.Tx) ([]entity.Order, error) {
	var orders []entity.Order
	query := "SELECT * FROM orders ORDER BY id DESC"
	var err error

	if tx != nil {
		err = tx.SelectContext(ctx, &orders, query)
	} else {
		err = repo.db.SelectContext(ctx, &orders, query)
	}

	if err != nil {
		return nil, err
	}

	if orders == nil {
		return []entity.Order{}, nil
	}

	return orders, nil
}

func (repo *OrderRepository) GetOneByIDQuery(ctx context.Context, id int, tx *sqlx.Tx) (*entity.Order, error) {
	var order entity.Order
	query := "SELECT * FROM orders WHERE id = ?"
	var err error

	if tx != nil {
		err = tx.GetContext(ctx, &order, query, id)
	} else {
		err = repo.db.GetContext(ctx, &order, query, id)
	}

	if err != nil {
		if err.Error() == error_utils.SystemErrorMessage.SqlxNoRow {
			return nil, nil
		}
		return nil, err
	}

	return &order, nil
}

func (repo *OrderRepository) CreateCommand(ctx context.Context, order *entity.Order, tx *sqlx.Tx) error {
	// First insert without code (code will be generated after getting ID)
	insertQuery := `INSERT INTO orders(code, customer_id, order_date, note, total_original_cost, total_sales_revenue, additional_cost, additional_cost_note, tax_percent) 
					VALUES ('TEMP', :customer_id, :order_date, :note, :total_original_cost, :total_sales_revenue, :additional_cost, :additional_cost_note, :tax_percent)`

	var result sql.Result
	var err error

	if tx != nil {
		result, err = tx.NamedExecContext(ctx, insertQuery, order)
	} else {
		result, err = repo.db.NamedExecContext(ctx, insertQuery, order)
	}

	if err != nil {
		return err
	}

	// Get the inserted ID
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	order.ID = int(id)

	// Generate code based on ID (DH + 5-digit format)
	code := fmt.Sprintf("DH%05d", order.ID)
	order.Code = code

	// Update the record with the generated code
	updateCodeQuery := `UPDATE orders SET code = ? WHERE id = ?`

	if tx != nil {
		_, err = tx.ExecContext(ctx, updateCodeQuery, code, order.ID)
	} else {
		_, err = repo.db.ExecContext(ctx, updateCodeQuery, code, order.ID)
	}

	return err
}

func (repo *OrderRepository) UpdateCommand(ctx context.Context, order *entity.Order, tx *sqlx.Tx) error {
	updateQuery := `UPDATE orders SET customer_id = :customer_id, order_date = :order_date, note = :note, 
					total_original_cost = :total_original_cost, total_sales_revenue = :total_sales_revenue, 
					additional_cost = :additional_cost, additional_cost_note = :additional_cost_note, 
					tax_percent = :tax_percent WHERE id = :id`

	if tx != nil {
		_, err := tx.NamedExecContext(ctx, updateQuery, order)
		return err
	}
	_, err := repo.db.NamedExecContext(ctx, updateQuery, order)
	return err
}

func (repo *OrderRepository) GetByCustomerIDQuery(ctx context.Context, customerID int, tx *sqlx.Tx) ([]entity.Order, error) {
	var orders []entity.Order
	query := "SELECT * FROM orders WHERE customer_id = ? ORDER BY id DESC"
	var err error

	if tx != nil {
		err = tx.SelectContext(ctx, &orders, query, customerID)
	} else {
		err = repo.db.SelectContext(ctx, &orders, query, customerID)
	}

	if err != nil {
		return nil, err
	}

	if orders == nil {
		return []entity.Order{}, nil
	}

	return orders, nil
}

func (repo *OrderRepository) GetAllWithFiltersQuery(ctx context.Context, customerID int, sortBy string, tx *sqlx.Tx) ([]entity.Order, error) {
	var orders []entity.Order
	query := "SELECT * FROM orders WHERE 1=1"
	var args []interface{}

	// Add customer filter
	if customerID > 0 {
		query += " AND customer_id = ?"
		args = append(args, customerID)
	}

	// Add sorting
	switch sortBy {
	case "order_date_asc":
		query += " ORDER BY order_date ASC"
	case "order_date_desc":
		query += " ORDER BY order_date DESC"
	default:
		query += " ORDER BY id DESC"
	}

	var err error
	if tx != nil {
		err = tx.SelectContext(ctx, &orders, query, args...)
	} else {
		err = repo.db.SelectContext(ctx, &orders, query, args...)
	}
	if err != nil {
		return nil, err
	}
	if orders == nil {
		return []entity.Order{}, nil
	}
	return orders, nil
}
