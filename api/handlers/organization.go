package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/sendinblue/dpe-insights/core/helpers"
	"github.com/sendinblue/dpe-insights/plugins/github/services"
)

// GetOrganizations Handler for fetching organizations data from github.
func GetOrganizations(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	organizations, err := services.GetOrganizations()
	if err != nil {
		helpers.GetError(err, w)
		return
	}

	data, err := json.Marshal(helpers.SuccessResponse{
		Message: "success",
		Data:    organizations,
	})
	if err != nil {
		log.Fatal(err)
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

// GetOrganization Handler for fetching organization data.
func GetOrganization(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	organizations, err := services.GetOrganization(1535969)
	if err != nil {
		helpers.GetError(err, w)
		return
	}

	data, err := json.Marshal(helpers.SuccessResponse{
		Message: "success",
		Data:    organizations,
	})
	if err != nil {
		log.Println(err)
		helpers.GetError(err, w)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}
