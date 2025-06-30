package service

import (
	"github.com/gin-gonic/gin"
	"github.com/pna/management-app-backend/internal/domain/model"
)

type ProductBomService interface {
	Create(ctx *gin.Context, request model.CreateProductBomRequest) (*model.ProductBomResponse, string)
	Update(ctx *gin.Context, request model.UpdateProductBomRequest) (*model.ProductBomResponse, string)
	GetAll(ctx *gin.Context) (*model.GetAllProductBomsResponse, string)
	GetOne(ctx *gin.Context, id int) (*model.GetOneProductBomResponse, string)
	GetByParentProductID(ctx *gin.Context, parentProductID int) (*model.GetProductBomsResponse, string)
	GetByComponentProductID(ctx *gin.Context, componentProductID int) (*model.GetProductBomsResponse, string)
	Delete(ctx *gin.Context, id int) string
}
