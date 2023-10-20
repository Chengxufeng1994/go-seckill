package plugin

import (
	"github.com/Chengxufeng1994/go-seckill/svc/sk-app-svc/model"
	"github.com/Chengxufeng1994/go-seckill/svc/sk-app-svc/service"
	"github.com/go-kit/kit/metrics"
	"time"
)

// metricMiddleware 定义监控中间件，嵌入Service
// 新增监控指标项：requestCount和requestLatency
type skAppMetricMiddleware struct {
	service.Service
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
}

// Metrics 封装监控方法
func SkAppMetrics(requestCount metrics.Counter, requestLatency metrics.Histogram) service.ServiceMiddleware {
	return func(next service.Service) service.Service {
		return skAppMetricMiddleware{
			next,
			requestCount,
			requestLatency}
	}
}

func (mw skAppMetricMiddleware) HealthCheck() (result bool) {

	defer func(begin time.Time) {
		lvs := []string{"method", "HealthCheck"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	result = mw.Service.HealthCheck()
	return
}

func (mw skAppMetricMiddleware) SecInfo(productId int) map[string]interface{} {

	defer func(begin time.Time) {
		lvs := []string{"method", "HealthCheck"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	ret := mw.Service.SecInfo(productId)
	return ret
}

func (mw skAppMetricMiddleware) SecInfoList() ([]map[string]interface{}, int, error) {

	defer func(begin time.Time) {
		lvs := []string{"method", "HealthCheck"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	data, num, error := mw.Service.SecInfoList()
	return data, num, error
}

func (mw skAppMetricMiddleware) SecKill(req *model.SecRequest) (map[string]interface{}, int, error) {

	defer func(begin time.Time) {
		lvs := []string{"method", "HealthCheck"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	result, num, error := mw.Service.SecKill(req)
	return result, num, error
}
