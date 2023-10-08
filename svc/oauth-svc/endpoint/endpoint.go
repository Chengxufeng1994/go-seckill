package endpoint

import (
	"context"
	"errors"
	"github.com/Chengxufeng1994/go-seckill/svc/oauth-svc/model"
	"github.com/Chengxufeng1994/go-seckill/svc/oauth-svc/service"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"net/http"
)

func MakeClientAuthorizationMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			if err, ok := ctx.Value(OAuth2ErrorKey).(error); ok {
				return nil, err
			}
			if _, ok := ctx.Value(OAuth2ClientDetailsKey).(*model.ClientDetails); !ok {
				return nil, ErrInvalidClientRequest
			}
			return next(ctx, request)
		}
	}
}

func MakeOAuth2AuthorizationMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {

		return func(ctx context.Context, request interface{}) (response interface{}, err error) {

			if err, ok := ctx.Value(OAuth2ErrorKey).(error); ok {
				return nil, err
			}
			if _, ok := ctx.Value(OAuth2DetailsKey).(*model.OAuth2Details); !ok {
				return nil, ErrInvalidUserRequest
			}
			return next(ctx, request)
		}
	}
}
func MakeAuthorityAuthorizationMiddleware(authority string, logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {

		return func(ctx context.Context, request interface{}) (response interface{}, err error) {

			if err, ok := ctx.Value(OAuth2ErrorKey).(error); ok {
				return nil, err
			}
			if details, ok := ctx.Value(OAuth2DetailsKey).(*model.OAuth2Details); !ok {
				return nil, ErrInvalidClientRequest
			} else {
				for _, value := range details.User.Authorities {
					if value == authority {
						return next(ctx, request)
					}
				}
				return nil, ErrNotPermit
			}
		}
	}
}

const (
	OAuth2DetailsKey       = "OAuth2Details"
	OAuth2ClientDetailsKey = "OAuth2ClientDetails"
	OAuth2ErrorKey         = "OAuth2Error"
)

var (
	ErrInvalidClientRequest = errors.New("invalid client message")
	ErrInvalidUserRequest   = errors.New("invalid user message")
	ErrNotPermit            = errors.New("not permit")
)

type OAuth2Endpoints struct {
	TokenEndpoint          endpoint.Endpoint
	CheckTokenEndpoint     endpoint.Endpoint
	GrpcCheckTokenEndpoint endpoint.Endpoint
	HealthCheckEndpoint    endpoint.Endpoint
}

type TokenRequest struct {
	GrantType string
	Reader    *http.Request
}

type TokenResponse struct {
	AccessToken *model.OAuth2Token `json:"access_token"`
	Error       string             `json:"error"`
}

func MakeTokenEndpoint(svc service.TokenGranter, clientService service.ClientDetailsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*TokenRequest)
		grantType := req.GrantType
		reader := req.Reader
		clientDetails := ctx.Value(OAuth2ClientDetailsKey).(*model.ClientDetails)
		token, err := svc.Grant(ctx, grantType, clientDetails, reader)
		var errString = ""
		if err != nil {
			errString = err.Error()
		}

		return &TokenResponse{
			AccessToken: token,
			Error:       errString,
		}, nil
	}
}

type CheckTokenRequest struct {
	Token         string
	ClientDetails model.ClientDetails
}

type CheckTokenResponse struct {
	OAuthDetails *model.OAuth2Details `json:"o_auth_details"`
	Error        string               `json:"error"`
}

func MakeCheckTokenEndpoint(svc service.TokenService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		checkTokenRequest := request.(*CheckTokenRequest)
		tokenDetails, err := svc.GetOAuth2DetailsByAccessToken(checkTokenRequest.Token)
		var errString = ""
		if err != nil {
			errString = err.Error()
		}

		return &CheckTokenResponse{
			OAuthDetails: tokenDetails,
			Error:        errString,
		}, nil
	}
}

type HealthCheckRequest struct {
}

type HealthCheckResponse struct {
	Status bool `json:"status"`
}

func MakeHealthCheckEndpoint(svc service.CommonService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		_ = request.(*HealthCheckRequest)
		status := svc.HealthCheck()
		return &HealthCheckResponse{Status: status}, nil
	}
}
