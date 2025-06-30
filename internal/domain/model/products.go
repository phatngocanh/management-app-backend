package model

type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required"` // Tên sản phẩm
	Cost        float64 `json:"cost" binding:"required"` // Giá vốn của sản phẩm (VND)
	CategoryID  *int    `json:"category_id"`             // ID danh mục sản phẩm
	UnitID      *int    `json:"unit_id"`                 // ID đơn vị tính
	Description string  `json:"description"`             // Mô tả chi tiết sản phẩm
}

type UpdateProductRequest struct {
	ID          int     `json:"id" binding:"required"`
	Name        string  `json:"name" binding:"required"` // Tên sản phẩm
	Cost        float64 `json:"cost"`                    // Giá vốn của sản phẩm (VND)
	CategoryID  *int    `json:"category_id"`             // ID danh mục sản phẩm
	UnitID      *int    `json:"unit_id"`                 // ID đơn vị tính
	Description string  `json:"description"`             // Mô tả chi tiết sản phẩm
}

type ProductResponse struct {
	ID          int                      `json:"id"`
	Name        string                   `json:"name"`                   // Tên sản phẩm
	Cost        float64                  `json:"cost"`                   // Giá vốn của sản phẩm (VND)
	CategoryID  *int                     `json:"category_id"`            // ID danh mục sản phẩm
	UnitID      *int                     `json:"unit_id"`                // ID đơn vị tính
	Description string                   `json:"description"`            // Mô tả chi tiết sản phẩm
	Category    *ProductCategoryResponse `json:"category,omitempty"`     // Thông tin danh mục
	Unit        *UnitOfMeasureResponse   `json:"unit,omitempty"`         // Thông tin đơn vị tính
	Inventory   *InventoryInfo           `json:"inventory,omitempty"`    // Thông tin tồn kho
	BOM         *ProductBOMInfo          `json:"bom,omitempty"`          // Thông tin công thức sản xuất (nếu có)
	UsedInBOMs  []ProductBOMUsage        `json:"used_in_boms,omitempty"` // Danh sách sản phẩm sử dụng sản phẩm này làm nguyên liệu
}

// BOM information when getting product details
type ProductBOMInfo struct {
	TotalComponents int                    `json:"total_components"` // Tổng số loại nguyên liệu
	Components      []BomComponentResponse `json:"components"`       // Danh sách nguyên liệu chi tiết
}

// Information about where this product is used as component
type ProductBOMUsage struct {
	ParentProductID   int     `json:"parent_product_id"`   // ID sản phẩm thành phẩm
	ParentProductName string  `json:"parent_product_name"` // Tên sản phẩm thành phẩm
	Quantity          float64 `json:"quantity"`            // Số lượng cần thiết
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
