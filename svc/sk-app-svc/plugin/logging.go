package plugin

import (
	"github.com/Chengxufeng1994/go-seckill/svc/sk-app-svc/service"
	"github.com/go-kit/kit/log"
	"time"
)

// loggingMiddleware Make a new type
// that contains Service interface and logger instance
type skAppLoggingMiddleware struct {
	service.Service
	logger log.Logger
}

func SkAppLoggingMiddleware(logger log.Logger) service.ServiceMiddleware {
	return func(next service.Service) service.Service {
		return skAppLoggingMiddleware{next, logger}
	}
}

func (mw skAppLoggingMiddleware) HealthCheck() bool {
	var ret bool
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"function", "HealthChcek",
			"result", ret,
			"took", time.Since(begin),
		)
	}(time.Now())

	ret = mw.Service.HealthCheck()
	return ret
}
