package importer

import (
	"context"
	"fmt"

	"github.com/sendinblue/dpe-insights/plugins/github/helpers"
	"github.com/sendinblue/dpe-insights/plugins/github/models/entities"
	"github.com/google/go-github/v32/github"
)

// TeamRepository interface.
type TeamRepository interface {
	Save(e entities.Team) error
	SaveTeamRepository(e entities.TeamRepository) error
	SaveTeamUser(e entities.TeamUser) error
}

// TeamClient GitHub Client Teams service interface.
type TeamClient interface {
	ListTeams(ctx context.Context, org string, opts *github.ListOptions) ([]*github.Team, *github.Response, error)
	ListTeamReposBySlug(ctx context.Context, org, slug string, opts *github.ListOptions) ([]*github.Repository, *github.Response, error)
	ListTeamMembersBySlug(ctx context.Context, org, slug string, opts *github.TeamListTeamMembersOptions) ([]*github.User, *github.Response, error)
}

// TeamImporter import pull request data from GitHub.
type TeamImporter struct {
	TeamClient             TeamClient
	Repository             TeamRepository
	GithubOrganizationName string
}

// ImportTeams Import pull request data.
func (ti *TeamImporter) ImportTeams() error {
	teams, res, err := ti.TeamClient.ListTeams(context.Background(), ti.GithubOrganizationName, nil)
	if err != nil {
		return err
	}

	helpers.RateLimitCoolDown(&res.Rate)

	for _, team := range teams {
		err = ti.saveTeam(team)

		if err != nil {
			fmt.Println(err)
			continue
		}

		err = ti.importTeamRepository(team)
		if err != nil {
			fmt.Println(err)
		}

		err = ti.importTeamUser(team)
		if err != nil {
			fmt.Println(err)
		}
	}

	return nil
}

// importTeamRepository Import team's repositories.
func (ti *TeamImporter) importTeamRepository(team *github.Team) error {
	page := 1
	for {
		opts := &github.ListOptions{
			Page:    page,
			PerPage: recordsPerPage,
		}

		repositories, res, err := ti.TeamClient.ListTeamReposBySlug(context.Background(), ti.GithubOrganizationName, team.GetSlug(), opts)
		if err != nil {
			return err
		}

		helpers.RateLimitCoolDown(&res.Rate)

		if len(repositories) == 0 {
			return nil
		}

		page++

		ti.saveTeamRepository(team, repositories)

		if len(repositories) < recordsPerPage {
			return nil
		}
	}
}

// importTeamUser Import team's users.
func (ti *TeamImporter) importTeamUser(team *github.Team) error {
	users, res, err := ti.TeamClient.ListTeamMembersBySlug(context.Background(), ti.GithubOrganizationName, team.GetSlug(), nil)
	if err != nil {
		return err
	}

	helpers.RateLimitCoolDown(&res.Rate)

	if len(users) == 0 {
		return nil
	}

	ti.saveTeamUser(team, users)

	return nil
}

// saveTeamUser Save team's users.
func (ti *TeamImporter) saveTeamUser(team *github.Team, users []*github.User) {
	for _, u := range users {
		e := entities.TeamUser{
			GithubUserID: int(u.GetID()),
			GithubTeamID: int(team.GetID()),
		}

		err := ti.Repository.SaveTeamUser(e)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (ti *TeamImporter) saveTeamRepository(team *github.Team, repos []*github.Repository) {
	for _, r := range repos {
		e := entities.TeamRepository{
			GithubRepositoryID: int(r.GetID()),
			GithubTeamID:       int(team.GetID()),
		}

		err := ti.Repository.SaveTeamRepository(e)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (ti *TeamImporter) saveTeam(team *github.Team) error {
	e := entities.Team{
		ID:          int(team.GetID()),
		Name:        team.GetName(),
		Slug:        team.GetSlug(),
		Description: team.GetDescription(),
	}

	err := ti.Repository.Save(e)
	if err != nil {
		return err
	}

	return nil
}
