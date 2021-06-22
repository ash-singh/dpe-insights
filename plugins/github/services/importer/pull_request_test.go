package importer

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/sendinblue/dpe-insights/core/constants"
	"github.com/sendinblue/dpe-insights/plugins/github/models/entities"
	"github.com/google/go-github/v32/github"
	"github.com/stretchr/testify/mock"
)

type mockPrRepository struct {
	mock.Mock
}

func (m *mockPrRepository) Save(e *entities.PullRequest) error {
	args := m.Called(e)
	return args.Error(0)
}

type mockSearchClient struct {
	mock.Mock
}

func (m *mockSearchClient) Issues(ctx context.Context, query string, opts *github.SearchOptions) (*github.IssuesSearchResult, *github.Response, error) {
	args := m.Called(ctx, query, opts)
	return args.Get(0).(*github.IssuesSearchResult), args.Get(1).(*github.Response), args.Error(2)
}

type mockPullRequestClient struct {
	mock.Mock
}

func (m *mockPullRequestClient) Get(ctx context.Context, owner string, repo string, number int) (*github.PullRequest, *github.Response, error) {
	args := m.Called()
	return args.Get(0).(*github.PullRequest), args.Get(1).(*github.Response), args.Error(2)
}

func (m *mockPullRequestClient) ListCommits(ctx context.Context, owner string, repo string, number int, opts *github.ListOptions) ([]*github.RepositoryCommit, *github.Response, error) {
	args := m.Called()
	return args.Get(0).([]*github.RepositoryCommit), args.Get(1).(*github.Response), args.Error(2)
}

func TestPullRequest_ImportPullRequestData(t *testing.T) {
	mRepository := new(mockPrRepository)
	mSearchClient := new(mockSearchClient)

	timeNow := time.Time{}
	endDate := timeNow
	startDate := timeNow.AddDate(0, 0, 1)

	labels := []*github.Label{
		{
			ID:          github.Int64(int64(1)),
			Name:        github.String("test-label"),
			Description: github.String("test label description"),
		},
	}

	testOrganizationName := "TestOrg"
	repositoryName := "dpe-insights"

	issueResult := &github.IssuesSearchResult{
		Total:             github.Int(1),
		IncompleteResults: github.Bool(false),
		Issues: []*github.Issue{
			{
				ID:            github.Int64(int64(111)),
				Number:        github.Int(222),
				Title:         github.String("DPE-1 Test Pull request"),
				Body:          github.String("This is test pull request body"),
				Labels:        labels,
				RepositoryURL: github.String("https://api.github.com/repos/" + testOrganizationName + "/" + repositoryName),
			},
		},
	}

	mPullRequestClient := new(mockPullRequestClient)

	pullRequest := &github.PullRequest{
		ID:     github.Int64(int64(1)),
		Number: github.Int(1),
		Title:  github.String("DPE-1 Test Pull request"),
		Body:   github.String("This is test pull request body"),
		Labels: labels,
	}

	e := &entities.PullRequest{
		PrID:            111,
		PrNumber:        222,
		Title:           "DPE-1 Test Pull request",
		BranchName:      "",
		Body:            "This is test pull request body",
		RepositoryName:  repositoryName,
		Comments:        0,
		Labels:          "[\"test-label\"]",
		TransformStatus: constants.TransformationStatusPending,
		FirstCommitAt:   timeNow,
		PrCreatedAt:     timeNow,
		PrUpdatedAt:     timeNow,
		PrClosedAt:      timeNow,
		PrMergedAt:      timeNow,
		TransformAt:     timeNow,
	}

	pullRequestCommits := &github.RepositoryCommit{
		Commit: &github.Commit{
			Author: &github.CommitAuthor{
				Date:  &timeNow,
				Name:  github.String("Test Author"),
				Email: github.String("test+author@example.com"),
				Login: nil,
			},
		},
	}

	openPullRequestQuery := fmt.Sprintf("org:%s type:pr is:open created:%s..%s",
		testOrganizationName,
		startDate.Format(constants.DateFormatISO),
		endDate.Format(constants.DateFormatISO))

	closedPullRequestQuery := fmt.Sprintf("org:%s type:pr closed:%s..%s",
		testOrganizationName,
		startDate.Format(constants.DateFormatISO),
		endDate.Format(constants.DateFormatISO))

	searchOpts := &github.SearchOptions{
		Sort:      "",
		Order:     "",
		TextMatch: false,
		ListOptions: github.ListOptions{
			Page:    1,
			PerPage: recordsPerPage,
		},
	}

	mGitResponse := &github.Response{
		Rate: github.Rate{
			Limit:     5000,
			Remaining: 1000,
			Reset:     github.Timestamp{},
		},
	}

	// setup expectations

	mPullRequestClient.
		On("Get").Return(pullRequest, mGitResponse, nil)

	mPullRequestClient.
		On("ListCommits").Return([]*github.RepositoryCommit{pullRequestCommits}, mGitResponse, nil)

	mSearchClient.
		On("Issues", context.Background(), openPullRequestQuery, searchOpts).Return(issueResult, mGitResponse, nil)

	mSearchClient.
		On("Issues", context.Background(), closedPullRequestQuery, searchOpts).Return(issueResult, mGitResponse, nil)

	mRepository.
		On("Save", e).Return(nil)

	pullRequestImporter := PullRequest{
		SearchClient:           mSearchClient,
		PullRequestClient:      mPullRequestClient,
		Repository:             mRepository,
		GithubOrganizationName: testOrganizationName,
	}

	err := pullRequestImporter.ImportPullRequestData(startDate, endDate)
	if err != nil {
		t.Fatal("Pull request import failed")
	}

	mPullRequestClient.AssertExpectations(t)
	mSearchClient.AssertExpectations(t)
	mRepository.AssertExpectations(t)
}

