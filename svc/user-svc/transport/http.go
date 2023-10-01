package transport

import (
	"context"
	"encoding/json"
	"github.com/Chengxufeng1994/go-seckill/svc/user-svc/endpoint"
	"github.com/go-kit/kit/log"
	kitzipkin "github.com/go-kit/kit/tracing/zipkin"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/openzipkin/zipkin-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	//"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func MakeHttpHandler(ctx context.Context, endpoints endpoint.UserEndpoints, zipkinTracer *zipkin.Tracer, logger log.Logger) http.Handler {
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

	r.Methods(http.MethodPost).
		Path("/check/valid").
		Handler(kithttp.NewServer(
			endpoints.UserEndpoint,
			decodeUserRequest,
			encodeUserResponse,
			options...,
		))

	r.Path("/metrics").
		Handler(promhttp.Handler())

	r.Methods(http.MethodGet).
		Path("/health").
		Handler(kithttp.NewServer(
			endpoints.HealthCheckEndpoint,
			decodeHealthCheckRequest,
			encodeUserResponse,
			options...,
		))

	return r
}

// decodeUserRequest decode request params to struct
func decodeUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var userRequest endpoint.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&userRequest); err != nil {
		return nil, err
	}

	return userRequest, nil
}

func encodeUserResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func decodeHealthCheckRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return endpoint.HealthCheckRequest{}, nil
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
