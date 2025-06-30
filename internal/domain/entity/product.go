package entity

type Product struct {
	ID          int     `db:"id"`
	Name        string  `db:"name"`        // Tên sản phẩm
	Cost        float64 `db:"cost"`        // Giá vốn của sản phẩm (VND)
	CategoryID  *int    `db:"category_id"` // ID danh mục sản phẩm
	UnitID      *int    `db:"unit_id"`     // ID đơn vị tính
	Description string  `db:"description"` // Mô tả chi tiết sản phẩm
}
