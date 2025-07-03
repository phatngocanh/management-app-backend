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
	ProductID       int `json:"product_id" binding:"required"`     // Sản phẩm (một order item chỉ có thể là 1 product)
	Quantity        int `json:"quantity" binding:"required"`       // Số lượng
	SellingPrice    int `json:"selling_price" binding:"required"`  // Giá bán
	OriginalPrice   int `json:"original_price" binding:"required"` // Giá vốn
	DiscountPercent int `json:"discount_percent"`                  // Chiết khấu
}

type OrderResponse struct {
	ID                 int                 `json:"id"`
	Code               string              `json:"code"`
	OrderDate          time.Time           `json:"order_date"`
	Note               *string             `json:"notes"`
	AdditionalCost     int                 `json:"additional_cost"`
	AdditionalCostNote *string             `json:"additional_cost_note"`
	Customer           CustomerResponse    `json:"customer"`
	OrderItems         []OrderItemResponse `json:"order_items,omitempty"`
	Images             []OrderImage        `json:"images,omitempty"`
	TotalAmount        *int                `json:"total_amount,omitempty"`
	ProductCount       *int                `json:"product_count,omitempty"`
	TaxPercent         int                 `json:"tax_percent"`
	// Profit/Loss fields for total order
	TotalProfitLoss           *int     `json:"total_profit_loss,omitempty"`            // Total profit/loss for the order
	TotalProfitLossPercentage *float64 `json:"total_profit_loss_percentage,omitempty"` // Total profit/loss percentage for the order
	TotalSalesRevenue         int      `json:"total_sales_revenue"`                    // Total sales revenue for the order
}

type OrderItemResponse struct {
	ID              int    `json:"id"`
	ProductName     string `json:"product_name"`
	OrderID         int    `json:"order_id"`
	ProductID       int    `json:"product_id"`
	Quantity        int    `json:"quantity"`
	SellingPrice    int    `json:"selling_price"`
	DiscountPercent int    `json:"discount_percent"`
	FinalAmount     *int   `json:"final_amount"`
	// Profit/Loss fields
	OriginalPrice        *int     `json:"original_price,omitempty"`         // Product's original price
	ProfitLoss           *int     `json:"profit_loss,omitempty"`            // Profit/Loss amount for this item
	ProfitLossPercentage *float64 `json:"profit_loss_percentage,omitempty"` // Profit/Loss percentage for this item
}

type GetOneOrderResponse struct {
	Order OrderResponse `json:"order"`
}
