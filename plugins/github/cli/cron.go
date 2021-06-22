package github

import (
	"log"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(cronCmd)
	cronCmd.Flags().IntP("days", "d", 1, "Look back days")
}

var cronCmd = &cobra.Command{
	Use:   "cron",
	Short: "Cron task to import/sync GitHub data",
	Long:  `Cron task to import/sync GitHub data`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Starting cron ...")

		days, _ := cmd.Flags().GetInt("days")
		cronTask(days)
	},
}

func cronTask(days int) {
	importTeam()

	log.Println("Importing Pull request data ...")
	importPullRequestData(days, "", "")

	log.Println("Importing Pull count request data ...")
	importPullRequestCountData(days)

	log.Println("Importing releases ...")
	importReleases(true)

	log.Println("Starting data transformation ...")
	transformPullRequestData()
	transformPullRequestCountData(days)
	transformTeamPullRequestCountData(days)
}
