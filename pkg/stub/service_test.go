package stub_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ch629/mockservice/pkg/recorder/recorder_mocks"
	"github.com/ch629/mockservice/pkg/stub"
	"github.com/ch629/mockservice/pkg/stub/matching"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func Test_service_Stubs(t *testing.T) {
	svc := stub.NewService(zap.NewNop(), nil)

	id := uuid.New()
	def := stub.Definition{
		Request: nil,
		Response: stub.Response{
			Headers: map[string]string{},
			Body:    []byte{},
			Status:  200,
		},
		ID: id,
	}

	// Add a stub
	returnedID := svc.AddStub(def)
	require.Equal(t, id, returnedID)
	require.Len(t, svc.Definitions(), 1)

	// Remove the stub
	returnedStub, err := svc.RemoveStub(id)
	require.NoError(t, err)
	require.Equal(t, def, *returnedStub)
	require.Empty(t, svc.Definitions())

	// Try to remove the stub again
	returnedStub, err = svc.RemoveStub(id)
	require.ErrorIs(t, err, stub.ErrNoDefinition)
	require.Nil(t, returnedStub)
}

func Test_service_Handler(t *testing.T) {
	for _, test := range []struct {
		setupMocks       func(recorderMock *recorder_mocks.MockService)
		name             string
		path             string
		method           string
		payload          string
		expectedResponse string
		stubs            []stub.Definition
		expectedStatus   int
	}{
		{
			name:           "No match",
			path:           "/",
			method:         http.MethodGet,
			payload:        ``,
			expectedStatus: http.StatusNoContent,
			setupMocks: func(mock *recorder_mocks.MockService) {
				mock.
					EXPECT().
					Record(gomock.Any()).
					Times(1)
			},
		},
		{
			name:    "Path match",
			path:    "/abc",
			method:  http.MethodGet,
			payload: ``,
			setupMocks: func(mockRecorder *recorder_mocks.MockService) {
				mockRecorder.
					EXPECT().
					Record(gomock.Any()).
					Times(1)
				mockRecorder.
					EXPECT().
					RecordStub(gomock.Any()).
					Times(1)
			},
			stubs: []stub.Definition{
				{
					Request: stub.NewPathMatcher(matching.EqualToMatcher("/abc")),
					Response: stub.Response{
						Headers: map[string]string{},
						Body:    []byte(`{"foo": "bar"}`),
						Status:  http.StatusOK,
					},
				},
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: `{"foo": "bar"}`,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			recorder := recorder_mocks.NewMockService(ctrl)

			if test.setupMocks != nil {
				test.setupMocks(recorder)
			}

			svc := stub.NewService(zap.NewNop(), recorder)

			for _, def := range test.stubs {
				svc.AddStub(def)
			}

			rr := httptest.NewRecorder()
			req, err := http.NewRequest(test.method, test.path, strings.NewReader(test.payload))
			require.NoError(t, err)

			svc.ServeHTTP(rr, req)
			if test.expectedResponse != "" {
				require.Equal(t, test.expectedResponse, rr.Body.String())
			}

			result := rr.Result()
			if result != nil {
				defer result.Body.Close()
				require.Equal(t, test.expectedStatus, result.StatusCode)
			}
		})
	}
}
