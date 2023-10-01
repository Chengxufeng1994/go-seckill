package loadbalance

import (
	"errors"
	"github.com/Chengxufeng1994/go-seckill/pkg/common"
	"math/rand"
	"time"
)

var ErrServiceInstanceNotExist = errors.New("service instances are not exist")

type LoadBalance interface {
	SelectService([]*common.ServiceInstance) (*common.ServiceInstance, error)
}

type RandomLoadBalance struct {
}

func NewRandomLoadBalance() LoadBalance {
	return &RandomLoadBalance{}
}

func (lb *RandomLoadBalance) SelectService(instances []*common.ServiceInstance) (*common.ServiceInstance, error) {
	n := len(instances)
	if instances == nil || n == 0 {
		return nil, ErrServiceInstanceNotExist
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return instances[r.Intn(n)], nil
}
