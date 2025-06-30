package service

import (
	"github.com/gin-gonic/gin"
	"github.com/pna/management-app-backend/internal/domain/model"
)

type ProductBomService interface {
	Create(ctx *gin.Context, request model.CreateProductBomRequest) (*model.ProductBomResponse, string)
	Update(ctx *gin.Context, request model.UpdateProductBomRequest) (*model.ProductBomResponse, string)
	GetAll(ctx *gin.Context) (*model.GetAllProductBomsResponse, string)
	GetByParentProductID(ctx *gin.Context, parentProductID int) (*model.GetOneProductBomResponse, string)
	GetByComponentProductID(ctx *gin.Context, componentProductID int) (*model.GetAllProductBomsResponse, string)
	DeleteByParentProductID(ctx *gin.Context, parentProductID int) string
	CalculateMaterialRequirements(ctx *gin.Context, request model.CalculateMaterialRequirementsRequest) (*model.MaterialRequirementsResponse, string)
}
