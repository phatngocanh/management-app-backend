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

type ProductCategoryHandler struct {
	categoryService service.ProductCategoryService
}

func NewProductCategoryHandler(categoryService service.ProductCategoryService) *ProductCategoryHandler {
	return &ProductCategoryHandler{
		categoryService: categoryService,
	}
}

// @Summary Create Product Category
// @Description Create a new product category
// @Tags Categories
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization: Bearer"
// @Param request body model.CreateProductCategoryRequest true "Category information"
// @Success 201 {object} httpcommon.HttpResponse[model.ProductCategoryResponse]
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
// @Router /categories [post]
func (h *ProductCategoryHandler) Create(ctx *gin.Context) {
	var request model.CreateProductCategoryRequest
	if err := validation.BindJsonAndValidate(ctx, &request); err != nil {
		return
	}

	response, errCode := h.categoryService.Create(ctx, request)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}

	ctx.JSON(http.StatusCreated, httpcommon.NewSuccessResponse(response))
}

// @Summary Update Product Category
// @Description Update an existing product category
// @Tags Categories
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization: Bearer"
// @Param request body model.UpdateProductCategoryRequest true "Updated category information"
// @Success 200 {object} httpcommon.HttpResponse[model.ProductCategoryResponse]
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 404 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
// @Router /categories [put]
func (h *ProductCategoryHandler) Update(ctx *gin.Context) {
	var request model.UpdateProductCategoryRequest
	if err := validation.BindJsonAndValidate(ctx, &request); err != nil {
		return
	}

	response, errCode := h.categoryService.Update(ctx, request)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}

	ctx.JSON(http.StatusOK, httpcommon.NewSuccessResponse(response))
}

// @Summary Get All Product Categories
// @Description Retrieve all product categories
// @Tags Categories
// @Produce json
// @Param Authorization header string true "Authorization: Bearer"
// @Success 200 {object} httpcommon.HttpResponse[model.GetAllProductCategoriesResponse]
// @Failure 500 {object} httpcommon.HttpResponse[any]
// @Router /categories [get]
func (h *ProductCategoryHandler) GetAll(ctx *gin.Context) {
	response, errCode := h.categoryService.GetAll(ctx)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}

	ctx.JSON(http.StatusOK, httpcommon.NewSuccessResponse(response))
}

// @Summary Get Product Category by ID
// @Description Retrieve a product category by its ID
// @Tags Categories
// @Produce json
// @Param Authorization header string true "Authorization: Bearer"
// @Param id path int true "Category ID"
// @Success 200 {object} httpcommon.HttpResponse[model.GetOneProductCategoryResponse]
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 404 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
// @Router /categories/{categoryId} [get]
func (h *ProductCategoryHandler) GetOne(ctx *gin.Context) {
	idStr := ctx.Param("categoryId")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(error_utils.ErrorCode.BAD_REQUEST, "id")
		ctx.JSON(statusCode, errResponse)
		return
	}

	response, errCode := h.categoryService.GetOne(ctx, id)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}

	ctx.JSON(http.StatusOK, httpcommon.NewSuccessResponse(response))
}

// @Summary Get Product Category by Code
// @Description Retrieve a product category by its code
// @Tags Categories
// @Produce json
// @Param Authorization header string true "Authorization: Bearer"
// @Param code path string true "Category Code"
// @Success 200 {object} httpcommon.HttpResponse[model.GetOneProductCategoryResponse]
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 404 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
// @Router /categories/code/{code} [get]
func (h *ProductCategoryHandler) GetByCode(ctx *gin.Context) {
	code := ctx.Param("code")
	if code == "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(error_utils.ErrorCode.BAD_REQUEST, "code")
		ctx.JSON(statusCode, errResponse)
		return
	}

	response, errCode := h.categoryService.GetByCode(ctx, code)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}

	ctx.JSON(http.StatusOK, httpcommon.NewSuccessResponse(response))
}
