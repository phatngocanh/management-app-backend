package serviceimplement

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pna/management-app-backend/internal/bean"
	"github.com/pna/management-app-backend/internal/domain/entity"
	"github.com/pna/management-app-backend/internal/domain/model"
	"github.com/pna/management-app-backend/internal/repository"
	"github.com/pna/management-app-backend/internal/service"
	"github.com/pna/management-app-backend/internal/utils/error_utils"
	log "github.com/sirupsen/logrus"
)

type OrderService struct {
	orderRepo            repository.OrderRepository
	orderItemRepo        repository.OrderItemRepository
	inventoryRepo        repository.InventoryRepository
	inventoryHistoryRepo repository.InventoryHistoryRepository
	productRepo          repository.ProductRepository
	bomRepo              repository.ProductBomRepository
	unitOfWork           repository.UnitOfWork
	userRepo             repository.UserRepository
	customerRepo         repository.CustomerRepository
	orderImageRepo       repository.OrderImageRepository
	s3Service            bean.S3Service
}

func NewOrderService(
	orderRepo repository.OrderRepository,
	inventoryRepo repository.InventoryRepository,
	inventoryHistoryRepo repository.InventoryHistoryRepository,
	orderItemRepo repository.OrderItemRepository,
	productRepo repository.ProductRepository,
	bomRepo repository.ProductBomRepository,
	unitOfWork repository.UnitOfWork,
	userRepo repository.UserRepository,
	orderImageRepo repository.OrderImageRepository,
	s3Service bean.S3Service,
	customerRepo repository.CustomerRepository,
) service.OrderService {
	return &OrderService{
		orderRepo:            orderRepo,
		inventoryRepo:        inventoryRepo,
		inventoryHistoryRepo: inventoryHistoryRepo,
		orderItemRepo:        orderItemRepo,
		productRepo:          productRepo,
		bomRepo:              bomRepo,
		unitOfWork:           unitOfWork,
		userRepo:             userRepo,
		customerRepo:         customerRepo,
		orderImageRepo:       orderImageRepo,
		s3Service:            s3Service,
	}
}

// RequiredMaterial represents the required quantity of a raw material
type RequiredMaterial struct {
	ProductID int
	Quantity  int
}

