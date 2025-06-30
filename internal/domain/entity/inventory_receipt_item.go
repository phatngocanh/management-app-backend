package entity

import "time"

type InventoryReceiptItem struct {
	ID                 int       `db:"id"`
	InventoryReceiptID int       `db:"inventory_receipt_id"`
	ProductID          int       `db:"product_id"`
	Quantity           int       `db:"quantity"`
	UnitCost           *float64  `db:"unit_cost"`
	Notes              *string   `db:"notes"`
	CreatedAt          time.Time `db:"created_at"`
	UpdatedAt          time.Time `db:"updated_at"`
}
