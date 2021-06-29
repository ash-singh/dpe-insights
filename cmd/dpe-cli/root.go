package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dpe-cli.yml)")
}

func initConfig() {
	// Don't forget to read config either from cfgFile or from home directory!
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".dpe-cli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".dpe-cli")
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Using system environment variables.")
	} else {
		fmt.Println("Using system environment variables defined in $HOME/.dpe-cli.yml")
	}

	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()
}

// rootCmd main command.
var rootCmd = &cobra.Command{
	Short: "DPE is a developer productivity engineer cli tool",
	Long: `Import data from multiple sources(Git, Pagerduty, etc)
 apply transformation and save to data store (MySQL).
 Complete documentation is available at https://github.com/sendinblue/dpe-insights#readme`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("dpe-cli v.0.1.0")
	},
}

// execute cli command.
func execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
