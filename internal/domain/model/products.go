package model

type CreateProductRequest struct {
	Name          string `json:"name" binding:"required"`           // Tên sản phẩm
	Spec          int    `json:"spec"`                              // Quy cách
	OriginalPrice int    `json:"original_price" binding:"required"` // Giá gốc của sản phẩm (VND)
	CategoryID    *int   `json:"category_id"`                       // ID danh mục sản phẩm
	UnitID        *int   `json:"unit_id"`                           // ID đơn vị tính
	Description   string `json:"description"`                       // Mô tả chi tiết sản phẩm
	UnitCode      string `json:"unit_code"`                         // Mã đơn vị tính
}

type UpdateProductRequest struct {
	ID            int    `json:"id" binding:"required"`
	Name          string `json:"name" binding:"required"`           // Tên sản phẩm
	Spec          int    `json:"spec"`                              // Quy cách
	OriginalPrice int    `json:"original_price" binding:"required"` // Giá gốc của sản phẩm (VND)
	CategoryID    *int   `json:"category_id"`                       // ID danh mục sản phẩm
	UnitID        *int   `json:"unit_id"`                           // ID đơn vị tính
	Description   string `json:"description"`                       // Mô tả chi tiết sản phẩm
}

type ProductResponse struct {
	ID            int                      `json:"id"`
	Name          string                   `json:"name"`                // Tên sản phẩm
	Spec          int                      `json:"spec"`                // Quy cách
	OriginalPrice int                      `json:"original_price"`      // Giá gốc của sản phẩm (VND)
	CategoryID    *int                     `json:"category_id"`         // ID danh mục sản phẩm
	UnitID        *int                     `json:"unit_id"`             // ID đơn vị tính
	Description   string                   `json:"description"`         // Mô tả chi tiết sản phẩm
	Category      *ProductCategoryResponse `json:"category,omitempty"`  // Thông tin danh mục
	Unit          *UnitOfMeasureResponse   `json:"unit,omitempty"`      // Thông tin đơn vị tính
	Inventory     *InventoryInfo           `json:"inventory,omitempty"` // Thông tin tồn kho
}

type InventoryInfo struct {
	Quantity int    `json:"quantity"` // Số lượng tồn kho
	Version  string `json:"version"`  // Version để optimistic lock
}

type GetAllProductsResponse struct {
	Products []ProductResponse `json:"products"`
}

type GetOneProductResponse struct {
	Product ProductResponse `json:"product"`
}
