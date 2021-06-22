package pullrequest

import (
	"fmt"
	"log"
	"time"

	mysqldb "github.com/sendinblue/dpe-insights/core/databases/mysql"
	"github.com/sendinblue/dpe-insights/plugins/github/models/entities"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// PullRequest repository for pull requests.
type PullRequest struct {
	Db sqlx.DB
}

const (
	tableName = "extracted_pull_requests"
)

// Insert pull request data.
func (pr *PullRequest) Insert(prEntity *entities.PullRequest) (*entities.PullRequest, error) {
	var err error

	query := "INSERT INTO " + tableName +
		" (pr_id, pr_number, title, branch_name, body, repository_name, owner_login, owner_id, " +
		"comments, review_comments, labels, commits, additions, deletions, changed_files, " +
		"first_commit_at, pr_created_at, pr_updated_at, pr_closed_at, transform_at, transform_status, pr_merged_at) " +
		"VALUES (:pr_id, :pr_number,:title, :branch_name, :body, :repository_name, :owner_login, :owner_id, " +
		":comments, :review_comments, :labels, :commits, :additions, :deletions, :changed_files, " +
		":first_commit_at, :pr_created_at, :pr_updated_at, :pr_closed_at, :transform_at, :transform_status, :pr_merged_at)"

	if _, err = pr.Db.NamedExec(query, prEntity); err != nil {
		return nil, err
	}

	return prEntity, nil
}

// Update pull request data.
func (pr *PullRequest) Update(prEntity *entities.PullRequest) error {
	var err error

	query := "UPDATE " + tableName +
		" SET transform_status = 'pending', title = ?, branch_name =?, body = ?, " +
		" comments = ?, review_comments = ?,commits = ?, additions = ?,deletions = ?, changed_files = ?, " +
		" first_commit_at = ?, pr_updated_at = ?, pr_closed_at = ?, pr_merged_at = ?, labels = ?" +
		" WHERE pr_id = ?"

	qry, args, err := sqlx.In(
		query,
		prEntity.Title,
		prEntity.BranchName,
		prEntity.Body,
		prEntity.Comments,
		prEntity.ReviewComments,
		prEntity.Commits,
		prEntity.Additions,
		prEntity.Deletions,
		prEntity.ChangedFiles,
		prEntity.FirstCommitAt,
		prEntity.PrUpdatedAt,
		prEntity.PrClosedAt,
		prEntity.PrMergedAt,
		prEntity.Labels,
		prEntity.PrID)
	if err != nil {
		return err
	}

	if _, err := pr.Db.Exec(qry, args...); err != nil {
		return err
	}

	return nil
}

// Fetch fetch records to transform.
func (pr *PullRequest) Fetch() ([]entities.PullRequest, error) {
	var pullRequestData []entities.PullRequest
	var err error

	query := "SELECT *  FROM " + tableName + " WHERE transform_status = 'pending' LIMIT 1000"

	if err = pr.Db.Select(&pullRequestData, query); err != nil {
		return nil, fmt.Errorf("%q: %w", "Failed to fetch pull request data from table", err)
	}

	return pullRequestData, nil
}

// MarkTransformed mark pull request as transformed.
func (pr *PullRequest) MarkTransformed(prIDs []int) {
	query := "UPDATE " + tableName + " SET transform_status = 'done', transform_at =?  WHERE pr_id in (?)"

	qry, args, err := sqlx.In(query, time.Now(), prIDs)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := pr.Db.Exec(qry, args...); err != nil {
		log.Fatal(err)
	}
}

// FetchRecordsForSync Fetch pull request record for sync.
func (pr *PullRequest) FetchRecordsForSync() ([]entities.SyncRecord, error) {
	var syncRecords []entities.SyncRecord
	var err error

	query := "SELECT pr_id, pr_number, repository_name  FROM " + tableName +
		" WHERE pr_closed_at = '0000-00-00 00:00:00' AND transform_status = 'done'" +
		" LIMIT 100"

	if err = pr.Db.Select(&syncRecords, query); err != nil {
		return nil, fmt.Errorf("%q: %w", "Failed to pull request records for sync process from table", err)
	}
	return syncRecords, nil
}

// Save Insert or update pull request data.
func (pr *PullRequest) Save(prEntity *entities.PullRequest) error {
	_, err := pr.Insert(prEntity)
	if err != nil {

		if driverErr, ok := err.(*mysql.MySQLError); ok {
			if driverErr.Number == mysqldb.ERDupEntry {
				_ = pr.Update(prEntity)
				return nil
			}
		}

		return err
	}
	return nil
}
