package setup

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	register "github.com/Chengxufeng1994/go-seckill/pkg/discover"
	"github.com/Chengxufeng1994/go-seckill/svc/sk-admin-svc/config"
	"github.com/Chengxufeng1994/go-seckill/svc/sk-admin-svc/db"
	"github.com/Chengxufeng1994/go-seckill/svc/sk-admin-svc/endpoint"
	"github.com/Chengxufeng1994/go-seckill/svc/sk-admin-svc/entity"
	"github.com/Chengxufeng1994/go-seckill/svc/sk-admin-svc/plugin"
	"github.com/Chengxufeng1994/go-seckill/svc/sk-admin-svc/repo/gorm"
	"github.com/Chengxufeng1994/go-seckill/svc/sk-admin-svc/service"
	"github.com/Chengxufeng1994/go-seckill/svc/sk-admin-svc/transport"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	kitzipkin "github.com/go-kit/kit/tracing/zipkin"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

func InitializeService(host string, port int) {
	ctx := context.Background()
	errChan := make(chan error)

	fieldKeys := []string{"method"}

	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "aoho",
		Subsystem: "user_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)

	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "aoho",
		Subsystem: "user_service",
		Name:      "request_latency",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)
	//ratebucket := rate.NewLimiter(rate.Every(time.Second*1), 100)

	var repo entity.Repository
	var activitySvc service.ActivityService
	var productSvc service.ProductService
	var commonSvc service.CommonService

	repo = gorm.NewGormRepository(db.Db)
	commonSvc = service.NewCommonService()
	commonSvc = plugin.SkAdminLoggingMiddleware(config.Logger)(commonSvc)
	commonSvc = plugin.SkAdminMetrics(requestCount, requestLatency)(commonSvc)

	activitySvc = service.NewActivityService(repo)
	activitySvc = plugin.ActivityLoggingMiddleware(config.Logger)(activitySvc)
	activitySvc = plugin.ActivityMetrics(requestCount, requestLatency)(activitySvc)

	productSvc = service.NewProductService(repo)
	productSvc = plugin.ProductLoggingMiddleware(config.Logger)(productSvc)
	productSvc = plugin.ProductMetrics(requestCount, requestLatency)(productSvc)

	getActivityEndpoint := endpoint.MakeGetActivityEndpoint(activitySvc)
	getActivityEndpoint = kitzipkin.TraceEndpoint(config.ZipkinTracer, "get-activity")(getActivityEndpoint)

	createActivityEndpoint := endpoint.MakeCreateActivityEndpoint(activitySvc)
	createActivityEndpoint = kitzipkin.TraceEndpoint(config.ZipkinTracer, "create-activity")(createActivityEndpoint)

	getProductEndpoint := endpoint.MakeGetProductEndpoint(productSvc)
	getProductEndpoint = kitzipkin.TraceEndpoint(config.ZipkinTracer, "get-product")(getProductEndpoint)

	createProductEndpoint := endpoint.MakeCreateProductEndpoint(productSvc)
	createProductEndpoint = kitzipkin.TraceEndpoint(config.ZipkinTracer, "create-product")(createProductEndpoint)

	healthCheckEndpoint := endpoint.MakeHealthCheckEndpoint(ctx, commonSvc)
	healthCheckEndpoint = kitzipkin.TraceEndpoint(config.ZipkinTracer, "health-endpoint")(healthCheckEndpoint)

	endpts := endpoint.SkAdminEndpoints{
		HealthCheckEndpoint:    healthCheckEndpoint,
		GetActivityEndpoint:    getActivityEndpoint,
		CreateActivityEndpoint: createActivityEndpoint,
		GetProductEndpoint:     getProductEndpoint,
		CreateProductEndpoint:  createProductEndpoint,
	}

	register.Register()

	go func() {
		config.Logger.Log("info", "Http server start at port:"+strconv.Itoa(port))
		handler := transport.MakeHttpHandler(ctx, endpts, config.ZipkinTracer, config.Logger)
		s := &http.Server{
			Addr:           fmt.Sprintf(":%d", port),
			Handler:        handler,
			ReadTimeout:    60 * time.Second,
			WriteTimeout:   60 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}
		errChan <- s.ListenAndServe()
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	err := <-errChan
	register.Deregister()
	config.Logger.Log("err", err)
}
