package setup

import (
	"fmt"
	conf "github.com/Chengxufeng1994/go-seckill/pkg/config"
	register "github.com/Chengxufeng1994/go-seckill/pkg/discover"
	"github.com/Chengxufeng1994/go-seckill/svc/sk-core-svc/config"
	"github.com/Chengxufeng1994/go-seckill/svc/sk-core-svc/worker"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func InitializeService() {
	errChan := make(chan error)
	config.Logger.Log("info", "initialize redis processor")
	config.Logger.Log("info", "CoreWriteToHandleGoroutineNum: "+strconv.Itoa(conf.Conf.SecKill.CoreWriteRedisGoroutineNum))
	dist := worker.NewDistributor(conf.Conf.Redis.RedisClient)
	dist.StartHandleWrite(conf.Conf.SecKill.CoreWriteRedisGoroutineNum)

	config.Logger.Log("info", "CoreHandleGoroutineNum: "+strconv.Itoa(conf.Conf.SecKill.CoreHandleGoroutineNum))
	dist.StartHandleUser(conf.Conf.SecKill.CoreHandleGoroutineNum)

	config.Logger.Log("info", "CoreReadToHandleGoroutineNum: "+strconv.Itoa(conf.Conf.SecKill.CoreReadRedisGoroutineNum))
	proc := worker.NewProcessor(conf.Conf.Redis.RedisClient)
	proc.Start(conf.Conf.SecKill.CoreReadRedisGoroutineNum)

	register.Register()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	err := <-errChan
	register.Deregister()
	fmt.Println(err)
}
