package entity

type Product struct {
	ID            int    `db:"id"`
	Name          string `db:"name"`           // Tên sản phẩm
	Spec          int    `db:"spec"`           // Quy cách
	OriginalPrice int    `db:"original_price"` // Giá gốc của sản phẩm (VND)
	CategoryID    *int   `db:"category_id"`    // ID danh mục sản phẩm
	UnitID        *int   `db:"unit_id"`        // ID đơn vị tính
	Description   string `db:"description"`    // Mô tả chi tiết sản phẩm
}
