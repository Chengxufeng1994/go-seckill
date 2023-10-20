package config

import (
	"fmt"
	"github.com/Chengxufeng1994/go-seckill/pkg/bootstrap"
	pkgconfig "github.com/Chengxufeng1994/go-seckill/pkg/config"
	"github.com/Chengxufeng1994/go-seckill/svc/sk-app-svc/model"
	"github.com/go-kit/kit/log"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"gorm.io/gorm"
	"os"
	"strconv"
	"sync"
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
	if err := initializeZipkin(zipkinUrl); err != nil {
		Logger.Log("failed to load initialize zipkin", err)
		os.Exit(1)
	}
}

func initializeZipkin(zipkinUrl string) error {
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

var SkAppContext = &SkAppCtx{
	UserConnMap: make(map[string]chan *model.SecResult, 1024),
	SecReqChan:  make(chan *model.SecRequest, 1024),
}

type SkAppCtx struct {
	SecReqChan       chan *model.SecRequest
	SecReqChanSize   int
	RWSecProductLock sync.RWMutex

	UserConnMap     map[string]chan *model.SecResult
	UserConnMapLock sync.Mutex
}

const (
	// ProductStatusNormal 商品狀態正常
	ProductStatusNormal = 0
	// ProductStatusSaleOut 商品狀態賣光
	ProductStatusSaleOut = 1
	// ProductStatusForceSaleOut 商品強制賣光
	ProductStatusForceSaleOut = 2
)
