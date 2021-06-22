package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// Config configuration for app.
type Config struct {
	Environment  string
	MysqlDSN     string
	DatabaseName string

	PluginGithubOauth2Token      string
	PluginGithubOrganizationName string
	PluginGithubOwnerName        string

	PluginPagerDutyAccessToken string
}

func getConfig(key string, defaultValue string) string {
	if value, ok := viper.Get(key).(string); ok {
		return value
	}

	return defaultValue
}

func loadEnvVariables() {
	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		// Find and read the config file
		err := viper.ReadInConfig()
		if err != nil {
			fmt.Printf("Error while reading config file %s", err)
		}
	}

	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()
}

// NewConfig New Config Object.
func NewConfig() *Config {
	loadEnvVariables()

	return &Config{
		Environment:  getConfig("APP_ENV", "dev"),
		MysqlDSN:     getConfig("MYSQL_DSN", "root:root@(mysql:3306)/dpe_insights?parseTime=true"),
		DatabaseName: getConfig("DB_NAME", "dpe_insights"),

		PluginGithubOauth2Token:      getConfig("PLUGIN_GITHUB_OAUTH2_TOKEN", ""),
		PluginGithubOrganizationName: getConfig("PLUGIN_GITHUB_ORGANIZATION_NAME", ""),
		PluginGithubOwnerName:        getConfig("PLUGIN_GITHUB_OWNER_NAME", ""),

		PluginPagerDutyAccessToken: getConfig("PLUGIN_PAGER_DUTY_ACCESS_TOKEN", ""),
	}
}
