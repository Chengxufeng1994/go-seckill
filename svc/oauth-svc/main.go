package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/Chengxufeng1994/go-seckill/pkg/bootstrap"
	"github.com/Chengxufeng1994/go-seckill/pkg/discover"
	"github.com/Chengxufeng1994/go-seckill/svc/oauth-svc/config"
	"github.com/Chengxufeng1994/go-seckill/svc/oauth-svc/endpoint"
	"github.com/Chengxufeng1994/go-seckill/svc/oauth-svc/entity"
	gormrepo "github.com/Chengxufeng1994/go-seckill/svc/oauth-svc/repo/gorm"
	memrepo "github.com/Chengxufeng1994/go-seckill/svc/oauth-svc/repo/mem"
	"github.com/Chengxufeng1994/go-seckill/svc/oauth-svc/service"
	"github.com/Chengxufeng1994/go-seckill/svc/oauth-svc/token"
	"github.com/Chengxufeng1994/go-seckill/svc/oauth-svc/transport"
	"github.com/go-kit/kit/log"
	kitzipkin "github.com/go-kit/kit/tracing/zipkin"
	"github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/propagation/b3"
	"google.golang.org/grpc/metadata"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
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

	var tokenEnhancer token.TokenEnhancer
	tokenEnhancer = token.NewJwtTokenEnhancer([]byte("secret"))

	var commonService service.CommonService
	commonService = service.NewCommonService()

	var clientDetailsRepo entity.ClientDetailsRepository
	clientDetailsRepo = gormrepo.NewClientDetailsRepository(config.Db)

	var inMemTokenRepo memrepo.TokenStore
	inMemTokenRepo = memrepo.NewInMemTokenRepository(tokenEnhancer)
	//clientDetailsRepo.CreateClientDetails(context.Background(), &entity.ClientDetails{
	//	ClientId:                    "clientid",
	//	ClientSecret:                "clientsecret",
	//	AccessTokenValiditySeconds:  1800,
	//	RefreshTokenValiditySeconds: 18000,
	//	RegisteredRedirectUri:       "http://127.0.0.1",
	//	AuthorizedGrantTypes:        "[\"password\"]",
	//})

	var clientDetailsService service.ClientDetailsService
	clientDetailsService = service.NewClientDetailsService(clientDetailsRepo)

	var userDetailsService service.UserDetailsService
	userDetailsService = service.NewRemoteUserDetailsService()

	var tokenService service.TokenService
	tokenService = service.NewTokenService(tokenEnhancer, inMemTokenRepo)

	tokenGrantDict := make(map[string]service.TokenGranter)
	tokenGrantDict["password"] = service.NewUsernamePasswordTokenGranter("password", userDetailsService, tokenService)
	tokenGrantDict["refresh_token"] = service.NewRefreshGranter("refresh_token", userDetailsService, tokenService)

	var tokenGranter service.TokenGranter
	tokenGranter = service.NewComposeTokenGranter(tokenGrantDict)

	tokenEndpoint := endpoint.MakeTokenEndpoint(tokenGranter, clientDetailsService)
	tokenEndpoint = endpoint.MakeClientAuthorizationMiddleware(config.Logger)(tokenEndpoint)
	tokenEndpoint = kitzipkin.TraceEndpoint(config.ZipkinTracer, "token-endpoint")(tokenEndpoint)

	checkTokenEndpoint := endpoint.MakeCheckTokenEndpoint(tokenService)
	checkTokenEndpoint = endpoint.MakeClientAuthorizationMiddleware(config.Logger)(checkTokenEndpoint)
	checkTokenEndpoint = kitzipkin.TraceEndpoint(config.ZipkinTracer, "check-endpoint")(checkTokenEndpoint)

	healthCheckEndpoint := endpoint.MakeHealthCheckEndpoint(commonService)
	healthCheckEndpoint = kitzipkin.TraceEndpoint(config.ZipkinTracer, "health-endpoint")(healthCheckEndpoint)

	endpts := endpoint.OAuth2Endpoints{
		TokenEndpoint:       tokenEndpoint,
		CheckTokenEndpoint:  checkTokenEndpoint,
		HealthCheckEndpoint: healthCheckEndpoint,
	}

	discover.Register()

	go runHttpSrv(ctx, endpts, clientDetailsService, config.ZipkinTracer, config.Logger, errChan)

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

func runHttpSrv(ctx context.Context, endpts endpoint.OAuth2Endpoints, clientDetailsService service.ClientDetailsService, tracer *zipkin.Tracer, logger log.Logger, errChan chan error) {
	config.Logger.Log("info", "Http server start at port:"+strconv.Itoa(*servicePort))
	handler := transport.MakeHttpHandler(ctx, endpts, clientDetailsService, tracer, logger)
	addr := fmt.Sprintf(":%d", *servicePort)
	errChan <- http.ListenAndServe(addr, handler)
}

func runGrpcSrv(ctx context.Context, endpts endpoint.OAuth2Endpoints, tracer *zipkin.Tracer, errChan chan error) {
	config.Logger.Log("info", "gRPC server start at port:"+strconv.Itoa(*grpcPort))
	//serverTracer := kitzipkin.GRPCServerTrace(tracer, kitzipkin.Name("grpc-transport"))
	tr := tracer
	md := metadata.MD{}
	parentSpan := tr.StartSpan("test")

	b3.InjectGRPC(&md)(parentSpan.Context())

	grpcAddr := fmt.Sprintf(":%d", *grpcPort)
	_, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		errChan <- err
		return
	}
}
