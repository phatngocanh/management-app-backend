package serviceimplement

import (
	"github.com/gin-gonic/gin"
	"github.com/pna/management-app-backend/internal/domain/entity"
	"github.com/pna/management-app-backend/internal/domain/model"
	"github.com/pna/management-app-backend/internal/repository"
	"github.com/pna/management-app-backend/internal/service"
	"github.com/pna/management-app-backend/internal/utils/error_utils"
	log "github.com/sirupsen/logrus"
)

type ProductCategoryService struct {
	categoryRepository repository.ProductCategoryRepository
}

func NewProductCategoryService(categoryRepository repository.ProductCategoryRepository) service.ProductCategoryService {
	return &ProductCategoryService{
		categoryRepository: categoryRepository,
	}
}

func (s *ProductCategoryService) Create(ctx *gin.Context, request model.CreateProductCategoryRequest) (*model.ProductCategoryResponse, string) {
	// Create category entity
	category := &entity.ProductCategory{
		Name:        request.Name,
		Code:        request.Code,
		Description: request.Description,
	}

	// Save to database
	err := s.categoryRepository.CreateCommand(ctx, category, nil)
	if err != nil {
		log.Error("ProductCategoryService.Create Error when create category: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	return &model.ProductCategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Code:        category.Code,
		Description: category.Description,
	}, ""
}

func (s *ProductCategoryService) Update(ctx *gin.Context, request model.UpdateProductCategoryRequest) (*model.ProductCategoryResponse, string) {
	// Check if category exists
	existingCategory, err := s.categoryRepository.GetOneByIDQuery(ctx, request.ID, nil)
	if err != nil {
		log.Error("ProductCategoryService.Update Error when get category: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	if existingCategory == nil {
		return nil, error_utils.ErrorCode.NOT_FOUND
	}

	// Update category entity
	category := &entity.ProductCategory{
		ID:          request.ID,
		Name:        request.Name,
		Code:        request.Code,
		Description: request.Description,
	}

	// Save to database
	err = s.categoryRepository.UpdateCommand(ctx, category, nil)
	if err != nil {
		log.Error("ProductCategoryService.Update Error when update category: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	return &model.ProductCategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Code:        category.Code,
		Description: category.Description,
	}, ""
}

func (s *ProductCategoryService) GetAll(ctx *gin.Context) (*model.GetAllProductCategoriesResponse, string) {
	// Get all categories
	categories, err := s.categoryRepository.GetAllQuery(ctx, nil)
	if err != nil {
		log.Error("ProductCategoryService.GetAll Error when get categories: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	// Convert to response models
	categoryResponses := make([]model.ProductCategoryResponse, len(categories))
	for i, category := range categories {
		categoryResponses[i] = model.ProductCategoryResponse{
			ID:          category.ID,
			Name:        category.Name,
			Code:        category.Code,
			Description: category.Description,
		}
	}

	return &model.GetAllProductCategoriesResponse{
		Categories: categoryResponses,
	}, ""
}

func (s *ProductCategoryService) GetOne(ctx *gin.Context, id int) (*model.GetOneProductCategoryResponse, string) {
	// Get category by ID
	category, err := s.categoryRepository.GetOneByIDQuery(ctx, id, nil)
	if err != nil {
		log.Error("ProductCategoryService.GetOne Error when get category: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	if category == nil {
		return nil, error_utils.ErrorCode.NOT_FOUND
	}

	return &model.GetOneProductCategoryResponse{
		Category: model.ProductCategoryResponse{
			ID:          category.ID,
			Name:        category.Name,
			Code:        category.Code,
			Description: category.Description,
		},
	}, ""
}

func (s *ProductCategoryService) GetByCode(ctx *gin.Context, code string) (*model.GetOneProductCategoryResponse, string) {
	// Get category by code
	category, err := s.categoryRepository.GetOneByCodeQuery(ctx, code, nil)
	if err != nil {
		log.Error("ProductCategoryService.GetByCode Error when get category: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	if category == nil {
		return nil, error_utils.ErrorCode.NOT_FOUND
	}

	return &model.GetOneProductCategoryResponse{
		Category: model.ProductCategoryResponse{
			ID:          category.ID,
			Name:        category.Name,
			Code:        category.Code,
			Description: category.Description,
		},
	}, ""
}
