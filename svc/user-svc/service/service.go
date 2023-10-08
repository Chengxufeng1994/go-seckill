package service

import (
	"context"
	"errors"
	"github.com/Chengxufeng1994/go-seckill/svc/user-svc/entity"
	"github.com/Chengxufeng1994/go-seckill/svc/user-svc/model"
	"log"
	"strings"
)

// Service Define a service interface
type Service interface {
	Check(ctx context.Context, username, password string) (int64, error)

	// HealthCheck check service health status
	HealthCheck() bool
}

// UserService implement Service interface
type userService struct {
	repo entity.Repository
}

func New(repo entity.Repository) Service {
	return &userService{
		repo: repo,
	}
}

// Check implement Service method
func (svc *userService) Check(ctx context.Context, username, password string) (int64, error) {
	userDao, err := svc.repo.GetUserByUsername(username)
	userDto := model.UserDao2Dto(userDao)

	if err != nil {
		log.Printf("Repository.GetUserByUsername, err : %v", err)
		return 0, err
	}

	if !strings.EqualFold(password, userDto.Password) {
		return 0, errors.New("password invalid")
	}

	return int64(userDto.ID), nil
}

// HealthCheck implement Service method
func (svc *userService) HealthCheck() bool {
	return true
}

// ServiceMiddleware define service middleware
type ServiceMiddleware func(Service) Service
