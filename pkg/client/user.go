package client

import (
	"context"
	"github.com/Chengxufeng1994/go-seckill/pb"
	"github.com/Chengxufeng1994/go-seckill/pkg/discover"
	"github.com/Chengxufeng1994/go-seckill/pkg/loadbalance"
	"github.com/opentracing/opentracing-go"
	"strings"
)

type UserClient interface {
	CheckUser(ctx context.Context, tracer opentracing.Tracer, request *pb.UserRequest) (*pb.UserResponse, error)
}

type UserClientImpl struct {
	manager     ClientManager
	serviceName string
	loadBalance loadbalance.LoadBalance
	tracer      opentracing.Tracer
}

func NewUserClient(serviceName string, lb loadbalance.LoadBalance, tracer opentracing.Tracer) (UserClient, error) {
	if strings.EqualFold(serviceName, "") {
		serviceName = "user"
	}
	if lb == nil {
		lb = defaultLoadBalance
	}

	return &UserClientImpl{
		manager: &DefaultClientManager{
			serviceName:     serviceName,
			loadBalance:     lb,
			discoveryClient: discover.ConsulService,
			logger:          discover.Logger,
		},
		serviceName: serviceName,
		loadBalance: lb,
		tracer:      tracer,
	}, nil
}

func (impl *UserClientImpl) CheckUser(ctx context.Context, tracer opentracing.Tracer, request *pb.UserRequest) (*pb.UserResponse, error) {
	response := new(pb.UserResponse)
	if err := impl.manager.DecoratorInvoke("/pb.UserService/Check", "user_check", tracer, ctx, request, response); err != nil {
		return nil, err
	}

	return response, nil
}
