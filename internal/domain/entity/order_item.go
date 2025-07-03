package entity

type OrderItem struct {
	ID              int  `db:"id"`
	OrderID         int  `db:"order_id"`         // Đơn hàng
	ProductID       *int `db:"product_id"`       // Sản phẩm (có thể là direct product hoặc parent product từ BOM)
	Quantity        int  `db:"quantity"`         // Số lượng
	SellingPrice    int  `db:"selling_price"`    // Giá bán
	OriginalPrice   int  `db:"original_price"`   // Giá vốn
	DiscountPercent int  `db:"discount_percent"` // Chiết khấu
	FinalAmount     int  `db:"final_amount"`     // Tổng tiền
}
