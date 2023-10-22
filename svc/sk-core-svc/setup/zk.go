package setup

import (
	"encoding/json"
	"fmt"
	conf "github.com/Chengxufeng1994/go-seckill/pkg/config"
	"github.com/Chengxufeng1994/go-seckill/svc/sk-app-svc/config"
	"github.com/go-zookeeper/zk"
	"log"
	"time"
)

func InitializeZk() error {
	hosts := []string{
		fmt.Sprintf("%v:%d", conf.Conf.Zk.Host, conf.Conf.Zk.Port),
	}
	option := zk.WithEventCallback(waitSecProductEvent)
	conn, _, err := zk.Connect(hosts, time.Second*5, option)
	if err != nil {
		config.Logger.Log("err", "connect zookeeper failed. svc_err: "+err.Error())
		return err
	}
	conf.Conf.Zk.ZkConn = conn
	conf.Conf.Zk.SecProductKey = "/product"
	loadSecConf(conn)

	return nil
}

func loadSecConf(conn *zk.Conn) {
	v, _, err := conn.Get(conf.Conf.Zk.SecProductKey) //conf.Etcd.EtcdSecProductKey
	if err != nil {
		config.Logger.Log("err", "get product info failed, err:"+err.Error())
		return
	}

	config.Logger.Log("info", "get product info")
	var secProductInfo []*conf.SecProductInfoConf
	err = json.Unmarshal(v, &secProductInfo)
	if err != nil {
		config.Logger.Log("err", "unmarshal second product info failed, err:"+err.Error())
	}
	updateSecProductInfo(secProductInfo)
}

func waitSecProductEvent(event zk.Event) {
	log.Print(">>>>>>>>>>>>>>>>>>>")
	log.Println("path:", event.Path)
	log.Println("type:", event.Type.String())
	log.Println("state:", event.State.String())
	log.Println("<<<<<<<<<<<<<<<<<<<")
	if event.Path == conf.Conf.Zk.SecProductKey {

	}
}

func updateSecProductInfo(secProductInfo []*conf.SecProductInfoConf) {
	tmp := make(map[int]*conf.SecProductInfoConf, 1024)
	for _, v := range secProductInfo {
		log.Printf("update sec product info %v\n", v)
		tmp[v.ProductId] = v
	}

	conf.Conf.SecKill.RWBlackLock.Lock()
	defer conf.Conf.SecKill.RWBlackLock.Unlock()
	conf.Conf.SecKill.SecProductInfoMap = tmp
}
