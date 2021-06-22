package github

import (
	"log"

	"github.com/spf13/cobra"
)

func init() {
	importCmd.AddCommand(releaseCmd)
	releaseCmd.Flags().BoolP("is-daily", "", true, "Enter true/false")
}

var releaseCmd = &cobra.Command{
	Use:   "releases",
	Short: "Import releases data from GitHub",
	Long:  `Import releases data from GitHub`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Importing releases data ...")

		isDaily, _ := cmd.Flags().GetBool("is-daily")
		importReleases(isDaily)
	},
}

func importReleases(isDaily bool) {
	dic := getDIContainer()

	im := dic.GetReleaseImporter()

	err := im.ImportReleases(isDaily)
	if err != nil {
		log.Fatal(err)
	}
}
