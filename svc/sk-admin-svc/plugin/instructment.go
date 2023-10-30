package plugin

import (
	"time"

	"github.com/Chengxufeng1994/go-seckill/svc/sk-admin-svc/model"
	"github.com/Chengxufeng1994/go-seckill/svc/sk-admin-svc/service"
	"github.com/go-kit/kit/metrics"
)

// metricMiddleware 定义监控中间件，嵌入Service
// 新增监控指标项：requestCount和requestLatency
type skAdminMetricMiddleware struct {
	service.CommonService
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
}

type activityMetricMiddleware struct {
	service.ActivityService
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
}

type productMetricMiddleware struct {
	service.ProductService
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
}

func SkAdminMetrics(requestCount metrics.Counter, requestLatency metrics.Histogram) service.CommonServiceMiddleware {
	return func(next service.CommonService) service.CommonService {
		return skAdminMetricMiddleware{
			next,
			requestCount,
			requestLatency,
		}
	}
}

// Metrics 封装监控方法
func ActivityMetrics(requestCount metrics.Counter, requestLatency metrics.Histogram) service.ActivityServiceMiddleware {
	return func(next service.ActivityService) service.ActivityService {
		return activityMetricMiddleware{
			next,
			requestCount,
			requestLatency}
	}
}

// Metrics 封装监控方法
func ProductMetrics(requestCount metrics.Counter, requestLatency metrics.Histogram) service.ProductServiceMiddleware {
	return func(next service.ProductService) service.ProductService {
		return productMetricMiddleware{
			next,
			requestCount,
			requestLatency}
	}
}

func (mw skAdminMetricMiddleware) HealthCheck() (result bool) {

	defer func(begin time.Time) {
		lvs := []string{"method", "HealthCheck"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	result = mw.CommonService.HealthCheck()
	return
}

func (mw activityMetricMiddleware) GetActivityList() ([]*model.Activity, error) {

	defer func(begin time.Time) {
		lvs := []string{"method", "HealthCheck"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	result, err := mw.ActivityService.GetActivityList()
	return result, err
}

func (mw activityMetricMiddleware) CreateActivity(activity *model.Activity) error {

	defer func(begin time.Time) {
		lvs := []string{"method", "HealthCheck"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	err := mw.ActivityService.CreateActivity(activity)
	return err
}

func (mw productMetricMiddleware) CreateProduct(product *model.Product) error {

	defer func(begin time.Time) {
		lvs := []string{"method", "HealthCheck"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	err := mw.ProductService.CreateProduct(product)
	return err
}

func (mw productMetricMiddleware) GetProductList() ([]*model.Product, error) {

	defer func(begin time.Time) {
		lvs := []string{"method", "HealthCheck"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	data, err := mw.ProductService.GetProductList()
	return data, err
}
