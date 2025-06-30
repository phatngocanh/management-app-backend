package service

import (
	"github.com/gin-gonic/gin"
	"github.com/pna/management-app-backend/internal/domain/model"
)

type InventoryHistoryService interface {
	GetAll(ctx *gin.Context, productID int) (*model.GetAllInventoryHistoriesResponse, string)
	Create(ctx *gin.Context, request model.CreateInventoryHistoryRequest) (*model.InventoryHistoryResponse, string)
}
