package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/Chengxufeng1994/go-seckill/pkg/bootstrap"
	"github.com/Chengxufeng1994/go-seckill/pkg/discover"
	"github.com/Chengxufeng1994/go-seckill/svc/sk-app-svc/config"
	"github.com/Chengxufeng1994/go-seckill/svc/sk-app-svc/endpoint"
	"github.com/Chengxufeng1994/go-seckill/svc/sk-app-svc/service"
	"github.com/Chengxufeng1994/go-seckill/svc/sk-app-svc/transport"
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

	// initialize redis

	initializeService(bootstrap.Conf.Http.Host, bootstrap.Conf.Http.Port)
}

func initializeService(host string, port int) {

	ctx := context.Background()
	errChan := make(chan error)

	svc := service.NewSkAppService()
	healthcheckEndpoint := endpoint.MakeHealthCheckEndpoint(svc)

	endpts := endpoint.SkAppEndpoints{
		HealthCheckEndpoint: healthcheckEndpoint,
	}

	handler := transport.MakeHttpHandler(ctx, endpts, config.ZipkinTracer, config.Logger)

	discover.Register()

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
	discover.Register()
	config.Logger.Log("err", err)
}
