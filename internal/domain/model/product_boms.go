package model

// Component for BOM - represents one component needed
type BomComponent struct {
	ComponentProductID int     `json:"component_product_id" binding:"required"` // ID sản phẩm nguyên liệu
	Quantity           float64 `json:"quantity" binding:"required"`             // Số lượng nguyên liệu cần thiết
}

// Component with full product info for response
type BomComponentResponse struct {
	ID                 int             `json:"id"`                          // ID của BOM entry
	ComponentProductID int             `json:"component_product_id"`        // ID sản phẩm nguyên liệu
	Quantity           float64         `json:"quantity"`                    // Số lượng nguyên liệu cần thiết
	ComponentProduct   *ProductBomInfo `json:"component_product,omitempty"` // Thông tin sản phẩm nguyên liệu
}

// Create BOM with multiple components
type CreateProductBomRequest struct {
	ParentProductID int            `json:"parent_product_id" binding:"required"` // ID sản phẩm thành phẩm
	Components      []BomComponent `json:"components" binding:"required,dive"`   // Danh sách nguyên liệu cần thiết
}

// Update entire BOM (replace all components)
type UpdateProductBomRequest struct {
	ParentProductID int            `json:"parent_product_id" binding:"required"` // ID sản phẩm thành phẩm
	Components      []BomComponent `json:"components" binding:"required,dive"`   // Danh sách nguyên liệu cần thiết mới
}

// BOM response - one product with all its components
type ProductBomResponse struct {
	ParentProductID int                    `json:"parent_product_id"`        // ID sản phẩm thành phẩm
	ParentProduct   *ProductBomInfo        `json:"parent_product,omitempty"` // Thông tin sản phẩm thành phẩm
	Components      []BomComponentResponse `json:"components"`               // Danh sách nguyên liệu
	TotalComponents int                    `json:"total_components"`         // Tổng số loại nguyên liệu
}

type GetAllProductBomsResponse struct {
	Boms []ProductBomResponse `json:"boms"`
}

type GetOneProductBomResponse struct {
	Bom ProductBomResponse `json:"bom"`
}

type ProductBomInfo struct {
	ID           int     `json:"id"`            // ID sản phẩm
	Name         string  `json:"name"`          // Tên sản phẩm
	Cost         float64 `json:"cost"`          // Giá vốn
	UnitCode     string  `json:"unit_code"`     // Mã đơn vị tính (VD: "thung", "cai", "ML")
	CategoryCode string  `json:"category_code"` // Mã danh mục (VD: "hoa-chat", "nhan")
}

// Material requirement calculation request
type CalculateMaterialRequirementsRequest struct {
	ParentProductID int     `json:"parent_product_id" binding:"required"` // ID sản phẩm thành phẩm
	Quantity        float64 `json:"quantity" binding:"required,gt=0"`     // Số lượng sản phẩm cần sản xuất
}

// Raw material requirement
type MaterialRequirement struct {
	ProductID        int             `json:"product_id"`        // ID nguyên liệu
	Product          *ProductBomInfo `json:"product"`           // Thông tin nguyên liệu
	RequiredQuantity float64         `json:"required_quantity"` // Tổng số lượng cần thiết
}

// Material requirements calculation response
type MaterialRequirementsResponse struct {
	ParentProductID      int                   `json:"parent_product_id"`     // ID sản phẩm thành phẩm
	ParentProduct        *ProductBomInfo       `json:"parent_product"`        // Thông tin sản phẩm thành phẩm
	RequestedQuantity    float64               `json:"requested_quantity"`    // Số lượng sản phẩm được yêu cầu
	MaterialRequirements []MaterialRequirement `json:"material_requirements"` // Danh sách nguyên liệu và số lượng cần thiết
	TotalMaterials       int                   `json:"total_materials"`       // Tổng số loại nguyên liệu
}
