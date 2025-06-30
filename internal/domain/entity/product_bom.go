package entity

import "time"

type ProductBom struct {
	ID                 int       `db:"id"`
	ParentProductID    int       `db:"parent_product_id"`    // ID sản phẩm thành phẩm
	ComponentProductID int       `db:"component_product_id"` // ID sản phẩm nguyên liệu
	Quantity           float64   `db:"quantity"`             // Số lượng nguyên liệu cần thiết
	CreatedAt          time.Time `db:"created_at"`           // Thời gian tạo
	UpdatedAt          time.Time `db:"updated_at"`           // Thời gian cập nhật
}
