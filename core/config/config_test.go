package config

import (
	"testing"

	_ "github.com/sendinblue/dpe-insights/testing"
)

// TestNewConfig test configuration values.
func TestNewConfig(t *testing.T) {
	conf := NewConfig()

	// assert databases name is correct
	if conf.DatabaseName != "dpe_insights" {
		t.Fatalf("DatabaseName %s != dpe_insights", conf.DatabaseName)
	}

	if conf.PluginGithubOauth2Token == "" {
		t.Fatal("PluginGithubOauth2Token is empty")
	}

	if conf.MysqlDSN == "" {
		t.Fatal("MysqlDSN is empty")
	}

	if conf.PluginPagerDutyAccessToken == "" {
		t.Fatal("PluginPagerDutyAccessToken is empty")
	}
}
