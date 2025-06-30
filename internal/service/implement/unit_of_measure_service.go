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

type UnitOfMeasureService struct {
	unitRepository repository.UnitOfMeasureRepository
	unitOfWork     repository.UnitOfWork
}

func NewUnitOfMeasureService(unitRepository repository.UnitOfMeasureRepository, unitOfWork repository.UnitOfWork) service.UnitOfMeasureService {
	return &UnitOfMeasureService{
		unitRepository: unitRepository,
		unitOfWork:     unitOfWork,
	}
}

func (s *UnitOfMeasureService) Create(ctx *gin.Context, request model.CreateUnitOfMeasureRequest) (*model.UnitOfMeasureResponse, string) {
	// Create unit entity
	unit := &entity.UnitOfMeasure{
		Name:        request.Name,
		Code:        request.Code,
		Description: request.Description,
	}

	// Save to database
	err := s.unitRepository.CreateCommand(ctx, unit, nil)
	if err != nil {
		log.Error("UnitOfMeasureService.Create Error when create unit: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	return &model.UnitOfMeasureResponse{
		ID:          unit.ID,
		Name:        unit.Name,
		Code:        unit.Code,
		Description: unit.Description,
	}, ""
}

func (s *UnitOfMeasureService) Update(ctx *gin.Context, request model.UpdateUnitOfMeasureRequest) (*model.UnitOfMeasureResponse, string) {
	// Check if unit exists
	existingUnit, err := s.unitRepository.GetOneByIDQuery(ctx, request.ID, nil)
	if err != nil {
		log.Error("UnitOfMeasureService.Update Error when get unit: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	if existingUnit == nil {
		return nil, error_utils.ErrorCode.NOT_FOUND
	}

	// Update unit entity
	unit := &entity.UnitOfMeasure{
		ID:          request.ID,
		Name:        request.Name,
		Code:        request.Code,
		Description: request.Description,
	}

	// Save to database
	err = s.unitRepository.UpdateCommand(ctx, unit, nil)
	if err != nil {
		log.Error("UnitOfMeasureService.Update Error when update unit: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	return &model.UnitOfMeasureResponse{
		ID:          unit.ID,
		Name:        unit.Name,
		Code:        unit.Code,
		Description: unit.Description,
	}, ""
}

func (s *UnitOfMeasureService) GetAll(ctx *gin.Context) (*model.GetAllUnitsOfMeasureResponse, string) {
	// Get all units
	units, err := s.unitRepository.GetAllQuery(ctx, nil)
	if err != nil {
		log.Error("UnitOfMeasureService.GetAll Error when get units: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	// Convert to response models
	unitResponses := make([]model.UnitOfMeasureResponse, len(units))
	for i, unit := range units {
		unitResponses[i] = model.UnitOfMeasureResponse{
			ID:          unit.ID,
			Name:        unit.Name,
			Code:        unit.Code,
			Description: unit.Description,
		}
	}

	return &model.GetAllUnitsOfMeasureResponse{
		Units: unitResponses,
	}, ""
}

func (s *UnitOfMeasureService) GetOne(ctx *gin.Context, id int) (*model.GetOneUnitOfMeasureResponse, string) {
	// Get unit by ID
	unit, err := s.unitRepository.GetOneByIDQuery(ctx, id, nil)
	if err != nil {
		log.Error("UnitOfMeasureService.GetOne Error when get unit: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	if unit == nil {
		return nil, error_utils.ErrorCode.NOT_FOUND
	}

	return &model.GetOneUnitOfMeasureResponse{
		Unit: model.UnitOfMeasureResponse{
			ID:          unit.ID,
			Name:        unit.Name,
			Code:        unit.Code,
			Description: unit.Description,
		},
	}, ""
}

func (s *UnitOfMeasureService) GetByCode(ctx *gin.Context, code string) (*model.GetOneUnitOfMeasureResponse, string) {
	// Get unit by code
	unit, err := s.unitRepository.GetOneByCodeQuery(ctx, code, nil)
	if err != nil {
		log.Error("UnitOfMeasureService.GetByCode Error when get unit: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	if unit == nil {
		return nil, error_utils.ErrorCode.NOT_FOUND
	}

	return &model.GetOneUnitOfMeasureResponse{
		Unit: model.UnitOfMeasureResponse{
			ID:          unit.ID,
			Name:        unit.Name,
			Code:        unit.Code,
			Description: unit.Description,
		},
	}, ""
}
