package transport

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Chengxufeng1994/go-seckill/svc/sk-admin-svc/endpoint"
	"github.com/Chengxufeng1994/go-seckill/svc/sk-admin-svc/model"
	"github.com/go-kit/kit/log"
	kitzipkin "github.com/go-kit/kit/tracing/zipkin"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/openzipkin/zipkin-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func MakeHttpHandler(ctx context.Context, endpts endpoint.SkAdminEndpoints, zipkinTracer *zipkin.Tracer, logger log.Logger) http.Handler {
	r := mux.NewRouter()

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

	r.Methods(http.MethodGet).Path("/activity/list").Handler(kithttp.NewServer(
		endpts.GetActivityEndpoint,
		decodeGetListRequest,
		encodeResponse,
		options...,
	))

	r.Methods(http.MethodPost).Path("/activity/create").Handler(kithttp.NewServer(
		endpts.CreateActivityEndpoint,
		decodeCreateActivityRequest,
		encodeResponse,
		options...,
	))

	r.Methods(http.MethodGet).Path("/product/list").Handler(kithttp.NewServer(
		endpts.GetProductEndpoint,
		decodeGetListRequest,
		encodeResponse,
		options...,
	))

	r.Methods(http.MethodPost).Path("/product/create").Handler(kithttp.NewServer(
		endpts.CreateProductEndpoint,
		decodeCreateProductCheckRequest,
		encodeResponse,
		options...,
	))

	r.Path("/metrics").Handler(promhttp.Handler())

	r.Methods(http.MethodGet).
		Path("/health").
		Handler(kithttp.NewServer(
			endpts.HealthCheckEndpoint,
			decodeHealthCheckRequest,
			encodeResponse,
			options...,
		))

	return r
}

func decodeGetListRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return &endpoint.GetListRequest{}, nil
}

func decodeCreateActivityRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var activity model.Activity
	if err := json.NewDecoder(r.Body).Decode(&activity); err != nil {
		return nil, err
	}

	return &activity, nil
}

func decodeCreateProductCheckRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var product model.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		return nil, err
	}
	return &product, nil
}

func decodeHealthCheckRequest(_ context.Context, req *http.Request) (interface{}, error) {
	return &endpoint.HealthCheckRequest{}, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, resp interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(resp)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"svc_err": err.Error(),
	})
}
