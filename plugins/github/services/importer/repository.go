package importer

import (
	"context"
	"fmt"

	mysqldb "github.com/sendinblue/dpe-insights/core/databases/mysql"
	"github.com/sendinblue/dpe-insights/plugins/github/helpers"
	"github.com/sendinblue/dpe-insights/plugins/github/models/entities"
	"github.com/go-sql-driver/mysql"
	"github.com/google/go-github/v32/github"
)

// Repository interface.
type Repository interface {
	Insert(e entities.Repository) error
	Update(e entities.Repository) error
	Fetch() ([]entities.Repository, error)
}

// RepositoryClient GitHub ReleaseClient Repositories service interface.
type RepositoryClient interface {
	ListByOrg(ctx context.Context, org string, opts *github.RepositoryListByOrgOptions) ([]*github.Repository, *github.Response, error)
}

// RepositoryImporter import pull request data from git.
type RepositoryImporter struct {
	RepositoryClient       RepositoryClient
	Repository             Repository
	GithubOrganizationName string
}

// ImportRepositories Import repository request data.
func (ri *RepositoryImporter) ImportRepositories() error {
	page := 1
	for {
		opts := &github.RepositoryListByOrgOptions{
			Sort: "",
			ListOptions: github.ListOptions{
				Page:    page,
				PerPage: recordsPerPage,
			},
		}

		repos, res, err := ri.RepositoryClient.ListByOrg(context.Background(), ri.GithubOrganizationName, opts)
		if err != nil {
			return err
		}

		helpers.RateLimitCoolDown(&res.Rate)

		page++

		ri.saveRepositories(repos)

		if len(repos) < recordsPerPage {
			break
		}
	}

	return nil
}

func (ri *RepositoryImporter) saveRepositories(repos []*github.Repository) {
	for _, r := range repos {
		err := ri.saveRepository(r)
		if err != nil {
			fmt.Println("Failed to save GitHub repository record!", err)
		}
	}
}

func (ri *RepositoryImporter) saveRepository(r *github.Repository) error {
	e := entities.Repository{
		ID:          int(r.GetID()),
		Name:        r.GetName(),
		Description: r.GetDescription(),
		Language:    r.GetLanguage(),
		Size:        r.GetSize(),
		OpenIssues:  r.GetOpenIssuesCount(),
		UpdatedAt:   r.GetUpdatedAt().Time,
		CreatedAt:   r.GetCreatedAt().Time,
	}

	err := ri.Repository.Insert(e)
	if err != nil {
		if driverErr, ok := err.(*mysql.MySQLError); ok {
			if driverErr.Number == mysqldb.ERDupEntry {
				return ri.Repository.Update(e)
			}
		}
		return err
	}

	return nil
}
