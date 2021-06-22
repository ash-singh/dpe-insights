package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	appTest "github.com/sendinblue/dpe-insights/testing"
)

// execute Test http request.
func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	Router().ServeHTTP(rr, req)
	return rr
}

func TestPing(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/v1/ping", nil)

	response := executeRequest(req)

	appTest.CheckResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body == "[]" {
		t.Errorf("Expected response. Got %s", body)
	}
}
