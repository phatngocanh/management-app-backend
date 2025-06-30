package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/pna/management-app-backend/internal/controller/http/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func MapRoutes(router *gin.Engine,
	healHandler *HealthHandler,
	helloWorldHandler *HelloWorldHandler,
	userHandler *UserHandler,
	productHandler *ProductHandler,
	productBomHandler *ProductBomHandler,
	productCategoryHandler *ProductCategoryHandler,
	unitOfMeasureHandler *UnitOfMeasureHandler,
	inventoryHandler *InventoryHandler,
	inventoryHistoryHandler *InventoryHistoryHandler,
	customerHandler *CustomerHandler,
	orderImageHandler *OrderImageHandler,
	statisticsHandler *StatisticsHandler,
	authMiddleware *middleware.AuthMiddleware,
) {
	// Apply CORS middleware to all routes
	router.Use(middleware.CorsMiddleware())

	v1 := router.Group("/api/v1")
	{
		health := v1.Group("/health")
		{
			health.GET("", healHandler.Check)
		}
		hello := v1.Group("/hello-world")
		{
			hello.GET("", helloWorldHandler.HelloWorld)
		}
		users := v1.Group("/users")
		{
			users.POST("/login", userHandler.Login)
		}
		products := v1.Group("/products")
		{
			products.POST("", authMiddleware.VerifyAccessToken, productHandler.Create)
			products.PUT("", authMiddleware.VerifyAccessToken, productHandler.Update)
			products.GET("", authMiddleware.VerifyAccessToken, productHandler.GetAll)
			products.GET("/:productId", authMiddleware.VerifyAccessToken, productHandler.GetOne)
			products.GET("/:productId/inventories", authMiddleware.VerifyAccessToken, inventoryHandler.GetByProductID)
			products.PUT("/:productId/inventories/quantity", authMiddleware.VerifyAccessToken, inventoryHandler.UpdateQuantity)
			products.GET("/:productId/inventories/histories", authMiddleware.VerifyAccessToken, inventoryHistoryHandler.GetAll)
		}
		boms := v1.Group("/boms")
		{
			boms.POST("", authMiddleware.VerifyAccessToken, productBomHandler.CreateProductBom)
			boms.PUT("", authMiddleware.VerifyAccessToken, productBomHandler.UpdateProductBom)
			boms.GET("", authMiddleware.VerifyAccessToken, productBomHandler.GetAllProductBoms)
			boms.GET("/parent/:parentProductId", authMiddleware.VerifyAccessToken, productBomHandler.GetProductBomByParentID)
			boms.GET("/component/:componentProductId", authMiddleware.VerifyAccessToken, productBomHandler.GetProductBomsByComponentID)
			boms.DELETE("/parent/:parentProductId", authMiddleware.VerifyAccessToken, productBomHandler.DeleteProductBom)
			boms.POST("/explosion", authMiddleware.VerifyAccessToken, productBomHandler.CalculateMaterialRequirements)
		}
		categories := v1.Group("/categories")
		{
			categories.POST("", authMiddleware.VerifyAccessToken, productCategoryHandler.Create)
			categories.PUT("", authMiddleware.VerifyAccessToken, productCategoryHandler.Update)
			categories.GET("", authMiddleware.VerifyAccessToken, productCategoryHandler.GetAll)
			categories.GET("/:categoryId", authMiddleware.VerifyAccessToken, productCategoryHandler.GetOne)
			categories.GET("/code/:code", authMiddleware.VerifyAccessToken, productCategoryHandler.GetByCode)
		}
		units := v1.Group("/units")
		{
			units.POST("", authMiddleware.VerifyAccessToken, unitOfMeasureHandler.Create)
			units.PUT("", authMiddleware.VerifyAccessToken, unitOfMeasureHandler.Update)
			units.GET("", authMiddleware.VerifyAccessToken, unitOfMeasureHandler.GetAll)
			units.GET("/:unitId", authMiddleware.VerifyAccessToken, unitOfMeasureHandler.GetOne)
			units.GET("/code/:code", authMiddleware.VerifyAccessToken, unitOfMeasureHandler.GetByCode)
		}
		customers := v1.Group("/customers")
		{
			customers.POST("", authMiddleware.VerifyAccessToken, customerHandler.Create)
			customers.PUT("/:customerId", authMiddleware.VerifyAccessToken, customerHandler.Update)
			customers.GET("", authMiddleware.VerifyAccessToken, customerHandler.GetAll)
			customers.GET("/:customerId", authMiddleware.VerifyAccessToken, customerHandler.GetOne)
		}
		inventory := v1.Group("/inventory")
		{
			inventory.GET("", authMiddleware.VerifyAccessToken, inventoryHandler.GetAll)
		}
		statistics := v1.Group("/statistics")
		{
			statistics.GET("/dashboard", authMiddleware.VerifyAccessToken, statisticsHandler.GetDashboardStats)
		}
	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
