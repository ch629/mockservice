package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/ch629/mockservice/pkg/api"
	"github.com/ch629/mockservice/pkg/config"
	"github.com/ch629/mockservice/pkg/recorder"
	"github.com/ch629/mockservice/pkg/stub"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Setup Logger
	var log *zap.Logger
	{
		var err error
		devCfg := zap.NewDevelopmentConfig()
		devCfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		devCfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		log, err = devCfg.Build()
		if err != nil {
			panic(fmt.Errorf("failed to create logger: %w", err))
		}
	}

	// Dependencies
	recorderService := recorder.New()
	stubService := stub.NewService(log, recorderService)

	// HTTP
	api := api.New(log, stubService, recorderService)
	log.Info("starting HTTP server")
	api.Start(ctx, config.API{
		Port: 8080,
	})
}
