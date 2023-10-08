package gorm

import (
	"context"
	"github.com/Chengxufeng1994/go-seckill/svc/oauth-svc/entity"
	"gorm.io/gorm"
)

type ClientDetailsRepository struct {
	db *gorm.DB
}

func NewClientDetailsRepository(db *gorm.DB) entity.ClientDetailsRepository {
	return &ClientDetailsRepository{
		db: db,
	}
}

func (r *ClientDetailsRepository) getTableName() string {
	return "client_details"
}

func (r *ClientDetailsRepository) GetClientDetailsByClientId(ctx context.Context, clientId string) (*entity.ClientDetails, error) {
	var clientDetails entity.ClientDetails
	err := r.db.WithContext(ctx).Table(r.getTableName()).Where("client_id = ?", clientId).Find(&clientDetails).Error
	if err != nil {
		return nil, err
	}

	return &clientDetails, nil
}

func (r *ClientDetailsRepository) CreateClientDetails(ctx context.Context, clientDetails *entity.ClientDetails) (uint, error) {
	result := r.db.WithContext(ctx).Table(r.getTableName()).Create(&clientDetails)
	if result.Error != nil {
		return 0, result.Error
	}
	if result.RowsAffected == 0 {
		return 0, result.Error
	}

	return clientDetails.ID, nil
}
