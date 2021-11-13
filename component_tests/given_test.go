package main_test

import (
	"github.com/cucumber/godog"
	"go.uber.org/zap"
)

func (s *scenarioContext) aDefinitionIsRegisteredWithPayload(body *godog.DocString) error {
	id, err := s.api.RegisterDefinition(body.Content)
	if err != nil {
		return err
	}
	s.log.Info("received", zap.Stringer("id", id))
	s.stubIDs = append(s.stubIDs, id)
	return nil
}
