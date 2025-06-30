package entity

import "time"

type InventoryReceipt struct {
	ID          int       `db:"id"`
	Code        string    `db:"code"`
	UserID      int       `db:"user_id"`
	ReceiptDate time.Time `db:"receipt_date"`
	Notes       *string   `db:"notes"`
	TotalItems  int       `db:"total_items"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
