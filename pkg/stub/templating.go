package stub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"text/template"

	"github.com/ch629/mockservice/pkg/domain"
	"github.com/google/uuid"
	"github.com/oliveagle/jsonpath"
)

var globalFuncs = template.FuncMap{
	//nolint:gocritic
	"randomUUID": func() string {
		return uuid.NewString()
	},
}

func requestFuncs(request domain.Request) template.FuncMap {
	return template.FuncMap{
		"json": func(jsonBs []byte, path string) string {
			// TODO: validate paths are valid in the stub registration?
			var jsonData interface{}
			if err := json.Unmarshal(jsonBs, &jsonData); err != nil {
				fmt.Println("INVALID JSON INPUT")
				return ""
			}
			res, err := jsonpath.JsonPathLookup(jsonData, path)
			if err != nil {
				fmt.Println("INVALID JSON PATH")
				return ""
			}
			return res.(string)
		},
	}
}

func templateRequest(stub Definition, request domain.Request) (*Response, error) {
	if !stub.Templated {
		return &stub.Response, nil
	}
	reqFuncs := requestFuncs(request)
	// Combine global and request func maps together
	for key, value := range globalFuncs {
		reqFuncs[key] = value
	}

	templ, err := template.New("stubTemplate").Funcs(reqFuncs).Parse(string(stub.Response.Body))
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	var buf bytes.Buffer
	if err := templ.Execute(&buf, request); err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	stub.Response.Body = buf.Bytes()
	return &stub.Response, nil
}
