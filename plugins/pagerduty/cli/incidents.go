package pagerduty

import (
	"log"
	"time"

	"github.com/sendinblue/dpe-insights/core/constants"
	"github.com/sendinblue/dpe-insights/plugins/pagerduty/models/repositories"
	"github.com/sendinblue/dpe-insights/plugins/pagerduty/services"
	"github.com/cheggaaa/pb/v3"
	"github.com/spf13/cobra"
)

func init() {
	importCmd.AddCommand(incidentsCmd)
	incidentsCmd.Flags().IntP("days", "d", 7, "Look back days")
	incidentsCmd.Flags().StringP("startDate", "", "", "Start date (yyyy-mm-dd)")
	incidentsCmd.Flags().StringP("endDate", "", "", "End date (yyyy-mm-dd)")
}

var incidentsCmd = &cobra.Command{
	Use:   "incidents",
	Short: "Import incident data",
	Long:  "Import incident data",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Importing incident data ...")
		days, _ := cmd.Flags().GetInt("days")
		startDate, _ := cmd.Flags().GetString("startDate")
		endDate, _ := cmd.Flags().GetString("endDate")
		importIncidents(days, startDate, endDate)
	},
}

func importIncidents(days int, startDateInput string, endDateInput string) {
	dic := getDIContainer()

	if startDateInput != "" && endDateInput != "" {

		startDate, _ := time.Parse(constants.DateFormatISO, startDateInput)
		endDate, _ := time.Parse(constants.DateFormatISO, endDateInput)

		err := importData(dic.GetPagerDutyClient(), dic.GetPagerDutyIncidentRepository(), startDate, endDate)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		bar := pb.Full.Start(days)

		endDate := time.Now()

		dateRange := 7

		if days < 7 {
			dateRange = days
		}

		for days > 0 {
			startDate := endDate.AddDate(0, 0, -dateRange)

			err := importData(dic.GetPagerDutyClient(), dic.GetPagerDutyIncidentRepository(), startDate, endDate)
			if err != nil {
				log.Fatal(err)
			}

			// date is inclusive so next end date is -1 day
			endDate = startDate.AddDate(0, 0, -1)

			bar.Add(dateRange)

			days -= 7
		}

		bar.Finish()
	}

	log.Println("Import complete.")
}

func importData(client *services.Client, repository *repositories.IncidentRepository, startDate time.Time, endDate time.Time) error {
	log.Printf("Fetching data between %v and %v", startDate.Format(time.RFC3339), endDate.Format(time.RFC3339))
	incidents, err := client.GetAllIncidents(startDate, endDate)
	if err != nil {
		return err
	}

	log.Printf("Found %v incidents, inserting", len(incidents))

	for _, incident := range incidents {
		err := repository.Insert(incident)
		if err != nil {
			return err
		}
	}

	return nil
}
