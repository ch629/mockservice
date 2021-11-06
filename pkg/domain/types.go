package domain

// TODO: Cookies
// TODO: Form data
// TODO: Auth
type Request struct {
	QueryParameters map[string][]string
	Headers         map[string][]string
	Path            string
	Method          string
	Body            []byte
}
