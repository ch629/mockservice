package main_test

import (
	"github.com/cucumber/godog"
)

func (s *scenarioContext) createSteps(ctx *godog.ScenarioContext) {
	// Given
	ctx.Step("^a definition is registered with payload$", s.aDefinitionIsRegisteredWithPayload)

	// When
	ctx.Step(`^a request is sent with method "([^"]*)", path "([^"]*)" and payload$`, s.aRequestIsSentWithMethodPathAndPayload)

	// Then
	ctx.Step("^the response body should match$", s.theResponseBodyShouldMatch)
	ctx.Step(`^the response should have status code (\d+)$`, s.theResponseShouldHaveStatusCode)
}
