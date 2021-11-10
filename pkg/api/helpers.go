package api

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

func (s *server) writeJSON(rw http.ResponseWriter, value interface{}, status int) {
	rw.WriteHeader(status)
	// TODO: Check this is sent
	rw.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(rw).Encode(value); err != nil {
		s.log.Error("failed to encode to ResponseWriter", zap.Error(err))
	}
}

func (s *server) writeError(rw http.ResponseWriter, err error, status int) {
	type response struct {
		Error string `json:"error"`
	}

	s.writeJSON(rw, response{err.Error()}, status)
}
