package repositories

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/sendinblue/dpe-insights/core/config"
	"github.com/sendinblue/dpe-insights/core/databases/mysql"
	"github.com/sendinblue/dpe-insights/plugins/pagerduty/models/entities"
	_ "github.com/sendinblue/dpe-insights/testing"
	"github.com/stretchr/testify/assert"
)

func TestIncidentRepository_Insert(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	db, _ := mysql.NewDB(config.NewConfig())
	repo := &IncidentRepository{Db: *db}
	duration, _ := time.ParseDuration("1s")
	now := time.Now().Truncate(time.Second)
	lastStatusChangedAt := now.Add(duration)
	incident := entities.Incident{
		ID:                  rand.Uint32(),
		Duration:            1,
		ServiceName:         "foo",
		CreatedAt:           now,
		LastStatusChangedAt: lastStatusChangedAt,
		Title:               "foo-title",
		Description:         "foo-description",
		Urgency:             "foo-urgency",
		Status:              "foo-status",
		FalsePositive:       true,
	}
	t.Cleanup(func() {
		_, _ = db.Exec("DELETE FROM "+tableNameUser+" WHERE id = ?", incident.ID)
	})

	err := repo.Insert(incident)
	assert.NoError(t, err)

	rowx := db.QueryRowx("SELECT * FROM "+tableNameUser+" WHERE id = ?", incident.ID)

	var actualIncident entities.Incident
	err = rowx.StructScan(&actualIncident)
	assert.NoError(t, err)
	assert.Equal(t, incident.ID, actualIncident.ID)
	assert.Equal(t, 1, actualIncident.Duration)
	assert.Equal(t, "foo", actualIncident.ServiceName)
	assert.Equal(t, now.Unix(), actualIncident.CreatedAt.Unix())
	assert.Equal(t, lastStatusChangedAt.Unix(), actualIncident.LastStatusChangedAt.Unix())
	assert.Equal(t, "foo-title", actualIncident.Title)
	assert.Equal(t, "foo-description", actualIncident.Description)
	assert.Equal(t, "foo-urgency", actualIncident.Urgency)
	assert.Equal(t, "foo-status", actualIncident.Status)
	assert.Equal(t, true, actualIncident.FalsePositive)

	fmt.Println(actualIncident)
}
