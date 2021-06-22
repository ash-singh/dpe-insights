package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/sendinblue/dpe-insights/plugins/pagerduty/models/entities"
)

var errorImporter = errors.New("importerError")

// HTTPClient interface.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client for access to PagerDuty API.
type Client struct {
	accessToken string
	httpClient  HTTPClient
}

// New PagerDuty client.
func New(accessToken string) *Client {
	return &Client{
		accessToken: accessToken,
		httpClient:  &http.Client{},
	}
}

type incidentResult struct {
	Incidents []entities.Incident `json:"incidents"`
	Limit     int                 `json:"limit"`
	Offset    int                 `json:"offset"`
	Total     interface{}         `json:"total"`
	More      bool                `json:"more"`
}

// GetAllIncidents retrieves all incidents for a given time frame.
func (client *Client) GetAllIncidents(since time.Time, until time.Time) ([]entities.Incident, error) {
	offset := 0
	var incidents []entities.Incident
	for {
		res, err := client.getIncidentsBetween(since, until, offset)
		if err != nil {
			return nil, err
		}

		for index := range res.Incidents {
			calculateDynamicFields(&res.Incidents[index])
		}

		incidents = append(incidents, res.Incidents...)
		if !res.More {
			break
		}

		offset += 100
		if offset >= 10000 {
			return nil, err // API will answer err 400 with offset >= 10000
		}
	}

	return incidents, nil
}

func calculateDynamicFields(incident *entities.Incident) {
	incident.ServiceName = incident.Service.Summary
	incident.Duration = int(incident.LastStatusChangedAt.Sub(incident.CreatedAt) / time.Second)
}

func (client *Client) getIncidentsBetween(start time.Time, until time.Time, offset int) (*incidentResult, error) {
	verb := "GET"
	requestStr := buildHTTPRequestURL(start, until, offset)
	req, err := http.NewRequest(verb, requestStr, nil)
	if err != nil {
		return nil, newImporterError(err.Error())
	}

	req.Header.Add("Authorization", fmt.Sprintf("Token token=%v", client.accessToken))
	req.Header.Set("Content-Type", "application/json")

	httpResp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, newImporterError(err.Error())
	}

	bodyBytes, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		log.Fatal(err)
	}

	if httpResp.StatusCode != http.StatusOK {
		return nil, newImporterError(fmt.Sprintf("HTTP request returned a non-success code (%v): %v", httpResp.StatusCode, string(bodyBytes)))
	}

	res, err := decodeHTTPResponseBody(bodyBytes)
	if err != nil {
		log.Fatal(err)
	}

	_ = httpResp.Body.Close()

	return res, nil
}

func newImporterError(s string) error {
	return fmt.Errorf("%w: %s", errorImporter, s)
}

func buildHTTPRequestURL(since time.Time, until time.Time, offset int) string {
	url := fmt.Sprintf("https://api.pagerduty.com/incidents?limit=100&since=%v&until=%v&time_zone=UTC&offset=%v", since.Format(time.RFC3339), until.Format(time.RFC3339), strconv.Itoa(offset))
	return url
}

func decodeHTTPResponseBody(body []byte) (*incidentResult, error) {
	var res *incidentResult
	err := json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
