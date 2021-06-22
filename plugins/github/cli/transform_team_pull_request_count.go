package github

import (
	"log"
	"time"

	"github.com/sendinblue/dpe-insights/core/constants"
	"github.com/cheggaaa/pb/v3"
	"github.com/spf13/cobra"
)

func init() {
	transformCmd.AddCommand(transformTeamPullRequestCountCmd)
	transformTeamPullRequestCountCmd.Flags().IntP("days", "d", 7, "Look back days")
}

var transformTeamPullRequestCountCmd = &cobra.Command{
	Use:   "team-pull-request-counts",
	Short: "Transform team pull request count data",
	Long:  `Transform team pull request count data`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Transform team pull request count data ...")

		days, _ := cmd.Flags().GetInt("days")

		transformTeamPullRequestCountData(days)
	},
}

func transformTeamPullRequestCountData(days int) {
	bar := pb.Full.Start(days)

	dic := getDIContainer()
	tprRepository := dic.GetTransformedPullRequestRepository()
	teamPrRepository := dic.GetTransformedTeamPullRequestCountRepository()

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
				log.Println("transformed Team Pull request count save error: ", err)
				continue
			}
		}

		bar.Increment()
		days--
	}

	bar.Finish()

	log.Println("Process completed.")
}
