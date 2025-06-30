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
	bomRepository      repository.ProductBomRepository
	productRepository  repository.ProductRepository
	categoryRepository repository.ProductCategoryRepository
	unitRepository     repository.UnitOfMeasureRepository
	unitOfWork         repository.UnitOfWork
}

func NewProductBomService(
	bomRepository repository.ProductBomRepository,
	productRepository repository.ProductRepository,
	categoryRepository repository.ProductCategoryRepository,
	unitRepository repository.UnitOfMeasureRepository,
	unitOfWork repository.UnitOfWork,
) service.ProductBomService {
	return &ProductBomService{
		bomRepository:      bomRepository,
		productRepository:  productRepository,
		categoryRepository: categoryRepository,
		unitRepository:     unitRepository,
		unitOfWork:         unitOfWork,
	}
}

// Helper function to create ProductBomInfo with unit and category codes
func (s *ProductBomService) buildProductBomInfo(ctx *gin.Context, product *entity.Product) *model.ProductBomInfo {
	if product == nil {
		return nil
	}

	unitCode := ""
	categoryCode := ""

	// Get unit code
	if product.UnitID != nil {
		unit, err := s.unitRepository.GetOneByIDQuery(ctx, *product.UnitID, nil)
		if err == nil && unit != nil {
			unitCode = unit.Code
		}
	}

	// Get category code
	if product.CategoryID != nil {
		category, err := s.categoryRepository.GetOneByIDQuery(ctx, *product.CategoryID, nil)
		if err == nil && category != nil {
			categoryCode = category.Code
		}
	}

	return &model.ProductBomInfo{
		ID:           product.ID,
		Name:         product.Name,
		Cost:         product.Cost,
		UnitCode:     unitCode,
		CategoryCode: categoryCode,
	}
}

