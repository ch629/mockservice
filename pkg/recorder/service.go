package recorder

import (
	"sync"

	"github.com/ch629/mockservice/pkg/domain"
)

const initialRequestCapacity = 100

type Service interface {
	Record(req domain.Request)
	Requests() []domain.Request

	Clear()
}

func New() Service {
	return &service{
		requests: make([]domain.Request, 0, initialRequestCapacity),
	}
}

type service struct {
	mux      sync.Mutex
	requests []domain.Request
}

func (s *service) Record(req domain.Request) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.requests = append(s.requests, req)
}

func (s *service) Requests() []domain.Request {
	return s.requests
}

func (s *service) Clear() {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.requests = make([]domain.Request, 0, initialRequestCapacity)
}
