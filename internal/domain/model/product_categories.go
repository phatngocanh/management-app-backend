package model

type CreateProductCategoryRequest struct {
	Name        string `json:"name" binding:"required"` // Tên danh mục
	Code        string `json:"code" binding:"required"` // Mã danh mục
	Description string `json:"description"`             // Mô tả danh mục
}

type UpdateProductCategoryRequest struct {
	ID          int    `json:"id" binding:"required"`
	Name        string `json:"name" binding:"required"` // Tên danh mục
	Code        string `json:"code" binding:"required"` // Mã danh mục
	Description string `json:"description"`             // Mô tả danh mục
}

type ProductCategoryResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`        // Tên danh mục
	Code        string `json:"code"`        // Mã danh mục
	Description string `json:"description"` // Mô tả danh mục
}

type GetAllProductCategoriesResponse struct {
	Categories []ProductCategoryResponse `json:"categories"`
}

type GetOneProductCategoryResponse struct {
	Category ProductCategoryResponse `json:"category"`
}
