package main_test

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
)

type api struct {
	client *resty.Client
	apiURI string
}

func (a *api) RegisterDefinition(payload string) (uuid.UUID, error) {
	type response struct {
		ID uuid.UUID `json:"id"`
	}
	resp, err := a.client.
		R().
		SetResult(response{}).
		SetBody(payload).
		Post("/admin/definition")
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to register definition: %w", err)
	}

	if resp.Error() != nil {
		return uuid.Nil, fmt.Errorf("error from response: %w", resp.Error())
	}

	if resp.IsError() {
		return uuid.Nil, fmt.Errorf("error status code from response, recieved: %d", resp.StatusCode())
	}

	return resp.Result().(*response).ID, nil
}
