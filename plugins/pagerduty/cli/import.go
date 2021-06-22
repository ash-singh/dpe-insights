package pagerduty

import (
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(importCmd)
}

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import data form PagerDuty",
	Long:  "Import data form PagerDuty",
}
