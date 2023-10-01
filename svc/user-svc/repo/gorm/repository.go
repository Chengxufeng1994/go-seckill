package gorm

import (
	"github.com/Chengxufeng1994/go-seckill/svc/user-svc/entity"
	"github.com/Chengxufeng1994/go-seckill/svc/user-svc/model"
	"gorm.io/gorm"
	"time"
)

type repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) entity.Repository {
	return &repository{
		db: db,
	}
}

func (repo *repository) GetUser(id uint) (model.User, error) {
	var user model.User
	result := repo.db.
		Model(&model.User{}).
		Where("id = ?", id).
		First(&user)
	if result.Error != nil {
		return user, result.Error
	}

	return user, nil
}

func (repo *repository) GetUserByUsername(username string) (model.User, error) {
	var user model.User
	result := repo.db.
		Model(&model.User{}).
		Where("username = ?", username).
		First(&user)
	if result.Error != nil {
		return user, result.Error
	}

	return user, nil
}

func (repo *repository) GetUsers() ([]model.User, error) {
	var dat []model.User
	err := repo.db.Model(&model.User{}).Find(&dat).Error
	if err != nil {
		return nil, err
	}

	users := make([]model.User, len(dat))
	for i, v := range dat {
		users[i] = v
	}

	return users, nil
}

func (repo *repository) CreateUser(user model.User) (uint, error) {
	user.CreatedAt = time.Now()
	result := repo.db.Model(&model.User{}).Create(&user)
	if result.Error != nil {
		return 0, result.Error
	}

	return user.ID, nil
}

func (repo *repository) UpdateUser(user model.User) error {
	result := repo.db.Save(&user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *repository) DeleteUser(id uint) error {
	result := repo.db.Delete(&model.User{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
