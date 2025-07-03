package service

import (
	"github.com/gin-gonic/gin"
	"github.com/pna/management-app-backend/internal/domain/model"
)

type OrderService interface {
	CreateOrder(ctx *gin.Context, orderRequest model.CreateOrderRequest, userId int) (*model.OrderResponse, string)
	GetOneOrder(ctx *gin.Context, orderID int) (model.GetOneOrderResponse, string)
}
