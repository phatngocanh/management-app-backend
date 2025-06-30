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

type InventoryReceiptHandler struct {
	inventoryReceiptService service.InventoryReceiptService
}

func NewInventoryReceiptHandler(inventoryReceiptService service.InventoryReceiptService) *InventoryReceiptHandler {
	return &InventoryReceiptHandler{
		inventoryReceiptService: inventoryReceiptService,
	}
}

// @Summary Create Inventory Receipt
// @Description Create a new inventory receipt with items
// @Tags Inventory Receipts
// @Accept json
// @Produce json
// @Param  Authorization header string true "Authorization: Bearer"
// @Param request body model.CreateInventoryReceiptRequest true "Inventory receipt information"
// @Success 201 {object} httpcommon.HttpResponse[model.InventoryReceiptResponse]
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 404 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
// @Router /inventory-receipts [post]
func (h *InventoryReceiptHandler) Create(ctx *gin.Context) {
	var request model.CreateInventoryReceiptRequest
	if err := validation.BindJsonAndValidate(ctx, &request); err != nil {
		return
	}

	response, errCode := h.inventoryReceiptService.Create(ctx, request)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}

	ctx.JSON(http.StatusCreated, httpcommon.NewSuccessResponse(response))
}

// @Summary Get All Inventory Receipts
// @Description Retrieve all inventory receipts (general data only)
// @Tags Inventory Receipts
// @Produce json
// @Param  Authorization header string true "Authorization: Bearer"
// @Success 200 {object} httpcommon.HttpResponse[model.GetAllInventoryReceiptsResponse]
// @Failure 500 {object} httpcommon.HttpResponse[any]
// @Router /inventory-receipts [get]
func (h *InventoryReceiptHandler) GetAll(ctx *gin.Context) {
	response, errCode := h.inventoryReceiptService.GetAll(ctx)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}

	ctx.JSON(http.StatusOK, httpcommon.NewSuccessResponse(response))
}

// @Summary Get Inventory Receipt by Code
// @Description Retrieve an inventory receipt by its code with all related items
// @Tags Inventory Receipts
// @Produce json
// @Param  Authorization header string true "Authorization: Bearer"
// @Param code path string true "Inventory Receipt Code"
// @Success 200 {object} httpcommon.HttpResponse[model.GetOneInventoryReceiptResponse]
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 404 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
// @Router /inventory-receipts/{receiptCode} [get]
func (h *InventoryReceiptHandler) GetOne(ctx *gin.Context) {
	code := ctx.Param("receiptCode")
	if code == "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(error_utils.ErrorCode.BAD_REQUEST, "receiptCode")
		ctx.JSON(statusCode, errResponse)
		return
	}

	response, errCode := h.inventoryReceiptService.GetByCode(ctx, code)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}

	ctx.JSON(http.StatusOK, httpcommon.NewSuccessResponse(response))
}
