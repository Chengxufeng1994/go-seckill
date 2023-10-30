package plugin

import (
	"time"

	"github.com/Chengxufeng1994/go-seckill/svc/sk-admin-svc/model"
	"github.com/Chengxufeng1994/go-seckill/svc/sk-admin-svc/service"
	"github.com/go-kit/kit/log"
)

// loggingMiddleware Make a new type
// that contains Service interface and logger instance
type skAdminLoggingMiddleware struct {
	service.CommonService
	logger log.Logger
}

type activityLoggingMiddleware struct {
	service.ActivityService
	logger log.Logger
}

type productLoggingMiddleware struct {
	service.ProductService
	logger log.Logger
}

// LoggingMiddleware make logging middleware
func SkAdminLoggingMiddleware(logger log.Logger) service.CommonServiceMiddleware {
	return func(next service.CommonService) service.CommonService {
		return skAdminLoggingMiddleware{next, logger}
	}
}

func ActivityLoggingMiddleware(logger log.Logger) service.ActivityServiceMiddleware {
	return func(next service.ActivityService) service.ActivityService {
		return activityLoggingMiddleware{next, logger}
	}
}

func ProductLoggingMiddleware(logger log.Logger) service.ProductServiceMiddleware {
	return func(next service.ProductService) service.ProductService {
		return productLoggingMiddleware{next, logger}
	}
}

func (mw productLoggingMiddleware) CreateProduct(product *model.Product) (err error) {

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"function", "Check",
			"product", product,
			"took", time.Since(begin),
		)
	}(time.Now())

	err = mw.ProductService.CreateProduct(product)
	return err
}

func (mw productLoggingMiddleware) GetProductList() ([]*model.Product, error) {

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"function", "Check",
			"took", time.Since(begin),
		)
	}(time.Now())

	data, err := mw.ProductService.GetProductList()
	return data, err
}

func (mw activityLoggingMiddleware) GetActivityList() ([]*model.Activity, error) {

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"function", "Check",
			"took", time.Since(begin),
		)
	}(time.Now())

	ret, err := mw.ActivityService.GetActivityList()
	return ret, err
}

func (mw activityLoggingMiddleware) CreateActivity(activity *model.Activity) error {

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"function", "Check",
			"activity", activity,
			"took", time.Since(begin),
		)
	}(time.Now())

	err := mw.ActivityService.CreateActivity(activity)
	return err
}

func (mw skAdminLoggingMiddleware) HealthCheck() (result bool) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"function", "HealthCheck",
			"result", result,
			"took", time.Since(begin),
		)
	}(time.Now())
	result = mw.CommonService.HealthCheck()
	return
}
