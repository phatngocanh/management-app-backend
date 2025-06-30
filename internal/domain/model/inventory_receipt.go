package model

import "time"

type CreateInventoryReceiptRequest struct {
	Notes string                              `json:"notes"`
	Items []CreateInventoryReceiptItemRequest `json:"items" binding:"required"`
}

type CreateInventoryReceiptItemRequest struct {
	ProductID int      `json:"product_id" binding:"required"`
	Quantity  int      `json:"quantity" binding:"required"`
	UnitCost  *float64 `json:"unit_cost"`
	Notes     string   `json:"notes"`
}

type InventoryReceiptResponse struct {
	ID          int                            `json:"id"`
	UserID      int                            `json:"user_id"`
	ReceiptDate time.Time                      `json:"receipt_date"`
	Notes       *string                        `json:"notes"`
	TotalItems  int                            `json:"total_items"`
	CreatedAt   time.Time                      `json:"created_at"`
	UpdatedAt   time.Time                      `json:"updated_at"`
	Items       []InventoryReceiptItemResponse `json:"items,omitempty"`
}

type InventoryReceiptItemResponse struct {
	ID                 int          `json:"id"`
	InventoryReceiptID int          `json:"inventory_receipt_id"`
	ProductID          int          `json:"product_id"`
	Quantity           int          `json:"quantity"`
	UnitCost           *float64     `json:"unit_cost"`
	Notes              *string      `json:"notes"`
	CreatedAt          time.Time    `json:"created_at"`
	UpdatedAt          time.Time    `json:"updated_at"`
	Product            *ProductInfo `json:"product,omitempty"`
}

type GetAllInventoryReceiptsResponse struct {
	InventoryReceipts []InventoryReceiptResponse `json:"inventory_receipts"`
}
