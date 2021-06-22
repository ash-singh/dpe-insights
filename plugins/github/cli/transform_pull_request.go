package github

import (
	"log"

	"github.com/spf13/cobra"
)

func init() {
	transformCmd.AddCommand(transformPullRequestCmd)
}

var transformPullRequestCmd = &cobra.Command{
	Use:   "pull-requests",
	Short: "Transform pull request data",
	Long:  `Transform pull request data`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Transform pull request data ...")

		transformPullRequestData()
	},
}

func transformPullRequestData() {
	dic := getDIContainer()

	prTransformer := dic.GetPullRequestTransformer()

	count, err := prTransformer.TransformRecords()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Process completed for total %d records \n", count)
}
