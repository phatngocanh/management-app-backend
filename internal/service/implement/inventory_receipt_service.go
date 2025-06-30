package serviceimplement

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pna/management-app-backend/internal/domain/entity"
	"github.com/pna/management-app-backend/internal/domain/model"
	"github.com/pna/management-app-backend/internal/repository"
	"github.com/pna/management-app-backend/internal/service"
	"github.com/pna/management-app-backend/internal/utils/error_utils"
	log "github.com/sirupsen/logrus"
)

type InventoryReceiptService struct {
	inventoryReceiptRepository     repository.InventoryReceiptRepository
	inventoryReceiptItemRepository repository.InventoryReceiptItemRepository
	inventoryRepository            repository.InventoryRepository
	inventoryHistoryRepository     repository.InventoryHistoryRepository
	userRepository                 repository.UserRepository
	productRepository              repository.ProductRepository
	unitOfWork                     repository.UnitOfWork
}

func NewInventoryReceiptService(
	inventoryReceiptRepository repository.InventoryReceiptRepository,
	inventoryReceiptItemRepository repository.InventoryReceiptItemRepository,
	inventoryRepository repository.InventoryRepository,
	inventoryHistoryRepository repository.InventoryHistoryRepository,
	userRepository repository.UserRepository,
	productRepository repository.ProductRepository,
	unitOfWork repository.UnitOfWork,
) service.InventoryReceiptService {
	return &InventoryReceiptService{
		inventoryReceiptRepository:     inventoryReceiptRepository,
		inventoryReceiptItemRepository: inventoryReceiptItemRepository,
		inventoryRepository:            inventoryRepository,
		inventoryHistoryRepository:     inventoryHistoryRepository,
		userRepository:                 userRepository,
		productRepository:              productRepository,
		unitOfWork:                     unitOfWork,
	}
}

