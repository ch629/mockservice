package stub

import (
	"fmt"
	"testing"

	"github.com/ch629/mockservice/pkg/domain"
	"github.com/ch629/mockservice/pkg/stub/field_matching"
	"github.com/ch629/mockservice/pkg/stub/request_matching"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func Test_templateRequest(t *testing.T) {
	staticUUID := uuid.NewString()
	globalFuncs["randomUUID"] = func() string {
		return staticUUID
	}
	stub := Definition{
		Request: request_matching.NewPathMatcher(field_matching.EndsWithMatcher("/abc")),
		Response: Response{
			Body: []byte(`
			{
				"id": "{{ randomUUID }}",
				"foo": "{{ json .Body "$.foo" }}"
			}
			`),
			Status: 200,
		},
		ID:        uuid.New(),
		Templated: true,
	}
	req := domain.Request{
		Path:   "/abc",
		Method: "GET",
		Body: []byte(`
		{
			"foo": "bar"
		}
		`),
	}

	resp, err := templateRequest(stub, req)
	require.NoError(t, err)

	require.JSONEq(t, fmt.Sprintf(`
	{
		"id": "%s",
		"foo": "bar"
	}
		`, staticUUID), string(resp.Body))
	require.Equal(t, 200, resp.Status)
}
