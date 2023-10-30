package endpoint

import (
	"context"
	"log"

	"github.com/Chengxufeng1994/go-seckill/svc/sk-admin-svc/model"
	"github.com/Chengxufeng1994/go-seckill/svc/sk-admin-svc/service"
	"github.com/go-kit/kit/endpoint"
)

type SkAdminEndpoints struct {
	HealthCheckEndpoint    endpoint.Endpoint
	GetActivityEndpoint    endpoint.Endpoint
	CreateActivityEndpoint endpoint.Endpoint
	GetProductEndpoint     endpoint.Endpoint
	CreateProductEndpoint  endpoint.Endpoint
}

type GetListRequest struct{}

type GetResponse struct {
	Result []*model.Activity `json:"result"`
	Error  error             `json:"error"`
}

func MakeGetActivityEndpoint(svc service.ActivityService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		log.Printf("GetActivityList")
		activityList, calError := svc.GetActivityList()
		if calError != nil {
			return GetResponse{Result: nil, Error: calError}, nil
		}
		return GetResponse{Result: activityList, Error: calError}, nil
	}
}

type GetProductListResponse struct {
	Result []*model.Product `json:"result"`
	Error  error            `json:"error"`
}

func MakeGetProductEndpoint(svc service.ProductService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		log.Printf("GetProductList")
		list, calError := svc.GetProductList()
		if calError != nil {
			return GetProductListResponse{Result: nil, Error: calError}, nil
		}
		return GetProductListResponse{Result: list, Error: calError}, nil
	}
}

type CreateResponse struct {
	Error error `json:"error"`
}

func MakeCreateActivityEndpoint(svc service.ActivityService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		activity := request.(*model.Activity)
		err := svc.CreateActivity(activity)
		return CreateResponse{Error: err}, nil
	}
}

func MakeCreateProductEndpoint(svc service.ProductService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		prod := request.(*model.Product)
		err := svc.CreateProduct(prod)
		return CreateResponse{Error: err}, nil
	}
}

type HealthCheckRequest struct {
}

type HealthCheckResponse struct {
	Status bool `json:"status"`
}

func MakeHealthCheckEndpoint(ctx context.Context, svc service.CommonService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		status := svc.HealthCheck()
		return &HealthCheckResponse{status}, nil
	}
}