func (s *InventoryReceiptService) Create(ctx *gin.Context, request model.CreateInventoryReceiptRequest) (*model.InventoryReceiptResponse, string) {
	// Get user details to get username for history
	user, err := s.userRepository.FindByIDQuery(ctx, request.UserID, nil)
	if err != nil {
		log.Error("InventoryReceiptService.Create Error when get user: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	if user == nil {
		log.Error("InventoryReceiptService.Create Error: user not found")
		return nil, error_utils.ErrorCode.NOT_FOUND
	}

	// Begin transaction
	tx, err := s.unitOfWork.Begin(ctx)
	if err != nil {
		log.Error("InventoryReceiptService.Create Error when begin transaction: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	// Defer rollback in case of error
	defer func() {
		if err != nil {
			if rollbackErr := s.unitOfWork.Rollback(tx); rollbackErr != nil {
				log.Error("InventoryReceiptService.Create Error when rollback transaction: " + rollbackErr.Error())
			}
		}
	}()

	// Set receipt date to now if not provided
	receiptDate := request.ReceiptDate
	if receiptDate.IsZero() {
		receiptDate = time.Now()
	}

	// Create inventory receipt entity
	inventoryReceipt := &entity.InventoryReceipt{
		UserID:      request.UserID,
		ReceiptDate: receiptDate,
		Notes:       request.Notes,
		TotalItems:  len(request.Items),
	}

	// Create inventory receipt (code will be generated in repository)
	err = s.inventoryReceiptRepository.CreateCommand(ctx, inventoryReceipt, tx)
	if err != nil {
		log.Error("InventoryReceiptService.Create Error when create inventory receipt: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	// Create receipt items and update inventory
	var itemResponses []model.InventoryReceiptItemResponse
	for _, itemRequest := range request.Items {
		// Validate product exists
		product, err := s.productRepository.GetOneByIDQuery(ctx, itemRequest.ProductID, tx)
		if err != nil {
			log.Error("InventoryReceiptService.Create Error when get product: " + err.Error())
			return nil, error_utils.ErrorCode.DB_DOWN
		}
		if product == nil {
			log.Error(fmt.Sprintf("InventoryReceiptService.Create Error: product with ID %d not found", itemRequest.ProductID))
			return nil, error_utils.ErrorCode.NOT_FOUND
		}

		// Create receipt item entity
		receiptItem := &entity.InventoryReceiptItem{
			InventoryReceiptID: inventoryReceipt.ID,
			ProductID:          itemRequest.ProductID,
			Quantity:           itemRequest.Quantity,
			UnitCost:           itemRequest.UnitCost,
			Notes:              itemRequest.Notes,
		}

		// Create receipt item
		err = s.inventoryReceiptItemRepository.CreateCommand(ctx, receiptItem, tx)
		if err != nil {
			log.Error("InventoryReceiptService.Create Error when create receipt item: " + err.Error())
			return nil, error_utils.ErrorCode.DB_DOWN
		}

		// Get current inventory for this product
		inventory, err := s.inventoryRepository.GetOneByProductIDQuery(ctx, itemRequest.ProductID, tx)
		if err != nil {
			log.Error("InventoryReceiptService.Create Error when get inventory: " + err.Error())
			return nil, error_utils.ErrorCode.DB_DOWN
		}

		var finalQuantity int
		if inventory == nil {
			// Create new inventory record if doesn't exist
			newInventory := &entity.Inventory{
				ProductID: itemRequest.ProductID,
				Quantity:  itemRequest.Quantity,
				Version:   uuid.New().String(),
			}
			err = s.inventoryRepository.CreateCommand(ctx, newInventory, tx)
			if err != nil {
				log.Error("InventoryReceiptService.Create Error when create inventory: " + err.Error())
				return nil, error_utils.ErrorCode.DB_DOWN
			}
			finalQuantity = itemRequest.Quantity
		} else {
			// Update existing inventory
			newVersion := uuid.New().String()
			s.inventoryRepository.GetOneByIDForUpdateQuery(ctx, itemRequest.ProductID, tx)
			err = s.inventoryRepository.UpdateQuantityWithVersionCommand(ctx, itemRequest.ProductID, itemRequest.Quantity, inventory.Version, newVersion, tx)
			if err != nil {
				log.Error("InventoryReceiptService.Create Error when update inventory: " + err.Error())
				return nil, error_utils.ErrorCode.DB_DOWN
			}
			finalQuantity = inventory.Quantity + itemRequest.Quantity
		}

		// Create inventory history record
		historyNote := fmt.Sprintf("Nhập kho từ phiếu nhập %s", inventoryReceipt.Code)
		if itemRequest.Notes != nil && *itemRequest.Notes != "" {
			historyNote += fmt.Sprintf(" - %s", *itemRequest.Notes)
		}

		inventoryHistory := &entity.InventoryHistory{
			ProductID:     itemRequest.ProductID,
			Quantity:      itemRequest.Quantity,
			FinalQuantity: finalQuantity,
			ImporterName:  user.Username,
			ImportedAt:    time.Now(),
			Note:          historyNote,
			ReferenceID:   &inventoryReceipt.ID,
		}

		err = s.inventoryHistoryRepository.CreateCommand(ctx, inventoryHistory, tx)
		if err != nil {
			log.Error("InventoryReceiptService.Create Error when create inventory history: " + err.Error())
			return nil, error_utils.ErrorCode.DB_DOWN
		}

		// Add to response items
		itemResponses = append(itemResponses, model.InventoryReceiptItemResponse{
			ID:                 receiptItem.ID,
			InventoryReceiptID: receiptItem.InventoryReceiptID,
			ProductID:          receiptItem.ProductID,
			Quantity:           receiptItem.Quantity,
			UnitCost:           receiptItem.UnitCost,
			Notes:              receiptItem.Notes,
			CreatedAt:          receiptItem.CreatedAt,
			UpdatedAt:          receiptItem.UpdatedAt,
		})
	}

	// Commit transaction
	err = s.unitOfWork.Commit(tx)
	if err != nil {
		log.Error("InventoryReceiptService.Create Error when commit transaction: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	// Return response
	return &model.InventoryReceiptResponse{
		ID:          inventoryReceipt.ID,
		Code:        inventoryReceipt.Code,
		UserID:      inventoryReceipt.UserID,
		ReceiptDate: inventoryReceipt.ReceiptDate,
		Notes:       inventoryReceipt.Notes,
		TotalItems:  inventoryReceipt.TotalItems,
		CreatedAt:   inventoryReceipt.CreatedAt,
		UpdatedAt:   inventoryReceipt.UpdatedAt,
		Items:       itemResponses,
	}, ""
}

func (s *InventoryReceiptService) GetAll(ctx *gin.Context) (*model.GetAllInventoryReceiptsResponse, string) {
	// Get all inventory receipts
	receipts, err := s.inventoryReceiptRepository.GetAllQuery(ctx, nil)
	if err != nil {
		log.Error("InventoryReceiptService.GetAll Error when get receipts: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	// Convert to response models (without items for GetAll)
	receiptResponses := make([]model.InventoryReceiptResponse, len(receipts))
	for i, receipt := range receipts {
		receiptResponses[i] = model.InventoryReceiptResponse{
			ID:          receipt.ID,
			Code:        receipt.Code,
			UserID:      receipt.UserID,
			ReceiptDate: receipt.ReceiptDate,
			Notes:       receipt.Notes,
			TotalItems:  receipt.TotalItems,
			CreatedAt:   receipt.CreatedAt,
			UpdatedAt:   receipt.UpdatedAt,
			// Items omitted for GetAll
		}
	}

	return &model.GetAllInventoryReceiptsResponse{
		InventoryReceipts: receiptResponses,
	}, ""
}

func (s *InventoryReceiptService) GetOne(ctx *gin.Context, id int) (*model.GetOneInventoryReceiptResponse, string) {
	// Get inventory receipt by ID
	receipt, err := s.inventoryReceiptRepository.GetOneByIDQuery(ctx, id, nil)
	if err != nil {
		log.Error("InventoryReceiptService.GetOne Error when get receipt: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	if receipt == nil {
		return nil, error_utils.ErrorCode.NOT_FOUND
	}

	// Get receipt items
	receiptItems, err := s.inventoryReceiptItemRepository.GetByInventoryReceiptIDQuery(ctx, id, nil)
	if err != nil {
		log.Error("InventoryReceiptService.GetOne Error when get receipt items: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	// Convert items to response models
	itemResponses := make([]model.InventoryReceiptItemResponse, len(receiptItems))
	for i, item := range receiptItems {
		itemResponses[i] = model.InventoryReceiptItemResponse{
			ID:                 item.ID,
			InventoryReceiptID: item.InventoryReceiptID,
			ProductID:          item.ProductID,
			Quantity:           item.Quantity,
			UnitCost:           item.UnitCost,
			Notes:              item.Notes,
			CreatedAt:          item.CreatedAt,
			UpdatedAt:          item.UpdatedAt,
		}
	}

	// Return response with items
	return &model.GetOneInventoryReceiptResponse{
		InventoryReceipt: model.InventoryReceiptResponse{
			ID:          receipt.ID,
			Code:        receipt.Code,
			UserID:      receipt.UserID,
			ReceiptDate: receipt.ReceiptDate,
			Notes:       receipt.Notes,
			TotalItems:  receipt.TotalItems,
			CreatedAt:   receipt.CreatedAt,
			UpdatedAt:   receipt.UpdatedAt,
			Items:       itemResponses,
		},
	}, ""
}

func (s *InventoryReceiptService) GetByCode(ctx *gin.Context, code string) (*model.GetOneInventoryReceiptResponse, string) {
	// Get inventory receipt by code
	receipt, err := s.inventoryReceiptRepository.GetOneByCodeQuery(ctx, code, nil)
	if err != nil {
		log.Error("InventoryReceiptService.GetByCode Error when get receipt: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	if receipt == nil {
		return nil, error_utils.ErrorCode.NOT_FOUND
	}

	// Get receipt items
	receiptItems, err := s.inventoryReceiptItemRepository.GetByInventoryReceiptIDQuery(ctx, receipt.ID, nil)
	if err != nil {
		log.Error("InventoryReceiptService.GetByCode Error when get receipt items: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	// Convert items to response models
	itemResponses := make([]model.InventoryReceiptItemResponse, len(receiptItems))
	for i, item := range receiptItems {
		itemResponses[i] = model.InventoryReceiptItemResponse{
			ID:                 item.ID,
			InventoryReceiptID: item.InventoryReceiptID,
			ProductID:          item.ProductID,
			Quantity:           item.Quantity,
			UnitCost:           item.UnitCost,
			Notes:              item.Notes,
			CreatedAt:          item.CreatedAt,
			UpdatedAt:          item.UpdatedAt,
		}
	}

	// Return response with items
	return &model.GetOneInventoryReceiptResponse{
		InventoryReceipt: model.InventoryReceiptResponse{
			ID:          receipt.ID,
			Code:        receipt.Code,
			UserID:      receipt.UserID,
			ReceiptDate: receipt.ReceiptDate,
			Notes:       receipt.Notes,
			TotalItems:  receipt.TotalItems,
			CreatedAt:   receipt.CreatedAt,
			UpdatedAt:   receipt.UpdatedAt,
			Items:       itemResponses,
		},
	}, ""
}
