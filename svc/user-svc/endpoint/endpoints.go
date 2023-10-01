package endpoint

import (
	"context"
	"errors"
	"github.com/Chengxufeng1994/go-seckill/svc/user-svc/service"
	"github.com/go-kit/kit/endpoint"
)

type UserEndpoints struct {
	UserEndpoint        endpoint.Endpoint
	HealthCheckEndpoint endpoint.Endpoint
}

var (
	ErrInvalidRequestType = errors.New("invalid username, password")
)

// UserRequest define request struct
type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// UserResponse define response struct
type UserResponse struct {
	Result bool   `json:"result"`
	UserId int64  `json:"user_id"`
	Error  string `json:"error"`
}

// MakeUserEndpoints
func MakeUserEndpoints(ctx context.Context, svc service.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(UserRequest)
		username := req.Username
		password := req.Password
		userId, err := svc.Check(ctx, username, password)
		if err != nil {
			return UserResponse{
				Result: false,
				Error:  ErrInvalidRequestType.Error(),
			}, nil
		}

		return UserResponse{
			Result: true,
			UserId: userId,
		}, nil
	}
}

// HealthCheckRequest define request struct
type HealthCheckRequest struct {
}

// HealthCheckResponse define response struct
type HealthCheckResponse struct {
	Status bool `json:"status"`
}

// MakeHealthCheckEndpoint make healthcheck endpoint
func MakeHealthCheckEndpoint(ctx context.Context, svc service.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		_ = request.(HealthCheckRequest)
		status := svc.HealthCheck()
		return &HealthCheckResponse{
			Status: status,
		}, nil
	}
}
