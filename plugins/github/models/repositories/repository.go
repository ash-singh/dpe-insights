package repositories

import (
	"fmt"

	"github.com/sendinblue/dpe-insights/plugins/github/models/entities"
	"github.com/jmoiron/sqlx"
)

// Repository model .
type Repository struct {
	Db sqlx.DB
}

const (
	tableNameRepository = "extracted_repositories"
)

// Insert repository record.
func (r *Repository) Insert(e entities.Repository) error {
	var err error

	query := "INSERT INTO " + tableNameRepository +
		"(id, name, size, description, open_issues, language, created_at, updated_at) " +
		"VALUES (:id, :name, :size, :description, :open_issues, :language, :updated_at, :created_at)"

	if _, err = r.Db.NamedExec(query, e); err != nil {
		return err
	}

	return nil
}

// Update repository.
func (r *Repository) Update(e entities.Repository) error {
	var err error

	query := "UPDATE " + tableNameTeams +
		" SET name = ? ," +
		" description = ?" +
		" WHERE id = ?"

	qry, args, err := sqlx.In(
		query,
		e.Name,
		e.Description,
		e.ID)
	if err != nil {
		return err
	}

	if _, err := r.Db.Exec(qry, args...); err != nil {
		return err
	}

	return nil
}

// Fetch repository.
func (r *Repository) Fetch() ([]entities.Repository, error) {
	var err error
	var data []entities.Repository

	query := "SELECT id, name  FROM " + tableNameRepository

	if err = r.Db.Select(&data, query); err != nil {
		return nil, fmt.Errorf("%q: %w", "Failed to fetch repository data from table", err)
	}

	return data, nil
}
