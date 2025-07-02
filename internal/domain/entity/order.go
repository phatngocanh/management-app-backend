package entity

import "time"

type Order struct {
	ID                 int       `db:"id"`
	Code               string    `db:"code"`                 // Mã đơn hàng (DH00001)
	CustomerID         int       `db:"customer_id"`          // Khách hàng
	OrderDate          time.Time `db:"order_date"`           // Ngày đặt hàng
	Note               *string   `db:"note"`                 // Ghi chú
	TotalOriginalCost  int       `db:"total_original_cost"`  // Tổng giá vốn sản phẩm
	TotalSalesRevenue  int       `db:"total_sales_revenue"`  // Tổng doanh thu
	AdditionalCost     int       `db:"additional_cost"`      // Chi phí phát sinh
	AdditionalCostNote *string   `db:"additional_cost_note"` // Ghi chú chi phí phát sinh
	TaxPercent         int       `db:"tax_percent"`          // Thuế suất
}
