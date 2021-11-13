package main_test

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/cucumber/godog"
)

func (s *scenarioContext) theResponseBodyShouldMatch(body *godog.DocString) error {
	var actualJSON, expectedJSON interface{}
	if err := json.Unmarshal([]byte(body.Content), &expectedJSON); err != nil {
		return fmt.Errorf("step doc string was not valid json: %w", err)
	}
	if err := json.Unmarshal([]byte(s.responseBody), &actualJSON); err != nil {
		return fmt.Errorf("response from the mockservice was not valid json: %w", err)
	}

	// TODO: cmp.Diff?
	if !reflect.DeepEqual(actualJSON, expectedJSON) {
		return fmt.Errorf("response did not match")
	}
	return nil
}

// TODO: Headers
func (s *scenarioContext) theResponseShouldHaveStatusCode(statusCode int) error {
	if s.responseStatus != statusCode {
		return fmt.Errorf("response code was %d", s.responseStatus)
	}
	return nil
}
