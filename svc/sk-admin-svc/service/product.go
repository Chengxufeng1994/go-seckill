package service

import (
	"log"

	"github.com/Chengxufeng1994/go-seckill/svc/sk-admin-svc/entity"
	"github.com/Chengxufeng1994/go-seckill/svc/sk-admin-svc/model"
)

type ProductService interface {
	CreateProduct(product *model.Product) error
	GetProductList() ([]*model.Product, error)
}

type ProductServiceImpl struct {
	repo entity.Repository
}

func NewProductService(repo entity.Repository) ProductService {
	return &ProductServiceImpl{
		repo: repo,
	}
}

func (p ProductServiceImpl) CreateProduct(product *model.Product) error {
	dao := model.ProductDtoToDao(product)
	if err := p.repo.CreateProduct(dao); err != nil {
		log.Printf("productEntity.CreateProduct failed, err: %v", err)
		return err
	}

	return nil
}

func (p ProductServiceImpl) GetProductList() ([]*model.Product, error) {
	list, err := p.repo.GetProductList()
	if err != nil {
		log.Printf("productEntity.GetProduct failed, err: %v", err)
		return nil, err
	}
	ret := make([]*model.Product, len(list))
	for i, v := range list {
		ret[i] = model.ProductDaoToDto(v)
	}

	return ret, nil
}

type ProductServiceMiddleware func(ProductService) ProductService
