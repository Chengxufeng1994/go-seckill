package worker

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	conf "github.com/Chengxufeng1994/go-seckill/pkg/config"
	"github.com/Chengxufeng1994/go-seckill/svc/sk-core-svc/config"
	"github.com/Chengxufeng1994/go-seckill/svc/sk-core-svc/model"
	srv_err "github.com/Chengxufeng1994/go-seckill/svc/sk-core-svc/service/svc_err"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

type Distributor interface {
	HandleWrite()
	StartHandleWrite(int)
	StartHandleUser(int)
}

type RedisDistributor struct {
	client *redis.Client
}

func NewDistributor(client *redis.Client) Distributor {
	return &RedisDistributor{
		client: client,
	}
}

func (dist *RedisDistributor) StartHandleWrite(num int) {
	for i := 0; i < num; i++ {
		go dist.HandleWrite()
	}
}

func (dist *RedisDistributor) StartHandleUser(num int) {
	for i := 0; i < num; i++ {
		go dist.HandleUser()
	}
}

func (dist *RedisDistributor) HandleWrite() {
	log.Println("[distributor] handle write running")

	for res := range config.SecLayerCtx.Handle2WriteChan {
		fmt.Println("===", res)
		err := dist.sendToRedis(res)
		if err != nil {
			log.Printf("[distributor] send to redis, err : %v, res : %v", err, res)
			continue
		}
	}
}

func (dist *RedisDistributor) sendToRedis(res *config.SecResult) (err error) {
	data, err := json.Marshal(res)
	if err != nil {
		log.Printf("[distributor] marshal failed, err : %v", err)
		return
	}

	fmt.Printf("[distributor] before lpush %v", conf.Conf.Redis.Layer2proxyQueueName)
	conn := conf.Conf.Redis.RedisClient
	err = conn.LPush(context.Background(), conf.Conf.Redis.Layer2proxyQueueName, string(data)).Err()
	fmt.Println("[distributor] after lpush")
	if err != nil {
		log.Printf("[distributor] lpush layer to proxy redis queue failed, err : %v", err)
		return
	}
	log.Printf("[distributor] lpush layer to proxy success. data[%v]", string(data))

	return
}

func (dist *RedisDistributor) HandleUser() {
	log.Println("[processor] handle user running")
	for req := range config.SecLayerCtx.Read2HandleChan {
		log.Printf("[processor] begin process request : %v", req)
		res, err := HandleSeckill(req)
		if err != nil {
			log.Printf("[processor] process request %v failed, err : %v", err)
			res = &config.SecResult{
				Code: srv_err.ErrServiceBusy,
			}
		}
		fmt.Println("process... ", res)
		timer := time.NewTicker(time.Millisecond * time.Duration(conf.Conf.SecKill.SendToWriteChanTimeout))
		select {
		case config.SecLayerCtx.Handle2WriteChan <- res:
		case <-timer.C:
			log.Printf("send to response chan timeout, res : %v", res)
			break
		}
	}
	return
}

func HandleSeckill(req *config.SecRequest) (res *config.SecResult, err error) {
	config.SecLayerCtx.RWSecProductLock.RLock()
	defer config.SecLayerCtx.RWSecProductLock.RUnlock()

	res = &config.SecResult{}
	res.ProductId = req.ProductId
	res.UserId = req.UserId

	product, ok := conf.Conf.SecKill.SecProductInfoMap[req.ProductId]
	if !ok {
		log.Printf("not found product : %v", req.ProductId)
		res.Code = srv_err.ErrNotFoundProduct
		return
	}

	if product.Status == srv_err.ProductStatusSoldout {
		res.Code = srv_err.ErrSoldout
		return
	}
	nowTime := time.Now().Unix()

	config.SecLayerCtx.HistoryMapLock.Lock()
	userHistory, ok := config.SecLayerCtx.HistoryMap[req.UserId]
	if !ok {
		userHistory = &model.UserBuyHistory{
			History: make(map[int]int, 16),
		}
		config.SecLayerCtx.HistoryMap[req.UserId] = userHistory
	}
	historyCount := userHistory.GetProductBuyCount(req.ProductId)
	config.SecLayerCtx.HistoryMapLock.Unlock()

	if historyCount >= product.OnePersonBuyLimit {
		res.Code = srv_err.ErrAlreadyBuy
		return
	}

	curSoldCount := config.SecLayerCtx.ProductCountMgr.Count(req.ProductId)

	if curSoldCount >= product.Total {
		res.Code = srv_err.ErrSoldout
		product.Status = srv_err.ProductStatusSoldout
		return
	}

	//curRate := rand.Float64()
	curRate := 0.1
	fmt.Println(curRate, product.BuyRate)
	if curRate > product.BuyRate {
		res.Code = srv_err.ErrRetry
		return
	}

	userHistory.Add(req.ProductId, 1)
	config.SecLayerCtx.ProductCountMgr.Add(req.ProductId, 1)

	// 用戶Id, 商品Id, 當前時間, 密鑰
	res.Code = srv_err.ErrSecKillSucc
	tokenData := fmt.Sprintf("userId=%d&productId=%d&timestamp=%d&security=%s", req.UserId, req.ProductId, nowTime, conf.Conf.SecKill.TokenPassWd)
	res.Token = fmt.Sprintf("%x", md5.Sum([]byte(tokenData))) //MD5加密
	res.TokenTime = nowTime

	return
}
