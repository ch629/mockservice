package main_test

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
)

type api struct {
	client *resty.Client
}

type apiResponse struct {
	Body       string
	StatusCode int
	Headers    map[string][]string
}

func (a *api) RegisterDefinition(payload string) (id uuid.UUID, err error) {
	type response struct {
		ID uuid.UUID `json:"id"`
	}
	resp, err := a.client.
		R().
		SetResult(response{}).
		SetBody(payload).
		Post("/admin/definition")
	if err != nil {
		return id, fmt.Errorf("failed to register definition: %w", err)
	}

	if resp.Error() != nil {
		return id, fmt.Errorf("error from response: %w", resp.Error())
	}

	if resp.IsError() {
		return id, fmt.Errorf("error status code from response, received: %d", resp.StatusCode())
	}

	return resp.Result().(*response).ID, nil
}

func (a *api) SendRequest(method, path, payload string) (*apiResponse, error) {
	resp, err := a.client.
		R().
		SetBody(payload).
		SetResult(json.RawMessage{}).
		Execute(method, path)
	if err != nil {
		return nil, fmt.Errorf("failed to register definition: %w", err)
	}

	if resp.Error() != nil {
		return nil, fmt.Errorf("error from response: %w", resp.Error())
	}

	return &apiResponse{
		Body:       string(resp.Body()),
		StatusCode: resp.StatusCode(),
		Headers:    resp.Header(),
	}, nil
}

func (a *api) DeleteStub(stubID uuid.UUID) error {
	_, err := a.client.R().SetPathParam("id", stubID.String()).Delete("/admin/definition/{id}")
	if err != nil {
		return fmt.Errorf("error from response: %w", err)
	}
	return nil
}
