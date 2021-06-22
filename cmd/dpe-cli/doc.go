package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func init() {
	rootCmd.AddCommand(cronCmd)
	cronCmd.Flags().StringP("outputDir", "o", "./docs", "Output directory")
}

var cronCmd = &cobra.Command{
	Use:   "doc",
	Short: "Generate documentation",
	Long:  `Generate documentation`,
	Run: func(cmd *cobra.Command, args []string) {
		outputDir, _ := cmd.Flags().GetString("outputDir")

		fmt.Println("Generating documents in ", outputDir)

		cronTask(outputDir)
	},
}

func cronTask(outputDir string) {
	err := doc.GenMarkdownTree(rootCmd, outputDir)
	if err != nil {
		log.Fatal(err)
	}
}
