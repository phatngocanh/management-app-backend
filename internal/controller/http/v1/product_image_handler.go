package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	httpcommon "github.com/pna/management-app-backend/internal/domain/http_common"
	"github.com/pna/management-app-backend/internal/service"
	"github.com/pna/management-app-backend/internal/utils/env"
	"github.com/pna/management-app-backend/internal/utils/error_utils"
)

type ProductImageHandler struct {
	productImageService service.ProductImageService
}

func NewProductImageHandler(productImageService service.ProductImageService) *ProductImageHandler {
	return &ProductImageHandler{
		productImageService: productImageService,
	}
}

// @Summary Generate Signed Upload URL
// @Description Generate a signed URL for uploading an image to S3
// @Tags Product Images
// @Produce json
// @Param  Authorization header string true "Authorization: Bearer"
// @Param productId path int true "Product ID"
// @Param fileName query string true "File name"
// @Param contentType query string true "Content type (e.g., image/jpeg)"
// @Success 200 {object} httpcommon.HttpResponse[model.GenerateProductImageSignedUploadURLResponse]
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 401 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
// @Router /products/{productId}/images/upload-url [post]
func (h *ProductImageHandler) GenerateSignedUploadURL(ctx *gin.Context) {
	// Get product ID from path parameter
	productID, err := strconv.Atoi(ctx.Param("productId"))
	if err != nil {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(error_utils.ErrorCode.BAD_REQUEST, "productId")
		ctx.JSON(statusCode, errResponse)
		return
	}

	// Get query parameters
	fileName := ctx.Query("fileName")
	if fileName == "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(error_utils.ErrorCode.BAD_REQUEST, "fileName is required")
		ctx.JSON(statusCode, errResponse)
		return
	}

	contentType := ctx.Query("contentType")
	if contentType == "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(error_utils.ErrorCode.BAD_REQUEST, "contentType is required")
		ctx.JSON(statusCode, errResponse)
		return
	}

	prefix, err := env.GetEnv("AWS_S3_PRODUCT_IMAGES_PREFIX")
	if err != nil {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(error_utils.ErrorCode.INTERNAL_SERVER_ERROR, "AWS_S3_PRODUCT_IMAGES_PREFIX is not set")
		ctx.JSON(statusCode, errResponse)
		return
	}
	// Generate signed upload URL
	response, errCode := h.productImageService.GenerateSignedUploadURL(ctx, productID, fileName, contentType, prefix)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}

	ctx.JSON(http.StatusOK, httpcommon.NewSuccessResponse(&response))
}

// @Summary Delete Product Image
// @Description Delete a specific image from a product
// @Tags Product Images
// @Produce json
// @Param  Authorization header string true "Authorization: Bearer"
// @Param productId path int true "Product ID"
// @Param imageId path int true "Image ID"
// @Success 200 {object} httpcommon.HttpResponse[any]
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 401 {object} httpcommon.HttpResponse[any]
// @Failure 404 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
// @Router /products/{productId}/images/{imageId} [delete]
func (h *ProductImageHandler) DeleteImage(ctx *gin.Context) {
	// Get image ID from path parameter
	imageID, err := strconv.Atoi(ctx.Param("imageId"))
	if err != nil {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(error_utils.ErrorCode.BAD_REQUEST, "imageId")
		ctx.JSON(statusCode, errResponse)
		return
	}

	// Delete the image
	errCode := h.productImageService.Delete(ctx, imageID)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}

	ctx.JSON(http.StatusOK, httpcommon.NewSuccessResponse[any](nil))
}
