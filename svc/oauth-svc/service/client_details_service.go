package service

import (
	"context"
	"errors"
	"github.com/Chengxufeng1994/go-seckill/svc/oauth-svc/entity"
	"github.com/Chengxufeng1994/go-seckill/svc/oauth-svc/model"
	"strings"
)

var (
	ErrClientMessage = errors.New("invalid client")
)

// ClientDetailsService Define a service interface
type ClientDetailsService interface {
	GetClientDetailByClientId(ctx context.Context, clientId string, clientSecret string) (*model.ClientDetails, error)
}

type ClientDetailsServiceImpl struct {
	repo entity.ClientDetailsRepository
}

func NewClientDetailsService(repo entity.ClientDetailsRepository) ClientDetailsService {
	return &ClientDetailsServiceImpl{
		repo: repo,
	}
}

func (svc *ClientDetailsServiceImpl) GetClientDetailByClientId(ctx context.Context, clientId string, clientSecret string) (*model.ClientDetails, error) {
	clientDetails, err := svc.repo.GetClientDetailsByClientId(ctx, clientId)
	if err != nil {
		return nil, err
	}

	if !strings.EqualFold(clientDetails.ClientSecret, clientSecret) {
		return nil, ErrClientMessage
	}

	return model.ClientDetailsEntity2Model(clientDetails)
}
