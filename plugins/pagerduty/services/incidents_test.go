package services

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/sendinblue/dpe-insights/core/config"
	_ "github.com/sendinblue/dpe-insights/testing"
	"github.com/stretchr/testify/assert"
)

const exampleResponse = `{"incidents":[{"incident_number":12311,"title":"[Warn] [DB] Memory usage on host:example.com","description":"[Warn] [DB] Memory usage on host:example.com","created_at":"2021-02-21T12:11:20Z","status":"resolved","incident_key":null,"service":{"id":"AQJQJ","type":"service_reference","summary":"Low","self":"https://api.pagerduty.com/services/example","html_url":"https://sendinblue.pagerduty.com/service-directory/example"},"assignments":[],"assigned_via":"escalation_policy","last_status_change_at":"2021-02-21T12:50:16Z","first_trigger_log_entry":{"id":"awqdasdas","type":"trigger_log_entry_reference","summary":"Triggered through the API","self":"https://api.pagerduty.com/log_entries/asdasda","html_url":"https://sendinblue.pagerduty.com/incidents/example/log_entries/EXMAPLE"},"alert_counts":{"all":1,"triggered":0,"resolved":1},"is_mergeable":true,"escalation_policy":{"id":"EXAMPLE","type":"escalation_policy_reference","summary":"Devops","self":"https://api.pagerduty.com/escalation_policies/example","html_url":"https://sendinblue.pagerduty.com/escalation_policies/example"},"teams":[],"pending_actions":[],"acknowledgements":[],"basic_alert_grouping":null,"alert_grouping":null,"last_status_change_by":{"id":"example","type":"service_reference","summary":"Low","self":"https://api.pagerduty.com/services/example","html_url":"https://sendinblue.pagerduty.com/service-directory/example"},"priority":null,"resolve_reason":null,"incidents_responders":[],"responder_requests":[],"subscriber_requests":[],"urgency":"low","id":"example","type":"incident","summary":"[#asdad] [Warn] [] Memory usage on host:example.com","self":"https://api.pagerduty.com/incidents/example","html_url":"https://sendinblue.pagerduty.com/incidents/example"}],"limit":100,"offset":0,"total":null,"more":false}`

type mockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func TestClient_GetAllIncidents_Integration(t *testing.T) {
	pdClient := New(config.NewConfig().PluginPagerDutyAccessToken)

	since, _ := time.Parse(time.RFC3339, "2021-02-20T12:11:19Z")
	until, _ := time.Parse(time.RFC3339, "2021-02-20T12:11:21Z")

	_, err := pdClient.GetAllIncidents(since, until)

	assert.NoError(t, err)
}

func TestDecodeHTTPResponseBody(t *testing.T) {
	responseBody, err := decodeHTTPResponseBody([]byte(exampleResponse))

	assert.NoError(t, err)
	assert.Equal(t, 1, len(responseBody.Incidents))
	incident := responseBody.Incidents[0]

	assert.Equal(t, uint32(12311), incident.ID)
	assert.Equal(t, time.Date(2021, 2, 21, 12, 11, 20, 0, time.UTC), incident.CreatedAt)
	assert.Equal(t, time.Date(2021, 2, 21, 12, 50, 16, 0, time.UTC), incident.LastStatusChangedAt)
}

func TestDecodeHTTPResponseBodyFail(t *testing.T) {
	_, err := decodeHTTPResponseBody([]byte("x"))
	assert.NotNil(t, err)
}

func TestBuildHTTPRequestURL(t *testing.T) {
	since := time.Date(2021, 3, 22, 15, 30, 0, 0, time.UTC)
	until := time.Date(2022, 3, 22, 15, 30, 0, 0, time.UTC)

	url := buildHTTPRequestURL(since, until, 30)

	assert.Equal(t, "https://api.pagerduty.com/incidents?limit=100&since=2021-03-22T15:30:00Z&until=2022-03-22T15:30:00Z&time_zone=UTC&offset=30", url)
}

func TestNew(t *testing.T) {
	client := New("x")

	assert.NotNil(t, client.accessToken)
	assert.NotNil(t, client.httpClient)
}

func TestClient_GetAllIncidents(t *testing.T) {
	client := New("x")
	client.httpClient = &mockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			assert.NotNil(t, req)
			assert.Equal(t, "Token token=x", req.Header.Get("Authorization"))
			assert.Equal(t, "application/json", req.Header.Get("Content-Type"))

			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(exampleResponse))),
			}, nil
		},
	}

	incidents, err := client.GetAllIncidents(time.Now(), time.Now())

	assert.NoError(t, err)
	assert.Equal(t, 1, len(incidents))
}
