package repositories

import (
	mysqldb "github.com/sendinblue/dpe-insights/core/databases/mysql"
	"github.com/sendinblue/dpe-insights/plugins/github/models/entities"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// Release model .
type Release struct {
	Db sqlx.DB
}

const (
	tableNameRelease = "extracted_releases"
)

// Insert Release record.
func (u *Release) Insert(releaseEntity *entities.Release) (*entities.Release, error) {
	var err error

	query := "INSERT INTO " + tableNameRelease +
		"(release_id, title, tag_name, body, repository_id, author_login, created_at, published_at) " +
		"VALUES (:release_id ,:title, :tag_name, :body, :repository_id, :author_login, :created_at, :published_at)"

	if _, err = u.Db.NamedExec(query, releaseEntity); err != nil {
		return nil, err
	}

	return releaseEntity, nil
}

// Update pull request data.
func (u *Release) Update(releaseEntity *entities.Release) error {
	var err error

	query := "UPDATE " + tableNameRelease +
		" SET title = ?, tag_name =?, body = ?, " +
		" repository_id = ?, author_login = ?,created_at = ?, published_at = ?" +
		" WHERE release_id = ?"

	qry, args, err := sqlx.In(
		query,
		releaseEntity.Title,
		releaseEntity.TagName,
		releaseEntity.Body,
		releaseEntity.RepositoryID,
		releaseEntity.AuthorLogin,
		releaseEntity.CreatedAt,
		releaseEntity.PublishedAt,
		releaseEntity.ID)
	if err != nil {
		return err
	}

	if _, err := u.Db.Exec(qry, args...); err != nil {
		return err
	}

	return nil
}

// Save insert or update pull request data.
func (u *Release) Save(releaseEntity *entities.Release) error {
	_, err := u.Insert(releaseEntity)
	if err != nil {

		if driverErr, ok := err.(*mysql.MySQLError); ok {
			if driverErr.Number == mysqldb.ERDupEntry {
				_ = u.Update(releaseEntity)
				return nil
			}
		}

		return err
	}
	return nil
}
