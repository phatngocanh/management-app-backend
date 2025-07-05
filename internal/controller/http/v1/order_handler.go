package v1

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pna/management-app-backend/internal/controller/http/middleware"
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

	userID := middleware.GetUserIdHelper(ctx)

	response, errorCode := h.orderService.CreateOrder(ctx, request, userID)
	if errorCode != "" {
		// Check for detailed error message from service
		detailedMessage, exists := ctx.Get("detailed_error_message")
		if exists && errorCode == error_utils.ErrorCode.INVENTORY_QUANTITY_EXCEEDED {
			// Use detailed message for inventory shortage
			statusCode := http.StatusBadRequest
			errResponse := httpcommon.NewErrorResponse(httpcommon.Error{
				Message: detailedMessage.(string),
				Field:   "",
				Code:    errorCode,
			})
			ctx.JSON(statusCode, errResponse)
		} else {
			statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errorCode, "")
			ctx.JSON(statusCode, errResponse)
		}
		return
	}

	ctx.JSON(http.StatusCreated, httpcommon.NewSuccessResponse(response))
}

// @Summary Get One Order
// @Description Get one order by id
// @Tags Orders
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization: Bearer"
// @Param id path int true "Order ID"
// @Success 200 {object} httpcommon.HttpResponse[model.OrderResponse]
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
// @Router /orders/{orderId} [get]
func (h *OrderHandler) GetOneOrder(ctx *gin.Context) {
	id := ctx.Param("orderId")
	orderID, err := strconv.Atoi(id)

	if err != nil {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(error_utils.ErrorCode.BAD_REQUEST, "customerId")
		ctx.JSON(statusCode, errResponse)
		return
	}

	response, errCode := h.orderService.GetOneOrder(ctx, orderID)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}

	ctx.JSON(http.StatusOK, httpcommon.NewSuccessResponse(&response))
}

// @Summary Update Order
// @Description Update an existing order
// @Tags Orders
// @Accept json
// @Produce json
// @Param  Authorization header string true "Authorization: Bearer"
// @Param request body model.UpdateOrderRequest true "Updated order information"
// @Success 200 {object} httpcommon.HttpResponse[any]
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 404 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
// @Router /orders/{orderId} [put]
func (h *OrderHandler) Update(ctx *gin.Context) {
	orderID, err := strconv.Atoi(ctx.Param("orderId"))
	if err != nil {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(error_utils.ErrorCode.BAD_REQUEST, "orderId")
		ctx.JSON(statusCode, errResponse)
		return
	}

	var request model.UpdateOrderRequest
	if err := validation.BindJsonAndValidate(ctx, &request); err != nil {
		return
	}

	request.ID = orderID
	errCode := h.orderService.Update(ctx, request)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}

	ctx.JSON(http.StatusOK, httpcommon.NewSuccessResponse[any](nil))
}

// @Summary Get All Orders
// @Description Retrieve all orders with optional filters and sorting
// @Tags Orders
// @Produce json
// @Param  Authorization header string true "Authorization: Bearer"
// @Param customer_id query int false "Filter by customer ID"
// @Param delivery_statuses query string false "Filter by delivery statuses (comma-separated, e.g., PENDING,DELIVERED)"
// @Param sort_by query string false "Sort by: order_date_asc, order_date_desc (default: id DESC)"
// @Success 200 {object} httpcommon.HttpResponse[model.GetAllOrdersResponse]
// @Failure 500 {object} httpcommon.HttpResponse[any]
// @Router /orders [get]
func (h *OrderHandler) GetAll(ctx *gin.Context) {
	// Get query parameters
	customerIDStr := ctx.Query("customer_id")
	sortBy := ctx.Query("sort_by")
	fromDateStr := ctx.Query("from_date")
	toDateStr := ctx.Query("to_date")

	// Parse customer ID if provided
	customerID := 0
	if customerIDStr != "" {
		if id, err := strconv.Atoi(customerIDStr); err == nil {
			customerID = id
		}
	}

	// Parse date filters
	var fromDate *time.Time
	var toDate *time.Time

	if fromDateStr != "" {
		if parsedDate, err := time.Parse("2006-01-02", fromDateStr); err == nil {
			fromDate = &parsedDate
		} else {
			statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(error_utils.ErrorCode.BAD_REQUEST, "from_date format should be YYYY-MM-DD")
			ctx.JSON(statusCode, errResponse)
			return
		}
	}

	if toDateStr != "" {
		if parsedDate, err := time.Parse("2006-01-02", toDateStr); err == nil {
			toDate = &parsedDate
		} else {
			statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(error_utils.ErrorCode.BAD_REQUEST, "to_date format should be YYYY-MM-DD")
			ctx.JSON(statusCode, errResponse)
			return
		}
	}

	response, errCode := h.orderService.GetAll(ctx, 0, customerID, sortBy, fromDate, toDate)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}

	ctx.JSON(http.StatusOK, httpcommon.NewSuccessResponse(&response))
}
