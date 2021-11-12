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

func (a *api) SendRequest(method, path, payload string) (string, error) {
	resp, err := a.client.R().SetBody(payload).SetResult(json.RawMessage{}).Execute(method, path)
	if err != nil {
		return "", fmt.Errorf("failed to register definition: %w", err)
	}

	if resp.Error() != nil {
		return "", fmt.Errorf("error from response: %w", resp.Error())
	}

	if resp.IsError() {
		return "", fmt.Errorf("error status code from response, received: %d", resp.StatusCode())
	}

	return string(resp.Body()), nil
}
