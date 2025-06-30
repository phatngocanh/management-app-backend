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

type ProductBomService struct {
	bomRepository     repository.ProductBomRepository
	productRepository repository.ProductRepository
	unitOfWork        repository.UnitOfWork
}

func NewProductBomService(bomRepository repository.ProductBomRepository, productRepository repository.ProductRepository, unitOfWork repository.UnitOfWork) service.ProductBomService {
	return &ProductBomService{
		bomRepository:     bomRepository,
		productRepository: productRepository,
		unitOfWork:        unitOfWork,
	}
}

func (s *ProductBomService) Create(ctx *gin.Context, request model.CreateProductBomRequest) (*model.ProductBomResponse, string) {
	// Create BOM entity
	bom := &entity.ProductBom{
		ParentProductID:    request.ParentProductID,
		ComponentProductID: request.ComponentProductID,
		Quantity:           request.Quantity,
	}

	// Save to database
	err := s.bomRepository.CreateCommand(ctx, bom, nil)
	if err != nil {
		log.Error("ProductBomService.Create Error when create bom: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	// Get product info for response
	parentProduct, _ := s.productRepository.GetOneByIDQuery(ctx, bom.ParentProductID, nil)
	componentProduct, _ := s.productRepository.GetOneByIDQuery(ctx, bom.ComponentProductID, nil)

	response := &model.ProductBomResponse{
		ID:                 bom.ID,
		ParentProductID:    bom.ParentProductID,
		ComponentProductID: bom.ComponentProductID,
		Quantity:           bom.Quantity,
	}

	if parentProduct != nil {
		response.ParentProduct = &model.ProductBomInfo{
			ID:            parentProduct.ID,
			Name:          parentProduct.Name,
			Spec:          parentProduct.Spec,
			OriginalPrice: parentProduct.OriginalPrice,
		}
	}

	if componentProduct != nil {
		response.ComponentProduct = &model.ProductBomInfo{
			ID:            componentProduct.ID,
			Name:          componentProduct.Name,
			Spec:          componentProduct.Spec,
			OriginalPrice: componentProduct.OriginalPrice,
		}
	}

	return response, ""
}

func (s *ProductBomService) Update(ctx *gin.Context, request model.UpdateProductBomRequest) (*model.ProductBomResponse, string) {
	// Check if BOM exists
	existingBom, err := s.bomRepository.GetOneByIDQuery(ctx, request.ID, nil)
	if err != nil {
		log.Error("ProductBomService.Update Error when get bom: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	if existingBom == nil {
		return nil, error_utils.ErrorCode.NOT_FOUND
	}

	// Update BOM entity
	bom := &entity.ProductBom{
		ID:                 request.ID,
		ParentProductID:    request.ParentProductID,
		ComponentProductID: request.ComponentProductID,
		Quantity:           request.Quantity,
	}

	// Save to database
	err = s.bomRepository.UpdateCommand(ctx, bom, nil)
	if err != nil {
		log.Error("ProductBomService.Update Error when update bom: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	// Get product info for response
	parentProduct, _ := s.productRepository.GetOneByIDQuery(ctx, bom.ParentProductID, nil)
	componentProduct, _ := s.productRepository.GetOneByIDQuery(ctx, bom.ComponentProductID, nil)

	response := &model.ProductBomResponse{
		ID:                 bom.ID,
		ParentProductID:    bom.ParentProductID,
		ComponentProductID: bom.ComponentProductID,
		Quantity:           bom.Quantity,
	}

	if parentProduct != nil {
		response.ParentProduct = &model.ProductBomInfo{
			ID:            parentProduct.ID,
			Name:          parentProduct.Name,
			Spec:          parentProduct.Spec,
			OriginalPrice: parentProduct.OriginalPrice,
		}
	}

	if componentProduct != nil {
		response.ComponentProduct = &model.ProductBomInfo{
			ID:            componentProduct.ID,
			Name:          componentProduct.Name,
			Spec:          componentProduct.Spec,
			OriginalPrice: componentProduct.OriginalPrice,
		}
	}

	return response, ""
}

func (s *ProductBomService) GetAll(ctx *gin.Context) (*model.GetAllProductBomsResponse, string) {
	// Get all BOMs
	boms, err := s.bomRepository.GetAllQuery(ctx, nil)
	if err != nil {
		log.Error("ProductBomService.GetAll Error when get boms: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	// Convert to response models
	bomResponses := make([]model.ProductBomResponse, len(boms))
	for i, bom := range boms {
		bomResponses[i] = model.ProductBomResponse{
			ID:                 bom.ID,
			ParentProductID:    bom.ParentProductID,
			ComponentProductID: bom.ComponentProductID,
			Quantity:           bom.Quantity,
		}

		// Optionally get product info
		parentProduct, _ := s.productRepository.GetOneByIDQuery(ctx, bom.ParentProductID, nil)
		componentProduct, _ := s.productRepository.GetOneByIDQuery(ctx, bom.ComponentProductID, nil)

		if parentProduct != nil {
			bomResponses[i].ParentProduct = &model.ProductBomInfo{
				ID:            parentProduct.ID,
				Name:          parentProduct.Name,
				Spec:          parentProduct.Spec,
				OriginalPrice: parentProduct.OriginalPrice,
			}
		}

		if componentProduct != nil {
			bomResponses[i].ComponentProduct = &model.ProductBomInfo{
				ID:            componentProduct.ID,
				Name:          componentProduct.Name,
				Spec:          componentProduct.Spec,
				OriginalPrice: componentProduct.OriginalPrice,
			}
		}
	}

	return &model.GetAllProductBomsResponse{
		Boms: bomResponses,
	}, ""
}

func (s *ProductBomService) GetOne(ctx *gin.Context, id int) (*model.GetOneProductBomResponse, string) {
	// Get BOM by ID
	bom, err := s.bomRepository.GetOneByIDQuery(ctx, id, nil)
	if err != nil {
		log.Error("ProductBomService.GetOne Error when get bom: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	if bom == nil {
		return nil, error_utils.ErrorCode.NOT_FOUND
	}

	// Get product info
	parentProduct, _ := s.productRepository.GetOneByIDQuery(ctx, bom.ParentProductID, nil)
	componentProduct, _ := s.productRepository.GetOneByIDQuery(ctx, bom.ComponentProductID, nil)

	response := model.ProductBomResponse{
		ID:                 bom.ID,
		ParentProductID:    bom.ParentProductID,
		ComponentProductID: bom.ComponentProductID,
		Quantity:           bom.Quantity,
	}

	if parentProduct != nil {
		response.ParentProduct = &model.ProductBomInfo{
			ID:            parentProduct.ID,
			Name:          parentProduct.Name,
			Spec:          parentProduct.Spec,
			OriginalPrice: parentProduct.OriginalPrice,
		}
	}

	if componentProduct != nil {
		response.ComponentProduct = &model.ProductBomInfo{
			ID:            componentProduct.ID,
			Name:          componentProduct.Name,
			Spec:          componentProduct.Spec,
			OriginalPrice: componentProduct.OriginalPrice,
		}
	}

	return &model.GetOneProductBomResponse{
		Bom: response,
	}, ""
}

func (s *ProductBomService) GetByParentProductID(ctx *gin.Context, parentProductID int) (*model.GetProductBomsResponse, string) {
	// Get BOMs by parent product ID
	boms, err := s.bomRepository.GetByParentProductIDQuery(ctx, parentProductID, nil)
	if err != nil {
		log.Error("ProductBomService.GetByParentProductID Error when get boms: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	// Convert to response models
	bomResponses := make([]model.ProductBomResponse, len(boms))
	for i, bom := range boms {
		bomResponses[i] = model.ProductBomResponse{
			ID:                 bom.ID,
			ParentProductID:    bom.ParentProductID,
			ComponentProductID: bom.ComponentProductID,
			Quantity:           bom.Quantity,
		}

		// Get component product info
		componentProduct, _ := s.productRepository.GetOneByIDQuery(ctx, bom.ComponentProductID, nil)
		if componentProduct != nil {
			bomResponses[i].ComponentProduct = &model.ProductBomInfo{
				ID:            componentProduct.ID,
				Name:          componentProduct.Name,
				Spec:          componentProduct.Spec,
				OriginalPrice: componentProduct.OriginalPrice,
			}
		}
	}

	return &model.GetProductBomsResponse{
		Boms: bomResponses,
	}, ""
}

func (s *ProductBomService) GetByComponentProductID(ctx *gin.Context, componentProductID int) (*model.GetProductBomsResponse, string) {
	// Get BOMs by component product ID
	boms, err := s.bomRepository.GetByComponentProductIDQuery(ctx, componentProductID, nil)
	if err != nil {
		log.Error("ProductBomService.GetByComponentProductID Error when get boms: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	// Convert to response models
	bomResponses := make([]model.ProductBomResponse, len(boms))
	for i, bom := range boms {
		bomResponses[i] = model.ProductBomResponse{
			ID:                 bom.ID,
			ParentProductID:    bom.ParentProductID,
			ComponentProductID: bom.ComponentProductID,
			Quantity:           bom.Quantity,
		}

		// Get parent product info
		parentProduct, _ := s.productRepository.GetOneByIDQuery(ctx, bom.ParentProductID, nil)
		if parentProduct != nil {
			bomResponses[i].ParentProduct = &model.ProductBomInfo{
				ID:            parentProduct.ID,
				Name:          parentProduct.Name,
				Spec:          parentProduct.Spec,
				OriginalPrice: parentProduct.OriginalPrice,
			}
		}
	}

	return &model.GetProductBomsResponse{
		Boms: bomResponses,
	}, ""
}

func (s *ProductBomService) Delete(ctx *gin.Context, id int) string {
	// Check if BOM exists
	existingBom, err := s.bomRepository.GetOneByIDQuery(ctx, id, nil)
	if err != nil {
		log.Error("ProductBomService.Delete Error when get bom: " + err.Error())
		return error_utils.ErrorCode.DB_DOWN
	}

	if existingBom == nil {
		return error_utils.ErrorCode.NOT_FOUND
	}

	// Delete from database
	err = s.bomRepository.DeleteCommand(ctx, id, nil)
	if err != nil {
		log.Error("ProductBomService.Delete Error when delete bom: " + err.Error())
		return error_utils.ErrorCode.DB_DOWN
	}

	return ""
}
