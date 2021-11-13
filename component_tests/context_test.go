package main_test

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/cucumber/godog"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	tc "github.com/testcontainers/testcontainers-go"
)

type (
	suiteContext struct {
		container tc.Container
		log       *zap.Logger
		api       api
	}

	scenarioContext struct {
		*suiteContext

		stubIDs []uuid.UUID

		responseBody   string
		responseStatus int
	}
)

func (c *suiteContext) InitializeTestSuite(ctx *godog.TestSuiteContext) {
	// Start the mockservice in docker
	ctx.BeforeSuite(func() {
		containerRequest := tc.ContainerRequest{
			FromDockerfile: tc.FromDockerfile{
				Context: "../",
			},
			ExposedPorts: []string{"8080:8080"},
		}

		container, err := tc.GenericContainer(context.TODO(), tc.GenericContainerRequest{
			ContainerRequest: containerRequest,
			Started:          true,
		})
		if err != nil {
			c.log.Fatal("Failed to create container", zap.Error(err))
		}

		ip, err := container.Host(context.TODO())
		if err != nil {
			c.log.Fatal("Failed to get container host", zap.Error(err))
		}
		mappedPort, err := container.MappedPort(context.TODO(), "8080")
		if err != nil {
			c.log.Fatal("Failed to get mapped port", zap.Error(err))
		}
		mockServiceURI := fmt.Sprintf("http://%s:%s", ip, mappedPort.Port())

		c.container = container
		c.api = api{
			client: resty.
				New().
				SetBaseURL(mockServiceURI).
				SetHeader("Content-Type", "application/json"),
		}

		c.log.Info("Started mockservice", zap.String("uri", mockServiceURI))
	})

	// Stop the mockservice docker container on finish
	ctx.AfterSuite(func() {
		if c.container != nil {
			c.log.Info("Terminating container", zap.String("container_id", c.container.GetContainerID()))
			if err := c.container.Terminate(context.TODO()); err != nil {
				c.log.Error("failed to terminate container", zap.Error(err))
			}
		}
	})
}

func (c *suiteContext) InitializeScenario(ctx *godog.ScenarioContext) {
	s := &scenarioContext{
		suiteContext: c,
	}

	// Delete all created stubs between scenarios
	ctx.After(func(ctx context.Context, _ *godog.Scenario, _ error) (context.Context, error) {
		s.clearRegisteredStubs()
		return ctx, nil
	})

	s.createSteps(ctx)
}

func (s *scenarioContext) clearRegisteredStubs() {
	for _, stubID := range s.stubIDs {
		if err := s.api.DeleteStub(stubID); err != nil {
			s.log.Warn("failed to delete stub", zap.Stringer("id", stubID), zap.Error(err))
		}
	}
}
