package transport

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Chengxufeng1994/go-seckill/svc/sk-app-svc/endpoint"
	"github.com/Chengxufeng1994/go-seckill/svc/sk-app-svc/model"
	"github.com/gin-gonic/gin"
	"github.com/go-kit/kit/log"
	kitzipkin "github.com/go-kit/kit/tracing/zipkin"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/openzipkin/zipkin-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"net/http"
)

var (
	ErrorBadRequest = errors.New("invalid request parameter")
)

func MakeHttpHandler(ctx context.Context, endpts endpoint.SkAppEndpoints, zipkinTracer *zipkin.Tracer, logger log.Logger) http.Handler {

	r := gin.Default()

	zipkinOptions := []kitzipkin.TracerOption{
		kitzipkin.Name("http-transport"),
	}
	zipkinServer := kitzipkin.HTTPServerTrace(zipkinTracer, zipkinOptions...)

	options := []kithttp.ServerOption{
		//kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		//kithttp.ServerErrorEncoder(kithttp.DefaultErrorEncoder),
		kithttp.ServerErrorEncoder(encodeError),
		zipkinServer,
	}

	secInfoHandler := kithttp.NewServer(
		endpts.GetSecInfoEndpoint,
		decodeGetSecInfoRequest,
		encodeJsonResponse,
		options...,
	)

	secInfoListHandler := kithttp.NewServer(
		endpts.GetSecInfoListEndpoint,
		decodeGetSecInfoListRequest,
		encodeJsonResponse,
		options...,
	)

	secKillHandler := kithttp.NewServer(
		endpts.SecKillEndpoint,
		decodeSecKillRequest,
		encodeJsonResponse,
		options...,
	)

	healthcheckHandler := kithttp.NewServer(
		endpts.HealthCheckEndpoint,
		decodeHealthCheckRequest,
		encodeJsonResponse,
		options...,
	)

	r.POST("/sec/info", gin.WrapH(secInfoHandler))

	r.GET("/sec/list", gin.WrapH(secInfoListHandler))

	r.POST("/sec/kill", gin.WrapH(secKillHandler))

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	r.GET("/health", gin.WrapH(healthcheckHandler))

	return h2c.NewHandler(r, &http2.Server{})
}

func decodeGetSecInfoRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var secInfoRequest endpoint.SecInfoRequest
	if err := json.NewDecoder(r.Body).Decode(&secInfoRequest); err != nil {
		return nil, err
	}
	return secInfoRequest, nil
}

func decodeGetSecInfoListRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func decodeSecKillRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var secRequest model.SecRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&secRequest); err != nil {
		return nil, err
	}

	return secRequest, nil
}

func decodeHealthCheckRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return endpoint.HealthCheckRequest{}, nil
}

func encodeJsonResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
