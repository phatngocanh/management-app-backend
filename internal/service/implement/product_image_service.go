package serviceimplement

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pna/management-app-backend/internal/bean"
	"github.com/pna/management-app-backend/internal/domain/entity"
	"github.com/pna/management-app-backend/internal/domain/model"
	"github.com/pna/management-app-backend/internal/repository"
	"github.com/pna/management-app-backend/internal/service"
	"github.com/pna/management-app-backend/internal/utils/error_utils"
	log "github.com/sirupsen/logrus"
)

type ProductImageService struct {
	imageRepository repository.ProductImageRepository
	unitOfWork      repository.UnitOfWork
	s3Service       bean.S3Service
}

func NewProductImageService(imageRepository repository.ProductImageRepository, unitOfWork repository.UnitOfWork, s3Service bean.S3Service) service.ProductImageService {
	return &ProductImageService{
		imageRepository: imageRepository,
		unitOfWork:      unitOfWork,
		s3Service:       s3Service,
	}
}

func (s *ProductImageService) GenerateSignedUploadURL(ctx *gin.Context, productID int, fileName string, contentType string) (model.GenerateProductImageSignedUploadURLResponse, string) {
	// Generate signed upload URL and S3 key with product images prefix
	signedURL, s3Key, err := s.s3Service.GenerateSignedUploadURLWithPrefix(ctx, fileName, contentType, "product-images/")
	if err != nil {
		log.Error("ProductImageService.GenerateSignedUploadURL Error generating signed upload URL: " + err.Error())
		return model.GenerateProductImageSignedUploadURLResponse{}, error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}

	// Create product image entity with S3 key (no signed URL yet since file hasn't been uploaded)
	productImage := &entity.ProductImage{
		ProductID: productID,
		ImageKey:  s3Key,
		IsPrimary: false, // Default to false, can be updated later
	}

	// Save to database
	err = s.imageRepository.CreateCommand(ctx, productImage, nil)
	if err != nil {
		log.Error("ProductImageService.GenerateSignedUploadURL Error saving to database: " + err.Error())
		return model.GenerateProductImageSignedUploadURLResponse{}, error_utils.ErrorCode.DB_DOWN
	}

	response := model.GenerateProductImageSignedUploadURLResponse{
		SignedURL: signedURL,
		ImageKey:  s3Key,
		ImageID:   productImage.ID,
	}

	return response, ""
}

