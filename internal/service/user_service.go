package service

import (
	"github.com/gin-gonic/gin"
	"github.com/pna/management-app-backend/internal/domain/model"
)

type UserService interface {
	Login(ctx *gin.Context, request model.LoginRequest) (*model.LoginResponse, string)
}
