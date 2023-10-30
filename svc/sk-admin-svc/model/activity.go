package model

import (
	"time"

	"github.com/Chengxufeng1994/go-seckill/svc/sk-admin-svc/entity"
)

const (
	ActivityStatusNormal  = 0
	ActivityStatusDisable = 1
	ActivityStatusExpire  = 2
)

type Activity struct {
	ActivityId   int    `json:"activity_id"`   //活动Id
	ActivityName string `json:"activity_name"` //活动名称
	ProductId    int    `json:"product_id"`    //商品Id
	StartTime    int64  `json:"start_time"`    //开始时间
	EndTime      int64  `json:"end_time"`      //结束时间
	Total        int    `json:"total"`         //商品总数
	Status       int    `json:"status"`        //状态

	StartTimeStr string  `json:"start_time_str"`
	EndTimeStr   string  `json:"end_time_str"`
	StatusStr    string  `json:"status_str"`
	Speed        int     `json:"speed"`
	BuyLimit     int     `json:"buy_limit"`
	BuyRate      float64 `json:"buy_rate"`
}

type SecProductInfoConf struct {
	ProductId         int     `json:"product_id"`           //商品Id
	StartTime         int64   `json:"start_time"`           //开始时间
	EndTime           int64   `json:"end_time"`             //结束时间
	Status            int     `json:"status"`               //状态
	Total             int     `json:"total"`                //商品总数
	Left              int     `json:"left"`                 //剩余商品数
	OnePersonBuyLimit int     `json:"one_person_buy_limit"` //一个人购买限制
	BuyRate           float64 `json:"buy_rate"`             //买中几率
	SoldMaxLimit      int     `json:"sold_max_limit"`       //每秒最多能卖多少个
}

func ActivityDaoToDto(dao entity.Activity) *Activity {
	startTime := dao.StartTime
	startUnixTimeUTC := time.Unix(startTime, 0)
	startTimeStr := startUnixTimeUTC.Format(time.RFC3339)

	endTime := dao.EndTime
	endUnixTimeUTC := time.Unix(endTime, 0)
	endTimeStr := endUnixTimeUTC.Format(time.RFC3339)

	status := dao.Status
	var statusStr string
	if status == ActivityStatusNormal {
		statusStr = "正常"
	} else if status == ActivityStatusDisable {
		statusStr = "已禁用"
	}

	return &Activity{
		ActivityId:   dao.ActivityId,
		ActivityName: dao.ActivityName,
		ProductId:    dao.ProductId,
		StartTime:    startTime,
		EndTime:      endTime,
		Total:        dao.Total,
		Status:       dao.Status,

		StartTimeStr: startTimeStr,
		EndTimeStr:   endTimeStr,
		StatusStr:    statusStr,
		Speed:        dao.SecSpeed,
		BuyLimit:     dao.BuyLimit,
		BuyRate:      dao.BuyRate,
	}
}

func ActivityDtoToDao(dao *Activity) *entity.Activity {
	return &entity.Activity{
		ActivityId:   dao.ActivityId,
		ActivityName: dao.ActivityName,
		ProductId:    dao.ProductId,
		StartTime:    dao.StartTime,
		EndTime:      dao.EndTime,
		Total:        dao.Total,
		Status:       dao.Status,
		SecSpeed:     dao.Speed,
		BuyLimit:     dao.BuyLimit,
		BuyRate:      dao.BuyRate,
	}
}
