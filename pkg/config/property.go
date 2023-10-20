package config

import (
	"github.com/go-zookeeper/zk"
	"github.com/redis/go-redis/v9"
	"sync"
)

var Conf Config

type Config struct {
	Postgres Postgres `mapstructure:"postgres"`
	Trace    Trace    `mapstructure:"trace"`
	Redis    Redis    `mapstructure:"redis"`
	Zk       Zk       `mapstructure:"zookeeper"`
	SecKill  SecKill  `mapstructure:"seckill"`
}

type Postgres struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Db       string `mapstructure:"db"`
}

type Trace struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	Url  string `mapstructure:"url"`
}

type Zk struct {
	ZkConn        *zk.Conn
	SecProductKey string
	Host          string `mapstructure:"host"`
	Port          int    `mapstructure:"port"`
}

type Redis struct {
	RedisClient *redis.Client
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
	Username    string `mapstructure:"username"`
	Password    string `mapstructure:"password"`
	Db          int    `mapstructure:"db"`

	Proxy2layerQueueName string `mapstructure:"Proxy2layerQueueName"`
	Layer2proxyQueueName string `mapstructure:"Layer2ProxyQueueName"`
	IdBlackListHash      string `mapstructure:"IdBlackListHash"`
	IpBlackListHash      string `mapstructure:"IpBlackListHash"`
	IdBlackListQueue     string `mapstructure:"IdBlackListQueue"`
	IpBlackListQueue     string `mapstructure:"IpBlackListQueue"`
}

type SecKill struct {
	RWBlackLock sync.RWMutex

	IPBlackMap map[string]struct{}
	IDBlackMap map[int]struct{}

	SecProductInfoMap map[int]*SecProductInfoConf

	AppWriteToHandleGoroutineNum int `mapstructure:"AppWriteToHandleGoroutineNum"`
	AppReadToHandleGoroutineNum  int `mapstructure:"AppReadToHandleGoroutineNum"`

	AppWaitResultTimeout int `mapstructure:"AppWaitResultTimeout"`
}

// SecProductInfoConf 商品資訊配置
type SecProductInfoConf struct {
	ProductId         int     `json:"product_id"`           //商品ID
	StartTime         int64   `json:"start_time"`           //开始时间
	EndTime           int64   `json:"end_time"`             //结束时间
	Status            int     `json:"status"`               //状态
	Total             int     `json:"total"`                //商品总数量
	Left              int     `json:"left"`                 //商品剩余数量
	OnePersonBuyLimit int     `json:"one_person_buy_limit"` //单个用户购买数量限制
	BuyRate           float64 `json:"buy_rate"`             //购买频率限制
	SoldMaxLimit      int     `json:"sold_max_limit"`
	// todo: svc_err
	SecLimit *SecLimit `json:"sec_limit"` //限速控制
}

// 每秒限制
type SecLimit struct {
	count   int   //次数
	preTime int64 //上一次记录时间
}
