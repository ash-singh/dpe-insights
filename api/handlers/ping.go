package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/sendinblue/dpe-insights/core/helpers"
)

// Ping Handles the ping request.
func Ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	data, err := json.Marshal(helpers.SuccessResponse{
		Message: "success",
	})
	if err != nil {
		log.Fatal(err)
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}
