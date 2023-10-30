package gorm

import (
	"github.com/Chengxufeng1994/go-seckill/svc/sk-admin-svc/entity"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) entity.Repository {
	return &Repository{
		db: db,
	}
}

func (r Repository) GetActivityList() ([]*entity.Activity, error) {
	var dat []*entity.Activity
	err := r.db.
		Table("activity").
		Model(&entity.Activity{}).
		Find(&dat).
		Order("activity_id desc").
		Error
	if err != nil {
		return nil, err
	}

	activities := make([]*entity.Activity, len(dat))
	for i, v := range dat {
		activities[i] = v
	}

	return activities, nil
}

func (r Repository) CreateActivity(activity *entity.Activity) error {
	err := r.db.
		Table("activity").
		Model(&entity.Activity{}).
		Create(activity).Error
	if err != nil {
		return err
	}

	return nil
}

func (r Repository) GetProductList() ([]*entity.Product, error) {
	var dat []*entity.Product
	err := r.db.
		Table("product").
		Model(&entity.Product{}).
		Find(&dat).
		Order("product_id desc").
		Error
	if err != nil {
		return nil, err
	}

	products := make([]*entity.Product, len(dat))
	for i, v := range dat {
		products[i] = v
	}

	return products, nil
}

func (r Repository) CreateProduct(product *entity.Product) error {
	err := r.db.
		Table("product").
		Model(&entity.Product{}).
		Create(product).Error
	if err != nil {
		return err
	}

	return nil
}
