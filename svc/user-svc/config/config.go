package config

import (
	"fmt"
	"github.com/Chengxufeng1994/go-seckill/pkg/bootstrap"
	pkgconfig "github.com/Chengxufeng1994/go-seckill/pkg/config"
	"github.com/Chengxufeng1994/go-seckill/svc/user-svc/model"
	"github.com/go-kit/kit/log"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
	"strconv"
)

var (
	Logger       log.Logger
	Db           *gorm.DB
	ZipkinTracer *zipkin.Tracer
)

func init() {
	{
		Logger = log.NewLogfmtLogger(os.Stderr)
		Logger = log.With(Logger, "ts", log.DefaultTimestampUTC)
		Logger = log.With(Logger, "caller", log.DefaultCaller)
	}

	// TODO: LoadRemoteConfig
	if err := pkgconfig.LoadLocalConfig(); err != nil {
		Logger.Log("failed to load local config", err)
		os.Exit(1)
	}

	zipkinUrl := fmt.Sprintf("http://%s:%d%s",
		pkgconfig.Conf.Trace.Host,
		pkgconfig.Conf.Trace.Port,
		pkgconfig.Conf.Trace.Url,
	)
	Logger.Log("zipkin url", zipkinUrl)
	if err := initZipkin(zipkinUrl); err != nil {
		Logger.Log("failed to load initialize zipkin", err)
		os.Exit(1)
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Taipei",
		pkgconfig.Conf.Postgres.Host,
		pkgconfig.Conf.Postgres.Username,
		pkgconfig.Conf.Postgres.Password,
		pkgconfig.Conf.Postgres.Db,
		pkgconfig.Conf.Postgres.Port)
	Logger.Log("dsn", dsn)
	if err := initDB(dsn); err != nil {
		Logger.Log("err", "failed to connect database", "reason", err)
		os.Exit(1)
	}

	if err := initMigrate(); err != nil {
		Logger.Log("err", "failed to migrate database", "reason", err)
		os.Exit(1)
	}
}

func initZipkin(zipkinUrl string) error {
	var (
		err           error
		useNoopTracer = zipkinUrl == ""
		reporter      = zipkinhttp.NewReporter(zipkinUrl)
	)
	//defer reporter.Close()
	zEP, _ := zipkin.NewEndpoint(bootstrap.Conf.Discover.ServiceName, strconv.Itoa(bootstrap.Conf.Discover.Port))
	ZipkinTracer, err = zipkin.NewTracer(
		reporter, zipkin.WithLocalEndpoint(zEP), zipkin.WithNoopTracer(useNoopTracer),
	)
	if err != nil {
		Logger.Log("failed to create zipkin", err)
		return err
	}
	if !useNoopTracer {
		Logger.Log("tracer", "Zipkin", "type", "Native", "URL", zipkinUrl)
	}
	return nil
}

func initDB(dsn string) error {
	var err error
	Db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Discard,
	})
	if err != nil {
		return err
	}
	return nil
}

func initMigrate() error {
	var err error
	err = Db.Table("users").AutoMigrate(&model.User{})
	if err != nil {
		return err
	}
	return nil
}
