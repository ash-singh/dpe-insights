package pullrequest

import (
	mysqldb "github.com/sendinblue/dpe-insights/core/databases/mysql"
	"github.com/sendinblue/dpe-insights/plugins/github/models/entities"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// TransformedTeamPRCountRepository repository for pull request count.
type TransformedTeamPRCountRepository struct {
	Db sqlx.DB
}

const (
	transformedTeamPullRequestCountTbl = "transformed_team_pull_request_count"
)

// Save Pull request count.
func (pc *TransformedTeamPRCountRepository) Save(e *entities.TransformedTeamPullRequestCount) error {
	var err error

	e.CloseTotalRatio = float64(e.Closed) / float64(e.Total)

	query := "INSERT INTO " + transformedTeamPullRequestCountTbl + " (total, open, closed, close_total_ratio, team_slug, date) VALUES (:total,:open,:closed, :close_total_ratio, :team_slug, :date)"

	if _, err = pc.Db.NamedExec(query, e); err != nil {
		if driverErr, ok := err.(*mysql.MySQLError); ok {
			if driverErr.Number == mysqldb.ERDupEntry {
				_ = pc.Update(e)
				return nil
			}
		}
		return err
	}

	return nil
}

// Update pull request data.
func (pc *TransformedTeamPRCountRepository) Update(e *entities.TransformedTeamPullRequestCount) error {
	var err error

	query := "UPDATE " + transformedTeamPullRequestCountTbl +
		" SET total = ?, open = ?, closed = ?, close_total_ratio=?" +
		" WHERE date = ? AND team_slug = ?"

	closeTotalRatio := float64(e.Closed) / float64(e.Total)

	qry, args, err := sqlx.In(
		query,
		e.Total,
		e.Open,
		e.Closed,
		closeTotalRatio,
		e.Date,
		e.TeamSlug)
	if err != nil {
		return err
	}

	if _, err := pc.Db.Exec(qry, args...); err != nil {
		return err
	}

	return nil
}
