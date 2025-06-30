package entity

import "time"

type ProductImage struct {
	ID                 int        `db:"id"`
	ProductID          int        `db:"product_id"`            // ID sản phẩm
	ImageKey           string     `db:"image_key"`             // S3 object key
	IsPrimary          bool       `db:"is_primary"`            // Hình ảnh chính hay không
	SignedURLExpiresAt *time.Time `db:"signed_url_expires_at"` // Thời gian hết hạn của signed URL
	CreatedAt          time.Time  `db:"created_at"`            // Thời gian tạo
}
