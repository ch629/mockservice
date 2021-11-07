package domain

// TODO: Cookies
// TODO: Form data
// TODO: Auth
type Request struct {
	QueryParameters map[string][]string `json:"query_parameters"`
	Headers         map[string][]string `json:"headers"`
	Path            string              `json:"path"`
	Method          string              `json:"method"`
	Body            []byte              `json:"body"`
}
