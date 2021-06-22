package repositories

import (
	"github.com/sendinblue/dpe-insights/plugins/pagerduty/models/entities"
	"github.com/jmoiron/sqlx"
)

// IncidentRepository repository for incidents.
type IncidentRepository struct {
	Db sqlx.DB
}

const (
	tableNameUser = "extracted_pagerduty_incidents"
)

// Insert incident data.
func (in *IncidentRepository) Insert(i entities.Incident) error {
	var err error

	query := "REPLACE INTO " + tableNameUser +
		"(id, created_at, last_status_change_at, service_name, duration, status, urgency, title, description, false_positive) " +
		"VALUES (:id, :created_at, :last_status_change_at, :service_name, :duration, :status, :urgency, :title, :description, :false_positive)"

	if _, err = in.Db.NamedExec(query, i); err != nil {
		return err
	}

	return nil
}
