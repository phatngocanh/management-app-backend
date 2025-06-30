package service

import (
	"context"
	"github.com/pna/management-app-backend/internal/domain/model"
)

type HelloWorldService interface {
	HelloWorld(ctx context.Context) (*model.HelloWorldResponse, string)
}
