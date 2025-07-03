package repositoryimplement

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pna/management-app-backend/internal/database"
	"github.com/pna/management-app-backend/internal/domain/entity"
	"github.com/pna/management-app-backend/internal/repository"
	"github.com/pna/management-app-backend/internal/utils/error_utils"
)

type OrderItemRepository struct {
	db *sqlx.DB
}

func NewOrderItemRepository(db database.Db) repository.OrderItemRepository {
	return &OrderItemRepository{db: db}
}

func (repo *OrderItemRepository) CreateCommand(ctx context.Context, item *entity.OrderItem, tx *sqlx.Tx) error {
	insertQuery := `INSERT INTO order_items(order_id, product_id, quantity, selling_price, original_price, discount_percent, final_amount)
					VALUES (:order_id, :product_id, :quantity, :selling_price, :original_price, :discount_percent, :final_amount)`

	var err error

	if tx != nil {
		_, err = tx.NamedExecContext(ctx, insertQuery, item)
	} else {
		_, err = repo.db.NamedExecContext(ctx, insertQuery, item)
	}

	return err
}

func (repo *OrderItemRepository) GetAllQuery(ctx context.Context, tx *sqlx.Tx) ([]entity.OrderItem, error) {
	var items []entity.OrderItem
	query := "SELECT * FROM order_items ORDER BY id"
	var err error

	if tx != nil {
		err = tx.SelectContext(ctx, &items, query)
	} else {
		err = repo.db.SelectContext(ctx, &items, query)
	}

	if err != nil {
		return nil, err
	}

	if items == nil {
		return []entity.OrderItem{}, nil
	}

	return items, nil
}

func (repo *OrderItemRepository) GetOneByIDQuery(ctx context.Context, id int, tx *sqlx.Tx) (*entity.OrderItem, error) {
	var item entity.OrderItem
	query := "SELECT * FROM order_items WHERE id = ?"
	var err error

	if tx != nil {
		err = tx.GetContext(ctx, &item, query, id)
	} else {
		err = repo.db.GetContext(ctx, &item, query, id)
	}

	if err != nil {
		if err.Error() == error_utils.SystemErrorMessage.SqlxNoRow {
			return nil, nil
		}
		return nil, err
	}

	return &item, nil
}

func (repo *OrderItemRepository) GetAllByOrderIDQuery(ctx context.Context, orderID int, tx *sqlx.Tx) ([]entity.OrderItem, error) {
	var items []entity.OrderItem
	query := "SELECT * FROM order_items WHERE order_id = ? ORDER BY id"
	var err error

	if tx != nil {
		err = tx.SelectContext(ctx, &items, query, orderID)
	} else {
		err = repo.db.SelectContext(ctx, &items, query, orderID)
	}

	if err != nil {
		return nil, err
	}

	if items == nil {
		return []entity.OrderItem{}, nil
	}

	return items, nil
}

func (repo *OrderItemRepository) UpdateCommand(ctx context.Context, item *entity.OrderItem, tx *sqlx.Tx) error {
	updateQuery := `UPDATE order_items SET order_id = :order_id, product_id = :product_id, 
					quantity = :quantity, selling_price = :selling_price, original_price = :original_price, 
					discount_percent = :discount_percent, final_amount = :final_amount WHERE id = :id`

	if tx != nil {
		_, err := tx.NamedExecContext(ctx, updateQuery, item)
		return err
	}
	_, err := repo.db.NamedExecContext(ctx, updateQuery, item)
	return err
}

func (repo *OrderItemRepository) DeleteCommand(ctx context.Context, id int, tx *sqlx.Tx) error {
	deleteQuery := "DELETE FROM order_items WHERE id = ?"

	if tx != nil {
		_, err := tx.ExecContext(ctx, deleteQuery, id)
		return err
	}
	_, err := repo.db.ExecContext(ctx, deleteQuery, id)
	return err
}

func (repo *OrderItemRepository) DeleteByOrderIDCommand(ctx context.Context, orderID int, tx *sqlx.Tx) error {
	deleteQuery := "DELETE FROM order_items WHERE order_id = ?"

	if tx != nil {
		_, err := tx.ExecContext(ctx, deleteQuery, orderID)
		return err
	}
	_, err := repo.db.ExecContext(ctx, deleteQuery, orderID)
	return err
}
