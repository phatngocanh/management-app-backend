package model

type CreateProductImageRequest struct {
	ProductID int    `json:"product_id" binding:"required"` // ID sản phẩm
	ImageKey  string `json:"image_key" binding:"required"`  // S3 object key
}

type UpdateProductImageRequest struct {
	ID        int    `json:"id" binding:"required"`
	ProductID int    `json:"product_id" binding:"required"` // ID sản phẩm
	ImageKey  string `json:"image_key" binding:"required"`  // S3 object key
}

type ProductImageResponse struct {
	ID        int    `json:"id"`
	ProductID int    `json:"product_id"` // ID sản phẩm
	ImageURL  string `json:"image_url"`  // URL hình ảnh signed (generated on demand)
	ImageKey  string `json:"image_key"`  // S3 object key
}

type GenerateProductImageSignedUploadURLResponse struct {
	SignedURL string `json:"signed_url"` // URL để upload
	ImageKey  string `json:"image_key"`  // S3 object key
	ImageID   int    `json:"image_id"`   // ID của hình ảnh được tạo
}

type GetAllProductImagesResponse struct {
	Images []ProductImageResponse `json:"images"`
}

type GetOneProductImageResponse struct {
	Image ProductImageResponse `json:"image"`
}

type GetProductImagesResponse struct {
	Images []ProductImageResponse `json:"images"`
}
