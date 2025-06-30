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

type ProductImageRepository struct {
	db *sqlx.DB
}

func NewProductImageRepository(db database.Db) repository.ProductImageRepository {
	return &ProductImageRepository{db: db}
}

func (repo *ProductImageRepository) GetAllQuery(ctx context.Context, tx *sqlx.Tx) ([]entity.ProductImage, error) {
	var images []entity.ProductImage
	query := "SELECT * FROM product_images ORDER BY id"
	var err error

	if tx != nil {
		err = tx.SelectContext(ctx, &images, query)
	} else {
		err = repo.db.SelectContext(ctx, &images, query)
	}

	if err != nil {
		return nil, err
	}

	if images == nil {
		return []entity.ProductImage{}, nil
	}

	return images, nil
}

func (repo *ProductImageRepository) GetOneByIDQuery(ctx context.Context, id int, tx *sqlx.Tx) (*entity.ProductImage, error) {
	var image entity.ProductImage
	query := "SELECT * FROM product_images WHERE id = ?"
	var err error

	if tx != nil {
		err = tx.GetContext(ctx, &image, query, id)
	} else {
		err = repo.db.GetContext(ctx, &image, query, id)
	}

	if err != nil {
		if err.Error() == error_utils.SystemErrorMessage.SqlxNoRow {
			return nil, nil
		}
		return nil, err
	}

	return &image, nil
}

func (repo *ProductImageRepository) GetByProductIDQuery(ctx context.Context, productID int, tx *sqlx.Tx) ([]entity.ProductImage, error) {
	var images []entity.ProductImage
	query := "SELECT * FROM product_images WHERE product_id = ? ORDER BY is_primary DESC, id"
	var err error

	if tx != nil {
		err = tx.SelectContext(ctx, &images, query, productID)
	} else {
		err = repo.db.SelectContext(ctx, &images, query, productID)
	}

	if err != nil {
		return nil, err
	}

	if images == nil {
		return []entity.ProductImage{}, nil
	}

	return images, nil
}

func (repo *ProductImageRepository) CreateCommand(ctx context.Context, image *entity.ProductImage, tx *sqlx.Tx) error {
	insertQuery := `INSERT INTO product_images(product_id, image_url, image_key, is_primary) VALUES (:product_id, :image_url, :image_key, :is_primary)`

	var result sql.Result
	var err error

	if tx != nil {
		result, err = tx.NamedExecContext(ctx, insertQuery, image)
	} else {
		result, err = repo.db.NamedExecContext(ctx, insertQuery, image)
	}

	if err != nil {
		return err
	}

	// Get the last inserted ID
	lastID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	// Set the ID to the image entity
	image.ID = int(lastID)
	return nil
}

func (repo *ProductImageRepository) UpdateCommand(ctx context.Context, image *entity.ProductImage, tx *sqlx.Tx) error {
	updateQuery := `UPDATE product_images SET product_id = :product_id, image_url = :image_url, image_key = :image_key, is_primary = :is_primary WHERE id = :id`

	if tx != nil {
		_, err := tx.NamedExecContext(ctx, updateQuery, image)
		return err
	}
	_, err := repo.db.NamedExecContext(ctx, updateQuery, image)
	return err
}

func (repo *ProductImageRepository) DeleteCommand(ctx context.Context, id int, tx *sqlx.Tx) error {
	deleteQuery := `DELETE FROM product_images WHERE id = ?`

	if tx != nil {
		_, err := tx.ExecContext(ctx, deleteQuery, id)
		return err
	}
	_, err := repo.db.ExecContext(ctx, deleteQuery, id)
	return err
}
