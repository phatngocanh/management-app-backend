package entity

import "time"

type ProductImage struct {
	ID        int       `db:"id"`
	ProductID int       `db:"product_id"` // ID sản phẩm
	ImageURL  string    `db:"image_url"`  // URL hình ảnh S3
	ImageKey  string    `db:"image_key"`  // S3 object key
	IsPrimary bool      `db:"is_primary"` // Hình ảnh chính hay không
	CreatedAt time.Time `db:"created_at"` // Thời gian tạo
}
