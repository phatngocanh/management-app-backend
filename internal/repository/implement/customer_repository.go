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

type CustomerRepository struct {
	db *sqlx.DB
}

func NewCustomerRepository(db database.Db) repository.CustomerRepository {
	return &CustomerRepository{db: db}
}

func (repo *CustomerRepository) GetAllQuery(ctx context.Context, tx *sqlx.Tx) ([]entity.Customer, error) {
	var customers []entity.Customer
	query := "SELECT * FROM customers ORDER BY id"
	var err error

	if tx != nil {
		err = tx.SelectContext(ctx, &customers, query)
	} else {
		err = repo.db.SelectContext(ctx, &customers, query)
	}

	if err != nil {
		return nil, err
	}

	if customers == nil {
		return []entity.Customer{}, nil
	}

	return customers, nil
}

func (repo *CustomerRepository) GetOneByIDQuery(ctx context.Context, id int, tx *sqlx.Tx) (*entity.Customer, error) {
	var customer entity.Customer
	query := "SELECT * FROM customers WHERE id = ?"
	var err error

	if tx != nil {
		err = tx.GetContext(ctx, &customer, query, id)
	} else {
		err = repo.db.GetContext(ctx, &customer, query, id)
	}

	if err != nil {
		if err.Error() == error_utils.SystemErrorMessage.SqlxNoRow {
			return nil, nil
		}
		return nil, err
	}

	return &customer, nil
}

func (repo *CustomerRepository) CreateCommand(ctx context.Context, customer *entity.Customer, tx *sqlx.Tx) error {
	// First insert without code (code will be generated after getting ID)
	insertQuery := `INSERT INTO customers(code, name, phone, address) VALUES ('TEMP', :name, :phone, :address)`

	var result sql.Result
	var err error

	if tx != nil {
		result, err = tx.NamedExecContext(ctx, insertQuery, customer)
	} else {
		result, err = repo.db.NamedExecContext(ctx, insertQuery, customer)
	}

	if err != nil {
		return err
	}

	// Get the inserted ID
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	customer.ID = int(id)

	// Generate code based on ID (KH + 5-digit format)
	code := fmt.Sprintf("KH%05d", customer.ID)
	customer.Code = code

	// Update the record with the generated code
	updateCodeQuery := `UPDATE customers SET code = ? WHERE id = ?`

	if tx != nil {
		_, err = tx.ExecContext(ctx, updateCodeQuery, code, customer.ID)
	} else {
		_, err = repo.db.ExecContext(ctx, updateCodeQuery, code, customer.ID)
	}

	return err
}

func (repo *CustomerRepository) UpdateCommand(ctx context.Context, customer *entity.Customer, tx *sqlx.Tx) error {
	updateQuery := `UPDATE customers SET name = :name, phone = :phone, address = :address WHERE id = :id`

	if tx != nil {
		_, err := tx.NamedExecContext(ctx, updateQuery, customer)
		return err
	}
	_, err := repo.db.NamedExecContext(ctx, updateQuery, customer)
	return err
}
