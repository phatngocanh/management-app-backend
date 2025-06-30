package service

import (
	"github.com/gin-gonic/gin"
	"github.com/pna/management-app-backend/internal/domain/model"
)

type ProductCategoryService interface {
	Create(ctx *gin.Context, request model.CreateProductCategoryRequest) (*model.ProductCategoryResponse, string)
	Update(ctx *gin.Context, request model.UpdateProductCategoryRequest) (*model.ProductCategoryResponse, string)
	GetAll(ctx *gin.Context) (*model.GetAllProductCategoriesResponse, string)
	GetOne(ctx *gin.Context, id int) (*model.GetOneProductCategoryResponse, string)
	GetByCode(ctx *gin.Context, code string) (*model.GetOneProductCategoryResponse, string)
}
