package entity

type OrderItem struct {
	ID              int  `db:"id"`
	OrderID         int  `db:"order_id"`         // Đơn hàng
	ProductID       *int `db:"product_id"`       // Sản phẩm (nullable - một order item chỉ có thể là 1 product hoặc 1 bom)
	BomID           *int `db:"bom_id"`           // BOM (nullable - một order item chỉ có thể là 1 product hoặc 1 bom)
	Quantity        int  `db:"quantity"`         // Số lượng
	SellingPrice    int  `db:"selling_price"`    // Giá bán
	OriginalPrice   int  `db:"original_price"`   // Giá vốn
	DiscountPercent int  `db:"discount_percent"` // Chiết khấu
	FinalAmount     int  `db:"final_amount"`     // Tổng tiền
}
