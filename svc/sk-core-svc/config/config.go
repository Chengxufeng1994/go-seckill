package config

import (
	"fmt"
	"github.com/Chengxufeng1994/go-seckill/pkg/bootstrap"
	pkgconfig "github.com/Chengxufeng1994/go-seckill/pkg/config"
	"github.com/Chengxufeng1994/go-seckill/svc/sk-core-svc/model"
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

var SecLayerCtx = &SecLayerContext{
	Read2HandleChan:  make(chan *SecRequest, 1024),
	Handle2WriteChan: make(chan *SecResult, 1024),
	HistoryMap:       make(map[int]*model.UserBuyHistory, 1024),
	ProductCountMgr:  model.NewProductCountMgr(),
}
var CoreCtx = &SkAppCtx{}

type SecResult struct {
	ProductId int    `json:"product_id"` //商品ID
	UserId    int    `json:"user_id"`    //用户ID
	Token     string `json:"token"`      //Token
	TokenTime int64  `json:"token_time"` //Token生成时间
	Code      int    `json:"code"`       //状态码
}

type SecRequest struct {
	ProductId     int             `json:"product_id"` //商品ID
	Source        string          `json:"source"`
	AuthCode      string          `json:"auth_code"`
	SecTime       int64           `json:"sec_time"`
	Nance         string          `json:"nance"`
	UserId        int             `json:"user_id"`
	UserAuthSign  string          `json:"user_auth_sign"` //用户授权签名
	ClientAddr    string          `json:"client_addr"`
	ClientRefence string          `json:"client_refence"`
	CloseNotify   <-chan bool     `json:"-"`
	ResultChan    chan *SecResult `json:"-"`
}

type SkAppCtx struct {
	SecReqChan       chan *SecRequest
	SecReqChanSize   int
	RWSecProductLock sync.RWMutex

	UserConnMap     map[string]chan *SecResult
	UserConnMapLock sync.Mutex
}

const (
	ProductStatusNormal       = 0 //商品状态正常
	ProductStatusSaleOut      = 1 //商品售罄
	ProductStatusForceSaleOut = 2 //商品强制售罄
)

type SecLayerContext struct {
	RWSecProductLock sync.RWMutex

	WaitGroup sync.WaitGroup

	Read2HandleChan  chan *SecRequest
	Handle2WriteChan chan *SecResult

	HistoryMap     map[int]*model.UserBuyHistory
	HistoryMapLock sync.Mutex

	ProductCountMgr *model.ProductCountMgr //商品计数
}
