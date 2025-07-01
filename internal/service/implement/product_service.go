package serviceimplement

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pna/management-app-backend/internal/bean"
	"github.com/pna/management-app-backend/internal/domain/entity"
	"github.com/pna/management-app-backend/internal/domain/model"
	"github.com/pna/management-app-backend/internal/repository"
	"github.com/pna/management-app-backend/internal/service"
	"github.com/pna/management-app-backend/internal/utils/error_utils"
	log "github.com/sirupsen/logrus"
)

type ProductService struct {
	productRepository      repository.ProductRepository
	inventoryRepository    repository.InventoryRepository
	categoryRepository     repository.ProductCategoryRepository
	unitRepository         repository.UnitOfMeasureRepository
	bomRepository          repository.ProductBomRepository
	unitOfWork             repository.UnitOfWork
	productImageRepository repository.ProductImageRepository
	s3Service              bean.S3Service
}

func NewProductService(
	productRepository repository.ProductRepository,
	inventoryRepository repository.InventoryRepository,
	categoryRepository repository.ProductCategoryRepository,
	unitRepository repository.UnitOfMeasureRepository,
	bomRepository repository.ProductBomRepository,
	unitOfWork repository.UnitOfWork,
	productImageRepository repository.ProductImageRepository,
	s3Service bean.S3Service,
) service.ProductService {
	return &ProductService{
		productRepository:      productRepository,
		inventoryRepository:    inventoryRepository,
		categoryRepository:     categoryRepository,
		unitRepository:         unitRepository,
		bomRepository:          bomRepository,
		unitOfWork:             unitOfWork,
		productImageRepository: productImageRepository,
		s3Service:              s3Service,
	}
}

// Helper function to build complete ProductResponse with all related information
func (s *ProductService) buildProductResponse(ctx *gin.Context, product *entity.Product) (*model.ProductResponse, string) {
	return s.buildProductResponseWithOptions(ctx, product, false)
}

func (s *ProductService) buildProductResponseWithOptions(ctx *gin.Context, product *entity.Product, noBom bool) (*model.ProductResponse, string) {
	response := &model.ProductResponse{
		ID:            product.ID,
		Name:          product.Name,
		Cost:          product.Cost,
		CategoryID:    product.CategoryID,
		UnitID:        product.UnitID,
		Description:   product.Description,
		OperationType: product.OperationType,
	}

	// Get inventory info
	inventory, err := s.inventoryRepository.GetOneByProductIDQuery(ctx, product.ID, nil)
	if err == nil && inventory != nil {
		response.Inventory = &model.InventoryInfo{
			Quantity: inventory.Quantity,
			Version:  inventory.Version,
		}
	}

	// Get category info
	if product.CategoryID != nil {
		category, err := s.categoryRepository.GetOneByIDQuery(ctx, *product.CategoryID, nil)
		if err == nil && category != nil {
			response.Category = &model.ProductCategoryResponse{
				ID:          category.ID,
				Name:        category.Name,
				Code:        category.Code,
				Description: category.Description,
			}
		}
	}

	// Get unit info
	if product.UnitID != nil {
		unit, err := s.unitRepository.GetOneByIDQuery(ctx, *product.UnitID, nil)
		if err == nil && unit != nil {
			response.Unit = &model.UnitOfMeasureResponse{
				ID:          unit.ID,
				Name:        unit.Name,
				Code:        unit.Code,
				Description: unit.Description,
			}
		}
	}

	// Skip BOM info if noBom is true (for performance optimization)
	if !noBom {
		// Get BOM info (if this product can be built from other products)
		bomEntries, err := s.bomRepository.GetByParentProductIDQuery(ctx, product.ID, nil)
		if err == nil && len(bomEntries) > 0 {
			bomComponents := make([]model.BomComponentResponse, len(bomEntries))
			for i, bomEntry := range bomEntries {
				// Get component product info
				componentProduct, _ := s.productRepository.GetOneByIDQuery(ctx, bomEntry.ComponentProductID, nil)

				bomComponents[i] = model.BomComponentResponse{
					ID:                 bomEntry.ID,
					ComponentProductID: bomEntry.ComponentProductID,
					Quantity:           bomEntry.Quantity,
				}

				if componentProduct != nil {
					// Get unit and category info for component
					unitName := ""
					categoryCode := ""

					if componentProduct.UnitID != nil {
						unit, _ := s.unitRepository.GetOneByIDQuery(ctx, *componentProduct.UnitID, nil)
						if unit != nil {
							unitName = unit.Name
						}
					}

					if componentProduct.CategoryID != nil {
						category, _ := s.categoryRepository.GetOneByIDQuery(ctx, *componentProduct.CategoryID, nil)
						if category != nil {
							categoryCode = category.Code
						}
					}

					bomComponents[i].ComponentProduct = &model.ProductBomInfo{
						ID:           componentProduct.ID,
						Name:         componentProduct.Name,
						Cost:         componentProduct.Cost,
						UnitName:     unitName,
						CategoryCode: categoryCode,
					}
				}
			}

			response.Bom = &model.ProductBOMInfo{
				TotalComponents: len(bomComponents),
				Components:      bomComponents,
			}
		}

		// Get usage info (where this product is used as component)
		usageEntries, err := s.bomRepository.GetByComponentProductIDQuery(ctx, product.ID, nil)
		if err == nil && len(usageEntries) > 0 {
			usageInfo := make([]model.ProductBOMUsage, len(usageEntries))
			for i, usageEntry := range usageEntries {
				parentProduct, _ := s.productRepository.GetOneByIDQuery(ctx, usageEntry.ParentProductID, nil)

				usageInfo[i] = model.ProductBOMUsage{
					ParentProductID: usageEntry.ParentProductID,
					Quantity:        usageEntry.Quantity,
				}

				if parentProduct != nil {
					usageInfo[i].ParentProductName = parentProduct.Name
				}
			}
			response.UsedInBoms = usageInfo
		}
	}

	// Get product images
	images, err := s.productImageRepository.GetByProductIDQuery(ctx, product.ID, nil)
	if err == nil && len(images) > 0 {
		imageResponses := make([]model.ProductImageResponse, len(images))
		for i, image := range images {
			// Generate signed URL for response
			signedURL, err := s.s3Service.GenerateSignedDownloadURL(ctx, image.ImageKey, 20*time.Minute)
			if err != nil {
				log.Error("ProductService.buildProductResponse Error generating signed URL: " + err.Error())
				signedURL = "" // Continue without signed URL
			}

			imageResponses[i] = model.ProductImageResponse{
				ID:        image.ID,
				ProductID: image.ProductID,
				ImageURL:  signedURL,
				ImageKey:  image.ImageKey,
			}
		}
		response.Images = imageResponses
	}

	return response, ""
}

