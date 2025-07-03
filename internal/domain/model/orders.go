package model

import "time"

type CreateOrderRequest struct {
	CustomerID         int                      `json:"customer_id" binding:"required"` // Khách hàng
	OrderDate          time.Time                `json:"order_date" binding:"required"`  // Ngày đặt hàng
	Note               *string                  `json:"note"`                           // Ghi chú
	AdditionalCost     int                      `json:"additional_cost"`                // Chi phí phát sinh
	AdditionalCostNote *string                  `json:"additional_cost_note"`           // Ghi chú chi phí phát sinh
	TaxPercent         int                      `json:"tax_percent"`                    // Thuế suất
	Items              []CreateOrderItemRequest `json:"items" binding:"required,dive"`  // Danh sách sản phẩm trong đơn hàng
}

type CreateOrderItemRequest struct {
	ProductID       *int `json:"product_id,omitempty"`              // Sản phẩm (nullable - một order item chỉ có thể là 1 product hoặc 1 bom parent)
	BomParentID     *int `json:"bom_parent_id,omitempty"`           // BOM Parent Product ID - sản phẩm cha cần sản xuất (nullable - một order item chỉ có thể là 1 product hoặc 1 bom parent)
	Quantity        int  `json:"quantity" binding:"required"`       // Số lượng
	SellingPrice    int  `json:"selling_price" binding:"required"`  // Giá bán
	OriginalPrice   int  `json:"original_price" binding:"required"` // Giá vốn
	DiscountPercent int  `json:"discount_percent"`                  // Chiết khấu
}

type UpdateOrderRequest struct {
	ID                 int                      `json:"id" binding:"required"`
	CustomerID         int                      `json:"customer_id" binding:"required"`         // Khách hàng
	OrderDate          time.Time                `json:"order_date" binding:"required"`          // Ngày đặt hàng
	Note               *string                  `json:"note"`                                   // Ghi chú
	TotalOriginalCost  int                      `json:"total_original_cost" binding:"required"` // Tổng giá vốn sản phẩm
	TotalSalesRevenue  int                      `json:"total_sales_revenue" binding:"required"` // Tổng doanh thu
	AdditionalCost     int                      `json:"additional_cost"`                        // Chi phí phát sinh
	AdditionalCostNote *string                  `json:"additional_cost_note"`                   // Ghi chú chi phí phát sinh
	TaxPercent         int                      `json:"tax_percent" binding:"required"`         // Thuế suất
	Items              []UpdateOrderItemRequest `json:"items" binding:"required,dive"`          // Danh sách sản phẩm trong đơn hàng
}

type UpdateOrderItemRequest struct {
	ID              int  `json:"id"`                                // ID của order item (0 nếu là item mới)
	ProductID       *int `json:"product_id"`                        // Sản phẩm (có thể là direct product hoặc parent product từ BOM)
	Quantity        int  `json:"quantity" binding:"required"`       // Số lượng
	SellingPrice    int  `json:"selling_price" binding:"required"`  // Giá bán
	OriginalPrice   int  `json:"original_price" binding:"required"` // Giá vốn
	DiscountPercent int  `json:"discount_percent"`                  // Chiết khấu
	FinalAmount     int  `json:"final_amount" binding:"required"`   // Tổng tiền
}

type OrderResponse struct {
	ID                 int                 `json:"id"`
	Code               string              `json:"code"`                 // Mã đơn hàng (DH00001)
	CustomerID         int                 `json:"customer_id"`          // Khách hàng
	OrderDate          time.Time           `json:"order_date"`           // Ngày đặt hàng
	Note               *string             `json:"note"`                 // Ghi chú
	TotalOriginalCost  int                 `json:"total_original_cost"`  // Tổng giá vốn sản phẩm
	TotalSalesRevenue  int                 `json:"total_sales_revenue"`  // Tổng doanh thu
	AdditionalCost     int                 `json:"additional_cost"`      // Chi phí phát sinh
	AdditionalCostNote *string             `json:"additional_cost_note"` // Ghi chú chi phí phát sinh
	TaxPercent         int                 `json:"tax_percent"`          // Thuế suất
	Customer           *CustomerResponse   `json:"customer,omitempty"`   // Thông tin khách hàng
	Items              []OrderItemResponse `json:"items,omitempty"`      // Danh sách sản phẩm trong đơn hàng
}

type OrderItemResponse struct {
	ID              int              `json:"id"`
	OrderID         int              `json:"order_id"`          // Đơn hàng
	ProductID       *int             `json:"product_id"`        // Sản phẩm (có thể là direct product hoặc parent product từ BOM)
	Quantity        int              `json:"quantity"`          // Số lượng
	SellingPrice    int              `json:"selling_price"`     // Giá bán
	OriginalPrice   int              `json:"original_price"`    // Giá vốn
	DiscountPercent int              `json:"discount_percent"`  // Chiết khấu
	FinalAmount     int              `json:"final_amount"`      // Tổng tiền
	Product         *ProductResponse `json:"product,omitempty"` // Thông tin sản phẩm
}

type GetAllOrdersResponse struct {
	Orders []OrderResponse `json:"orders"`
}

type GetOneOrderResponse struct {
	Order OrderResponse `json:"order"`
}
