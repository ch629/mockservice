package main_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/spf13/pflag"
	testcontainers "github.com/testcontainers/testcontainers-go"
)

var opts = godog.Options{
	Output: colors.Colored(os.Stdout),
}

func init() {
	godog.BindCommandLineFlags("godog.", &opts)
}

func TestMain(m *testing.M) {
	pflag.Parse()
	opts.Paths = pflag.Args()

	suiteCtx := &suiteContext{}
	status := godog.TestSuite{
		Name:                 "mockservice",
		TestSuiteInitializer: suiteCtx.InitializeTestSuite,
		ScenarioInitializer:  suiteCtx.InitializeScenario,
		Options:              &opts,
	}.Run()

	os.Exit(status)
}

type suiteContext struct {
	mockServiceUri string
}

func (c *suiteContext) InitializeTestSuite(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(func() {
		// TODO: use compose instead? - https://golang.testcontainers.org/features/docker_compose/
		containerRequest := testcontainers.ContainerRequest{
			FromDockerfile: testcontainers.FromDockerfile{
				Context: "../",
			},
			ExposedPorts: []string{"8080:8080"},
			// WaitingFor:   wait.ForHTTP("/admin/definition"),
		}

		container, err := testcontainers.GenericContainer(context.TODO(), testcontainers.GenericContainerRequest{
			ContainerRequest: containerRequest,
			Started:          true,
		})
		if err != nil {
			panic(err)
		}

		ip, err := container.Host(context.TODO())
		if err != nil {
			panic(err)
		}
		mappedPort, err := container.MappedPort(context.TODO(), "8080")
		if err != nil {
			panic(err)
		}
		uri := fmt.Sprintf("http://%s:%s", ip, mappedPort.Port())
		c.mockServiceUri = uri
		fmt.Println("Started on uri:", uri)
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
	ctx.Step("^something happens$", func() error {
		fmt.Println("fail")
		resp, err := http.DefaultClient.Get(s.mockServiceUri + "/admin/definition")
		if err != nil {
			return fmt.Errorf("get request to definitions: %w", err)
		}

		defer resp.Body.Close()
		bs, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("reading body: %w", err)
		}
		fmt.Println("Got body", string(bs))
		return nil
	})
}
