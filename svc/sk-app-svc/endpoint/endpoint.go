package endpoint

import (
	"context"
	"github.com/Chengxufeng1994/go-seckill/svc/sk-app-svc/service"
	"github.com/go-kit/kit/endpoint"
)

type SkAppEndpoints struct {
	HealthCheckEndpoint endpoint.Endpoint
}

type HealthCheckRequest struct {
}

type HealthCheckResponse struct {
	Status bool `json:"status"`
}

func MakeHealthCheckEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		//_ = request.(*HealthCheckRequest)
		status := svc.HealthCheck()
		return HealthCheckResponse{
			Status: status,
		}, nil
	}
}
