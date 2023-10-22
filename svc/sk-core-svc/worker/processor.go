package worker

import (
	"context"
	"encoding/json"
	"fmt"
	conf "github.com/Chengxufeng1994/go-seckill/pkg/config"
	"github.com/Chengxufeng1994/go-seckill/svc/sk-core-svc/config"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

type Processor interface {
	HandleRead()
	Start(num int)
}

type RedisProcessor struct {
	client *redis.Client
}

func NewProcessor(client *redis.Client) Processor {
	return &RedisProcessor{
		client: client,
	}
}

func (proc *RedisProcessor) HandleRead() {
	log.Println("[processor] read goroutine running")
	for {
		for {
			data, err := proc.client.BRPop(context.Background(), time.Second, conf.Conf.Redis.Proxy2layerQueueName).Result()
			if err != nil {
				continue
			}
			log.Printf("[processor] brpop from proxy to layer queue, data : %s\n", data)

			// 解析 SecRequest
			var req config.SecRequest
			err = json.Unmarshal([]byte(data[1]), &req)
			if err != nil {
				log.Printf("[processor] unmarshal to secrequest failed, err : %v", err)
				continue
			}

			// 判斷是否超時
			nowTime := time.Now().Unix()
			//int64(config.SecLayerCtx.SecLayerConf.MaxRequestWaitTimeout)
			fmt.Println(nowTime, " ", req.SecTime, " ", 100)
			if nowTime-req.SecTime >= int64(conf.Conf.SecKill.MaxRequestWaitTimeout) {
				log.Printf("[processor] req[%v] is expire", req)
				continue
			}

			// 設置超時時間
			timer := time.NewTicker(time.Millisecond * time.Duration(conf.Conf.SecKill.CoreWaitResultTimeout))
			select {
			case config.SecLayerCtx.Read2HandleChan <- &req:
			case <-timer.C:
				log.Printf("[processor] send to handle chan timeout, req : %v", req)
				break
			}
		}
	}
}

func (proc *RedisProcessor) Start(num int) {
	for i := 0; i < num; i++ {
		go proc.HandleRead()
	}
}
