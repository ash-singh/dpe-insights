package github

import (
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(importCmd)
}

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import data form GitHub",
	Long:  `Import data form GitHub`,
}
