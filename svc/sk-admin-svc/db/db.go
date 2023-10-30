package db

import (
	"github.com/Chengxufeng1994/go-seckill/svc/sk-admin-svc/config"
	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	Db *gorm.DB
)

func init() {
}

func InitializeDB(dsn string) error {
	config.Logger.Log("dsn", dsn)
	var err error
	Db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Discard,
	})

	if err != nil {
		config.Logger.Log("err", "failed to connect database", "reason", err)
		return errors.Wrap(err, "initializeDB")
	}

	return err
}