// calculateRequiredMaterialsForProduct calculates the required raw materials for a direct product order
func (s *OrderService) calculateRequiredMaterialsForProduct(ctx *gin.Context, productID int, quantity int, tx *sqlx.Tx) ([]RequiredMaterial, error) {
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

func (s *OrderService) CreateOrder(ctx *gin.Context, orderRequest model.CreateOrderRequest, userId int) (*model.OrderResponse, string) {
	// Begin transaction
	tx, err := s.unitOfWork.Begin(ctx)
	if err != nil {
		log.Error("OrderService.CreateOrder Error when begin transaction: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	// Defer rollback in case of error
	defer func() {
		if rollbackErr := s.unitOfWork.Rollback(tx); rollbackErr != nil {
			log.Error("OrderService.CreateOrder Error when rollback transaction: " + rollbackErr.Error())
		}
	}()

	// Calculate total required materials from all order items
	requiredMaterials := make(map[int]int) // productID -> total quantity needed

	for _, item := range orderRequest.Items {
		var materials []RequiredMaterial

		// Product order (can be direct product or parent product from BOM)
		materials, err = s.calculateRequiredMaterialsForProduct(ctx, item.ProductID, item.Quantity, tx)

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
		if inventory.Quantity < requiredQty {
			log.Error(fmt.Sprintf("OrderService.CreateOrder Error: insufficient inventory for product ID %d: required %d, available %d",
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

	user, err := s.userRepo.FindByIDQuery(ctx, userId, tx)
	if err != nil {
		log.Error("OrderService.CreateOrder Error when get user: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	// Deduct inventory for all required materials
	for productID, requiredQty := range requiredMaterials {
		inventory := inventoryMap[productID]
		newQuantity := inventory.Quantity - int(requiredQty)

		uuid := uuid.New()
		err = s.inventoryRepo.UpdateQuantityCommand(ctx, productID, newQuantity, uuid.String(), tx)
		if err != nil {
			log.Error(fmt.Sprintf("OrderService.CreateOrder Error when update inventory for product ID %d: %s", productID, err.Error()))
			return nil, error_utils.ErrorCode.DB_DOWN
		}

		inventoryHistory := &entity.InventoryHistory{
			ProductID:     productID,
			Quantity:      -requiredQty,
			FinalQuantity: newQuantity,
			ImporterName:  user.Username,
			ImportedAt:    time.Now(),
			Note:          "Xuất cho đơn hàng: " + order.Code,
			ReferenceID:   &order.ID,
		}

		err = s.inventoryHistoryRepo.CreateCommand(ctx, inventoryHistory, tx)
		if err != nil {
			log.Error("OrderService.CreateOrder Error when create inventory history: " + err.Error())
			return nil, error_utils.ErrorCode.DB_DOWN
		}

	}

	// Commit transaction
	err = s.unitOfWork.Commit(tx)
	if err != nil {
		log.Error("OrderService.CreateOrder Error when commit transaction: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	customer, err := s.customerRepo.GetOneByIDQuery(ctx, order.CustomerID, tx)
	if err != nil {
		log.Error("OrderService.CreateOrder Error when get customer: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	// Return response
	return &model.OrderResponse{
		ID:        order.ID,
		Code:      order.Code,
		OrderDate: order.OrderDate,
		Customer: model.CustomerResponse{
			ID:      order.CustomerID,
			Name:    customer.Name,
			Phone:   customer.Phone,
			Address: customer.Address,
		},
		Note:               order.Note,
		TotalSalesRevenue:  order.TotalSalesRevenue,
		AdditionalCost:     order.AdditionalCost,
		AdditionalCostNote: order.AdditionalCostNote,
		TaxPercent:         order.TaxPercent,
	}, ""
}

// Helper to calculate total amount and product count from order items
func calculateOrderAmountsAndProductCount(orderItems []entity.OrderItem) (totalAmount int, productCount int) {
	productIDSet := make(map[int]struct{})
	totalAmount = 0
	for _, item := range orderItems {
		finalAmount := item.FinalAmount
		if finalAmount == 0 {
			itemTotal := item.Quantity * item.SellingPrice
			discountAmount := (itemTotal * item.DiscountPercent) / 100
			calculated := itemTotal - discountAmount
			finalAmount = calculated
		}
		totalAmount += finalAmount
		productIDSet[item.ProductID] = struct{}{}
	}
	productCount = len(productIDSet)
	return
}

func (s *OrderService) GetOneOrder(ctx *gin.Context, id int) (model.GetOneOrderResponse, string) {
	order, err := s.orderRepo.GetOneByIDQuery(ctx, id, nil)
	if err != nil {
		log.Error("OrderService.GetOne Error: " + err.Error())
		return model.GetOneOrderResponse{}, error_utils.ErrorCode.DB_DOWN
	}
	if order == nil {
		return model.GetOneOrderResponse{}, error_utils.ErrorCode.NOT_FOUND
	}

	// Fetch customer information
	customer, err := s.customerRepo.GetOneByIDQuery(ctx, order.CustomerID, nil)
	if err != nil {
		log.Error("OrderService.GetOne Error fetching customer: " + err.Error())
		return model.GetOneOrderResponse{}, error_utils.ErrorCode.DB_DOWN
	}

	// Fetch order items
	orderItems, err := s.orderItemRepo.GetAllByOrderIDQuery(ctx, order.ID, nil)
	if err != nil {
		log.Error("OrderService.GetOne Error fetching order items: " + err.Error())
		return model.GetOneOrderResponse{}, error_utils.ErrorCode.DB_DOWN
	}

	orderItemResponses := make([]model.OrderItemResponse, 0, len(orderItems))
	totalOriginalCost := 0
	totalProfitLoss := 0

	for _, item := range orderItems {
		product, err := s.productRepo.GetOneByIDQuery(ctx, item.ProductID, nil)
		if err != nil {
			log.Error("OrderService.GetOne Error fetching product: " + err.Error())
			return model.GetOneOrderResponse{}, error_utils.ErrorCode.DB_DOWN
		}

		finalAmount := item.FinalAmount
		if finalAmount == 0 {
			itemTotal := item.Quantity * item.SellingPrice
			discountAmount := (itemTotal * item.DiscountPercent) / 100
			calculated := itemTotal - discountAmount
			finalAmount = calculated
		}

		// Calculate profit/loss for this item
		originalCost := item.Quantity * item.OriginalPrice
		sellingRevenue := item.Quantity * item.SellingPrice
		discountAmount := (sellingRevenue * item.DiscountPercent) / 100
		finalRevenue := sellingRevenue - discountAmount
		profitLoss := finalRevenue - originalCost
		profitLossPercentage := 0.0
		if originalCost > 0 {
			profitLossPercentage = float64(profitLoss) / float64(originalCost) * 100
		}

		// Accumulate totals
		totalOriginalCost += originalCost
		totalProfitLoss += profitLoss

		orderItemResponses = append(orderItemResponses, model.OrderItemResponse{
			ID:              item.ID,
			OrderID:         item.OrderID,
			ProductID:       item.ProductID,
			ProductName:     product.Name,
			Quantity:        item.Quantity,
			SellingPrice:    item.SellingPrice,
			DiscountPercent: item.DiscountPercent,
			FinalAmount:     &finalAmount,
			// Profit/Loss fields
			OriginalPrice:        &item.OriginalPrice,
			ProfitLoss:           &profitLoss,
			ProfitLossPercentage: &profitLossPercentage,
		})
	}

	totalAmount, productCount := calculateOrderAmountsAndProductCount(orderItems)
	totalAmount += order.AdditionalCost
	totalAmount += int(float64(totalAmount) * float64(order.TaxPercent) / 100)

	// Use stored values for total order profit/loss
	totalProfitLoss = order.TotalSalesRevenue - order.TotalOriginalCost + order.AdditionalCost
	totalProfitLossPercentage := 0.0
	if order.TotalOriginalCost > 0 {
		totalProfitLossPercentage = float64(totalProfitLoss) / float64(order.TotalOriginalCost) * 100
	}

	// Fetch order images and generate signed URLs
	orderImages, err := s.orderImageRepo.GetAllByOrderIDQuery(ctx, order.ID, nil)
	if err != nil {
		log.Error("OrderService.GetOne Error fetching order images: " + err.Error())
		// Continue without images rather than failing the entire request
		orderImages = make([]entity.OrderImage, 0)
	}

	// Convert images to response model with signed URLs
	var imageResponses []model.OrderImage
	if len(orderImages) > 0 {
		for _, img := range orderImages {
			// Generate a fresh signed URL for each image
			signedURL, err := s.s3Service.GenerateSignedDownloadURL(ctx, img.S3Key, 20*time.Second)
			if err != nil {
				log.Error("OrderService.GetOne Error generating signed URL for image: " + err.Error())
				// Continue with other images even if one fails
				signedURL = ""
			}

			imageResponses = append(imageResponses, model.OrderImage{
				ID:       img.ID,
				OrderID:  img.OrderID,
				ImageURL: signedURL,
				ImageKey: img.S3Key,
			})
		}
	}

	resp := model.GetOneOrderResponse{Order: model.OrderResponse{
		ID:                 order.ID,
		OrderDate:          order.OrderDate,
		Note:               order.Note,
		AdditionalCost:     order.AdditionalCost,
		AdditionalCostNote: order.AdditionalCostNote,
		Customer: model.CustomerResponse{
			ID:      customer.ID,
			Name:    customer.Name,
			Phone:   customer.Phone,
			Address: customer.Address,
		},
		OrderItems:   orderItemResponses,
		Images:       imageResponses,
		TotalAmount:  &totalAmount,
		ProductCount: &productCount,
		TaxPercent:   order.TaxPercent,
		// Profit/Loss fields for total order
		TotalProfitLoss:           &totalProfitLoss,
		TotalProfitLossPercentage: &totalProfitLossPercentage,
		TotalSalesRevenue:         order.TotalSalesRevenue,
	}}
	return resp, ""
}

func (s *OrderService) Update(ctx context.Context, req model.UpdateOrderRequest) string {
	existing, err := s.orderRepo.GetOneByIDQuery(ctx, req.ID, nil)
	if err != nil {
		log.Error("OrderService.Update Error: " + err.Error())
		return error_utils.ErrorCode.DB_DOWN
	}
	if existing == nil {
		return error_utils.ErrorCode.NOT_FOUND
	}

	if req.CustomerID != 0 {
		existing.CustomerID = req.CustomerID
	}
	if !req.OrderDate.IsZero() {
		existing.OrderDate = req.OrderDate
	}
	if req.Note != nil {
		existing.Note = req.Note
	}
	if req.AdditionalCost != nil {
		existing.AdditionalCost = *req.AdditionalCost
	}
	if req.AdditionalCostNote != nil {
		existing.AdditionalCostNote = req.AdditionalCostNote
	}
	if req.TaxPercent != nil {
		existing.TaxPercent = *req.TaxPercent
	}

	err = s.orderRepo.UpdateCommand(ctx, existing, nil)
	if err != nil {
		log.Error("OrderService.Update Error when update order: " + err.Error())
		return error_utils.ErrorCode.DB_DOWN
	}

	return ""
}

func (s *OrderService) GetAll(ctx context.Context, userID int, customerID int, sortBy string) (model.GetAllOrdersResponse, string) {
	orders, err := s.orderRepo.GetAllWithFiltersQuery(ctx, customerID, sortBy, nil)
	if err != nil {
		log.Error("OrderService.GetAll Error: " + err.Error())
		return model.GetAllOrdersResponse{}, error_utils.ErrorCode.DB_DOWN
	}

	resp := model.GetAllOrdersResponse{Orders: make([]model.OrderResponse, 0, len(orders))}
	for _, o := range orders {
		// Fetch customer information
		customer, err := s.customerRepo.GetOneByIDQuery(ctx, o.CustomerID, nil)
		if err != nil {
			log.Error("OrderService.GetAll Error fetching customer: " + err.Error())
			continue
		}

		// Fetch order items to calculate total amount and product count
		orderItems, err := s.orderItemRepo.GetAllByOrderIDQuery(ctx, o.ID, nil)
		if err != nil {
			log.Error("OrderService.GetAll Error fetching order items: " + err.Error())
			continue
		}
		totalAmount, productCount := calculateOrderAmountsAndProductCount(orderItems)
		totalAmount += o.AdditionalCost
		totalAmount += int(float64(totalAmount) * float64(o.TaxPercent) / 100)
		// Calculate profit/loss from stored cost and revenue values
		totalProfitLoss := o.TotalSalesRevenue - o.TotalOriginalCost + o.AdditionalCost
		totalProfitLossPercentage := 0.0
		if o.TotalOriginalCost > 0 {
			totalProfitLossPercentage = float64(totalProfitLoss) / float64(o.TotalOriginalCost) * 100
		}

		resp.Orders = append(resp.Orders, model.OrderResponse{
			ID:                 o.ID,
			OrderDate:          o.OrderDate,
			Note:               o.Note,
			AdditionalCost:     o.AdditionalCost,
			AdditionalCostNote: o.AdditionalCostNote,
			Customer: model.CustomerResponse{
				ID:      customer.ID,
				Name:    customer.Name,
				Phone:   customer.Phone,
				Address: customer.Address,
			},
			OrderItems:                nil, // Omit order items in GetAll
			TaxPercent:                o.TaxPercent,
			TotalAmount:               &totalAmount,
			ProductCount:              &productCount,
			TotalProfitLoss:           &totalProfitLoss,
			TotalProfitLossPercentage: &totalProfitLossPercentage,
		})
	}
	return resp, ""
}
