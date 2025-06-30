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

type InventoryReceiptRepository struct {
	db *sqlx.DB
}

func NewInventoryReceiptRepository(db database.Db) repository.InventoryReceiptRepository {
	return &InventoryReceiptRepository{db: db}
}

func (repo *InventoryReceiptRepository) CreateCommand(ctx context.Context, receipt *entity.InventoryReceipt, tx *sqlx.Tx) error {
	// First insert without code (code will be generated after getting ID)
	insertQuery := `INSERT INTO inventory_receipts(code, user_id, receipt_date, notes, total_items) 
					VALUES ('TEMP', :user_id, :receipt_date, :notes, :total_items)`

	var result sql.Result
	var err error

	if tx != nil {
		result, err = tx.NamedExecContext(ctx, insertQuery, receipt)
	} else {
		result, err = repo.db.NamedExecContext(ctx, insertQuery, receipt)
	}

	if err != nil {
		return err
	}

	// Get the inserted ID
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	receipt.ID = int(id)

	// Generate code based on ID (NH + 5-digit format)
	code := fmt.Sprintf("NK%05d", receipt.ID)
	receipt.Code = code

	// Update the record with the generated code
	updateCodeQuery := `UPDATE inventory_receipts SET code = ? WHERE id = ?`

	if tx != nil {
		_, err = tx.ExecContext(ctx, updateCodeQuery, code, receipt.ID)
	} else {
		_, err = repo.db.ExecContext(ctx, updateCodeQuery, code, receipt.ID)
	}

	return err
}

func (repo *InventoryReceiptRepository) GetAllQuery(ctx context.Context, tx *sqlx.Tx) ([]entity.InventoryReceipt, error) {
	var receipts []entity.InventoryReceipt
	query := "SELECT * FROM inventory_receipts ORDER BY created_at DESC"
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
	query := "SELECT * FROM inventory_receipts WHERE id = ?"
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

func (repo *InventoryReceiptRepository) GetOneByCodeQuery(ctx context.Context, code string, tx *sqlx.Tx) (*entity.InventoryReceipt, error) {
	var receipt entity.InventoryReceipt
	query := "SELECT * FROM inventory_receipts WHERE code = ?"
	var err error

	if tx != nil {
		err = tx.GetContext(ctx, &receipt, query, code)
	} else {
		err = repo.db.GetContext(ctx, &receipt, query, code)
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
	query := "SELECT * FROM inventory_receipts WHERE user_id = ? ORDER BY created_at DESC"
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
	updateQuery := `UPDATE inventory_receipts SET code = :code, user_id = :user_id, receipt_date = :receipt_date, 
					notes = :notes, total_items = :total_items WHERE id = :id`

	if tx != nil {
		_, err := tx.NamedExecContext(ctx, updateQuery, receipt)
		return err
	}

	_, err := repo.db.NamedExecContext(ctx, updateQuery, receipt)
	return err
}

func (repo *InventoryReceiptRepository) DeleteCommand(ctx context.Context, id int, tx *sqlx.Tx) error {
	deleteQuery := "DELETE FROM inventory_receipts WHERE id = ?"
	var err error

	if tx != nil {
		_, err = tx.ExecContext(ctx, deleteQuery, id)
	} else {
		_, err = repo.db.ExecContext(ctx, deleteQuery, id)
	}

	return err
}
