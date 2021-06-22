package api

import (
	"net/http"

	"github.com/sendinblue/dpe-insights/api/handlers"
	"github.com/gorilla/mux"
)

// Router Configured Router Object.
func Router() *mux.Router {
	r := mux.NewRouter()

	// Endpoint is required by GKE health check
	r.HandleFunc("/", handlers.Ping).Methods(http.MethodGet)

	api := r.PathPrefix("/api/v1").Subrouter()

	// Ping test
	api.HandleFunc("/ping", handlers.Ping).Methods(http.MethodGet)

	api.HandleFunc("/org/", handlers.GetOrganizations).Methods(http.MethodGet)
	api.HandleFunc("/org/{id:[0-9]+}", handlers.GetOrganization).Methods(http.MethodGet)

	// import team and related data
	api.HandleFunc("/org/{id:[0-9]+}/team/import", handlers.ImportTeamData).Methods(http.MethodGet)

	// pull request count data import
	api.HandleFunc("/org/{id:[0-9]+}/pull-request/count/{days:[0-9]+}", handlers.ImportDailyPullRequestCountData).Methods(http.MethodGet)

	// pull request data import
	api.HandleFunc("/org/{id:[0-9]+}/pull-request/{days:[0-9]+}", handlers.ImportPullRequestData).Methods(http.MethodGet)

	// sync extracted pull request data
	api.HandleFunc("/org/{id:[0-9]+}/pull-request/sync", handlers.SyncExtractedPullRequestData).Methods(http.MethodGet)

	// transform extracted pull request data
	api.HandleFunc("/org/{id:[0-9]+}/pull-request/transform", handlers.TransformExtractedPullRequestData).Methods(http.MethodGet)

	// transform extracted pull request count data
	api.HandleFunc("/org/{id:[0-9]+}/pull-request/count/transform/{days:[0-9]+}", handlers.TransformExtractedPullRequestCountData).Methods(http.MethodGet)

	// transform team pull request weekly count data
	api.HandleFunc("/org/{id:[0-9]+}/pull-request/team/count/transform/{weeks:[0-9]+}", handlers.TransformTeamPullRequestWeeklyCountData).Methods(http.MethodGet)

	return r
}
