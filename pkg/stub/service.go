package stub

import (
	"errors"
	"net/http"
	"sync"

	"github.com/ch629/mockservice/pkg/recorder"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

var ErrNoDefinition = errors.New("no definition")

//go:generate mockgen -destination=stub_mocks/mock_stub.go -package=stub_mocks . Service
// TODO: Separate stubs from HTTP logic?
type Service interface {
	http.Handler

	AddStub(def Definition) uuid.UUID
	RemoveStub(id uuid.UUID) (*Definition, error)

	Definitions() []Definition
}

// TODO: Pull this into a domain pkg
type Definition struct {
	// TODO: Multiple request matchers so we can say which was closest?
	Request  RequestMatcher
	Response Response
	ID       uuid.UUID
}

func NewService(log *zap.Logger, recorder recorder.Service) Service {
	return &service{
		definitions: make(map[uuid.UUID]Definition),
		log:         log,
		recorder:    recorder,
	}
}

type service struct {
	mux sync.RWMutex
	log *zap.Logger

	// TODO: This has to be ordered or pulled out as a priority list
	definitions map[uuid.UUID]Definition
	recorder    recorder.Service
}

func (s *service) AddStub(def Definition) uuid.UUID {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.log.Info("adding stub", zap.Stringer("id", def.ID))
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
	s.recorder.Record(req)
	s.mux.RLock()
	defer s.mux.RUnlock()
	for id, def := range s.definitions {
		if def.Request.Matches(req) {
			s.recorder.RecordStub(id)
			if err := def.Response.WriteTo(rw); err != nil {
				s.log.Error("failed to write response back", zap.Error(err))
			}
			s.log.Debug("request matched", zap.Stringer("stub", id))
			return
		}
	}
	// TODO: No stub found response
	rw.WriteHeader(http.StatusNoContent)
	s.log.Info("didn't find any stub for", zap.String("path", req.Path))
}
