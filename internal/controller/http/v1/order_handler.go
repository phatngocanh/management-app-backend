package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	httpcommon "github.com/pna/management-app-backend/internal/domain/http_common"
	"github.com/pna/management-app-backend/internal/domain/model"
	"github.com/pna/management-app-backend/internal/service"
	"github.com/pna/management-app-backend/internal/utils/error_utils"
	"github.com/pna/management-app-backend/internal/utils/validation"
)

type OrderHandler struct {
	orderService service.OrderService
}

func NewOrderHandler(orderService service.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

// @Summary Create Order
// @Description Create a new order
// @Tags Orders
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization: Bearer"
// @Param request body model.CreateOrderRequest true "Order creation information"
// @Success 201 {object} httpcommon.HttpResponse[model.OrderResponse]
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 401 {object} httpcommon.HttpResponse[any]
// @Failure 404 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
// @Router /orders [post]
func (h *OrderHandler) CreateOrder(ctx *gin.Context) {
	var request model.CreateOrderRequest
	if err := validation.BindJsonAndValidate(ctx, &request); err != nil {
		return
	}

	response, errorCode := h.orderService.CreateOrder(ctx, request)
	if errorCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errorCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}

	ctx.JSON(http.StatusCreated, httpcommon.NewSuccessResponse(response))
}
