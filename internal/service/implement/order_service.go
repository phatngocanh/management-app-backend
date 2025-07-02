package serviceimplement

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/pna/management-app-backend/internal/domain/entity"
	"github.com/pna/management-app-backend/internal/domain/model"
	"github.com/pna/management-app-backend/internal/repository"
	"github.com/pna/management-app-backend/internal/service"
	"github.com/pna/management-app-backend/internal/utils/error_utils"
	log "github.com/sirupsen/logrus"
)

type OrderService struct {
	orderRepo     repository.OrderRepository
	orderItemRepo repository.OrderItemRepository
	inventoryRepo repository.InventoryRepository
	productRepo   repository.ProductRepository
	bomRepo       repository.ProductBomRepository
	unitOfWork    repository.UnitOfWork
}

func NewOrderService(
	orderRepo repository.OrderRepository,
	inventoryRepo repository.InventoryRepository,
	orderItemRepo repository.OrderItemRepository,
	productRepo repository.ProductRepository,
	bomRepo repository.ProductBomRepository,
	unitOfWork repository.UnitOfWork,
) service.OrderService {
	return &OrderService{
		orderRepo:     orderRepo,
		inventoryRepo: inventoryRepo,
		orderItemRepo: orderItemRepo,
		productRepo:   productRepo,
		bomRepo:       bomRepo,
		unitOfWork:    unitOfWork,
	}
}

// RequiredMaterial represents the required quantity of a raw material
type RequiredMaterial struct {
	ProductID int
	Quantity  float64
}

