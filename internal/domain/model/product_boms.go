package model

type CreateProductBomRequest struct {
	ParentProductID    int     `json:"parent_product_id" binding:"required"`    // ID sản phẩm thành phẩm
	ComponentProductID int     `json:"component_product_id" binding:"required"` // ID sản phẩm nguyên liệu
	Quantity           float64 `json:"quantity" binding:"required"`             // Số lượng nguyên liệu cần thiết
}

type UpdateProductBomRequest struct {
	ID                 int     `json:"id" binding:"required"`
	ParentProductID    int     `json:"parent_product_id" binding:"required"`    // ID sản phẩm thành phẩm
	ComponentProductID int     `json:"component_product_id" binding:"required"` // ID sản phẩm nguyên liệu
	Quantity           float64 `json:"quantity" binding:"required"`             // Số lượng nguyên liệu cần thiết
}

type ProductBomResponse struct {
	ID                 int             `json:"id"`
	ParentProductID    int             `json:"parent_product_id"`           // ID sản phẩm thành phẩm
	ComponentProductID int             `json:"component_product_id"`        // ID sản phẩm nguyên liệu
	Quantity           float64         `json:"quantity"`                    // Số lượng nguyên liệu cần thiết
	ParentProduct      *ProductBomInfo `json:"parent_product,omitempty"`    // Thông tin sản phẩm thành phẩm
	ComponentProduct   *ProductBomInfo `json:"component_product,omitempty"` // Thông tin sản phẩm nguyên liệu
}

type GetAllProductBomsResponse struct {
	Boms []ProductBomResponse `json:"boms"`
}

type GetOneProductBomResponse struct {
	Bom ProductBomResponse `json:"bom"`
}

type GetProductBomsResponse struct {
	Boms []ProductBomResponse `json:"boms"`
}

type ProductBomInfo struct {
	ID            int    `json:"id"`             // ID sản phẩm
	Name          string `json:"name"`           // Tên sản phẩm
	Spec          int    `json:"spec"`           // Qui cách
	OriginalPrice int    `json:"original_price"` // Giá gốc
}
