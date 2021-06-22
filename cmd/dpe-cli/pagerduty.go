package main

import (
	pagerduty "github.com/sendinblue/dpe-insights/plugins/pagerduty/cli"
)

func init() {
	rootCmd.AddCommand(pagerduty.RootCmd)
}
