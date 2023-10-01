package discover

import (
	"github.com/Chengxufeng1994/go-seckill/pkg/common"
	"github.com/go-kit/kit/sd/consul"
	capi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
	"log"
	"strconv"
	"sync"
)

type KitDiscoverClient struct {
	Host       string // Consul Host
	Port       int    // Consul Port
	InstanceId string // Consul InstanceId
	client     consul.Client
	// 连接 consul 的配置
	config *capi.Config
	mutex  sync.Mutex
	// 服务实例缓存字段
	instancesMap sync.Map
}

func New(consulHost string, consulPort int, instanceId string) *KitDiscoverClient {
	consulConfig := capi.DefaultConfig()
	consulConfig.Address = consulHost + ":" + strconv.Itoa(consulPort)
	apiClient, err := capi.NewClient(consulConfig)
	if err != nil {
		return nil
	}

	client := consul.NewClient(apiClient)

	return &KitDiscoverClient{
		Host:       consulHost,
		Port:       consulPort,
		InstanceId: instanceId,
		client:     client,
		config:     consulConfig,
	}
}

func (cli *KitDiscoverClient) Register(serviceName, instanceId, healthCheckUrl string, instanceHost string, instancePort int, weight int, meta map[string]string, tags []string, logger *log.Logger) bool {
	registration := &capi.AgentServiceRegistration{
		ID:      cli.InstanceId,
		Name:    serviceName,
		Address: instanceHost,
		Port:    instancePort,
		Meta:    meta,
		Tags:    tags,
		Weights: &capi.AgentWeights{
			Passing: weight,
		},
		Check: &capi.AgentServiceCheck{
			DeregisterCriticalServiceAfter: "30s",
			HTTP:                           "http://" + instanceHost + ":" + strconv.Itoa(instancePort) + healthCheckUrl,
			Interval:                       "15s",
		},
	}
	err := cli.client.Register(registration)
	if err != nil {
		if logger != nil {
			logger.Println("register service failed.", err)
		}
		return false
	}
	if logger != nil {
		logger.Println("register service success.")
	}
	return true
}

func (cli *KitDiscoverClient) DeRegister(instanceId string, logger *log.Logger) bool {
	registration := &capi.AgentServiceRegistration{
		ID: cli.InstanceId,
	}

	err := cli.client.Deregister(registration)
	if err != nil {
		if logger != nil {
			logger.Println("deregister service failed.", err)
		}
		return false
	}
	if logger != nil {
		logger.Println("deregister service success.")
	}
	return true
}

func (cli *KitDiscoverClient) DiscoverServices(serviceName string, logger *log.Logger) []*common.ServiceInstance {
	instancesList, ok := cli.instancesMap.Load(serviceName)
	if ok {
		return instancesList.([]*common.ServiceInstance)
	}

	cli.mutex.Lock()
	defer cli.mutex.Unlock()
	instancesList, ok = cli.instancesMap.Load(serviceName)
	if ok {
		return instancesList.([]*common.ServiceInstance)
	} else {
		go func() {
			params := make(map[string]interface{})
			params["type"] = "service"
			params["service"] = serviceName
			wp, _ := watch.Parse(params)
			wp.Handler = func(idx uint64, data interface{}) {
				switch v := data.(type) {
				case []*capi.ServiceEntry:
					if len(v) == 0 {
						cli.instancesMap.Store(serviceName, []*common.ServiceInstance{})
						return
					}
					healthService := make([]*common.ServiceInstance, 0)
					for _, service := range v {
						if service.Checks.AggregatedStatus() == capi.HealthPassing {
							healthService = append(healthService, entryService2serviceInstance(service.Service))
						}
					}
					cli.instancesMap.Store(serviceName, healthService)
				}
			}

			defer wp.Stop()
			wp.Run(cli.config.Address)
		}()
	}

	entries, _, err := cli.client.Service(serviceName, "", false, nil)
	if err != nil {
		cli.instancesMap.Store(serviceName, []interface{}{})
		logger.Printf("discover service %s failed", serviceName)
		return nil
	}
	instances := make([]*common.ServiceInstance, len(entries))
	for i := range entries {
		instances[i] = entryService2serviceInstance(entries[i].Service)
	}

	cli.instancesMap.Store(serviceName, instances)
	return instances
}

func entryService2serviceInstance(svc *capi.AgentService) *common.ServiceInstance {
	rpcPort := svc.Port - 1
	if svc.Meta != nil {
		if rpcPortStr, ok := svc.Meta["rpcPort"]; ok {
			rpcPort, _ = strconv.Atoi(rpcPortStr)
		}
	}
	return &common.ServiceInstance{
		Host:     svc.Address,
		Port:     svc.Port,
		GrpcPort: rpcPort,
		Weight:   svc.Weights.Passing,
	}
}
