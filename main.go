package main

import (
	"context"
	"runtime"

	"net/http"
	"os"

	"github.com/smartpet/websocket/business"
	"github.com/smartpet/websocket/constant"
	"github.com/smartpet/websocket/metrics"
	"github.com/smartpet/websocket/utils/configs"
	"github.com/smartpet/websocket/utils/flags"
	log "github.com/smartpet/websocket/utils/logger"
)

func setupRoutes() {
	http.HandleFunc("/ws", business.WsEndpoint)
}
func main() {

	Initialization()

	setupRoutes()
	http.ListenAndServe(":8001", nil)
}
func initAWS() {

	if os.Getenv(constant.ModeKey) == "local" {
		return
	}
	log.Info(context.Background()).Msg("AWS config initilization started")
	//log.Fatal(ctx).Err(err).Stack().Msg("error initializaing AWS config")

	ctx := context.Background()
	err := configs.InitAWS(ctx)
	if err != nil {
		log.Fatal(ctx).Err(err).Stack().Msg("error initializaing AWS config")
	}
}

func initMetrics() {
	// Below are value for time interval for time metric. Start corresponds to startin value in seconds,
	// it'll start from 10ms, with increment of 20ms, and go till 300ms
	metrics.Init(metrics.BucketConfig{Start: 0.01, Width: 0.03, Count: 4})
}
func initConfigs() {
	// init configs
	configs.Init(flags.BaseConfigPath())
	var err error
	if flags.ApplicationMode() == constant.TestMode {
		err = configs.InitTestModeConfigs(flags.BaseConfigPath(), constant.DatabaseConfig, constant.LoggerConfig, constant.ApplicationConfig, constant.ExternalConfig)
	} else if flags.ApplicationMode() == constant.ReleaseMode {
		err = configs.InitReleaseModeConfigs(constant.DatabaseConfig, constant.LoggerConfig)
	}
	if err != nil {
		log.ApplicationFatal(context.Background()).Err(err).Msg("error loading configs")
	}
}

func startLogger() {
	// start logger
	loggerConfig, err := configs.Get(constant.LoggerConfig)
	if err != nil {
		log.ApplicationFatal(context.Background()).Err(err).Msg("error getting logger config")
	}
	//log.InitLoggerWithFileEmit(log.Level(loggerConfig.GetString(constant.LogLevelConfigKey)))
	log.InitLogger(log.Level(loggerConfig.GetString(
		constant.LogLevelConfigKey)))

}

func Initialization() {
	initAWS()
	initMetrics()
	initConfigs()
	startLogger()
	log.ApplicationInfo(context.Background()).Int("numCPUs", runtime.NumCPU()).Int("maxProcs", runtime.GOMAXPROCS(0)).Send()

}
