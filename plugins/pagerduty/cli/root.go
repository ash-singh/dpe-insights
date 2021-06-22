package pagerduty

import (
	"github.com/spf13/cobra"
)

// RootCmd PagerDuty main command.
var RootCmd = &cobra.Command{
	Use:   "pagerduty",
	Short: "Import data from PagerDuty",
	Long:  `Import data form PagerDuty`,
}
