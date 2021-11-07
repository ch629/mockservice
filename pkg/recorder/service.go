package recorder

import (
	"sync"

	"github.com/ch629/mockservice/pkg/domain"
	"github.com/google/uuid"
)

const initialRequestCapacity = 100

//go:generate mockgen -destination=recorder_mocks/mock_recorder.go -package=recorder_mocks . Service
type Service interface {
	Record(req domain.Request)
	Requests() []domain.Request

	RecordStub(id uuid.UUID)
	Stubs() []StubRecord

	Clear()
}

func New() Service {
	return &service{
		requests: make([]domain.Request, 0, initialRequestCapacity),
		stubs:    make(map[uuid.UUID]int),
	}
}

type service struct {
	requestMux sync.Mutex
	stubsMux   sync.RWMutex
	stubs      map[uuid.UUID]int
	requests   []domain.Request
}

func (s *service) Record(req domain.Request) {
	s.requestMux.Lock()
	defer s.requestMux.Unlock()
	// TODO: Should this be a sized queue instead & then configurable size?
	s.requests = append(s.requests, req)
}

func (s *service) Requests() []domain.Request {
	return s.requests
}

func (s *service) Clear() {
	s.requestMux.Lock()
	defer s.requestMux.Unlock()
	s.requests = make([]domain.Request, 0, initialRequestCapacity)
}

// RecordStub records that a stub was hit & stores the count
func (s *service) RecordStub(id uuid.UUID) {
	s.stubsMux.Lock()
	defer s.stubsMux.Unlock()
	s.stubs[id]++
}

// Stubs returns all of the recorded stubs with their count
func (s *service) Stubs() []StubRecord {
	s.stubsMux.RLock()
	defer s.stubsMux.RUnlock()

	stubs := make([]StubRecord, 0, len(s.stubs))

	for id, count := range s.stubs {
		stubs = append(stubs, StubRecord{
			ID:    id,
			Count: count,
		})
	}

	return stubs
}
