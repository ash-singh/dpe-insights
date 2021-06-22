package github

import (
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(transformCmd)
}

var transformCmd = &cobra.Command{
	Use:   "transform",
	Short: "Transform imported data",
	Long:  `Transform imported data`,
}
