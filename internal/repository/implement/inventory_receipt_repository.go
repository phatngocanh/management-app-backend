package repositoryimplement

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pna/management-app-backend/internal/database"
	"github.com/pna/management-app-backend/internal/domain/entity"
	"github.com/pna/management-app-backend/internal/repository"
	"github.com/pna/management-app-backend/internal/utils/error_utils"
)

type InventoryReceiptRepository struct {
	db *sqlx.DB
}

func NewInventoryReceiptRepository(db database.Db) repository.InventoryReceiptRepository {
	return &InventoryReceiptRepository{db: db}
}

func (repo *InventoryReceiptRepository) CreateCommand(ctx context.Context, receipt *entity.InventoryReceipt, tx *sqlx.Tx) error {
	insertQuery := `INSERT INTO inventory_receipt(user_id, receipt_date, notes, total_items) 
					VALUES (:user_id, :receipt_date, :notes, :total_items)`

	if tx != nil {
		result, err := tx.NamedExecContext(ctx, insertQuery, receipt)
		if err != nil {
			return err
		}
		id, err := result.LastInsertId()
		if err != nil {
			return err
		}
		receipt.ID = int(id)
		return nil
	}

	result, err := repo.db.NamedExecContext(ctx, insertQuery, receipt)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	receipt.ID = int(id)
	return nil
}

func (repo *InventoryReceiptRepository) GetAllQuery(ctx context.Context, tx *sqlx.Tx) ([]entity.InventoryReceipt, error) {
	var receipts []entity.InventoryReceipt
	query := "SELECT * FROM inventory_receipt ORDER BY created_at DESC"
	var err error

	if tx != nil {
		err = tx.SelectContext(ctx, &receipts, query)
	} else {
		err = repo.db.SelectContext(ctx, &receipts, query)
	}

	if err != nil {
		return nil, err
	}

	return receipts, nil
}

func (repo *InventoryReceiptRepository) GetOneByIDQuery(ctx context.Context, id int, tx *sqlx.Tx) (*entity.InventoryReceipt, error) {
	var receipt entity.InventoryReceipt
	query := "SELECT * FROM inventory_receipt WHERE id = ?"
	var err error

	if tx != nil {
		err = tx.GetContext(ctx, &receipt, query, id)
	} else {
		err = repo.db.GetContext(ctx, &receipt, query, id)
	}

	if err != nil {
		if err.Error() == error_utils.SystemErrorMessage.SqlxNoRow {
			return nil, nil
		}
		return nil, err
	}

	return &receipt, nil
}

func (repo *InventoryReceiptRepository) GetByUserIDQuery(ctx context.Context, userID int, tx *sqlx.Tx) ([]entity.InventoryReceipt, error) {
	var receipts []entity.InventoryReceipt
	query := "SELECT * FROM inventory_receipt WHERE user_id = ? ORDER BY created_at DESC"
	var err error

	if tx != nil {
		err = tx.SelectContext(ctx, &receipts, query, userID)
	} else {
		err = repo.db.SelectContext(ctx, &receipts, query, userID)
	}

	if err != nil {
		return nil, err
	}

	return receipts, nil
}

func (repo *InventoryReceiptRepository) UpdateCommand(ctx context.Context, receipt *entity.InventoryReceipt, tx *sqlx.Tx) error {
	updateQuery := `UPDATE inventory_receipt SET user_id = :user_id, receipt_date = :receipt_date, 
					notes = :notes, total_items = :total_items WHERE id = :id`

	if tx != nil {
		_, err := tx.NamedExecContext(ctx, updateQuery, receipt)
		return err
	}

	_, err := repo.db.NamedExecContext(ctx, updateQuery, receipt)
	return err
}

func (repo *InventoryReceiptRepository) DeleteCommand(ctx context.Context, id int, tx *sqlx.Tx) error {
	deleteQuery := "DELETE FROM inventory_receipt WHERE id = ?"
	var err error

	if tx != nil {
		_, err = tx.ExecContext(ctx, deleteQuery, id)
	} else {
		_, err = repo.db.ExecContext(ctx, deleteQuery, id)
	}

	return err
}
