package http

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/pna/management-app-backend/internal/controller/http/middleware"

	"github.com/gin-gonic/gin"

	v1 "github.com/pna/management-app-backend/internal/controller/http/v1"
)

type Server struct {
	healthHandler           *v1.HealthHandler
	helloWorldHandler       *v1.HelloWorldHandler
	authMiddleware          *middleware.AuthMiddleware
	userHandler             *v1.UserHandler
	productHandler          *v1.ProductHandler
	productBomHandler       *v1.ProductBomHandler
	productCategoryHandler  *v1.ProductCategoryHandler
	unitOfMeasureHandler    *v1.UnitOfMeasureHandler
	inventoryHandler        *v1.InventoryHandler
	inventoryHistoryHandler *v1.InventoryHistoryHandler
	inventoryReceiptHandler *v1.InventoryReceiptHandler
	customerHandler         *v1.CustomerHandler
	statisticsHandler       *v1.StatisticsHandler
	productImageHandler     *v1.ProductImageHandler
}

func NewServer(
	healthHandler *v1.HealthHandler,
	helloWorldHandler *v1.HelloWorldHandler,
	authMiddleware *middleware.AuthMiddleware,
	userHandler *v1.UserHandler,
	productHandler *v1.ProductHandler,
	productBomHandler *v1.ProductBomHandler,
	productCategoryHandler *v1.ProductCategoryHandler,
	unitOfMeasureHandler *v1.UnitOfMeasureHandler,
	inventoryHandler *v1.InventoryHandler,
	inventoryHistoryHandler *v1.InventoryHistoryHandler,
	inventoryReceiptHandler *v1.InventoryReceiptHandler,
	customerHandler *v1.CustomerHandler,
	statisticsHandler *v1.StatisticsHandler,
	productImageHandler *v1.ProductImageHandler,
) *Server {
	return &Server{
		healthHandler:           healthHandler,
		helloWorldHandler:       helloWorldHandler,
		authMiddleware:          authMiddleware,
		userHandler:             userHandler,
		productHandler:          productHandler,
		productBomHandler:       productBomHandler,
		productCategoryHandler:  productCategoryHandler,
		unitOfMeasureHandler:    unitOfMeasureHandler,
		inventoryHandler:        inventoryHandler,
		inventoryHistoryHandler: inventoryHistoryHandler,
		inventoryReceiptHandler: inventoryReceiptHandler,
		customerHandler:         customerHandler,
		statisticsHandler:       statisticsHandler,
		productImageHandler:     productImageHandler,
	}
}

func (s *Server) Run() {
	router := gin.New()
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	httpServerInstance := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}
	fmt.Println("Server running at " + httpServerInstance.Addr)

	v1.MapRoutes(
		router,
		s.healthHandler,
		s.helloWorldHandler,
		s.userHandler,
		s.productHandler,
		s.productBomHandler,
		s.productCategoryHandler,
		s.unitOfMeasureHandler,
		s.inventoryHandler,
		s.inventoryHistoryHandler,
		s.inventoryReceiptHandler,
		s.customerHandler,
		s.statisticsHandler,
		s.productImageHandler,
		s.authMiddleware,
	)
	err := httpServerInstance.ListenAndServe()
	if err != nil {
		fmt.Println("There is error: " + err.Error())
		return
	}
}
