package github

import (
	"log"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/spf13/cobra"
)

func init() {
	importCmd.AddCommand(pullRequestCountCmd)
	pullRequestCountCmd.Flags().IntP("days", "d", 7, "Look back days")
}

var pullRequestCountCmd = &cobra.Command{
	Use:   "pull-request-counts",
	Short: "Import pull request count data from GitHub",
	Long:  `Import daily pull request count data i.e number of open, total and closed PRs`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Importing pull count request ...")
		days, _ := cmd.Flags().GetInt("days")

		importPullRequestCountData(days)
	},
}

func importPullRequestCountData(days int) {
	dic := getDIContainer()
	bar := pb.Full.Start(days)

	importer := dic.GetPullRequestImporter()
	repository := dic.GetPullRequestCountRepository()

	for days > 0 {
		bar.Increment()

		importDate := time.Now().AddDate(0, 0, -days)

		pullRequestCount, err := importer.GetPullRequestStatsWithFilter(importDate)
		if err != nil {
			log.Println(err)
			continue
		}

		err = repository.Save(pullRequestCount)

		if err != nil {
			log.Println(err)
			continue
		}

		days--
	}

	bar.Finish()
}
