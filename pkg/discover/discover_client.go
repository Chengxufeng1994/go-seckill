package discover

import (
	"github.com/Chengxufeng1994/go-seckill/pkg/common"
	"log"
)

type DiscoveryClient interface {
	/**
	 * 服务注册接口
	 * @param serviceName 服务名
	 * @param instanceId 服务实例Id
	 * @param instancePort 服务实例端口
	 * @param healthCheckUrl 健康检查地址
	 * @param weight 权重
	 * @param meta 服务实例元数据
	 */
	Register(svcName, instanceId, healthcheckUrl, svcHost string, svcPort int, weight int, meta map[string]string, tags []string, logger *log.Logger) bool

	/**
	 * 服务注销接口
	 * @param instanceId 服务实例Id
	 */
	DeRegister(instanceId string, logger *log.Logger) bool

	/**
	 * 发现服务实例接口
	 * @param serviceName 服务名
	 */
	DiscoverServices(serviceName string, logger *log.Logger) []*common.ServiceInstance
}