func (s *ProductImageService) Create(ctx *gin.Context, request model.CreateProductImageRequest) (*model.ProductImageResponse, string) {
	// Create image entity
	image := &entity.ProductImage{
		ProductID: request.ProductID,
		ImageKey:  request.ImageKey,
		IsPrimary: request.IsPrimary,
	}

	// Save to database
	err := s.imageRepository.CreateCommand(ctx, image, nil)
	if err != nil {
		log.Error("ProductImageService.Create Error when create image: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	// Generate signed URL for response
	signedURL, err := s.s3Service.GenerateSignedDownloadURL(ctx, image.ImageKey, 20*time.Minute)
	if err != nil {
		log.Error("ProductImageService.Create Error generating signed URL: " + err.Error())
		signedURL = "" // Continue without signed URL
	}

	return &model.ProductImageResponse{
		ID:        image.ID,
		ProductID: image.ProductID,
		ImageURL:  signedURL,
		ImageKey:  image.ImageKey,
		IsPrimary: image.IsPrimary,
	}, ""
}

func (s *ProductImageService) Update(ctx *gin.Context, request model.UpdateProductImageRequest) (*model.ProductImageResponse, string) {
	// Check if image exists
	existingImage, err := s.imageRepository.GetOneByIDQuery(ctx, request.ID, nil)
	if err != nil {
		log.Error("ProductImageService.Update Error when get image: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	if existingImage == nil {
		return nil, error_utils.ErrorCode.NOT_FOUND
	}

	// Update image entity
	image := &entity.ProductImage{
		ID:        request.ID,
		ProductID: request.ProductID,
		ImageKey:  request.ImageKey,
		IsPrimary: request.IsPrimary,
	}

	// Save to database
	err = s.imageRepository.UpdateCommand(ctx, image, nil)
	if err != nil {
		log.Error("ProductImageService.Update Error when update image: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	// Generate signed URL for response
	signedURL, err := s.s3Service.GenerateSignedDownloadURL(ctx, image.ImageKey, 20*time.Minute)
	if err != nil {
		log.Error("ProductImageService.Update Error generating signed URL: " + err.Error())
		signedURL = "" // Continue without signed URL
	}

	return &model.ProductImageResponse{
		ID:        image.ID,
		ProductID: image.ProductID,
		ImageURL:  signedURL,
		ImageKey:  image.ImageKey,
		IsPrimary: image.IsPrimary,
	}, ""
}

func (s *ProductImageService) GetAll(ctx *gin.Context) (*model.GetAllProductImagesResponse, string) {
	// Get all images
	images, err := s.imageRepository.GetAllQuery(ctx, nil)
	if err != nil {
		log.Error("ProductImageService.GetAll Error when get images: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	// Convert to response models with signed URLs
	imageResponses := make([]model.ProductImageResponse, len(images))
	for i, image := range images {
		// Generate signed URL for each image
		signedURL, err := s.s3Service.GenerateSignedDownloadURL(ctx, image.ImageKey, 20*time.Minute)
		if err != nil {
			log.Error("ProductImageService.GetAll Error generating signed URL for image " + string(rune(image.ID)) + ": " + err.Error())
			signedURL = "" // Continue without signed URL
		}

		imageResponses[i] = model.ProductImageResponse{
			ID:        image.ID,
			ProductID: image.ProductID,
			ImageURL:  signedURL,
			ImageKey:  image.ImageKey,
			IsPrimary: image.IsPrimary,
		}
	}

	return &model.GetAllProductImagesResponse{
		Images: imageResponses,
	}, ""
}

func (s *ProductImageService) GetOne(ctx *gin.Context, id int) (*model.GetOneProductImageResponse, string) {
	// Get image by ID
	image, err := s.imageRepository.GetOneByIDQuery(ctx, id, nil)
	if err != nil {
		log.Error("ProductImageService.GetOne Error when get image: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	if image == nil {
		return nil, error_utils.ErrorCode.NOT_FOUND
	}

	// Generate signed URL
	signedURL, err := s.s3Service.GenerateSignedDownloadURL(ctx, image.ImageKey, 20*time.Minute)
	if err != nil {
		log.Error("ProductImageService.GetOne Error generating signed URL: " + err.Error())
		signedURL = "" // Continue without signed URL
	}

	return &model.GetOneProductImageResponse{
		Image: model.ProductImageResponse{
			ID:        image.ID,
			ProductID: image.ProductID,
			ImageURL:  signedURL,
			ImageKey:  image.ImageKey,
			IsPrimary: image.IsPrimary,
		},
	}, ""
}

func (s *ProductImageService) GetByProductID(ctx *gin.Context, productID int) (*model.GetProductImagesResponse, string) {
	// Get images by product ID
	images, err := s.imageRepository.GetByProductIDQuery(ctx, productID, nil)
	if err != nil {
		log.Error("ProductImageService.GetByProductID Error when get images: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	// Convert to response models with signed URLs
	imageResponses := make([]model.ProductImageResponse, len(images))
	for i, image := range images {
		// Generate signed URL for each image
		signedURL, err := s.s3Service.GenerateSignedDownloadURL(ctx, image.ImageKey, 20*time.Minute)
		if err != nil {
			log.Error("ProductImageService.GetByProductID Error generating signed URL for image " + string(rune(image.ID)) + ": " + err.Error())
			signedURL = "" // Continue without signed URL
		}

		imageResponses[i] = model.ProductImageResponse{
			ID:        image.ID,
			ProductID: image.ProductID,
			ImageURL:  signedURL,
			ImageKey:  image.ImageKey,
			IsPrimary: image.IsPrimary,
		}
	}

	return &model.GetProductImagesResponse{
		Images: imageResponses,
	}, ""
}

func (s *ProductImageService) Delete(ctx *gin.Context, id int) string {
	// Check if image exists
	existingImage, err := s.imageRepository.GetOneByIDQuery(ctx, id, nil)
	if err != nil {
		log.Error("ProductImageService.Delete Error when get image: " + err.Error())
		return error_utils.ErrorCode.DB_DOWN
	}

	if existingImage == nil {
		return error_utils.ErrorCode.NOT_FOUND
	}

	// Delete from S3 first
	err = s.s3Service.DeleteImage(ctx, existingImage.ImageKey)
	if err != nil {
		log.Error("ProductImageService.Delete Error when delete from S3: " + err.Error())
		// Continue with database deletion even if S3 deletion fails
	}

	// Delete from database
	err = s.imageRepository.DeleteCommand(ctx, id, nil)
	if err != nil {
		log.Error("ProductImageService.Delete Error when delete image: " + err.Error())
		return error_utils.ErrorCode.DB_DOWN
	}

	return ""
}
