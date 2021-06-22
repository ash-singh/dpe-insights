package repositories

import (
	"fmt"

	mysqldb "github.com/sendinblue/dpe-insights/core/databases/mysql"
	"github.com/sendinblue/dpe-insights/plugins/github/models/entities"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// Team repository for pull requests.
type Team struct {
	Db sqlx.DB
}

const (
	tableNameTeams = "extracted_teams"
)

// Save team record.
func (t *Team) Save(e entities.Team) error {
	var err error

	query := "INSERT INTO " + tableNameTeams +
		"(id, name, slug, description) " +
		"VALUES (:id, :name, :slug, :description)"

	if _, err = t.Db.NamedExec(query, e); err != nil {
		if driverErr, ok := err.(*mysql.MySQLError); ok {
			if driverErr.Number == mysqldb.ERDupEntry {
				return t.Update(e)
			}
		}
		return err
	}

	return nil
}

// SaveTeamRepository team-repository record.
func (t *Team) SaveTeamRepository(e entities.TeamRepository) error {
	var err error

	query := "INSERT INTO extracted_team_repositories " +
		"(github_team_id, github_repository_id)" +
		"VALUES (:github_team_id, :github_repository_id)"

	if _, err = t.Db.NamedExec(query, e); err != nil {
		if driverErr, ok := err.(*mysql.MySQLError); ok {
			if driverErr.Number == mysqldb.ERDupEntry {
				return nil
			}
		}
		return err
	}

	return nil
}

// SaveTeamUser team-user record.
func (t *Team) SaveTeamUser(e entities.TeamUser) error {
	var err error

	query := "INSERT INTO extracted_team_users " +
		"(github_team_id, github_user_id)" +
		"VALUES (:github_team_id, :github_user_id)"

	if _, err = t.Db.NamedExec(query, e); err != nil {
		if driverErr, ok := err.(*mysql.MySQLError); ok {
			if driverErr.Number == mysqldb.ERDupEntry {
				return nil
			}
		}
		return err
	}

	return nil
}

// Update team.
func (t *Team) Update(e entities.Team) error {
	var err error

	query := "UPDATE " + tableNameTeams +
		" SET description = ?" +
		" WHERE slug = ?"

	qry, args, err := sqlx.In(
		query,
		e.Description,
		e.Slug)
	if err != nil {
		return err
	}

	if _, err := t.Db.Exec(qry, args...); err != nil {
		return err
	}

	return nil
}

// Fetch team records.
func (t *Team) Fetch() ([]entities.Team, error) {
	var TeamData []entities.Team
	var err error

	query := "SELECT * FROM " + tableNameTeams + " LIMIT 1000"

	if err = t.Db.Select(&TeamData, query); err != nil {
		return nil, fmt.Errorf("%q: %w", "Failed to fetch team data from table", err)
	}

	return TeamData, nil
}
