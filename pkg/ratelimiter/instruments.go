package ratelimiter

import (
	"context"
	"errors"
	"github.com/go-kit/kit/endpoint"
	"go.uber.org/ratelimit"
	"golang.org/x/time/rate"
	"time"
)

var ErrLimitExceed = errors.New("rate limit exceed!")

func NewTokenBucketLimiterWithBuildIn(bkt *rate.Limiter) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if !bkt.Allow() {
				return nil, ErrLimitExceed
			}
			return next(ctx, request)
		}
	}
}

func NewTokenBucketLimiterWithUber(limit ratelimit.Limiter) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			limit.Take()
			return next(ctx, request)
		}
	}
}

func DynamicLimitter(interval int, burst int) endpoint.Middleware {
	bucket := rate.NewLimiter(rate.Every(time.Second*time.Duration(interval)), burst)
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if !bucket.Allow() {
				return nil, ErrLimitExceed
			}
			return next(ctx, request)
		}
	}
}
