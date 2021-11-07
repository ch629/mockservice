package api

import "net/http"

func (s *server) registerRoutes() {
	adminRouter := s.router.PathPrefix("/admin").Subrouter()

	// Admin routes
	for _, route := range []struct {
		Path    string
		Handler func() http.HandlerFunc
		Method  string
	}{
		{
			Path:    "/definition",
			Handler: s.listDefinitions,
			Method:  http.MethodGet,
		},
		{
			Path:    "/definition",
			Handler: s.registerDefinition,
			Method:  http.MethodPost,
		},
		{
			Path:    "/definition/{id}",
			Handler: s.deleteDefinition,
			Method:  http.MethodDelete,
		},
		{
			Path:    "/request",
			Handler: s.requests,
			Method:  http.MethodGet,
		},
		{
			Path:    "/stubs",
			Handler: s.stubCounts,
			Method:  http.MethodGet,
		},
	} {
		adminRouter.
			HandleFunc(route.Path, route.Handler()).
			Methods(route.Method)
	}

	s.router.PathPrefix("/").Handler(s.stubService)
}
