package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	conf "github.com/Chengxufeng1994/go-seckill/pkg/config"
	"github.com/Chengxufeng1994/go-seckill/svc/sk-app-svc/config"
	"github.com/redis/go-redis/v9"
)

type Distributor interface {
	WriteHandle()
	Start(int)
}

type RedisDistributor struct {
	client *redis.Client
}

func NewDistributor(client *redis.Client) Distributor {
	return &RedisDistributor{
		client: client,
	}
}

func (dist *RedisDistributor) WriteHandle() {
	for {
		config.Logger.Log("info", "write data to redis")

		req := <-config.SkAppContext.SecReqChan
		log.Printf(fmt.Sprintf("access time: %v", req.AccessTime))
		data, err := json.Marshal(req)
		if err != nil {
			log.Printf("json.Marashal req failed. svc_err: %v, req: %v", err, req)
			continue
		}

		err = dist.client.LPush(context.Background(), conf.Conf.Redis.Proxy2layerQueueName, string(data)).Err()
		if err != nil {
			log.Printf("dist lpush req failed. svc_err: %v, req: %v", err, req)
			continue
		}

		log.Printf("lpush req success. req: %v", string(data))
	}
}

func (dist *RedisDistributor) Start(num int) {
	for i := 0; i < num; i++ {
		go dist.WriteHandle()
	}
}
