package github

import (
	"log"
	"time"

	"github.com/sendinblue/dpe-insights/core/constants"
	"github.com/cheggaaa/pb/v3"
	"github.com/spf13/cobra"
)

func init() {
	importCmd.AddCommand(pullRequestCmd)
	pullRequestCmd.Flags().IntP("days", "d", 7, "Look back days")
	pullRequestCmd.Flags().StringP("startDate", "", "", "Start date (yyyy-mm-dd)")
	pullRequestCmd.Flags().StringP("endDate", "", "", "End date (yyyy-mm-dd)")
}

var pullRequestCmd = &cobra.Command{
	Use:   "pull-requests",
	Short: "Import pull request data from GitHub",
	Long:  `Import pull request data`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Importing pull request data ...")

		days, _ := cmd.Flags().GetInt("days")
		startDate, _ := cmd.Flags().GetString("startDate")
		endDate, _ := cmd.Flags().GetString("endDate")

		importPullRequestData(days, startDate, endDate)
	},
}

func importPullRequestData(days int, startDateInput, endDateInput string) {
	dic := getDIContainer()

	im := dic.GetPullRequestImporter()

	if startDateInput != "" && endDateInput != "" {

		startDate, _ := time.Parse(constants.DateFormatISO, startDateInput)
		endDate, _ := time.Parse(constants.DateFormatISO, endDateInput)

		log.Printf("Fetching data between %v and %v", startDate, endDate)

		err := im.ImportPullRequestData(startDate, endDate)
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

			err := im.ImportPullRequestData(startDate, endDate)
			if err != nil {
				log.Println(err)
			}

			// date is inclusive so next end date is -1 day
			endDate = startDate.AddDate(0, 0, -1)

			bar.Add(dateRange)

			days -= 7
		}

		bar.Finish()
	}

	log.Println("Import Complete.")
}
