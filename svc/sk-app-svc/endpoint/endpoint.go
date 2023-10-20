package endpoint

import (
	"context"
	"github.com/Chengxufeng1994/go-seckill/svc/sk-app-svc/model"
	"github.com/Chengxufeng1994/go-seckill/svc/sk-app-svc/service"
	"github.com/go-kit/kit/endpoint"
)

type SkAppEndpoints struct {
	HealthCheckEndpoint    endpoint.Endpoint
	GetSecInfoEndpoint     endpoint.Endpoint
	GetSecInfoListEndpoint endpoint.Endpoint
	SecKillEndpoint        endpoint.Endpoint
}

type SecInfoRequest struct {
	ProductId int `json:"id"`
}

type Response struct {
	Result map[string]interface{} `json:"result"`
	Error  error                  `json:"svc_err"`
	Code   int                    `json:"code"`
}

func MakeGetSecInfoEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		infoRequest := request.(SecInfoRequest)
		info := svc.SecInfo(infoRequest.ProductId)
		return Response{
			Result: info,
			Error:  nil,
		}, nil
	}
}

type SecInfoListResponse struct {
	Result []map[string]interface{} `json:"result"`
	Num    int                      `json:"num"`
	Error  error                    `json:"error"`
}

func MakeGetSecInfoListEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		ret, num, err := svc.SecInfoList()
		return SecInfoListResponse{
			Result: ret,
			Num:    num,
			Error:  err,
		}, nil
	}
}

func MakeSecKillEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(model.SecRequest)
		ret, code, calError := svc.SecKill(&req)
		return Response{Result: ret, Code: code, Error: calError}, nil
	}
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
