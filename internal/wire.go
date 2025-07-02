//go:build wireinject
// +build wireinject

package internal

import (
	"github.com/google/wire"
	beanimplement "github.com/pna/management-app-backend/internal/bean/implement"
	"github.com/pna/management-app-backend/internal/controller"
	"github.com/pna/management-app-backend/internal/controller/http"
	"github.com/pna/management-app-backend/internal/controller/http/middleware"
	v1 "github.com/pna/management-app-backend/internal/controller/http/v1"
	"github.com/pna/management-app-backend/internal/database"
	repositoryimplement "github.com/pna/management-app-backend/internal/repository/implement"
	serviceimplement "github.com/pna/management-app-backend/internal/service/implement"
)

var container = wire.NewSet(
	controller.NewApiContainer,
)

// may have grpc server in the future
var serverSet = wire.NewSet(
	http.NewServer,
)

// handler === controller | with service and repository layers to form 3 layers architecture
var handlerSet = wire.NewSet(
	v1.NewHealthHandler,
	v1.NewHelloWorldHandler,
	v1.NewUserHandler,
	v1.NewProductHandler,
	v1.NewProductBomHandler,
	v1.NewProductCategoryHandler,
	v1.NewUnitOfMeasureHandler,
	v1.NewInventoryHandler,
	v1.NewInventoryHistoryHandler,
	v1.NewCustomerHandler,
	v1.NewStatisticsHandler,
	v1.NewInventoryReceiptHandler,
	v1.NewProductImageHandler,
	v1.NewOrderHandler,
)

var serviceSet = wire.NewSet(
	serviceimplement.NewHelloWorldService,
	serviceimplement.NewUserService,
	serviceimplement.NewProductService,
	serviceimplement.NewInventoryService,
	serviceimplement.NewInventoryHistoryService,
	serviceimplement.NewCustomerService,
	serviceimplement.NewStatisticsService,
	serviceimplement.NewUnitOfMeasureService,
	serviceimplement.NewProductCategoryService,
	serviceimplement.NewProductImageService,
	serviceimplement.NewProductBomService,
	serviceimplement.NewInventoryReceiptService,
	serviceimplement.NewOrderService,
)

var repositorySet = wire.NewSet(
	repositoryimplement.NewHelloWorldRepository,
	repositoryimplement.NewUserRepository,
	repositoryimplement.NewProductRepository,
	repositoryimplement.NewInventoryRepository,
	repositoryimplement.NewInventoryHistoryRepository,
	repositoryimplement.NewUnitOfWork,
	repositoryimplement.NewCustomerRepository,
	repositoryimplement.NewUnitOfMeasureRepository,
	repositoryimplement.NewProductCategoryRepository,
	repositoryimplement.NewProductImageRepository,
	repositoryimplement.NewProductBomRepository,
	repositoryimplement.NewInventoryReceiptRepository,
	repositoryimplement.NewInventoryReceiptItemRepository,
	repositoryimplement.NewOrderRepository,
	repositoryimplement.NewOrderItemRepository,
)

var middlewareSet = wire.NewSet(
	middleware.NewAuthMiddleware,
)

var beanSet = wire.NewSet(
	beanimplement.NewBcryptPasswordEncoder,
	beanimplement.NewS3Service,
)

func InitializeContainer(
	db database.Db,
) *controller.ApiContainer {
	wire.Build(serverSet, handlerSet, serviceSet, repositorySet, middlewareSet, beanSet, container)
	return &controller.ApiContainer{}
}
