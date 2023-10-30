package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	conf "github.com/Chengxufeng1994/go-seckill/pkg/config"
	"github.com/Chengxufeng1994/go-seckill/svc/sk-admin-svc/entity"
	"github.com/Chengxufeng1994/go-seckill/svc/sk-admin-svc/model"
	"github.com/go-zookeeper/zk"
)

type ActivityService interface {
	GetActivityList() ([]*model.Activity, error)
	CreateActivity(*model.Activity) error
}

type ActivityServiceImpl struct {
	repo entity.Repository
}

func NewActivityService(repo entity.Repository) ActivityService {
	return &ActivityServiceImpl{
		repo: repo,
	}
}

func (svc ActivityServiceImpl) GetActivityList() ([]*model.Activity, error) {
	list, err := svc.repo.GetActivityList()
	if err != nil {
		log.Printf("ActivityEntity.GetActivityList, err : %v", err)
		return nil, err
	}

	ret := make([]*model.Activity, len(list))
	for i, item := range list {
		ret[i] = model.ActivityDaoToDto(*item)
	}

	log.Printf("get activity success, activity list is [%v]", ret)
	return ret, nil
}

func (svc ActivityServiceImpl) CreateActivity(activity *model.Activity) error {
	var err error
	dao := model.ActivityDtoToDao(activity)
	err = svc.repo.CreateActivity(dao)
	if err != nil {
		log.Printf("ActivityModel.CreateActivity, err : %v", err)
		return err
	}

	log.Printf("start sycn to zookeeper")
	err = svc._syncToZk(activity)
	if err != nil {
		log.Printf("activity product info sync to zookeeperj failed, err : %v", err)
		return err
	}

	return nil
}

func (svc ActivityServiceImpl) _syncToZk(activity *model.Activity) error {
	key := conf.Conf.Zk.SecProductKey
	secProductInfoList, err := svc._loadProductFromZk(key)
	if err != nil {
		secProductInfoList = []*model.SecProductInfoConf{}
	}
	secProductInfo := &model.SecProductInfoConf{}
	secProductInfo.EndTime = activity.EndTime
	secProductInfo.OnePersonBuyLimit = activity.BuyLimit
	secProductInfo.ProductId = activity.ProductId
	secProductInfo.SoldMaxLimit = activity.Speed
	secProductInfo.StartTime = activity.StartTime
	secProductInfo.Status = activity.Status
	secProductInfo.Total = activity.Total
	secProductInfo.BuyRate = activity.BuyRate
	secProductInfoList = append(secProductInfoList, secProductInfo)
	data, err := json.Marshal(secProductInfoList)
	if err != nil {
		log.Printf("_syncToZk json marshal failed, err : %v", err)
		return err
	}

	conn := conf.Conf.Zk.ZkConn
	{
		var byt = []byte(string(data))
		var flags int32 = 0
		var acls = zk.WorldACL(zk.PermAll)
		fmt.Println(acls)

		exists, _, _ := conn.Exists(key)
		if exists {
			if _, err := conn.Set(key, byt, flags); err != nil {
				log.Printf("_syncToZk set to zk failed, data = [%v]", err)
			}
		} else {
			if _, err := conn.Create(key, byt, flags, acls); err != nil {
				log.Printf("_syncToZk create to zk failed, data = [%v]", err)
			}
		}
	}

	log.Printf("_syncToZk put to zk success, data = [%v]", string(data))
	return nil
}

func (svc ActivityServiceImpl) _loadProductFromZk(key string) ([]*model.SecProductInfoConf, error) {
	_, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	conn := conf.Conf.Zk.ZkConn
	v, s, err := conn.Get(key)
	if err != nil {
		log.Printf("_loadProductFromZk get [%s] from zk failed, err : %v\n", key, err)
		return nil, err
	}

	log.Printf("_loadProductFromZk get from zk success, rsp : %v\n", s)
	var secProductInfo []*model.SecProductInfoConf
	fmt.Printf("value of path[%s]=[%s].\n", key, v)
	if err := json.Unmarshal(v, &secProductInfo); err != nil {
		log.Printf("_loadProductFromZk json unmarshal failed, err : %v\n", err)
		return nil, err
	}

	return secProductInfo, nil
}

type ActivityServiceMiddleware func(ActivityService) ActivityService
