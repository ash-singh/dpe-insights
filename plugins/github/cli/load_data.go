package github

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(loadDataCmd)
	loadDataCmd.Flags().IntP("days", "d", 7, "Look back days")
	loadDataCmd.Flags().StringP("startDate", "", "", "Start date (yyyy-mm-dd)")
	loadDataCmd.Flags().StringP("endDate", "", "", "End date (yyyy-mm-dd)")
}

var loadDataCmd = &cobra.Command{
	Use:   "load-data",
	Short: "Load all data from GitHub",
	Long:  `Load all data (pull request, releases, teams etc) from GitHub`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Load data from GitHub ...")

		days, _ := cmd.Flags().GetInt("days")
		startDate, _ := cmd.Flags().GetString("startDate")
		endDate, _ := cmd.Flags().GetString("endDate")

		loadData(days, startDate, endDate)
	},
}

func loadData(days int, startDate, endDate string) {
	importTeam()

	fmt.Println("Importing pull request data.")
	importPullRequestData(days, startDate, endDate)
	importPullRequestCountData(days)

	fmt.Println("Importing releases data.")
	importReleases(false)

	fmt.Println("Transforming data.")
	transformPullRequestData()
	transformPullRequestCountData(days)
	transformTeamPullRequestCountData(days)

	fmt.Println("Data load complete.")
}
