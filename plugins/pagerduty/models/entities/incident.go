package entities

import (
	"time"
)

// Incident entity.
type Incident struct {
	ID                  uint32          `json:"incident_number" db:"id"`
	CreatedAt           time.Time       `json:"created_at" db:"created_at"`
	Urgency             string          `json:"urgency" db:"urgency"`
	Status              string          `json:"status" db:"status"`
	Title               string          `json:"title" db:"title"`
	Description         string          `json:"description" db:"description"`
	LastStatusChangedAt time.Time       `json:"last_status_change_at" db:"last_status_change_at"`
	Service             IncidentService `json:"service"`
	FalsePositive       bool            `json:"false_positive" db:"false_positive"`

	ServiceName string `db:"service_name"`
	Duration    int    `db:"duration"`
}

// IncidentService entity.
type IncidentService struct {
	ID      string `json:"id"`
	Summary string `json:"summary"`
}