func (s *ProductService) Create(ctx *gin.Context, request model.CreateProductRequest) (*model.ProductResponse, string) {
	// Begin transaction
	tx, err := s.unitOfWork.Begin(ctx)
	if err != nil {
		log.Error("ProductService.Create Error when begin transaction: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	// Defer rollback in case of error
	defer func() {
		if err != nil {
			if rollbackErr := s.unitOfWork.Rollback(tx); rollbackErr != nil {
				log.Error("ProductService.Create Error when rollback transaction: " + rollbackErr.Error())
			}
		}
	}()

	// Create product entity
	product := &entity.Product{
		Name:          request.Name,
		Cost:          request.Cost,
		CategoryID:    request.CategoryID,
		UnitID:        request.UnitID,
		Description:   request.Description,
		OperationType: request.OperationType,
	}

	// Save product to database
	err = s.productRepository.CreateCommand(ctx, product, tx)
	if err != nil {
		log.Error("ProductService.Create Error when create product: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	// Create inventory for the product
	inventory := &entity.Inventory{
		ProductID: product.ID,
		Quantity:  0, // Start with 0 quantity
		Version:   uuid.New().String(),
	}

	err = s.inventoryRepository.CreateCommand(ctx, inventory, tx)
	if err != nil {
		log.Error("ProductService.Create Error when create inventory: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	// Commit transaction
	err = s.unitOfWork.Commit(tx)
	if err != nil {
		log.Error("ProductService.Create Error when commit transaction: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	// Return complete response with all related info
	response, errCode := s.buildProductResponse(ctx, product)
	if errCode != "" {
		return nil, errCode
	}
	return response, ""
}

func (s *ProductService) Update(ctx *gin.Context, request model.UpdateProductRequest) (*model.ProductResponse, string) {
	// Check if product exists
	existingProduct, err := s.productRepository.GetOneByIDQuery(ctx, request.ID, nil)
	if err != nil {
		log.Error("ProductService.Update Error when get product: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	if existingProduct == nil {
		return nil, error_utils.ErrorCode.NOT_FOUND
	}

	// Update product entity
	product := &entity.Product{
		ID:            request.ID,
		Name:          request.Name,
		Cost:          request.Cost,
		CategoryID:    request.CategoryID,
		UnitID:        request.UnitID,
		Description:   request.Description,
		OperationType: request.OperationType,
	}

	// Save to database
	err = s.productRepository.UpdateCommand(ctx, product, nil)
	if err != nil {
		log.Error("ProductService.Update Error when update product: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	// Return complete response with all related info
	response, errCode := s.buildProductResponse(ctx, product)
	if errCode != "" {
		return nil, errCode
	}
	return response, ""
}

func (s *ProductService) GetAll(ctx *gin.Context, categoryFilter string, operationTypeFilter string, noBom bool) (*model.GetAllProductsResponse, string) {
	// Get all products with category filter
	products, err := s.productRepository.GetAllQuery(ctx, categoryFilter, operationTypeFilter, nil)
	if err != nil {
		log.Error("ProductService.GetAll Error when get products: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	// Convert to response models with complete info
	productResponses := make([]model.ProductResponse, len(products))
	for i, product := range products {
		response, errCode := s.buildProductResponseWithOptions(ctx, &product, noBom)
		if errCode != "" {
			log.Error("ProductService.GetAll Error when build product response for product " + string(rune(product.ID)) + ": " + errCode)
			// Continue with basic info if detailed info fails
			productResponses[i] = model.ProductResponse{
				ID:            product.ID,
				Name:          product.Name,
				Cost:          product.Cost,
				CategoryID:    product.CategoryID,
				UnitID:        product.UnitID,
				Description:   product.Description,
				OperationType: product.OperationType,
			}
			continue
		}
		productResponses[i] = *response
	}

	return &model.GetAllProductsResponse{
		Products: productResponses,
	}, ""
}

func (s *ProductService) GetOne(ctx *gin.Context, id int) (*model.GetOneProductResponse, string) {
	// Get product by ID
	product, err := s.productRepository.GetOneByIDQuery(ctx, id, nil)
	if err != nil {
		log.Error("ProductService.GetOne Error when get product: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	if product == nil {
		return nil, error_utils.ErrorCode.NOT_FOUND
	}

	// Return complete response with all related info
	response, errCode := s.buildProductResponse(ctx, product)
	if errCode != "" {
		return nil, errCode
	}

	return &model.GetOneProductResponse{
		Product: *response,
	}, ""
}