func TestPullRequest_GetPullRequestStatsWithFilter(t *testing.T) {
	mRepository := new(mockPrRepository)
	mSearchClient := new(mockSearchClient)

	timeNow := time.Time{}

	labels := []*github.Label{
		{
			ID:          github.Int64(int64(1)),
			Name:        github.String("test-label"),
			Description: github.String("test label description"),
		},
	}

	issueResult := &github.IssuesSearchResult{
		Total:             github.Int(1),
		IncompleteResults: github.Bool(false),
		Issues: []*github.Issue{
			{
				ID:            github.Int64(int64(111)),
				Number:        github.Int(222),
				Title:         github.String("DPE-1 Test Pull request"),
				Body:          github.String("This is test pull request body"),
				Labels:        labels,
				RepositoryURL: github.String("https://api.github.com/repos/DTSL/account-symfony"),
			},
		},
	}

	mPullRequestClient := new(mockPullRequestClient)

	testOrganizationName := "TestOrg"

	openPullRequestQuery := fmt.Sprintf("org:%s is:pr is:open created:%s", testOrganizationName, timeNow.Format(constants.DateFormatISO))
	closedPullRequestQuery := fmt.Sprintf("org:%s is:pr closed:%s", testOrganizationName, timeNow.Format(constants.DateFormatISO))

	mGitResponse := &github.Response{
		Rate: github.Rate{
			Limit:     5000,
			Remaining: 1000,
			Reset:     github.Timestamp{},
		},
	}

	// setup expectations

	mSearchClient.
		On("Issues", context.Background(), openPullRequestQuery, (*github.SearchOptions)(nil)).Return(issueResult, mGitResponse, nil)
	mSearchClient.
		On("Issues", context.Background(), closedPullRequestQuery, (*github.SearchOptions)(nil)).Return(issueResult, mGitResponse, nil)

	pullRequestImporter := PullRequest{
		SearchClient:           mSearchClient,
		PullRequestClient:      mPullRequestClient,
		Repository:             mRepository,
		GithubOrganizationName: testOrganizationName,
	}

	_, err := pullRequestImporter.GetPullRequestStatsWithFilter(timeNow)
	if err != nil {
		t.Fatal("Pull request import failed")
	}

	mPullRequestClient.AssertExpectations(t)
	mSearchClient.AssertExpectations(t)
	mRepository.AssertExpectations(t)
}