func (s *OrderService) CreateOrder(ctx *gin.Context, orderRequest model.CreateOrderRequest) (*model.OrderResponse, string) {
	// Begin transaction
	tx, err := s.unitOfWork.Begin(ctx)
	if err != nil {
		log.Error("OrderService.CreateOrder Error when begin transaction: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	// Defer rollback in case of error
	defer func() {
		if err != nil {
			if rollbackErr := s.unitOfWork.Rollback(tx); rollbackErr != nil {
				log.Error("OrderService.CreateOrder Error when rollback transaction: " + rollbackErr.Error())
			}
		}
	}()

	// Calculate total required materials from all order items
	requiredMaterials := make(map[int]float64) // productID -> total quantity needed

	for _, item := range orderRequest.Items {
		var materials []RequiredMaterial

		if item.ProductID != nil {
			// Direct product order
			materials, err = s.calculateRequiredMaterialsForProduct(ctx, *item.ProductID, float64(item.Quantity), tx)
		} else if item.BomParentID != nil {
			// BOM parent product order - directly calculate materials for the parent product
			materials, err = s.calculateRequiredMaterialsForProduct(ctx, *item.BomParentID, float64(item.Quantity), tx)
		} else {
			log.Error("OrderService.CreateOrder Error: order item must have either product_id or bom_parent_id")
			return nil, error_utils.ErrorCode.BAD_REQUEST
		}

		if err != nil {
			log.Error("OrderService.CreateOrder Error calculating required materials: " + err.Error())
			return nil, error_utils.ErrorCode.DB_DOWN
		}

		// Aggregate the required materials
		for _, material := range materials {
			requiredMaterials[material.ProductID] += material.Quantity
		}
	}

	// Get all required product IDs for inventory checking
	var productIDs []int
	for productID := range requiredMaterials {
		productIDs = append(productIDs, productID)
	}

	fmt.Println("requiredMaterials", requiredMaterials)

	return nil, error_utils.ErrorCode.DB_DOWN

	// Lock the inventories to prevent concurrent access
	inventories, err := s.inventoryRepo.SelectManyForUpdate(ctx, productIDs, tx)
	if err != nil {
		log.Error("OrderService.CreateOrder Error when lock inventories: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	// Check if we have enough inventory for all required materials
	inventoryMap := make(map[int]*entity.Inventory)
	for i := range inventories {
		inventoryMap[inventories[i].ProductID] = &inventories[i]
	}

	for productID, requiredQty := range requiredMaterials {
		inventory, exists := inventoryMap[productID]
		if !exists {
			log.Error(fmt.Sprintf("OrderService.CreateOrder Error: inventory not found for product ID %d", productID))
			return nil, error_utils.ErrorCode.NOT_FOUND
		}
		if float64(inventory.Quantity) < requiredQty {
			log.Error(fmt.Sprintf("OrderService.CreateOrder Error: insufficient inventory for product ID %d: required %.3f, available %d",
				productID, requiredQty, inventory.Quantity))
			return nil, error_utils.ErrorCode.INVENTORY_QUANTITY_EXCEEDED
		}
	}

	// Create the order
	order := &entity.Order{
		CustomerID:         orderRequest.CustomerID,
		OrderDate:          orderRequest.OrderDate,
		Note:               orderRequest.Note,
		TotalOriginalCost:  0, // Will be calculated
		TotalSalesRevenue:  0, // Will be calculated
		AdditionalCost:     orderRequest.AdditionalCost,
		AdditionalCostNote: orderRequest.AdditionalCostNote,
		TaxPercent:         orderRequest.TaxPercent,
	}

	err = s.orderRepo.CreateCommand(ctx, order, tx)
	if err != nil {
		log.Error("OrderService.CreateOrder Error when create order: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	// Calculate totals and create order items
	var totalOriginalCost, totalSalesRevenue int
	for _, itemRequest := range orderRequest.Items {
		finalAmount := s.calculateFinalAmount(itemRequest.SellingPrice, itemRequest.Quantity, itemRequest.DiscountPercent)

		orderItem := &entity.OrderItem{
			OrderID:         order.ID,
			ProductID:       itemRequest.ProductID,
			BomParentID:     itemRequest.BomParentID,
			Quantity:        itemRequest.Quantity,
			SellingPrice:    itemRequest.SellingPrice,
			OriginalPrice:   itemRequest.OriginalPrice,
			DiscountPercent: itemRequest.DiscountPercent,
			FinalAmount:     finalAmount,
		}

		err = s.orderItemRepo.CreateCommand(ctx, orderItem, tx)
		if err != nil {
			log.Error("OrderService.CreateOrder Error when create order item: " + err.Error())
			return nil, error_utils.ErrorCode.DB_DOWN
		}

		totalOriginalCost += itemRequest.OriginalPrice * itemRequest.Quantity
		totalSalesRevenue += finalAmount
	}

	// Update order totals
	order.TotalOriginalCost = totalOriginalCost
	order.TotalSalesRevenue = totalSalesRevenue
	err = s.orderRepo.UpdateCommand(ctx, order, tx)
	if err != nil {
		log.Error("OrderService.CreateOrder Error when update order totals: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	// Deduct inventory for all required materials
	for productID, requiredQty := range requiredMaterials {
		inventory := inventoryMap[productID]
		newQuantity := inventory.Quantity - int(requiredQty)

		err = s.inventoryRepo.UpdateQuantityCommand(ctx, productID, newQuantity, inventory.Version, tx)
		if err != nil {
			log.Error(fmt.Sprintf("OrderService.CreateOrder Error when update inventory for product ID %d: %s", productID, err.Error()))
			return nil, error_utils.ErrorCode.DB_DOWN
		}
	}

	// Commit transaction
	err = s.unitOfWork.Commit(tx)
	if err != nil {
		log.Error("OrderService.CreateOrder Error when commit transaction: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	// Return response
	return &model.OrderResponse{
		ID:                 order.ID,
		Code:               order.Code,
		CustomerID:         order.CustomerID,
		OrderDate:          order.OrderDate,
		Note:               order.Note,
		TotalOriginalCost:  order.TotalOriginalCost,
		TotalSalesRevenue:  order.TotalSalesRevenue,
		AdditionalCost:     order.AdditionalCost,
		AdditionalCostNote: order.AdditionalCostNote,
		TaxPercent:         order.TaxPercent,
	}, ""
}

// calculateRequiredMaterialsForProduct calculates the required raw materials for a direct product order
func (s *OrderService) calculateRequiredMaterialsForProduct(ctx *gin.Context, productID int, quantity float64, tx *sqlx.Tx) ([]RequiredMaterial, error) {
	product, err := s.productRepo.GetOneByIDQuery(ctx, productID, tx)
	if err != nil || product == nil {
		return nil, fmt.Errorf("failed to get product %d: %w", productID, err)
	}

	// If it's a PURCHASE type product, it's a raw material itself
	if product.OperationType == "PURCHASE" {
		return []RequiredMaterial{{ProductID: productID, Quantity: quantity}}, nil
	}

	// For PACKAGING and MANUFACTURING, expand their BOMs
	boms, err := s.bomRepo.GetByParentProductIDQuery(ctx, productID, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to get BOMs for product %d: %w", productID, err)
	}

	var allMaterials []RequiredMaterial
	for _, bom := range boms {
		// Recursively calculate materials for each component
		componentMaterials, err := s.calculateRequiredMaterialsForProduct(ctx, bom.ComponentProductID, bom.Quantity*quantity, tx)
		if err != nil {
			return nil, err
		}
		allMaterials = append(allMaterials, componentMaterials...)
	}

	return allMaterials, nil
}

// calculateFinalAmount calculates the final amount after applying discount
func (s *OrderService) calculateFinalAmount(sellingPrice, quantity, discountPercent int) int {
	subtotal := sellingPrice * quantity
	discount := (subtotal * discountPercent) / 100
	return subtotal - discount
}
