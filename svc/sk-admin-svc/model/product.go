package model

import (
	"github.com/Chengxufeng1994/go-seckill/svc/sk-admin-svc/entity"
)

type Product struct {
	ProductId   int    `json:"product_id"`   //商品Id
	ProductName string `json:"product_name"` //商品名称
	Total       int    `json:"total"`        //商品数量
	Status      int    `json:"status"`       //商品状态
}

func ProductDaoToDto(dao *entity.Product) *Product {
	return &Product{
		ProductId:   dao.ProductId,
		ProductName: dao.ProductName,
		Total:       dao.Total,
		Status:      dao.Status,
	}
}

func ProductDtoToDao(dao *Product) *entity.Product {
	return &entity.Product{
		ProductId:   dao.ProductId,
		ProductName: dao.ProductName,
		Total:       dao.Total,
		Status:      dao.Status,
	}
}
