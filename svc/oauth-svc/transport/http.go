package transport

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Chengxufeng1994/go-seckill/svc/oauth-svc/endpoint"
	"github.com/Chengxufeng1994/go-seckill/svc/oauth-svc/service"
	"github.com/go-kit/kit/log"
	kitzipkin "github.com/go-kit/kit/tracing/zipkin"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/openzipkin/zipkin-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var (
	ErrorBadRequest         = errors.New("invalid request parameter")
	ErrorGrantTypeRequest   = errors.New("invalid request grant type")
	ErrorTokenRequest       = errors.New("invalid request token")
	ErrInvalidClientRequest = errors.New("invalid client message")
)

func MakeHttpHandler(ctx context.Context, endpoints endpoint.OAuth2Endpoints, clientDetailsService service.ClientDetailsService, zipkinTracer *zipkin.Tracer, logger log.Logger) http.Handler {
	r := mux.NewRouter()

	zipkinOptions := []kitzipkin.TracerOption{
		kitzipkin.Name("http-transport"),
	}
	zipkinServer := kitzipkin.HTTPServerTrace(zipkinTracer, zipkinOptions...)

	options := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encodeError),
		zipkinServer,
	}

	clientAuthorizationOptions := []kithttp.ServerOption{
		kithttp.ServerBefore(makeClientAuthorizationContext(clientDetailsService, logger)),
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encodeError),
		zipkinServer,
	}

	r.Path("/metrics").
		Handler(promhttp.Handler())

	r.Methods(http.MethodPost).
		Path("/oauth/token").
		Handler(kithttp.NewServer(
			endpoints.TokenEndpoint,
			decodeTokenRequest,
			encodeJsonResponse,
			clientAuthorizationOptions...,
		))

	r.Methods(http.MethodPost).
		Path("/oauth/check_token").
		Handler(kithttp.NewServer(
			endpoints.CheckTokenEndpoint,
			decodeCheckTokenRequest,
			encodeJsonResponse,
			clientAuthorizationOptions...,
		))

	r.Methods(http.MethodGet).
		Path("/health").
		Handler(kithttp.NewServer(
			endpoints.HealthCheckEndpoint,
			decodeHealthCheckRequest,
			encodeJsonResponse,
			options...,
		))

	return r
}

func makeClientAuthorizationContext(clientDetailsService service.ClientDetailsService, logger log.Logger) kithttp.RequestFunc {
	return func(ctx context.Context, r *http.Request) context.Context {
		clientId, clientSecret, ok := r.BasicAuth()
		if !ok {
			return context.WithValue(ctx, endpoint.OAuth2ErrorKey, ErrInvalidClientRequest)
		}

		clientDetails, err := clientDetailsService.GetClientDetailByClientId(ctx, clientId, clientSecret)
		if err != nil {
			return context.WithValue(ctx, endpoint.OAuth2ErrorKey, ErrInvalidClientRequest)
		}

		return context.WithValue(ctx, endpoint.OAuth2ClientDetailsKey, clientDetails)
	}
}

func decodeTokenRequest(_ context.Context, r *http.Request) (any, error) {
	grantType := r.URL.Query().Get("grant_type")
	if grantType == "" {
		return nil, ErrorGrantTypeRequest
	}
	return &endpoint.TokenRequest{
		GrantType: grantType,
		Reader:    r,
	}, nil
}

func decodeCheckTokenRequest(_ context.Context, r *http.Request) (any, error) {
	tokenValue := r.URL.Query().Get("token")
	if tokenValue == "" {
		return nil, ErrorTokenRequest
	}

	return &endpoint.CheckTokenRequest{
		Token: tokenValue,
	}, nil
}

func decodeHealthCheckRequest(_ context.Context, r *http.Request) (any, error) {
	return &endpoint.HealthCheckRequest{}, nil
}

func encodeJsonResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

// encode errors from business-logic
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
