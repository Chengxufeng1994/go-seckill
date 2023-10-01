package discover

import (
	"errors"
	"github.com/Chengxufeng1994/go-seckill/pkg/bootstrap"
	"github.com/Chengxufeng1994/go-seckill/pkg/common"
	"github.com/Chengxufeng1994/go-seckill/pkg/loadbalance"
	uuid "github.com/satori/go.uuid"
	"log"
	"os"
	"strconv"
)

var (
	ConsulService DiscoveryClient
	LoadBalance   loadbalance.LoadBalance
	Logger        *log.Logger
)

var ErrNoInstanceExisted = errors.New("no available client")

func init() {
	Logger = log.New(os.Stdout, "", log.LstdFlags)
	Logger.Println("[pkg] discover init.")
	instanceId := bootstrap.Conf.Discover.InstanceId + "-" + uuid.NewV4().String()
	ConsulService = New(bootstrap.Conf.Discover.Host, bootstrap.Conf.Discover.Port, instanceId)
	LoadBalance = loadbalance.NewRandomLoadBalance()
}

func DiscoveryService(serviceName string) (*common.ServiceInstance, error) {
	instances := ConsulService.DiscoverServices(serviceName, Logger)
	if len(instances) < 1 {
		Logger.Printf("[discover] no available client for %s.", serviceName)
		return nil, ErrNoInstanceExisted
	}

	return LoadBalance.SelectService(instances)
}

func Register() {
	// if consul service create failed, stop this service
	if ConsulService == nil {
		panic(0)
	}

	meta := make(map[string]string)
	meta["rpc"] = strconv.Itoa(bootstrap.Conf.Rpc.Port)

	if register := ConsulService.Register(
		bootstrap.Conf.Discover.ServiceName,
		bootstrap.Conf.Discover.InstanceId,
		"/health",
		bootstrap.Conf.Http.Host,
		bootstrap.Conf.Http.Port,
		bootstrap.Conf.Discover.Weight,
		meta,
		nil,
		Logger,
	); !register {
		Logger.Printf("[discover] register service %s-service failed.", bootstrap.Conf.Discover.ServiceName)
		panic(0)
	}

	Logger.Printf("[discover] register service %s-service success.", bootstrap.Conf.Discover.ServiceName)
}

func Deregister() {
	if ConsulService == nil {
		panic(0)
	}

	if deregister := ConsulService.DeRegister(
		bootstrap.Conf.Discover.InstanceId,
		Logger,
	); !deregister {
		Logger.Printf("[discover] deregister service %s-service failed.", bootstrap.Conf.Discover.ServiceName)
		panic(0)
	}

	Logger.Printf("[discover] deregister service %s-service success.", bootstrap.Conf.Discover.ServiceName)
}
