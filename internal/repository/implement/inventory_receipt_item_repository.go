package repositoryimplement

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pna/management-app-backend/internal/database"
	"github.com/pna/management-app-backend/internal/domain/entity"
	"github.com/pna/management-app-backend/internal/repository"
	"github.com/pna/management-app-backend/internal/utils/error_utils"
)

type InventoryReceiptItemRepository struct {
	db *sqlx.DB
}

func NewInventoryReceiptItemRepository(db database.Db) repository.InventoryReceiptItemRepository {
	return &InventoryReceiptItemRepository{db: db}
}

func (repo *InventoryReceiptItemRepository) CreateCommand(ctx context.Context, item *entity.InventoryReceiptItem, tx *sqlx.Tx) error {
	insertQuery := `INSERT INTO inventory_receipt_items(inventory_receipt_id, product_id, quantity, unit_cost, notes) 
					VALUES (:inventory_receipt_id, :product_id, :quantity, :unit_cost, :notes)`

	if tx != nil {
		result, err := tx.NamedExecContext(ctx, insertQuery, item)
		if err != nil {
			return err
		}
		id, err := result.LastInsertId()
		if err != nil {
			return err
		}
		item.ID = int(id)
		return nil
	}

	result, err := repo.db.NamedExecContext(ctx, insertQuery, item)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	item.ID = int(id)
	return nil
}

func (repo *InventoryReceiptItemRepository) GetAllQuery(ctx context.Context, tx *sqlx.Tx) ([]entity.InventoryReceiptItem, error) {
	var items []entity.InventoryReceiptItem
	query := "SELECT * FROM inventory_receipt_items ORDER BY created_at DESC"
	var err error

	if tx != nil {
		err = tx.SelectContext(ctx, &items, query)
	} else {
		err = repo.db.SelectContext(ctx, &items, query)
	}

	if err != nil {
		return nil, err
	}

	return items, nil
}

func (repo *InventoryReceiptItemRepository) GetOneByIDQuery(ctx context.Context, id int, tx *sqlx.Tx) (*entity.InventoryReceiptItem, error) {
	var item entity.InventoryReceiptItem
	query := "SELECT * FROM inventory_receipt_items WHERE id = ?"
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

func (repo *InventoryReceiptItemRepository) GetByInventoryReceiptIDQuery(ctx context.Context, inventoryReceiptID int, tx *sqlx.Tx) ([]entity.InventoryReceiptItem, error) {
	var items []entity.InventoryReceiptItem
	query := "SELECT * FROM inventory_receipt_items WHERE inventory_receipt_id = ? ORDER BY created_at DESC"
	var err error

	if tx != nil {
		err = tx.SelectContext(ctx, &items, query, inventoryReceiptID)
	} else {
		err = repo.db.SelectContext(ctx, &items, query, inventoryReceiptID)
	}

	if err != nil {
		return nil, err
	}

	return items, nil
}

func (repo *InventoryReceiptItemRepository) GetByProductIDQuery(ctx context.Context, productID int, tx *sqlx.Tx) ([]entity.InventoryReceiptItem, error) {
	var items []entity.InventoryReceiptItem
	query := "SELECT * FROM inventory_receipt_items WHERE product_id = ? ORDER BY created_at DESC"
	var err error

	if tx != nil {
		err = tx.SelectContext(ctx, &items, query, productID)
	} else {
		err = repo.db.SelectContext(ctx, &items, query, productID)
	}

	if err != nil {
		return nil, err
	}

	return items, nil
}

func (repo *InventoryReceiptItemRepository) UpdateCommand(ctx context.Context, item *entity.InventoryReceiptItem, tx *sqlx.Tx) error {
	updateQuery := `UPDATE inventory_receipt_items SET inventory_receipt_id = :inventory_receipt_id, 
					product_id = :product_id, quantity = :quantity, unit_cost = :unit_cost, 
					notes = :notes WHERE id = :id`

	if tx != nil {
		_, err := tx.NamedExecContext(ctx, updateQuery, item)
		return err
	}

	_, err := repo.db.NamedExecContext(ctx, updateQuery, item)
	return err
}

func (repo *InventoryReceiptItemRepository) DeleteCommand(ctx context.Context, id int, tx *sqlx.Tx) error {
	deleteQuery := "DELETE FROM inventory_receipt_items WHERE id = ?"
	var err error

	if tx != nil {
		_, err = tx.ExecContext(ctx, deleteQuery, id)
	} else {
		_, err = repo.db.ExecContext(ctx, deleteQuery, id)
	}

	return err
}

func (repo *InventoryReceiptItemRepository) DeleteByInventoryReceiptIDCommand(ctx context.Context, inventoryReceiptID int, tx *sqlx.Tx) error {
	deleteQuery := "DELETE FROM inventory_receipt_items WHERE inventory_receipt_id = ?"
	var err error

	if tx != nil {
		_, err = tx.ExecContext(ctx, deleteQuery, inventoryReceiptID)
	} else {
		_, err = repo.db.ExecContext(ctx, deleteQuery, inventoryReceiptID)
	}

	return err
}
