package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	httpcommon "github.com/pna/management-app-backend/internal/domain/http_common"
	"github.com/pna/management-app-backend/internal/domain/model"
	"github.com/pna/management-app-backend/internal/service"
	"github.com/pna/management-app-backend/internal/utils/error_utils"
	"github.com/pna/management-app-backend/internal/utils/validation"
)

type UnitOfMeasureHandler struct {
	unitService service.UnitOfMeasureService
}

func NewUnitOfMeasureHandler(unitService service.UnitOfMeasureService) *UnitOfMeasureHandler {
	return &UnitOfMeasureHandler{
		unitService: unitService,
	}
}

// @Summary Create Unit of Measure
// @Description Create a new unit of measure
// @Tags Units
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization: Bearer"
// @Param request body model.CreateUnitOfMeasureRequest true "Unit information"
// @Success 201 {object} httpcommon.HttpResponse[model.UnitOfMeasureResponse]
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
// @Router /units [post]
func (h *UnitOfMeasureHandler) Create(ctx *gin.Context) {
	var request model.CreateUnitOfMeasureRequest
	if err := validation.BindJsonAndValidate(ctx, &request); err != nil {
		return
	}

	response, errCode := h.unitService.Create(ctx, request)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}

	ctx.JSON(http.StatusCreated, httpcommon.NewSuccessResponse(response))
}

// @Summary Update Unit of Measure
// @Description Update an existing unit of measure
// @Tags Units
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization: Bearer"
// @Param request body model.UpdateUnitOfMeasureRequest true "Updated unit information"
// @Success 200 {object} httpcommon.HttpResponse[model.UnitOfMeasureResponse]
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 404 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
// @Router /units [put]
func (h *UnitOfMeasureHandler) Update(ctx *gin.Context) {
	var request model.UpdateUnitOfMeasureRequest
	if err := validation.BindJsonAndValidate(ctx, &request); err != nil {
		return
	}

	response, errCode := h.unitService.Update(ctx, request)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}

	ctx.JSON(http.StatusOK, httpcommon.NewSuccessResponse(response))
}

// @Summary Get All Units of Measure
// @Description Retrieve all units of measure
// @Tags Units
// @Produce json
// @Param Authorization header string true "Authorization: Bearer"
// @Success 200 {object} httpcommon.HttpResponse[model.GetAllUnitsOfMeasureResponse]
// @Failure 500 {object} httpcommon.HttpResponse[any]
// @Router /units [get]
func (h *UnitOfMeasureHandler) GetAll(ctx *gin.Context) {
	response, errCode := h.unitService.GetAll(ctx)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}

	ctx.JSON(http.StatusOK, httpcommon.NewSuccessResponse(response))
}

// @Summary Get Unit of Measure by ID
// @Description Retrieve a unit of measure by its ID
// @Tags Units
// @Produce json
// @Param Authorization header string true "Authorization: Bearer"
// @Param id path int true "Unit ID"
// @Success 200 {object} httpcommon.HttpResponse[model.GetOneUnitOfMeasureResponse]
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 404 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
// @Router /units/{unitId} [get]
func (h *UnitOfMeasureHandler) GetOne(ctx *gin.Context) {
	idStr := ctx.Param("unitId")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(error_utils.ErrorCode.BAD_REQUEST, "id")
		ctx.JSON(statusCode, errResponse)
		return
	}

	response, errCode := h.unitService.GetOne(ctx, id)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}

	ctx.JSON(http.StatusOK, httpcommon.NewSuccessResponse(response))
}

// @Summary Get Unit of Measure by Code
// @Description Retrieve a unit of measure by its code
// @Tags Units
// @Produce json
// @Param Authorization header string true "Authorization: Bearer"
// @Param code path string true "Unit Code"
// @Success 200 {object} httpcommon.HttpResponse[model.GetOneUnitOfMeasureResponse]
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 404 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
// @Router /units/code/{code} [get]
func (h *UnitOfMeasureHandler) GetByCode(ctx *gin.Context) {
	code := ctx.Param("code")
	if code == "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(error_utils.ErrorCode.BAD_REQUEST, "code")
		ctx.JSON(statusCode, errResponse)
		return
	}

	response, errCode := h.unitService.GetByCode(ctx, code)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}

	ctx.JSON(http.StatusOK, httpcommon.NewSuccessResponse(response))
}
