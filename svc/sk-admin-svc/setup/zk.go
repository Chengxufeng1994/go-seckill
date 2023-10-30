package setup

import (
	"fmt"
	conf "github.com/Chengxufeng1994/go-seckill/pkg/config"
	"github.com/go-zookeeper/zk"
	"time"
)

// InitializeZk initialize zookeeper
func InitializeZk() error {
	hosts := []string{
		fmt.Sprintf("%v:%d", conf.Conf.Zk.Host, conf.Conf.Zk.Port),
	}
	conn, _, err := zk.Connect(hosts, time.Second*5)
	if err != nil {
		return err
	}
	conf.Conf.Zk.ZkConn = conn
	conf.Conf.Zk.SecProductKey = "/product"
	return nil
}
