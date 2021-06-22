package github

import (
	"log"

	"github.com/spf13/cobra"
)

func init() {
	importCmd.AddCommand(teamCmd)
}

var teamCmd = &cobra.Command{
	Use:   "teams",
	Short: "Import teams data from GitHub",
	Long:  `Import teams data from GitHub`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Importing teams data ...")
		importTeam()
	},
}

func importTeam() {
	dic := getDIContainer()

	uImporter := dic.GetUserImporter()
	err := uImporter.ImportUsers()
	if err != nil {
		log.Fatal("Users import failed!", err)
	}
	log.Println("Users import is complete.")

	rImporter := dic.GetRepositoryImporter()
	err = rImporter.ImportRepositories()

	if err != nil {
		log.Fatal("Repositories import failed!", err)
	}
	log.Println("Repositories import is complete.")

	tImporter := dic.GetTeamImporter()
	err = tImporter.ImportTeams()

	if err != nil {
		log.Fatal("Teams import failed!", err)
	}
	log.Println("Teams data import is complete.")
}
