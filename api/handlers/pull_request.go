package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/sendinblue/dpe-insights/core/constants"
	"github.com/sendinblue/dpe-insights/core/helpers"
	"github.com/gorilla/mux"
)

var di *diContainer

func init() {
	di = getDIContainer()
}

// ImportDailyPullRequestCountData Import daily pull request count data.
func ImportDailyPullRequestCountData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	days, _ := strconv.Atoi(mux.Vars(r)["days"])

	now := time.Now()

	prImporter := di.GetPullRequestImporter()
	prcRepository := di.GetPullRequestCountRepository()

	for days > 0 {
		importDate := now.AddDate(0, 0, -days)
		pullRequestCount, err := prImporter.GetPullRequestStatsWithFilter(importDate)
		if err != nil {
			log.Println(err)
			continue
		}

		err = prcRepository.Save(pullRequestCount)

		if err != nil {
			log.Println(err)
			continue
		}

		time.Sleep(3 * time.Second)

		days--
	}

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

// ImportPullRequestData Import pull request data.
func ImportPullRequestData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	days, _ := strconv.Atoi(mux.Vars(r)["days"])

	now := time.Now()
	startDate := now.AddDate(0, 0, -days)
	endDate := now

	prImporter := di.GetPullRequestImporter()

	err := prImporter.ImportPullRequestData(startDate, endDate)
	if err != nil {
		log.Println(err)
		helpers.GetError(err, w)
		return
	}

	data, err := json.Marshal(helpers.SuccessResponse{
		Message: "success",
	})
	if err != nil {
		log.Println(err)
		helpers.GetError(err, w)
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

// SyncExtractedPullRequestData Sync pull request data.
func SyncExtractedPullRequestData(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	prSyncData := di.GetSyncDataService()
	prSyncData.Sync()

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

// TransformExtractedPullRequestData Transform extracted pull request data.
func TransformExtractedPullRequestData(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	prTransformer := di.GetPullRequestTransformer()

	_, _ = prTransformer.TransformRecords()

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

// TransformExtractedPullRequestCountData Transform pull request count data.
func TransformExtractedPullRequestCountData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	days, _ := strconv.Atoi(mux.Vars(r)["days"])

	now := time.Now()

	tprRepository := di.GetTransformedPullRequestRepository()
	tprCountRepository := di.GetTransformedPullRequestCountRepository()

	for days > 0 {
		transformDate := now.AddDate(0, 0, -days)

		result, err := tprRepository.FetchCountStats(transformDate)
		if err != nil {
			log.Println("FetchCountStats", err)
			continue
		}
		err = tprCountRepository.Save(result)

		if err != nil {
			log.Println("TransformedCount Save: ", err)
		}

		days--
	}

	fmt.Println("Transformation Pull request count is complete.")

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

// TransformTeamPullRequestWeeklyCountData Transform team pull request weekly count data.
func TransformTeamPullRequestWeeklyCountData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	days, _ := strconv.Atoi(mux.Vars(r)["days"])

	tprRepository := di.GetTransformedPullRequestRepository()
	teamPrRepository := di.GetTransformedTeamPullRequestCountRepository()

	for days > 0 {
		startTime := time.Now().AddDate(0, 0, -days)
		endTime := startTime.AddDate(0, 0, 1)

		startDate := startTime.Format(constants.DateFormatISO)
		endDate := endTime.Format(constants.DateFormatISO)

		records, err := tprRepository.FetchTeamsPullRequestCount(startDate, endDate)
		if err != nil {
			log.Println("FetchTeamsPullRequestCount", err)
			continue
		}

		for _, r := range *records {
			err = teamPrRepository.Save(&r)
			if err != nil {
				log.Println("transformed team pull request count save error: ", err)
			}
		}

		days--
	}

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
