package model

type CreateUnitOfMeasureRequest struct {
	Name        string `json:"name" binding:"required"` // Tên đơn vị
	Code        string `json:"code" binding:"required"` // Mã đơn vị
	Description string `json:"description"`             // Mô tả đơn vị
}

type UpdateUnitOfMeasureRequest struct {
	ID          int    `json:"id" binding:"required"`
	Name        string `json:"name" binding:"required"` // Tên đơn vị
	Code        string `json:"code" binding:"required"` // Mã đơn vị
	Description string `json:"description"`             // Mô tả đơn vị
}

type UnitOfMeasureResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`        // Tên đơn vị
	Code        string `json:"code"`        // Mã đơn vị
	Description string `json:"description"` // Mô tả đơn vị
}

type GetAllUnitsOfMeasureResponse struct {
	Units []UnitOfMeasureResponse `json:"units"`
}

type GetOneUnitOfMeasureResponse struct {
	Unit UnitOfMeasureResponse `json:"unit"`
}
