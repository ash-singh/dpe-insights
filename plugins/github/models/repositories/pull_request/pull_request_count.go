package pullrequest

import (
	mysqldb "github.com/sendinblue/dpe-insights/core/databases/mysql"
	"github.com/sendinblue/dpe-insights/plugins/github/models/entities"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// CountRepository repository for pull request count.
type CountRepository struct {
	Db sqlx.DB
}

const (
	pullRequestCountTbl = "extracted_pull_request_count"
)

// Insert Pull request count.
func (pc *CountRepository) Insert(pullRequestCount *entities.PullRequestCount) error {
	var err error
	query := "INSERT INTO " + pullRequestCountTbl + " (total, open, closed, date) VALUES (:total,:open,:closed,:date)"

	if _, err = pc.Db.NamedExec(query, pullRequestCount); err != nil {
		return err
	}

	return nil
}

// Update pull request data.
func (pc *CountRepository) Update(pullRequestCount *entities.PullRequestCount) error {
	var err error

	query := "UPDATE " + pullRequestCountTbl +
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

// Save Pull request count.
func (pc *CountRepository) Save(pullRequestCount *entities.PullRequestCount) error {
	err := pc.Insert(pullRequestCount)
	if err != nil {
		if driverErr, ok := err.(*mysql.MySQLError); ok {
			if driverErr.Number == mysqldb.ERDupEntry {
				_ = pc.Update(pullRequestCount)
			}
		} else {
			return err
		}
	}

	return nil
}
