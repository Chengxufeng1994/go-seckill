package service

import (
	"context"
	"errors"
	"github.com/Chengxufeng1994/go-seckill/pb"
	"github.com/Chengxufeng1994/go-seckill/pkg/client"
	"github.com/Chengxufeng1994/go-seckill/svc/oauth-svc/model"
)

var (
	InvalidAuthentication = errors.New("invalid auth")
	InvalidUserInfo       = errors.New("invalid user info")
)

// Service Define a service interface
type UserDetailsService interface {

	// Get UserDetails By username
	GetUserDetailByUsername(ctx context.Context, username, password string) (*model.UserDetails, error)
}

type RemoteUserDetailsServiceImpl struct {
	userClient client.UserClient
}

func NewRemoteUserDetailsService() UserDetailsService {
	userClient, _ := client.NewUserClient("user", nil, nil)
	return &RemoteUserDetailsServiceImpl{
		userClient: userClient,
	}
}

func (svc *RemoteUserDetailsServiceImpl) GetUserDetailByUsername(ctx context.Context, username, password string) (*model.UserDetails, error) {
	response, err := svc.userClient.CheckUser(ctx, nil, &pb.UserRequest{
		Username: username,
		Password: password,
	})

	if err != nil {
		return nil, err
	}

	if response.UserId == 0 {
		return nil, InvalidUserInfo
	}

	return &model.UserDetails{
		UserId:   response.UserId,
		Username: username,
		Password: password,
	}, nil
}
