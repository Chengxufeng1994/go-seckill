package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/Chengxufeng1994/go-seckill/pkg/bootstrap"
	conf "github.com/Chengxufeng1994/go-seckill/pkg/config"
	"github.com/Chengxufeng1994/go-seckill/pkg/discover"
	"github.com/Chengxufeng1994/go-seckill/svc/sk-app-svc/config"
	"github.com/Chengxufeng1994/go-seckill/svc/sk-app-svc/endpoint"
	"github.com/Chengxufeng1994/go-seckill/svc/sk-app-svc/plugin"
	"github.com/Chengxufeng1994/go-seckill/svc/sk-app-svc/service"
	"github.com/Chengxufeng1994/go-seckill/svc/sk-app-svc/setup"
	"github.com/Chengxufeng1994/go-seckill/svc/sk-app-svc/transport"
	"github.com/Chengxufeng1994/go-seckill/svc/sk-app-svc/worker"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {

	flag.Parse()
	// initialize zookeeper
	err := setup.InitializeZk()
	if err != nil {
		log.Fatalf("%s", err)
	}
	// initialize redis
	err = setup.InitializeRedis()
	if err != nil {
		log.Fatalf("%s", err)
	}

	initializeService(bootstrap.Conf.Http.Host, bootstrap.Conf.Http.Port)
}

func runProcessor(client *redis.Client) {
	config.Logger.Log("info", "initialize redis processor")
	config.Logger.Log("info", "AppWriteToHandleGoroutineNum: "+strconv.Itoa(conf.Conf.SecKill.AppWriteToHandleGoroutineNum))
	distributor := worker.NewDistributor(client)
	distributor.Start(conf.Conf.SecKill.AppWriteToHandleGoroutineNum)

	config.Logger.Log("info", "AppReadToHandleGoroutineNum: "+strconv.Itoa(conf.Conf.SecKill.AppReadToHandleGoroutineNum))
	processor := worker.NewProcessor(client)
	processor.Start(conf.Conf.SecKill.AppReadToHandleGoroutineNum)
}

func initializeService(host string, port int) {

	ctx := context.Background()
	errChan := make(chan error)

	fieldKeys := []string{"method"}

	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "default",
		Subsystem: "sk_app",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)

	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "default",
		Subsystem: "sk_app",
		Name:      "request_latency",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)

	svc := service.NewSkAppService()
	svc = plugin.SkAppLoggingMiddleware(config.Logger)(svc)
	svc = plugin.SkAppMetrics(requestCount, requestLatency)(svc)

	getSecInfoEndpoint := endpoint.MakeGetSecInfoEndpoint(svc)
	getSecInfoListEndpoint := endpoint.MakeGetSecInfoListEndpoint(svc)
	secKillEndpoint := endpoint.MakeSecKillEndpoint(svc)
	healthcheckEndpoint := endpoint.MakeHealthCheckEndpoint(svc)

	endpts := endpoint.SkAppEndpoints{
		HealthCheckEndpoint:    healthcheckEndpoint,
		GetSecInfoEndpoint:     getSecInfoEndpoint,
		GetSecInfoListEndpoint: getSecInfoListEndpoint,
		SecKillEndpoint:        secKillEndpoint,
	}

	handler := transport.MakeHttpHandler(ctx, endpts, config.ZipkinTracer, config.Logger)

	discover.Register()

	go runProcessor(conf.Conf.Redis.RedisClient)

	go func() {
		config.Logger.Log("info", "Http server start at port:"+strconv.Itoa(port))
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
	discover.Deregister()
	config.Logger.Log("err", err)
}
