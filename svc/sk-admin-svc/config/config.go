package config

import (
	"fmt"
	"github.com/Chengxufeng1994/go-seckill/pkg/bootstrap"
	conf "github.com/Chengxufeng1994/go-seckill/pkg/config"
	"github.com/go-kit/kit/log"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"os"
	"strconv"
)

var (
	Logger       log.Logger
	ZipkinTracer *zipkin.Tracer
)

func init() {
	{
		Logger = log.NewLogfmtLogger(os.Stderr)
		Logger = log.With(Logger, "ts", log.DefaultTimestampUTC)
		Logger = log.With(Logger, "caller", log.DefaultCaller)
	}

	// TODO: LoadRemoteConfig
	if err := conf.LoadLocalConfig(); err != nil {
		Logger.Log("failed to load local config", err)
		os.Exit(1)
	}

	zipkinUrl := fmt.Sprintf("http://%s:%d%s",
		conf.Conf.Trace.Host,
		conf.Conf.Trace.Port,
		conf.Conf.Trace.Url,
	)
	Logger.Log("zipkin url", zipkinUrl)
	if err := initZipkin(zipkinUrl); err != nil {
		Logger.Log("failed to load initialize zipkin", err)
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
