package github

import (
	"log"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/spf13/cobra"
)

func init() {
	transformCmd.AddCommand(transformPullRequestCountCmd)
	transformPullRequestCountCmd.Flags().IntP("days", "d", 7, "Look back days")
}

var transformPullRequestCountCmd = &cobra.Command{
	Use:   "pull-request-counts",
	Short: "Transform pull request count data",
	Long:  `Transform pull request count data`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Transform pull request count data ...")

		days, _ := cmd.Flags().GetInt("days")

		transformPullRequestCountData(days)
	},
}

func transformPullRequestCountData(days int) {
	bar := pb.Full.Start(days)

	dic := getDIContainer()
	tprCountRepository := dic.GetTransformedPullRequestCountRepository()
	tprRepository := dic.GetTransformedPullRequestRepository()

	for days > 0 {
		transformDate := time.Now().AddDate(0, 0, -days)

		result, err := tprRepository.FetchCountStats(transformDate)
		if err != nil {
			log.Println("FetchCountStats:", err)
			continue
		}
		err = tprCountRepository.Save(result)
		if err != nil {
			log.Println("TransformedCount Insert: ", err)
		}

		bar.Increment()

		days--
	}

	bar.Finish()

	log.Println("Process completed.")
}
