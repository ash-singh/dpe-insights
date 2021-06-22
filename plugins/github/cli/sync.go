package github

import (
	"log"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(syncCmd)
}

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync data from GitHub",
	Long:  `Sync data form GitHub`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Syncing Github Data")
		syncData()
	},
}

func syncData() {
	dic := getDIContainer()

	prSyncService := dic.GetSyncDataService()
	prTransformer := dic.GetPullRequestTransformer()

	total := prSyncService.Sync()

	log.Println("Synced ", total, " records")

	log.Println("Transforming synced records.")
	count, _ := prTransformer.TransformRecords()

	log.Println("Transformed ", count, " records")

	log.Println("Process completed")
}
