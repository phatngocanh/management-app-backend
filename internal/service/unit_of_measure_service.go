package service

import (
	"github.com/gin-gonic/gin"
	"github.com/pna/management-app-backend/internal/domain/model"
)

type UnitOfMeasureService interface {
	Create(ctx *gin.Context, request model.CreateUnitOfMeasureRequest) (*model.UnitOfMeasureResponse, string)
	Update(ctx *gin.Context, request model.UpdateUnitOfMeasureRequest) (*model.UnitOfMeasureResponse, string)
	GetAll(ctx *gin.Context) (*model.GetAllUnitsOfMeasureResponse, string)
	GetOne(ctx *gin.Context, id int) (*model.GetOneUnitOfMeasureResponse, string)
	GetByCode(ctx *gin.Context, code string) (*model.GetOneUnitOfMeasureResponse, string)
}
