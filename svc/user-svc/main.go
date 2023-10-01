package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/Chengxufeng1994/go-seckill/pb"
	"github.com/Chengxufeng1994/go-seckill/pkg/bootstrap"
	"github.com/Chengxufeng1994/go-seckill/pkg/discover"
	"github.com/Chengxufeng1994/go-seckill/pkg/ratelimiter"
	"github.com/Chengxufeng1994/go-seckill/svc/user-svc/config"
	"github.com/Chengxufeng1994/go-seckill/svc/user-svc/endpoint"
	"github.com/Chengxufeng1994/go-seckill/svc/user-svc/entity"
	"github.com/Chengxufeng1994/go-seckill/svc/user-svc/plugin"
	repo "github.com/Chengxufeng1994/go-seckill/svc/user-svc/repo/gorm"
	"github.com/Chengxufeng1994/go-seckill/svc/user-svc/service"
	"github.com/Chengxufeng1994/go-seckill/svc/user-svc/transport"
	"github.com/go-kit/kit/log"
	kitzipkin "github.com/go-kit/kit/tracing/zipkin"
	"github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/propagation/b3"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

var (
	servicePort *int
	grpcPort    *int
)

func init() {
	servicePort = flag.Int("service.port", bootstrap.Conf.Http.Port, "service port.")
	grpcPort = flag.Int("grpc.port", bootstrap.Conf.Rpc.Port, "gRPC port.")
}

func main() {

	flag.Parse()

	ctx := context.Background()
	errChan := make(chan error)

	var repository entity.Repository
	repository = repo.New(config.Db)

	var svc service.Service
	svc = service.New(repository)
	// add logging middleware
	svc = plugin.LoggingMiddleware(config.Logger)(svc)

	var endpts endpoint.UserEndpoints
	userEndpoint := endpoint.MakeUserEndpoints(ctx, svc)
	//rl := ratelimit.New(100) // per second
	//userEndpoint = ratelimiter.NewTokenBucketLimiterWithUber(rl)(userEndpoint)
	rl := rate.NewLimiter(rate.Every(time.Second*1), 100)
	userEndpoint = ratelimiter.NewTokenBucketLimiterWithBuildIn(rl)(userEndpoint)
	userEndpoint = kitzipkin.TraceEndpoint(config.ZipkinTracer, "user-endpoint")(userEndpoint)

	healthcheckEndpoint := endpoint.MakeHealthCheckEndpoint(ctx, svc)
	healthcheckEndpoint = kitzipkin.TraceEndpoint(config.ZipkinTracer, "health-endpoint")(healthcheckEndpoint)

	endpts = endpoint.UserEndpoints{
		UserEndpoint:        userEndpoint,
		HealthCheckEndpoint: healthcheckEndpoint,
	}

	discover.Register()

	go runHttpSrv(ctx, endpts, config.ZipkinTracer, config.Logger, errChan)

	go runGrpcSrv(ctx, endpts, config.ZipkinTracer, errChan)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	err := <-errChan
	discover.Deregister()
	config.Logger.Log("err", err)
}

func runHttpSrv(ctx context.Context, endpts endpoint.UserEndpoints, tracer *zipkin.Tracer, logger log.Logger, errChan chan error) {
	config.Logger.Log("info", "Http server start at port:"+strconv.Itoa(*servicePort))
	handler := transport.MakeHttpHandler(ctx, endpts, tracer, logger)
	addr := fmt.Sprintf(":%d", *servicePort)
	errChan <- http.ListenAndServe(addr, handler)
}

func runGrpcSrv(ctx context.Context, endpts endpoint.UserEndpoints, tracer *zipkin.Tracer, errChan chan error) {
	config.Logger.Log("info", "gRPC server start at port:"+strconv.Itoa(*grpcPort))
	serverTracer := kitzipkin.GRPCServerTrace(tracer, kitzipkin.Name("grpc-transport"))
	tr := tracer
	md := metadata.MD{}
	parentSpan := tr.StartSpan("test")

	b3.InjectGRPC(&md)(parentSpan.Context())

	grpcAddr := fmt.Sprintf(":%d", *grpcPort)
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		errChan <- err
		return
	}

	handler := transport.NewGrpcServer(ctx, endpts, serverTracer)
	grpcSrv := grpc.NewServer()
	//grpcSrv := grpc.NewServer(grpc.StatsHandler(zipkingrpc.NewServerHandler(tracer)))
	pb.RegisterUserServiceServer(grpcSrv, handler)
	errChan <- grpcSrv.Serve(lis)
}
