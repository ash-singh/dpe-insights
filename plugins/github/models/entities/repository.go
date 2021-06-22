package entities

import "time"

// Repository  Model, table: extracted_repository.
type Repository struct {
	ID          int       `json:"id,omitempty"`
	Size        int       `json:"size" db:"size"`
	OpenIssues  int       `json:"open_issues" db:"open_issues"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Language    string    `json:"language"`
	IsPrivate   bool      `json:"is_private" db:"is_private"`
	IsDisabled  bool      `json:"is_disabled" db:"is_disabled"`
	IsArchived  bool      `json:"is_archived" db:"is_archived"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
