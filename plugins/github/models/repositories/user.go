package repositories

import (
	mysqldb "github.com/sendinblue/dpe-insights/core/databases/mysql"
	"github.com/sendinblue/dpe-insights/plugins/github/models/entities"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// User model .
type User struct {
	Db sqlx.DB
}

const (
	tableNameUser = "extracted_users"
)

// Save User record.
func (u *User) Save(e *entities.User) error {
	var err error

	query := "INSERT INTO " + tableNameUser +
		"(github_id, github_login_name, github_user_type, is_site_admin) " +
		"VALUES (:github_id, :github_login_name, :github_user_type, :is_site_admin)"

	if _, err = u.Db.NamedExec(query, e); err != nil {
		if driverErr, ok := err.(*mysql.MySQLError); ok {
			if driverErr.Number == mysqldb.ERDupEntry {
				return u.Update(e)
			}
		}
		return err
	}

	return nil
}

// Update User.
func (u *User) Update(e *entities.User) error {
	var err error

	query := "UPDATE " + tableNameUser +
		" SET github_login_name = ?, " +
		" github_user_type = ? " +
		" WHERE id = ?"

	qry, args, err := sqlx.In(
		query,
		e.GithubLoginName,
		e.GithubUserType,
		e.GithubID)
	if err != nil {
		return err
	}

	if _, err := u.Db.Exec(qry, args...); err != nil {
		return err
	}

	return nil
}
