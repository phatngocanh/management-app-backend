package entity

import "time"

type ProductImage struct {
	ID        int       `db:"id"`
	ProductID int       `db:"product_id"` // ID sản phẩm
	ImageKey  string    `db:"image_key"`  // S3 object key
	CreatedAt time.Time `db:"created_at"` // Thời gian tạo
}
