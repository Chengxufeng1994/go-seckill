package main

import (
	"fmt"
	"log"

	"github.com/Chengxufeng1994/go-seckill/pkg/bootstrap"
	_ "github.com/Chengxufeng1994/go-seckill/pkg/bootstrap"
	pkgconfig "github.com/Chengxufeng1994/go-seckill/pkg/config"
	"github.com/Chengxufeng1994/go-seckill/svc/sk-admin-svc/db"
	"github.com/Chengxufeng1994/go-seckill/svc/sk-admin-svc/setup"
)

func main() {
	var err error

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable search_path=%s TimeZone=Asia/Taipei",
		pkgconfig.Conf.Postgres.Host,
		pkgconfig.Conf.Postgres.Username,
		pkgconfig.Conf.Postgres.Password,
		pkgconfig.Conf.Postgres.Db,
		pkgconfig.Conf.Postgres.Port,
		"sec_kill",
	)
	err = db.InitializeDB(dsn)
	if err != nil {
		log.Fatal(err)
	}

	err = setup.InitializeZk()
	if err != nil {
		log.Fatal(err)
	}

	setup.InitializeService(bootstrap.Conf.Http.Host, bootstrap.Conf.Http.Port)
}
