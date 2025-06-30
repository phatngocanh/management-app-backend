package entity

import "time"

type UnitOfMeasure struct {
	ID          int       `db:"id"`
	Name        string    `db:"name"`        // Tên đơn vị (VD: Thùng, Cái, ML)
	Code        string    `db:"code"`        // Mã đơn vị (VD: THUNG, CAI, ML)
	Description string    `db:"description"` // Mô tả đơn vị
	CreatedAt   time.Time `db:"created_at"`  // Thời gian tạo
}

type unitOfMeasureCode struct {
	THUNG string
	CAI   string
	ML    string
	KG    string
	M     string
}

var UnitOfMeasureCode unitOfMeasureCode = unitOfMeasureCode{
	THUNG: "THUNG",
	CAI:   "CAI",
	ML:    "ML",
	KG:    "KG",
	M:     "M",
}