func (s *ProductBomService) Create(ctx *gin.Context, request model.CreateProductBomRequest) (*model.ProductBomResponse, string) {
	// Begin transaction
	tx, err := s.unitOfWork.Begin(ctx)
	if err != nil {
		log.Error("ProductBomService.Create Error when begin transaction: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	// Defer rollback in case of error
	defer func() {
		if err != nil {
			if rollbackErr := s.unitOfWork.Rollback(tx); rollbackErr != nil {
				log.Error("ProductBomService.Create Error when rollback transaction: " + rollbackErr.Error())
			}
		}
	}()

	// Verify parent product exists
	parentProduct, err := s.productRepository.GetOneByIDQuery(ctx, request.ParentProductID, tx)
	if err != nil {
		log.Error("ProductBomService.Create Error when get parent product: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}
	if parentProduct == nil {
		return nil, error_utils.ErrorCode.NOT_FOUND
	}

	// Create BOM entries for each component
	bomComponents := make([]model.BomComponentResponse, len(request.Components))
	for i, component := range request.Components {
		// Create BOM entity
		bom := &entity.ProductBom{
			ParentProductID:    request.ParentProductID,
			ComponentProductID: component.ComponentProductID,
			Quantity:           component.Quantity,
		}

		// Save to database
		err = s.bomRepository.CreateCommand(ctx, bom, tx)
		if err != nil {
			log.Error("ProductBomService.Create Error when create bom: " + err.Error())
			return nil, error_utils.ErrorCode.DB_DOWN
		}

		// Get component product info
		componentProduct, _ := s.productRepository.GetOneByIDQuery(ctx, component.ComponentProductID, tx)

		bomComponents[i] = model.BomComponentResponse{
			ID:                 bom.ID,
			ComponentProductID: component.ComponentProductID,
			Quantity:           component.Quantity,
		}

		if componentProduct != nil {
			bomComponents[i].ComponentProduct = s.buildProductBomInfo(ctx, componentProduct)
		}
	}

	// Commit transaction
	err = s.unitOfWork.Commit(tx)
	if err != nil {
		log.Error("ProductBomService.Create Error when commit transaction: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	// Prepare response
	response := &model.ProductBomResponse{
		ParentProductID: request.ParentProductID,
		Components:      bomComponents,
		TotalComponents: len(bomComponents),
	}

	if parentProduct != nil {
		response.ParentProduct = s.buildProductBomInfo(ctx, parentProduct)
	}

	return response, ""
}

func (s *ProductBomService) Update(ctx *gin.Context, request model.UpdateProductBomRequest) (*model.ProductBomResponse, string) {
	// Begin transaction
	tx, err := s.unitOfWork.Begin(ctx)
	if err != nil {
		log.Error("ProductBomService.Update Error when begin transaction: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	// Defer rollback in case of error
	defer func() {
		if err != nil {
			if rollbackErr := s.unitOfWork.Rollback(tx); rollbackErr != nil {
				log.Error("ProductBomService.Update Error when rollback transaction: " + rollbackErr.Error())
			}
		}
	}()

	// Verify parent product exists
	parentProduct, err := s.productRepository.GetOneByIDQuery(ctx, request.ParentProductID, tx)
	if err != nil {
		log.Error("ProductBomService.Update Error when get parent product: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}
	if parentProduct == nil {
		return nil, error_utils.ErrorCode.NOT_FOUND
	}

	// Delete all existing BOM entries for this parent product
	existingBoms, err := s.bomRepository.GetByParentProductIDQuery(ctx, request.ParentProductID, tx)
	if err != nil {
		log.Error("ProductBomService.Update Error when get existing boms: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	for _, existingBom := range existingBoms {
		err = s.bomRepository.DeleteCommand(ctx, existingBom.ID, tx)
		if err != nil {
			log.Error("ProductBomService.Update Error when delete existing bom: " + err.Error())
			return nil, error_utils.ErrorCode.DB_DOWN
		}
	}

	// Create new BOM entries for each component
	bomComponents := make([]model.BomComponentResponse, len(request.Components))
	for i, component := range request.Components {
		// Create BOM entity
		bom := &entity.ProductBom{
			ParentProductID:    request.ParentProductID,
			ComponentProductID: component.ComponentProductID,
			Quantity:           component.Quantity,
		}

		// Save to database
		err = s.bomRepository.CreateCommand(ctx, bom, tx)
		if err != nil {
			log.Error("ProductBomService.Update Error when create bom: " + err.Error())
			return nil, error_utils.ErrorCode.DB_DOWN
		}

		// Get component product info
		componentProduct, _ := s.productRepository.GetOneByIDQuery(ctx, component.ComponentProductID, tx)

		bomComponents[i] = model.BomComponentResponse{
			ID:                 bom.ID,
			ComponentProductID: component.ComponentProductID,
			Quantity:           component.Quantity,
		}

		if componentProduct != nil {
			bomComponents[i].ComponentProduct = s.buildProductBomInfo(ctx, componentProduct)
		}
	}

	// Commit transaction
	err = s.unitOfWork.Commit(tx)
	if err != nil {
		log.Error("ProductBomService.Update Error when commit transaction: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	// Prepare response
	response := &model.ProductBomResponse{
		ParentProductID: request.ParentProductID,
		Components:      bomComponents,
		TotalComponents: len(bomComponents),
	}

	if parentProduct != nil {
		response.ParentProduct = s.buildProductBomInfo(ctx, parentProduct)
	}

	return response, ""
}

func (s *ProductBomService) GetAll(ctx *gin.Context) (*model.GetAllProductBomsResponse, string) {
	// Get all BOM entries
	allBoms, err := s.bomRepository.GetAllQuery(ctx, nil)
	if err != nil {
		log.Error("ProductBomService.GetAll Error when get boms: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	// Group BOM entries by parent product ID
	bomMap := make(map[int][]entity.ProductBom)
	for _, bom := range allBoms {
		bomMap[bom.ParentProductID] = append(bomMap[bom.ParentProductID], bom)
	}

	// Convert to response models
	bomResponses := make([]model.ProductBomResponse, 0, len(bomMap))
	for parentProductID, components := range bomMap {
		// Get parent product info
		parentProduct, _ := s.productRepository.GetOneByIDQuery(ctx, parentProductID, nil)

		// Convert components
		bomComponents := make([]model.BomComponentResponse, len(components))
		for i, component := range components {
			// Get component product info
			componentProduct, _ := s.productRepository.GetOneByIDQuery(ctx, component.ComponentProductID, nil)

			bomComponents[i] = model.BomComponentResponse{
				ID:                 component.ID,
				ComponentProductID: component.ComponentProductID,
				Quantity:           component.Quantity,
			}

			if componentProduct != nil {
				bomComponents[i].ComponentProduct = s.buildProductBomInfo(ctx, componentProduct)
			}
		}

		bomResponse := model.ProductBomResponse{
			ParentProductID: parentProductID,
			Components:      bomComponents,
			TotalComponents: len(bomComponents),
		}

		if parentProduct != nil {
			bomResponse.ParentProduct = s.buildProductBomInfo(ctx, parentProduct)
		}

		bomResponses = append(bomResponses, bomResponse)
	}

	return &model.GetAllProductBomsResponse{
		Boms: bomResponses,
	}, ""
}

func (s *ProductBomService) GetByParentProductID(ctx *gin.Context, parentProductID int) (*model.GetOneProductBomResponse, string) {
	// Get BOM entries for parent product
	boms, err := s.bomRepository.GetByParentProductIDQuery(ctx, parentProductID, nil)
	if err != nil {
		log.Error("ProductBomService.GetByParentProductID Error when get boms: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	// Get parent product info
	parentProduct, err := s.productRepository.GetOneByIDQuery(ctx, parentProductID, nil)
	if err != nil {
		log.Error("ProductBomService.GetByParentProductID Error when get parent product: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	if parentProduct == nil {
		return nil, error_utils.ErrorCode.NOT_FOUND
	}

	// Convert components
	bomComponents := make([]model.BomComponentResponse, len(boms))
	for i, bom := range boms {
		// Get component product info
		componentProduct, _ := s.productRepository.GetOneByIDQuery(ctx, bom.ComponentProductID, nil)

		bomComponents[i] = model.BomComponentResponse{
			ID:                 bom.ID,
			ComponentProductID: bom.ComponentProductID,
			Quantity:           bom.Quantity,
		}

		if componentProduct != nil {
			bomComponents[i].ComponentProduct = s.buildProductBomInfo(ctx, componentProduct)
		}
	}

	response := model.ProductBomResponse{
		ParentProductID: parentProductID,
		Components:      bomComponents,
		TotalComponents: len(bomComponents),
	}

	if parentProduct != nil {
		response.ParentProduct = s.buildProductBomInfo(ctx, parentProduct)
	}

	return &model.GetOneProductBomResponse{
		Bom: response,
	}, ""
}

func (s *ProductBomService) GetByComponentProductID(ctx *gin.Context, componentProductID int) (*model.GetAllProductBomsResponse, string) {
	// Get BOM entries where this product is used as component
	boms, err := s.bomRepository.GetByComponentProductIDQuery(ctx, componentProductID, nil)
	if err != nil {
		log.Error("ProductBomService.GetByComponentProductID Error when get boms: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	// Group by parent product ID
	bomMap := make(map[int][]entity.ProductBom)
	for _, bom := range boms {
		bomMap[bom.ParentProductID] = append(bomMap[bom.ParentProductID], bom)
	}

	// For each parent product, get its complete BOM
	bomResponses := make([]model.ProductBomResponse, 0, len(bomMap))
	for parentProductID := range bomMap {
		bomResponse, errCode := s.GetByParentProductID(ctx, parentProductID)
		if errCode != "" {
			log.Error("ProductBomService.GetByComponentProductID Error when get BOM for parent " + string(rune(parentProductID)))
			continue
		}
		bomResponses = append(bomResponses, bomResponse.Bom)
	}

	return &model.GetAllProductBomsResponse{
		Boms: bomResponses,
	}, ""
}

func (s *ProductBomService) DeleteByParentProductID(ctx *gin.Context, parentProductID int) string {
	// Begin transaction
	tx, err := s.unitOfWork.Begin(ctx)
	if err != nil {
		log.Error("ProductBomService.DeleteByParentProductID Error when begin transaction: " + err.Error())
		return error_utils.ErrorCode.DB_DOWN
	}

	// Defer rollback in case of error
	defer func() {
		if err != nil {
			if rollbackErr := s.unitOfWork.Rollback(tx); rollbackErr != nil {
				log.Error("ProductBomService.DeleteByParentProductID Error when rollback transaction: " + rollbackErr.Error())
			}
		}
	}()

	// Get all BOM entries for this parent product
	boms, err := s.bomRepository.GetByParentProductIDQuery(ctx, parentProductID, tx)
	if err != nil {
		log.Error("ProductBomService.DeleteByParentProductID Error when get boms: " + err.Error())
		return error_utils.ErrorCode.DB_DOWN
	}

	if len(boms) == 0 {
		return error_utils.ErrorCode.NOT_FOUND
	}

	// Delete all BOM entries
	for _, bom := range boms {
		err = s.bomRepository.DeleteCommand(ctx, bom.ID, tx)
		if err != nil {
			log.Error("ProductBomService.DeleteByParentProductID Error when delete bom: " + err.Error())
			return error_utils.ErrorCode.DB_DOWN
		}
	}

	// Commit transaction
	err = s.unitOfWork.Commit(tx)
	if err != nil {
		log.Error("ProductBomService.DeleteByParentProductID Error when commit transaction: " + err.Error())
		return error_utils.ErrorCode.DB_DOWN
	}

	return ""
}

func (s *ProductBomService) CalculateMaterialRequirements(ctx *gin.Context, request model.CalculateMaterialRequirementsRequest) (*model.MaterialRequirementsResponse, string) {
	// Verify parent product exists
	parentProduct, err := s.productRepository.GetOneByIDQuery(ctx, request.ParentProductID, nil)
	if err != nil {
		log.Error("ProductBomService.CalculateMaterialRequirements Error when get parent product: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}
	if parentProduct == nil {
		return nil, error_utils.ErrorCode.NOT_FOUND
	}

	// Map to accumulate material requirements by product ID
	materialMap := make(map[int]float64)

	// Recursively calculate material requirements
	err = s.calculateRequirementsRecursive(ctx, request.ParentProductID, request.Quantity, materialMap)
	if err != nil {
		log.Error("ProductBomService.CalculateMaterialRequirements Error in recursive calculation: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	// Convert map to slice and get product info
	var materialRequirements []model.MaterialRequirement
	for productID, quantity := range materialMap {
		// Get product info
		product, err := s.productRepository.GetOneByIDQuery(ctx, productID, nil)
		if err != nil {
			log.Error("ProductBomService.CalculateMaterialRequirements Error when get material product: " + err.Error())
			continue // Skip this material if we can't get its info
		}

		requirement := model.MaterialRequirement{
			ProductID:        productID,
			RequiredQuantity: quantity,
		}

		if product != nil {
			requirement.Product = s.buildProductBomInfo(ctx, product)
		}

		materialRequirements = append(materialRequirements, requirement)
	}

	// Prepare response
	response := &model.MaterialRequirementsResponse{
		ParentProductID:      request.ParentProductID,
		RequestedQuantity:    request.Quantity,
		MaterialRequirements: materialRequirements,
		TotalMaterials:       len(materialRequirements),
	}

	if parentProduct != nil {
		response.ParentProduct = s.buildProductBomInfo(ctx, parentProduct)
	}

	return response, ""
}

// Helper function to recursively calculate material requirements
func (s *ProductBomService) calculateRequirementsRecursive(ctx *gin.Context, productID int, quantity float64, materialMap map[int]float64) error {
	// Get BOM for this product
	boms, err := s.bomRepository.GetByParentProductIDQuery(ctx, productID, nil)
	if err != nil {
		return err
	}

	// If no BOM found, this is a raw material
	if len(boms) == 0 {
		// Accumulate quantity for this raw material
		materialMap[productID] += quantity
		return nil
	}

	// If BOM exists, recursively calculate for each component
	for _, bom := range boms {
		componentQuantity := bom.Quantity * quantity
		err = s.calculateRequirementsRecursive(ctx, bom.ComponentProductID, componentQuantity, materialMap)
		if err != nil {
			return err
		}
	}

	return nil
}
