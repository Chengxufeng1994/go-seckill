package client

import (
	"context"
	"errors"
	"github.com/Chengxufeng1994/go-seckill/pkg/bootstrap"
	conf "github.com/Chengxufeng1994/go-seckill/pkg/config"
	"github.com/Chengxufeng1994/go-seckill/pkg/discover"
	"github.com/Chengxufeng1994/go-seckill/pkg/loadbalance"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"google.golang.org/grpc"
	"log"
	"strconv"
	"time"
)

var (
	ErrRPCService = errors.New("no rpc service")
)

var defaultLoadBalance loadbalance.LoadBalance = &loadbalance.RandomLoadBalance{}

type ClientManager interface {
	DecoratorInvoke(path string, hystrixName string, tracer opentracing.Tracer,
		ctx context.Context, inputVal interface{}, outVal interface{}) (err error)
}

type DefaultClientManager struct {
	serviceName     string
	logger          *log.Logger
	discoveryClient discover.DiscoveryClient
	loadBalance     loadbalance.LoadBalance
	after           []InvokerAfterFunc
	before          []InvokerBeforeFunc
}

type InvokerAfterFunc func() (err error)

type InvokerBeforeFunc func() (err error)

func (manager *DefaultClientManager) DecoratorInvoke(path string, hystrixName string,
	tracer opentracing.Tracer, ctx context.Context, inputVal interface{}, outVal interface{}) (err error) {

	for _, fn := range manager.before {
		if err = fn(); err != nil {
			return err
		}
	}

	if err = hystrix.Do(hystrixName, func() error {

		instances := manager.discoveryClient.DiscoverServices(manager.serviceName, manager.logger)
		if instance, err := manager.loadBalance.SelectService(instances); err == nil {
			if instance.GrpcPort > 0 {
				if instance.Host == "host.docker.internal" {
					instance.Host = "127.0.0.1"
				}

				if conn, err := grpc.Dial(instance.Host+":"+strconv.Itoa(instance.GrpcPort), grpc.WithInsecure(),
					grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(genTracer(tracer), otgrpc.LogPayloads())), grpc.WithTimeout(1*time.Second)); err == nil {
					if err = conn.Invoke(ctx, path, inputVal, outVal); err != nil {
						return err
					}
				} else {
					return err
				}
			} else {
				return ErrRPCService
			}
		} else {
			return err
		}
		return nil
	}, func(e error) error {
		return e
	}); err != nil {
		return err
	} else {
		for _, fn := range manager.after {
			if err = fn(); err != nil {
				return err
			}
		}
		return nil
	}
}

func genTracer(tracer opentracing.Tracer) opentracing.Tracer {
	if tracer != nil {
		return tracer
	}
	zipkinUrl := "http://" + conf.Conf.Trace.Host + ":" + strconv.Itoa(conf.Conf.Trace.Port) + conf.Conf.Trace.Url
	useNoopTracer := zipkinUrl == ""

	reporter := zipkinhttp.NewReporter(zipkinUrl)
	// create our local service endpoint
	zEP, err := zipkin.NewEndpoint(bootstrap.Conf.Discover.ServiceName, strconv.Itoa(bootstrap.Conf.Discover.Port))
	if err != nil {
		log.Fatalf("zipkin.NewEndpoint unable to create local endpoint: %+v\n", err)
	}
	nativeTracer, err := zipkin.NewTracer(
		reporter, zipkin.WithLocalEndpoint(zEP), zipkin.WithNoopTracer(useNoopTracer),
	)
	if err != nil {
		log.Fatalf("zipkin.NewTracer unable to create tracer: %+v\n", err)
	}

	return zipkinot.Wrap(nativeTracer)
}
