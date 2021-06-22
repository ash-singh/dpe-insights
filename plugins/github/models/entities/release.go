package entities

import "time"

// Release model, table: extracted_releases.
type Release struct {
	ID           int       `json:"id,omitempty"`
	ReleaseID    int       `json:"release_id" db:"release_id"`
	Title        string    `json:"title"`
	TagName      string    `json:"tag_name" db:"tag_name"`
	Body         string    `json:"body"`
	RepositoryID int       `json:"repository_id" db:"repository_id"`
	AuthorLogin  string    `json:"author_login" db:"author_login"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	PublishedAt  time.Time `json:"published_at" db:"published_at"`
}
