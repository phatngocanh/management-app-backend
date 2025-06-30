package model

import "time"

type InventoryReceiptItemRequest struct {
	ProductID int      `json:"product_id" binding:"required"`
	Quantity  int      `json:"quantity" binding:"required"`
	UnitCost  *float64 `json:"unit_cost"`
	Notes     *string  `json:"notes"`
}

type CreateInventoryReceiptRequest struct {
	UserID      int                           `json:"user_id" binding:"required"`
	ReceiptDate time.Time                     `json:"receipt_date"`
	Notes       *string                       `json:"notes"`
	Items       []InventoryReceiptItemRequest `json:"items" binding:"required,dive"`
}

type InventoryReceiptItemResponse struct {
	ID                 int       `json:"id"`
	InventoryReceiptID int       `json:"inventory_receipt_id"`
	ProductID          int       `json:"product_id"`
	Quantity           int       `json:"quantity"`
	UnitCost           *float64  `json:"unit_cost"`
	Notes              *string   `json:"notes"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type InventoryReceiptResponse struct {
	ID          int                            `json:"id"`
	Code        string                         `json:"code"`
	UserID      int                            `json:"user_id"`
	ReceiptDate time.Time                      `json:"receipt_date"`
	Notes       *string                        `json:"notes"`
	TotalItems  int                            `json:"total_items"`
	CreatedAt   time.Time                      `json:"created_at"`
	UpdatedAt   time.Time                      `json:"updated_at"`
	Items       []InventoryReceiptItemResponse `json:"items,omitempty"`
}

type GetAllInventoryReceiptsResponse struct {
	InventoryReceipts []InventoryReceiptResponse `json:"inventory_receipts"`
}

type GetOneInventoryReceiptResponse struct {
	InventoryReceipt InventoryReceiptResponse `json:"inventory_receipt"`
}
