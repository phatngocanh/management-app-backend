package service

import (
	"github.com/gin-gonic/gin"
	"github.com/pna/management-app-backend/internal/domain/model"
)

type InventoryReceiptService interface {
	Create(ctx *gin.Context, request model.CreateInventoryReceiptRequest) (*model.InventoryReceiptResponse, string)
	GetAll(ctx *gin.Context) (*model.GetAllInventoryReceiptsResponse, string)
	GetOne(ctx *gin.Context, id int) (*model.GetOneInventoryReceiptResponse, string)
}
