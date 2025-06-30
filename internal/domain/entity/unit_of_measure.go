package entity

import "time"

type UnitOfMeasure struct {
	ID          int       `db:"id"`
	Name        string    `db:"name"`        // Tên đơn vị (VD: Thùng, Cái, ML)
	Code        string    `db:"code"`        // Mã đơn vị (VD: THUNG, CAI, ML)
	Description string    `db:"description"` // Mô tả đơn vị
	CreatedAt   time.Time `db:"created_at"`  // Thời gian tạo
}
