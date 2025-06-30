package service

import (
	"github.com/gin-gonic/gin"
	"github.com/pna/management-app-backend/internal/domain/model"
)

type ProductImageService interface {
	Create(ctx *gin.Context, request model.CreateProductImageRequest) (*model.ProductImageResponse, string)
	Update(ctx *gin.Context, request model.UpdateProductImageRequest) (*model.ProductImageResponse, string)
	GetAll(ctx *gin.Context) (*model.GetAllProductImagesResponse, string)
	GetOne(ctx *gin.Context, id int) (*model.GetOneProductImageResponse, string)
	GetByProductID(ctx *gin.Context, productID int) (*model.GetProductImagesResponse, string)
	Delete(ctx *gin.Context, id int) string
	GenerateSignedUploadURL(ctx *gin.Context, productID int, fileName string, contentType string) (model.GenerateProductImageSignedUploadURLResponse, string)
}
