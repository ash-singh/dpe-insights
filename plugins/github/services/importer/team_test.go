package importer

import (
	"context"
	"testing"

	"github.com/sendinblue/dpe-insights/plugins/github/models/entities"
	"github.com/go-sql-driver/mysql"
	"github.com/google/go-github/v32/github"
	"github.com/stretchr/testify/mock"
)

type mockTeamRepository struct {
	mock.Mock
}

func (m *mockTeamRepository) Save(e entities.Team) error {
	args := m.Called(e)
	return args.Error(0)
}

func (m *mockTeamRepository) SaveTeamRepository(e entities.TeamRepository) error {
	args := m.Called(e)
	return args.Error(0)
}

func (m *mockTeamRepository) SaveTeamUser(e entities.TeamUser) error {
	args := m.Called(e)
	return args.Error(0)
}

type mockTeamClient struct {
	mock.Mock
}

func (c *mockTeamClient) ListTeams(ctx context.Context, org string, opts *github.ListOptions) ([]*github.Team, *github.Response, error) {
	args := c.Called(ctx, org)
	return args.Get(0).([]*github.Team), args.Get(1).(*github.Response), args.Error(2)
}

func (c *mockTeamClient) ListTeamReposBySlug(ctx context.Context, org, slug string, opts *github.ListOptions) ([]*github.Repository, *github.Response, error) {
	args := c.Called(ctx, org, slug, opts)
	return args.Get(0).([]*github.Repository), args.Get(1).(*github.Response), args.Error(2)
}

func (c *mockTeamClient) ListTeamMembersBySlug(ctx context.Context, org, slug string, opts *github.TeamListTeamMembersOptions) ([]*github.User, *github.Response, error) {
	args := c.Called(ctx, org, slug, opts)
	return args.Get(0).([]*github.User), args.Get(1).(*github.Response), args.Error(2)
}

func TestTeamImporter_ImportTeams(t *testing.T) {
	mTeamClient := new(mockTeamClient)
	mTeamRepository := new(mockTeamRepository)

	githubTeam := &github.Team{
		ID:          github.Int64(int64(1)),
		Name:        github.String("Test Team"),
		Slug:        github.String("test-team"),
		Description: github.String("this is unit test team"),
	}

	var repositories []*github.Repository

	repository := &github.Repository{
		ID: github.Int64(int64(1)),
	}

	repositories = append(repositories, repository)

	var users []*github.User

	user := &github.User{
		ID:   github.Int64(int64(1)),
		Name: github.String("Test User"),
	}

	users = append(users, user)

	teamEntity := entities.Team{
		ID:          1,
		Name:        "Test Team",
		Slug:        "test-team",
		Description: "this is unit test team",
	}

	teamRepoEntity := entities.TeamRepository{
		GithubRepositoryID: 1,
		GithubTeamID:       teamEntity.ID,
	}

	teamUserEntity := entities.TeamUser{
		GithubUserID: 1,
		GithubTeamID: teamEntity.ID,
	}

	ctx := context.Background()
	org := "TestOrg"

	opts := &github.ListOptions{
		Page:    1,
		PerPage: recordsPerPage,
	}

	var opt *github.TeamListTeamMembersOptions

	mGitResponse := &github.Response{
		Rate: github.Rate{
			Limit:     5000,
			Remaining: 1000,
			Reset:     github.Timestamp{},
		},
	}

	// setup expectations
	mTeamClient.
		On("ListTeams", ctx, org).Return([]*github.Team{githubTeam}, mGitResponse, nil).
		On("ListTeamReposBySlug", ctx, org, teamEntity.Slug, opts).Return(repositories, mGitResponse, nil).
		On("ListTeamMembersBySlug", ctx, org, teamEntity.Slug, opt).Return(users, mGitResponse, nil)

	mTeamRepository.
		On("Save", teamEntity).Return(nil).
		On("SaveTeamRepository", teamRepoEntity).Return(nil).
		On("SaveTeamUser", teamUserEntity).Return(nil)

	teamImporter := &TeamImporter{
		Repository:             mTeamRepository,
		TeamClient:             mTeamClient,
		GithubOrganizationName: org,
	}

	err := teamImporter.ImportTeams()
	if err != nil {
		t.Fatal("Team import failed")
	}

	mTeamRepository.AssertExpectations(t)
}

func TestTeamImporter_ImportTeams_Failed_Github_Fetch(t *testing.T) {
	mTeamClient := new(mockTeamClient)
	mTeamRepository := new(mockTeamRepository)

	ctx := context.Background()
	org := "TestOrg"

	githubError := &github.Error{
		Resource: "",
		Field:    "",
		Code:     "121",
		Message:  "User fetch failed",
	}

	mGitResponse := &github.Response{
		Rate: github.Rate{
			Limit:     5000,
			Remaining: 1000,
			Reset:     github.Timestamp{},
		},
	}

	// setup expectations
	mTeamClient.
		On("ListTeams", ctx, org).Return([]*github.Team{}, mGitResponse, githubError)

	teamImporter := &TeamImporter{
		Repository:             mTeamRepository,
		TeamClient:             mTeamClient,
		GithubOrganizationName: org,
	}

	err := teamImporter.ImportTeams()

	if err == nil {
		t.Fatal("Github fetch error is not managed")
	}

	mTeamRepository.AssertExpectations(t)
}

func TestTeamImporter_ImportTeams_Failed_Team_Save(t *testing.T) {
	mTeamClient := new(mockTeamClient)
	mTeamRepository := new(mockTeamRepository)

	ctx := context.Background()
	org := "TestOrg"

	team := &github.Team{
		ID:          github.Int64(int64(1)),
		Name:        github.String("Test Team"),
		Slug:        github.String("test-team"),
		Description: github.String("this is unit test team"),
	}

	teamEntity := entities.Team{
		ID:          1,
		Name:        "Test Team",
		Slug:        "test-team",
		Description: "this is unit test team",
	}

	insertErr := &mysql.MySQLError{
		Number:  11212,
		Message: "Save Error",
	}

	mGitResponse := &github.Response{
		Rate: github.Rate{
			Limit:     5000,
			Remaining: 1000,
			Reset:     github.Timestamp{},
		},
	}

	// setup expectations
	mTeamClient.
		On("ListTeams", ctx, org).Return([]*github.Team{team}, mGitResponse, nil)

	mTeamRepository.
		On("Save", teamEntity).Return(insertErr)

	teamImporter := &TeamImporter{
		Repository:             mTeamRepository,
		TeamClient:             mTeamClient,
		GithubOrganizationName: org,
	}

	err := teamImporter.ImportTeams()
	if err != nil {
		t.Fatal("Team import save error not managed")
	}

	mTeamClient.AssertNotCalled(t, "ListTeamReposBySlug")
	mTeamClient.AssertNotCalled(t, "ListTeamMembersBySlug")
	mTeamRepository.AssertExpectations(t)
}
