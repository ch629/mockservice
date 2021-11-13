package main_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var opts = godog.Options{
	Output: colors.Colored(os.Stdout),
	Format: "pretty",
}

// TODO: Pass in a flag to point to an already-running app
func init() {
	godog.BindCommandLineFlags("godog.", &opts)
}

func TestMain(m *testing.M) {
	pflag.Parse()
	opts.Paths = pflag.Args()

	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	log, err := cfg.Build()
	if err != nil {
		panic(fmt.Errorf("failed to create logger: %w", err))
	}
	suiteCtx := &suiteContext{
		log: log,
	}
	status := godog.TestSuite{
		Name:                 "mockservice",
		TestSuiteInitializer: suiteCtx.InitializeTestSuite,
		ScenarioInitializer:  suiteCtx.InitializeScenario,
		Options:              &opts,
	}.Run()

	os.Exit(status)
}
