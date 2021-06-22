package helpers

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"
)

type mockResponseWriter struct {
	mock.Mock
}

func (m *mockResponseWriter) Header() http.Header {
	args := m.Called()
	return args.Get(0).(http.Header)
}

func (m *mockResponseWriter) Write(bytes []byte) (int, error) {
	args := m.Called(bytes)
	return args.Get(0).(int), args.Error(1)
}

func (m *mockResponseWriter) WriteHeader(statusCode int) {
	m.Called(statusCode)
}

type testError struct{}

func (receiver *testError) Error() string {
	return "api error"
}

func TestGetError(t *testing.T) {
	err := new(testError)
	mWriter := new(mockResponseWriter)

	header := http.Header{}

	response := ErrorResponse{
		StatusCode:   http.StatusInternalServerError,
		ErrorMessage: err.Error(),
	}

	message, _ := json.Marshal(response)

	mWriter.
		On("Header").Return(header)

	mWriter.
		On("WriteHeader", http.StatusInternalServerError).Return()

	mWriter.
		On("Write", message).Return(8, nil)

	GetError(err, mWriter)

	mWriter.AssertExpectations(t)
}
