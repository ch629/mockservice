package main_test

import (
	"github.com/cucumber/godog"
	"go.uber.org/zap"
)

func (s *scenarioContext) aRequestIsSentWithMethodPathAndPayload(method, path string, body *godog.DocString) error {
	response, err := s.api.SendRequest(method, path, body.Content)
	if err != nil {
		return err
	}
	s.responseBody = response.Body
	s.responseStatus = response.StatusCode
	s.log.Info("received", zap.String("payload", response.Body))
	return nil
}
