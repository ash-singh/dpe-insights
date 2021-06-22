package main

import (
	github "github.com/sendinblue/dpe-insights/plugins/github/cli"
)

func init() {
	rootCmd.AddCommand(github.RootCmd)
}
