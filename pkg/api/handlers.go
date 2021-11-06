package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ch629/mockservice/pkg/stub"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// GET /definition
func (s *server) listDefinitions() http.HandlerFunc {
	type definitionDto struct {
		ID uuid.UUID `json:"id"`
	}
	type response struct {
		Definitions []definitionDto `json:"definitions"`
	}

	return func(rw http.ResponseWriter, _ *http.Request) {
		definitions := s.stubService.Definitions()

		resp := response{
			Definitions: make([]definitionDto, len(definitions)),
		}
		for idx, def := range definitions {
			resp.Definitions[idx] = definitionDto{
				ID: def.ID,
			}
		}

		rw.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(rw).Encode(resp); err != nil {
			s.logger.Error("failed to encode definitions to ResponseWriter", zap.Error(err))
		}
	}
}

// POST /definition
func (s *server) registerDefinition() http.HandlerFunc {
	type request struct {
		Path string `json:"path"`
	}
	type response struct {
		ID uuid.UUID `json:"id"`
	}

	return func(rw http.ResponseWriter, r *http.Request) {
		var req request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		id := s.stubService.AddStub(stub.Definition{
			Request: stub.NewLoggedMatcher(s.logger, stub.NewPathMatcher(req.Path)),
			ID:      uuid.New(),
			Response: stub.Response{
				Headers: map[string]string{},
				Body:    []byte(`{"foo": "bar"}`),
				Status:  http.StatusOK,
			},
		})
		rw.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(rw).Encode(response{
			ID: id,
		}); err != nil {
			s.logger.Error("failed to encode id to ResponseWriter", zap.Error(err))
		}
	}
}

// DELETE /definition/{id}
func (s *server) deleteDefinition() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id, err := uuid.Parse(params["id"])
		if err != nil {
			// TODO : Invalid uuid
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		_, err = s.stubService.RemoveStub(id)
		if err != nil {
			if errors.Is(err, stub.ErrNoDefinition) {
				rw.WriteHeader(http.StatusNotFound)
				return
			}
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		rw.WriteHeader(http.StatusOK)
	}
}
