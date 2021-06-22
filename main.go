package main

import (
	"log"
	"net/http"
	"os"

	"github.com/sendinblue/dpe-insights/api"
	"github.com/sendinblue/dpe-insights/core/logger"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

// app Application .
type app struct {
	Router *mux.Router
}

// Initialize the application.
func (a *app) Initialize() {
	setConfiguration()

	// Configure application log file
	logger.Configure()

	a.Router = api.Router()
}

// setConfiguration set env and configuration values.
func setConfiguration() {
	viper.SetDefault("APP_PORT", "8080")

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		// Find and read the config file
		err := viper.ReadInConfig()
		if err != nil {
			log.Fatalf("Error while reading config file %s", err)
		}
	}

	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()
}

// Run the application.
func (a *app) Run(port string) {
	err := http.ListenAndServe(":"+port, a.Router)
	if err != nil {
		log.Fatal()
	}
}

func main() {
	a := app{}
	a.Initialize()

	port := viper.GetString("APP_PORT")
	a.Run(port)
}
