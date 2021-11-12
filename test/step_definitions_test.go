package main_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/spf13/pflag"
	tc "github.com/testcontainers/testcontainers-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var opts = godog.Options{
	Output: colors.Colored(os.Stdout),
	Format: "pretty",
}

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

type suiteContext struct {
	container      tc.Container
	log            *zap.Logger
	mockServiceURI string
}

func (c *suiteContext) InitializeTestSuite(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(func() {
		containerRequest := tc.ContainerRequest{
			FromDockerfile: tc.FromDockerfile{
				Context: "../",
			},
			ExposedPorts: []string{"8080:8080"},
			// WaitingFor:   wait.ForHTTP("/admin/definition"),
		}

		container, err := tc.GenericContainer(context.TODO(), tc.GenericContainerRequest{
			ContainerRequest: containerRequest,
			Started:          true,
		})
		if err != nil {
			c.log.Fatal("Failed to create container", zap.Error(err))
		}

		c.container = container

		ip, err := container.Host(context.TODO())
		if err != nil {
			c.log.Fatal("Failed to get container host", zap.Error(err))
		}
		mappedPort, err := container.MappedPort(context.TODO(), "8080")
		if err != nil {
			c.log.Fatal("Failed to get mapped port", zap.Error(err))
		}
		c.mockServiceURI = fmt.Sprintf("http://%s:%s", ip, mappedPort.Port())
		c.log.Info("Started mockservice", zap.String("uri", c.mockServiceURI))
	})

	ctx.AfterSuite(func() {
		if c.container != nil {
			c.log.Info("Terminating container", zap.String("container_id", c.container.GetContainerID()))
			if err := c.container.Terminate(context.TODO()); err != nil {
				c.log.Error("failed to terminate container", zap.Error(err))
			}
		}
	})
}

type scenarioContext struct {
	*suiteContext
}

func (c *suiteContext) InitializeScenario(ctx *godog.ScenarioContext) {
	s := &scenarioContext{
		suiteContext: c,
	}
	s.createSteps(ctx)
}

func (s *scenarioContext) createSteps(ctx *godog.ScenarioContext) {
	ctx.Step("^a definition is registered with payload$", s.aDefinitionIsRegisteredWithPayload)
}

func (s *scenarioContext) aDefinitionIsRegisteredWithPayload(body *godog.DocString) error {
	// TODO: Pull some of this logic out to a service
	resp, err := http.DefaultClient.Post(s.mockServiceURI+"/admin/definition", "application/json", strings.NewReader(body.Content))
	if err != nil {
		return fmt.Errorf("POST to /admin/definition: %w", err)
	}
	defer resp.Body.Close()
	respBs, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received a non 200 status code: %d", resp.StatusCode)
	}
	s.log.Info("received", zap.ByteString("body", respBs))
	return nil
}
