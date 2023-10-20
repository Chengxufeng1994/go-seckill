package setup

import (
	"context"
	"fmt"
	conf "github.com/Chengxufeng1994/go-seckill/pkg/config"
	"github.com/Chengxufeng1994/go-seckill/svc/sk-app-svc/config"
	"github.com/redis/go-redis/v9"
	"log"
	"strconv"
	"time"
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

	loadBlackList(conf.Conf.Redis.RedisClient)

	return nil
}

func loadBlackList(conn *redis.Client) {
	conf.Conf.SecKill.IPBlackMap = make(map[string]struct{})
	conf.Conf.SecKill.IDBlackMap = make(map[int]struct{})

	result, err := conn.HGetAll(context.Background(), conf.Conf.Redis.IdBlackListHash).Result()
	if err != nil {
		config.Logger.Log("svc_err", "hget all failed, svc_err: "+err.Error())
		return
	}

	for _, v := range result {
		id, err := strconv.Atoi(v)
		if err != nil {
			config.Logger.Log("svc_err", "invalid user id: "+v)
			continue
		}
		conf.Conf.SecKill.IDBlackMap[id] = struct{}{}
	}

	result, err = conn.HGetAll(context.Background(), conf.Conf.Redis.IpBlackListHash).Result()
	if err != nil {
		config.Logger.Log("svc_err", "hget all failed, svc_err: "+err.Error())
		return
	}

	for _, v := range result {
		conf.Conf.SecKill.IPBlackMap[v] = struct{}{}
	}

	//go syncIdBlackList(conn)
	//go syncIpBlackList(conn)
}

// 同步用户ID黑名单
func syncIdBlackList(conn *redis.Client) {
	for {
		idArr, err := conn.BRPop(context.Background(), time.Minute, conf.Conf.Redis.IdBlackListQueue).Result()
		if err != nil {
			log.Printf("brpop id failed, err : %v", err)
			continue
		}
		id, _ := strconv.Atoi(idArr[1])
		conf.Conf.SecKill.RWBlackLock.Lock()
		conf.Conf.SecKill.IDBlackMap[id] = struct{}{}
		conf.Conf.SecKill.RWBlackLock.Unlock()
	}
}

// 同步用户IP黑名单
func syncIpBlackList(conn *redis.Client) {
	var ipList []string
	lastTime := time.Now().Unix()

	for {
		ipArr, err := conn.BRPop(context.Background(), time.Minute, conf.Conf.Redis.IpBlackListQueue).Result()
		if err != nil {
			log.Printf("brpop ip failed, err : %v", err)
			continue
		}

		ip := ipArr[1]
		curTime := time.Now().Unix()
		ipList = append(ipList, ip)

		if len(ipList) > 100 || curTime-lastTime > 5 {
			conf.Conf.SecKill.RWBlackLock.Lock()
			{
				for _, v := range ipList {
					conf.Conf.SecKill.IPBlackMap[v] = struct{}{}
				}
			}
			conf.Conf.SecKill.RWBlackLock.Unlock()

			lastTime = curTime
			log.Printf("sync ip list from redis success, ip[%v]", ipList)
		}
	}
}
