package service

import (
	"github.com/gin-gonic/gin"
	"github.com/pna/management-app-backend/internal/domain/model"
)

type ProductImageService interface {
	Create(ctx *gin.Context, request model.CreateProductImageRequest) (*model.ProductImageResponse, string)
	Delete(ctx *gin.Context, id int) string
	GenerateSignedUploadURL(ctx *gin.Context, productID int, fileName string, contentType string, prefix string) (model.GenerateProductImageSignedUploadURLResponse, string)
	GetByProductID(ctx *gin.Context, productID int) ([]model.ProductImageResponse, string)
}
