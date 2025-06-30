package entity

import "time"

type ProductCategory struct {
	ID          int       `db:"id"`
	Name        string    `db:"name"`        // Tên danh mục
	Code        string    `db:"code"`        // Mã danh mục
	Description string    `db:"description"` // Mô tả danh mục
	CreatedAt   time.Time `db:"created_at"`  // Thời gian tạo
}
