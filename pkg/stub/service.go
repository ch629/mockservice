package stub

import (
	"errors"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

var ErrNoDefinition = errors.New("no definition")

type Service interface {
	http.Handler

	AddStub(def Definition) uuid.UUID
	RemoveStub(id uuid.UUID) (*Definition, error)

	Definitions() []Definition
}

type Definition struct {
	// TODO: Multiple request matchers so we can say which was closest?
	Request  RequestMatcher
	Response Response
	ID       uuid.UUID
}

func NewService(logger *zap.Logger) Service {
	return &service{
		definitions: make(map[uuid.UUID]Definition),
		logger:      logger,
	}
}

type service struct {
	mux    sync.RWMutex
	logger *zap.Logger

	definitions map[uuid.UUID]Definition
}

func (s *service) AddStub(def Definition) uuid.UUID {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.logger.Info("adding stub", zap.Stringer("id", def.ID))
	s.definitions[def.ID] = def
	return def.ID
}

func (s *service) RemoveStub(id uuid.UUID) (*Definition, error) {
	s.mux.Lock()
	defer s.mux.Unlock()
	def, ok := s.definitions[id]
	if !ok {
		return nil, ErrNoDefinition
	}
	delete(s.definitions, id)
	return &def, nil
}

func (s *service) Definitions() []Definition {
	s.mux.RLock()
	defer s.mux.RUnlock()
	defs := make([]Definition, 0, len(s.definitions))
	for _, d := range s.definitions {
		defs = append(defs, d)
	}
	return defs
}

// TODO: Closest stub
func (s *service) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r == nil {
		return
	}
	req := RequestFromHTTP(*r)
	s.mux.RLock()
	defer s.mux.RUnlock()
	for _, def := range s.definitions {
		if def.Request.Matches(req) {
			if err := def.Response.WriteTo(rw); err != nil {
				s.logger.Error("failed to write response back", zap.Error(err))
			}
			return
		}
	}
	s.logger.Info("didn't find any stub for", zap.String("path", req.Path))
}
