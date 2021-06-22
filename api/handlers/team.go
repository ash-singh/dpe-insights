package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/sendinblue/dpe-insights/core/helpers"
)

func init() {
	di = getDIContainer()
}

// ImportTeamData Handler for importing team data from github.
func ImportTeamData(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	uImporter := di.GetUserImporter()
	rImporter := di.GetRepositoryImporter()
	tImporter := di.GetTeamImporter()

	log.Println("Importing Github users.")
	_ = uImporter.ImportUsers()
	log.Println("User import complete.")

	log.Println("Importing Github repositories.")
	_ = rImporter.ImportRepositories()
	log.Println("Repositories import complete.")

	log.Println("Importing Github teams.")
	_ = tImporter.ImportTeams()
	log.Println("Team import complete.")

	data, err := json.Marshal(helpers.SuccessResponse{
		Message: "success",
	})
	if err != nil {
		log.Println(err)
		helpers.GetError(err, w)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}
