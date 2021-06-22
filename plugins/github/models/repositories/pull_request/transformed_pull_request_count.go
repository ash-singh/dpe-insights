package pullrequest

import (
	mysqldb "github.com/sendinblue/dpe-insights/core/databases/mysql"
	"github.com/sendinblue/dpe-insights/plugins/github/models/entities"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// TransformedPRCountRepository repository for pull request count.
type TransformedPRCountRepository struct {
	Db sqlx.DB
}

const (
	transformedPullRequestCountTbl = "transformed_pull_request_count"
)

// Insert Pull request count.
func (pc *TransformedPRCountRepository) Insert(pullRequestCount *entities.TransformedPullRequestCount) error {
	var err error

	query := "INSERT INTO " + transformedPullRequestCountTbl + " (total, open, closed, date) VALUES (:total,:open,:closed,:date)"

	if _, err = pc.Db.NamedExec(query, pullRequestCount); err != nil {
		return err
	}

	return nil
}

// Update pull request data.
func (pc *TransformedPRCountRepository) Update(pullRequestCount *entities.TransformedPullRequestCount) error {
	var err error

	query := "UPDATE " + transformedPullRequestCountTbl +
		" SET total = ?, open = ?, closed = ?" +
		" WHERE date = ?"

	qry, args, err := sqlx.In(
		query,
		pullRequestCount.Total,
		pullRequestCount.Open,
		pullRequestCount.Closed,
		pullRequestCount.Date)
	if err != nil {
		return err
	}

	if _, err := pc.Db.Exec(qry, args...); err != nil {
		return err
	}

	return nil
}

// Save insert or update pull request data.
func (pc *TransformedPRCountRepository) Save(e *entities.TransformedPullRequestCount) error {
	err := pc.Insert(e)
	if err != nil {
		if driverErr, ok := err.(*mysql.MySQLError); ok {
			if driverErr.Number == mysqldb.ERDupEntry {
				return pc.Update(e)
			}
		}
		return err
	}
	return nil
}
