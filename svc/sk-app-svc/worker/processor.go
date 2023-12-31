package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	conf "github.com/Chengxufeng1994/go-seckill/pkg/config"
	"github.com/Chengxufeng1994/go-seckill/svc/sk-app-svc/config"
	"github.com/Chengxufeng1994/go-seckill/svc/sk-app-svc/model"
	"github.com/redis/go-redis/v9"
)

type Processor interface {
	ReadHandle()
	Start(int)
}

type RedisProcessor struct {
	client *redis.Client
}

func NewProcessor(client *redis.Client) Processor {
	return &RedisProcessor{
		client: client,
	}
}

func (proc *RedisProcessor) ReadHandle() {
	for {
		data, err := proc.client.BRPop(context.Background(), time.Second, conf.Conf.Redis.Layer2proxyQueueName).Result()
		if err != nil {
			continue
		}

		var result *model.SecResult
		err = json.Unmarshal([]byte(data[1]), &result)
		if err != nil {
			config.Logger.Log("svc_err", fmt.Sprintf("json.Unmarshal req failed. svc_err: %v", err))
			continue
		}

		userKey := fmt.Sprintf("%d_%d", result.UserId, result.ProductId)
		fmt.Println(userKey)
		config.SkAppContext.UserConnMapLock.Lock()
		resultChan, ok := config.SkAppContext.UserConnMap[userKey]
		config.SkAppContext.UserConnMapLock.Unlock()
		if !ok {
			log.Printf("user not found: %v\n", userKey)
			continue
		}
		log.Printf("request result send to chan\n")

		resultChan <- result
		log.Printf("request result send to chan success, userKey:%v\n", userKey)
	}
}

func (proc *RedisProcessor) Start(num int) {
	for i := 0; i < num; i++ {
		go proc.ReadHandle()
	}
}
