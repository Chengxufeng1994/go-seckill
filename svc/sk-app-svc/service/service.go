package service

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	conf "github.com/Chengxufeng1994/go-seckill/pkg/config"
	"github.com/Chengxufeng1994/go-seckill/svc/sk-app-svc/config"
	"github.com/Chengxufeng1994/go-seckill/svc/sk-app-svc/model"
	"github.com/Chengxufeng1994/go-seckill/svc/sk-app-svc/service/svc_err"
)

// Service Define a service interface
type Service interface {
	// HealthCheck check service health status
	HealthCheck() bool
	// SecInfo get product information by productId
	SecInfo(int) map[string]interface{}
	// SecInfoList get sec information
	SecInfoList() ([]map[string]interface{}, int, error)

	SecKill(req *model.SecRequest) (map[string]interface{}, int, error)
}

type SkAppService struct {
}

func NewSkAppService() Service {
	return &SkAppService{}
}

func (svc *SkAppService) HealthCheck() bool {
	return true
}

func (svc *SkAppService) SecInfo(productId int) map[string]interface{} {
	config.SkAppContext.RWSecProductLock.RLock()
	defer config.SkAppContext.RWSecProductLock.RUnlock()

	v, ok := conf.Conf.SecKill.SecProductInfoMap[productId]
	if !ok {
		return nil
	}

	data := make(map[string]interface{})
	data["product_id"] = productId
	data["start_time"] = v.StartTime
	data["end_time"] = v.EndTime
	data["status"] = v.Status

	return data
}

func (svc *SkAppService) SecInfoList() ([]map[string]interface{}, int, error) {
	config.SkAppContext.RWSecProductLock.RLock()
	defer config.SkAppContext.RWSecProductLock.RUnlock()

	var data []map[string]interface{}
	for _, v := range conf.Conf.SecKill.SecProductInfoMap {
		item, _, err := SecInfoById(v.ProductId)
		if err != nil {
			log.Printf("get sec info, err : %v", err)
			continue
		}
		data = append(data, item)
	}
	return data, 0, nil
}

func SecInfoById(productId int) (map[string]interface{}, int, error) {
	v, ok := conf.Conf.SecKill.SecProductInfoMap[productId]
	if !ok {
		return nil, svc_err.ErrNotFoundProductId, fmt.Errorf("not found product id: %d", productId)
	}

	start := false
	end := false
	status := "status"
	var err error
	var code int
	data := make(map[string]interface{})
	nowTime := time.Now().Unix()

	if nowTime-v.StartTime < 0 {
		start = false
		end = false
		status = "second kill not start"
		code = svc_err.ErrActiveNotStart
		err = fmt.Errorf(status)
	}

	if nowTime-v.EndTime > 0 {
		start = false
		end = false
		status = "second kill is already end"
		code = svc_err.ErrActiveAlreadyEnd
		err = fmt.Errorf(status)
	}

	if v.Status == config.ProductStatusForceSaleOut || v.Status == config.ProductStatusSaleOut {
		start = false
		end = false
		status = "product is sale out"
		code = svc_err.ErrActiveAlreadyEnd
		err = fmt.Errorf(status)
	}

	curRate := rand.Float64()
	/**
	* 放大于购买比率的1.5倍的请求进入core层
	 */
	if curRate > v.BuyRate*1.5 {
		start = false
		end = false
		status = "retry"
		code = svc_err.ErrRetry
		err = fmt.Errorf(status)
	}

	data = map[string]interface{}{
		"product_id": productId,
		"start":      start,
		"end":        end,
		"status":     status,
	}

	return data, code, err
}

func (svc *SkAppService) SecKill(req *model.SecRequest) (map[string]interface{}, int, error) {
	var code int
	// TODO: svc_limit

	userKey := fmt.Sprintf("%d_%d", req.UserId, req.ProductId)
	resultChan := make(chan *model.SecResult, 1)
	config.SkAppContext.UserConnMapLock.Lock()
	config.SkAppContext.UserConnMap[userKey] = resultChan
	config.SkAppContext.UserConnMapLock.Unlock()

	// 將請求送入通道並且推入到 redis queue 中
	config.SkAppContext.SecReqChan <- req

	ticker := time.NewTicker(time.Second + time.Duration(conf.Conf.SecKill.AppWaitResultTimeout))
	defer func() {
		ticker.Stop()
		config.SkAppContext.UserConnMapLock.Lock()
		delete(config.SkAppContext.UserConnMap, userKey)
		config.SkAppContext.UserConnMapLock.Unlock()
	}()

	data, code, err := SecInfoById(req.ProductId)
	if err != nil {
		log.Printf("userId[%d] secInfoById Id failed, req[%v]", req.UserId, req)
		return nil, code, err
	}

	select {
	case <-ticker.C:
		code = svc_err.ErrProcessTimeout
		err = fmt.Errorf("request timeout")
		return nil, code, err
	case <-req.CloseNotify:
		code = svc_err.ErrClientClosed
		err = fmt.Errorf("client already closed")
		return nil, code, err
	case ret := <-resultChan:
		code = ret.Code
		if code != 1002 {
			return data, code, svc_err.GetErrMsg(code)
		}
		log.Printf("sec kill success\n")
		data["product_id"] = ret.ProductId
		data["token"] = ret.Token
		data["user_id"] = ret.UserId
		return data, code, nil
	}
}

type ServiceMiddleware func(Service) Service
