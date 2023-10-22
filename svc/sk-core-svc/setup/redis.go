package setup

import (
	"context"
	"fmt"
	conf "github.com/Chengxufeng1994/go-seckill/pkg/config"
	"github.com/Chengxufeng1994/go-seckill/svc/sk-app-svc/config"
	"github.com/redis/go-redis/v9"
)

func InitializeRedis() error {
	config.Logger.Log("info", "initialize redis")
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.Conf.Redis.Host, conf.Conf.Redis.Port),
		Password: conf.Conf.Redis.Password,
		DB:       conf.Conf.Redis.Db,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		config.Logger.Log("err", "connect redis failed. svc_err: "+err.Error())
		return err
	}

	conf.Conf.Redis.RedisClient = client
	config.Logger.Log("info", "initialize redis success")

	return nil
}
