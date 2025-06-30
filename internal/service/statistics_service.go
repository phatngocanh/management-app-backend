package service

import (
	"context"

	"github.com/pna/management-app-backend/internal/domain/model"
)

type StatisticsService interface {
	GetDashboardStats(ctx context.Context) (model.DashboardStatsResponse, string)
}
