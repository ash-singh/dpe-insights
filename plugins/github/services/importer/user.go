package importer

import (
	"context"
	"fmt"

	"github.com/sendinblue/dpe-insights/plugins/github/helpers"
	"github.com/sendinblue/dpe-insights/plugins/github/models/entities"
	"github.com/google/go-github/v32/github"
)

// UserRepository interface.
type UserRepository interface {
	Save(e *entities.User) error
}

// UserClient GitHub Client Organizations service interface.
type UserClient interface {
	ListMembers(ctx context.Context, org string, opts *github.ListMembersOptions) ([]*github.User, *github.Response, error)
}

// UserImporter import pull request data from git.
type UserImporter struct {
	UserClient             UserClient
	Repository             UserRepository
	GithubOrganizationName string
}

// ImportUsers Import git users.
func (u *UserImporter) ImportUsers() error {
	page := 1
	for {
		opts := &github.ListMembersOptions{
			ListOptions: github.ListOptions{
				Page:    page,
				PerPage: recordsPerPage,
			},
		}

		users, res, err := u.UserClient.ListMembers(context.Background(), u.GithubOrganizationName, opts)

		helpers.RateLimitCoolDown(&res.Rate)

		if err != nil {
			return err
		}

		page++

		u.saveUsers(users)

		if len(users) < recordsPerPage {
			break
		}
	}

	return nil
}

func (u *UserImporter) saveUsers(users []*github.User) {
	for _, user := range users {
		err := u.saveUser(user)
		if err != nil {
			fmt.Println("Failed to save GitHub user: ", err)
		}
	}
}

func (u *UserImporter) saveUser(user *github.User) error {
	e := &entities.User{
		GithubID:        int(user.GetID()),
		GithubLoginName: user.GetLogin(),
		IsSiteAdmin:     user.GetSiteAdmin(),
		GithubUserType:  user.GetType(),
	}

	err := u.Repository.Save(e)
	if err != nil {
		return err
	}

	return nil
}
