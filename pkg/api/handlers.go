package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/ch629/mockservice/pkg/domain"
	"github.com/ch629/mockservice/pkg/recorder"
	"github.com/ch629/mockservice/pkg/stub"
	"github.com/ch629/mockservice/pkg/stub/matching"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// GET /definition
func (s *server) listDefinitions() http.HandlerFunc {
	type definitionDto struct {
		Request stub.RequestMatcher `json:"request"`
		ID      uuid.UUID           `json:"id"`
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
				ID:      def.ID,
				Request: def.Request,
			}
		}

		s.writeJSON(rw, resp, http.StatusOK)
	}
}

// POST /definition
func (s *server) registerDefinition() http.HandlerFunc {
	// TODO: validation
	type request struct {
		Request struct {
			Path   json.RawMessage `json:"path"`
			Method json.RawMessage `json:"method"`
		} `json:"request"`
		Response struct {
			Headers map[string]string `json:"headers"`
			Body    json.RawMessage   `json:"body"`
			Status  int               `json:"status"`
		} `json:"response"`
	}
	type response struct {
		ID uuid.UUID `json:"id"`
	}

	return func(rw http.ResponseWriter, r *http.Request) {
		var req request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.writeError(rw, fmt.Errorf("received invalid payload: %s", err), http.StatusBadRequest)
			return
		}

		fieldMatcher, err := matching.UnmarshalJSONToFieldMatcher(req.Request.Path)
		if err != nil {
			s.writeError(rw, err, http.StatusBadRequest)
			return
		}

		id := s.stubService.AddStub(stub.Definition{
			Request: stub.NewLoggedMatcher(s.log, stub.NewPathMatcher(fieldMatcher)),
			ID:      uuid.New(),
			Response: stub.Response{
				Headers: req.Response.Headers,
				Body:    req.Response.Body,
				Status:  req.Response.Status,
			},
		})
		s.writeJSON(rw, response{id}, http.StatusOK)
	}
}

// DELETE /definition/{id}
func (s *server) deleteDefinition() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id, err := uuid.Parse(params["id"])
		if err != nil {
			s.writeError(rw, fmt.Errorf("received invalid id: %s", err), http.StatusBadRequest)
			return
		}
		_, err = s.stubService.RemoveStub(id)
		if err != nil {
			if errors.Is(err, stub.ErrNoDefinition) {
				s.writeError(rw, fmt.Errorf("no stub found with ID: '%s'", id), http.StatusNotFound)
				return
			}
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		rw.WriteHeader(http.StatusOK)
	}
}

// GET /request
func (s *server) requests() http.HandlerFunc {
	type response struct {
		Requests []domain.Request `json:"requests"`
	}

	return func(rw http.ResponseWriter, _ *http.Request) {
		s.writeJSON(rw, response{s.recorderService.Requests()}, http.StatusOK)
	}
}

// GET /stubs
func (s *server) stubCounts() http.HandlerFunc {
	type response struct {
		Stubs []recorder.StubRecord `json:"stubs"`
	}

	return func(rw http.ResponseWriter, _ *http.Request) {
		s.writeJSON(rw, response{s.recorderService.Stubs()}, http.StatusOK)
	}
}
