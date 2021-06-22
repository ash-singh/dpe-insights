package github

import (
	"github.com/spf13/cobra"
)

// RootCmd git hub main command.
var RootCmd = &cobra.Command{
	Use:   "github",
	Short: "Import, Transform and Save data from GitHub",
	Long:  `Import data form GitHub`,
}
