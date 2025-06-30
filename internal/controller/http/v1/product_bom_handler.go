package v1

import (
	"github.com/pna/management-app-backend/internal/utils/validation"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	httpcommon "github.com/pna/management-app-backend/internal/domain/http_common"
	"github.com/pna/management-app-backend/internal/domain/model"
	"github.com/pna/management-app-backend/internal/service"
	"github.com/pna/management-app-backend/internal/utils/error_utils"
)

type ProductBomHandler struct {
	bomService service.ProductBomService
}

func NewProductBomHandler(bomService service.ProductBomService) *ProductBomHandler {
	return &ProductBomHandler{
		bomService: bomService,
	}
}

// @Summary Create Product BOM
// @Description Create a Bill of Materials (BOM) for a product with multiple components
// @Tags BOMs
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization: Bearer"
// @Param request body model.CreateProductBomRequest true "BOM creation information with parent product and components"
// @Success 201 {object} httpcommon.HttpResponse[model.ProductBomResponse]
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 404 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
// @Router /boms [post]
func (h *ProductBomHandler) CreateProductBom(ctx *gin.Context) {
	var request model.CreateProductBomRequest
	if err := validation.BindJsonAndValidate(ctx, &request); err != nil {
		return
	}

	result, errorCode := h.bomService.Create(ctx, request)
	if errorCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errorCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}

	ctx.JSON(http.StatusCreated, httpcommon.NewSuccessResponse(result))
}

// @Summary Update Product BOM
// @Description Update a complete Bill of Materials (BOM) for a product, replacing all existing components
// @Tags BOMs
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization: Bearer"
// @Param request body model.UpdateProductBomRequest true "BOM update information with parent product and new components"
// @Success 200 {object} httpcommon.HttpResponse[model.ProductBomResponse]
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 404 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
// @Router /boms [put]
func (h *ProductBomHandler) UpdateProductBom(ctx *gin.Context) {
	var request model.UpdateProductBomRequest
	if err := validation.BindJsonAndValidate(ctx, &request); err != nil {
		return
	}

	result, errorCode := h.bomService.Update(ctx, request)
	if errorCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errorCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}

	ctx.JSON(http.StatusOK, httpcommon.NewSuccessResponse(result))
}

// @Summary Get All Product BOMs
// @Description Retrieve all Bills of Materials (BOMs) in the system, grouped by parent product
// @Tags BOMs
// @Produce json
// @Param Authorization header string true "Authorization: Bearer"
// @Success 200 {object} httpcommon.HttpResponse[model.GetAllProductBomsResponse]
// @Failure 500 {object} httpcommon.HttpResponse[any]
// @Router /boms [get]
func (h *ProductBomHandler) GetAllProductBoms(ctx *gin.Context) {
	result, errorCode := h.bomService.GetAll(ctx)
	if errorCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errorCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}

	ctx.JSON(http.StatusOK, httpcommon.NewSuccessResponse(result))
}

// @Summary Get BOM by Parent Product ID
// @Description Retrieve the complete Bill of Materials (BOM) for a specific parent product
// @Tags BOMs
// @Produce json
// @Param Authorization header string true "Authorization: Bearer"
// @Param parentProductId path int true "Parent Product ID"
// @Success 200 {object} httpcommon.HttpResponse[model.GetOneProductBomResponse]
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 404 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
// @Router /boms/parent/{parentProductId} [get]
func (h *ProductBomHandler) GetProductBomByParentID(ctx *gin.Context) {
	parentProductIDParam := ctx.Param("parentProductId")
	parentProductID, err := strconv.Atoi(parentProductIDParam)
	if err != nil {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(error_utils.ErrorCode.BAD_REQUEST, "parentProductId")
		ctx.JSON(statusCode, errResponse)
		return
	}

	result, errorCode := h.bomService.GetByParentProductID(ctx, parentProductID)
	if errorCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errorCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}

	ctx.JSON(http.StatusOK, httpcommon.NewSuccessResponse(result))
}

// @Summary Get BOMs by Component Product ID
// @Description Find all Bills of Materials (BOMs) that use a specific product as a component/ingredient
// @Tags BOMs
// @Produce json
// @Param Authorization header string true "Authorization: Bearer"
// @Param componentProductId path int true "Component Product ID"
// @Success 200 {object} httpcommon.HttpResponse[model.GetAllProductBomsResponse]
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
// @Router /boms/component/{componentProductId} [get]
func (h *ProductBomHandler) GetProductBomsByComponentID(ctx *gin.Context) {
	componentProductIDParam := ctx.Param("componentProductId")
	componentProductID, err := strconv.Atoi(componentProductIDParam)
	if err != nil {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(error_utils.ErrorCode.BAD_REQUEST, "componentProductId")
		ctx.JSON(statusCode, errResponse)
		return
	}

	result, errorCode := h.bomService.GetByComponentProductID(ctx, componentProductID)
	if errorCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errorCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}

	ctx.JSON(http.StatusOK, httpcommon.NewSuccessResponse(result))
}

// @Summary Delete Product BOM
// @Description Delete the complete Bill of Materials (BOM) for a specific parent product
// @Tags BOMs
// @Produce json
// @Param Authorization header string true "Authorization: Bearer"
// @Param parentProductId path int true "Parent Product ID"
// @Success 200 {object} httpcommon.HttpResponse[string]
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 404 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
// @Router /boms/parent/{parentProductId} [delete]
func (h *ProductBomHandler) DeleteProductBom(ctx *gin.Context) {
	parentProductIDParam := ctx.Param("parentProductId")
	parentProductID, err := strconv.Atoi(parentProductIDParam)
	if err != nil {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(error_utils.ErrorCode.BAD_REQUEST, "parentProductId")
		ctx.JSON(statusCode, errResponse)
		return
	}

	errorCode := h.bomService.DeleteByParentProductID(ctx, parentProductID)
	if errorCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errorCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}

	message := "Xóa BOM thành công"
	ctx.JSON(http.StatusOK, httpcommon.NewSuccessResponse(&message))
}
