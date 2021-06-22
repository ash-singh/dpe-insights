package importer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/sendinblue/dpe-insights/core/constants"
	"github.com/sendinblue/dpe-insights/plugins/github/helpers"
	"github.com/sendinblue/dpe-insights/plugins/github/models/entities"
	"github.com/google/go-github/v32/github"
)

// PullRequestRepository interface.
type PullRequestRepository interface {
	Save(prEntity *entities.PullRequest) error
}

// SearchClient github issue search client.
type SearchClient interface {
	Issues(ctx context.Context, query string, opts *github.SearchOptions) (*github.IssuesSearchResult, *github.Response, error)
}

// PullRequestClient github pull request client.
type PullRequestClient interface {
	Get(ctx context.Context, owner string, repo string, number int) (*github.PullRequest, *github.Response, error)
	ListCommits(ctx context.Context, owner string, repo string, number int, opts *github.ListOptions) ([]*github.RepositoryCommit, *github.Response, error)
}

// PullRequest import pull request data from git.
type PullRequest struct {
	PullRequestClient      PullRequestClient
	SearchClient           SearchClient
	Repository             PullRequestRepository
	GithubOrganizationName string
}

const (
	recordsPerPage = 100
)

// PullRequestStats Pull Request Stats.
type PullRequestStats struct {
	Total  int
	Open   int
	Closed int
}

// ImportPullRequestData Import pull request data.
func (pr *PullRequest) ImportPullRequestData(startDate time.Time, endDate time.Time) error {
	openPullRequestQuery := fmt.Sprintf("org:%s type:pr is:open created:%s..%s",
		pr.GithubOrganizationName,
		startDate.Format(constants.DateFormatISO),
		endDate.Format(constants.DateFormatISO))

	closedPullRequestQuery := fmt.Sprintf("org:%s type:pr closed:%s..%s",
		pr.GithubOrganizationName,
		startDate.Format(constants.DateFormatISO),
		endDate.Format(constants.DateFormatISO))

	err := pr.importDataFromQuery(openPullRequestQuery)
	if err != nil {
		return err
	}

	err = pr.importDataFromQuery(closedPullRequestQuery)

	if err != nil {
		return err
	}
	return nil
}

func (pr *PullRequest) importDataFromQuery(query string) error {
	page := 1
	for {
		searchOpts := &github.SearchOptions{
			Sort:      "",
			Order:     "",
			TextMatch: false,
			ListOptions: github.ListOptions{
				Page:    page,
				PerPage: recordsPerPage,
			},
		}

		result, res, err := pr.SearchClient.Issues(context.Background(), query, searchOpts)
		if err != nil {
			return err
		}

		helpers.RateLimitCoolDown(&res.Rate)

		page++

		pr.savePullRequestData(result)

		if len(result.Issues) < recordsPerPage {
			return nil
		}

	}
}

// GetPullRequestStatsWithFilter Fetch pull request stats with filter.
func (pr *PullRequest) GetPullRequestStatsWithFilter(prTime time.Time) (*entities.PullRequestCount, error) {
	openQuery := fmt.Sprintf("org:%s is:pr is:open created:%s", pr.GithubOrganizationName, prTime.Format(constants.DateFormatISO))
	closedQuery := fmt.Sprintf("org:%s is:pr closed:%s", pr.GithubOrganizationName, prTime.Format(constants.DateFormatISO))

	openPullRequestResult, res, err := pr.SearchClient.Issues(context.Background(), openQuery, nil)
	if err != nil {
		return nil, err
	}

	helpers.RateLimitCoolDown(&res.Rate)

	closedPullRequestResult, res, err := pr.SearchClient.Issues(context.Background(), closedQuery, nil)
	if err != nil {
		return nil, err
	}

	helpers.RateLimitCoolDown(&res.Rate)

	return &entities.PullRequestCount{
		Total:  closedPullRequestResult.GetTotal() + openPullRequestResult.GetTotal(),
		Closed: closedPullRequestResult.GetTotal(),
		Open:   openPullRequestResult.GetTotal(),
		Date:   prTime.Format(constants.DateFormatISO),
	}, nil
}

func (pr *PullRequest) savePullRequestData(prData *github.IssuesSearchResult) {
	for _, issue := range prData.Issues {

		pullRequest := *issue

		repoName := strings.SplitN(pullRequest.GetRepositoryURL(), pr.GithubOrganizationName+"/", 2)[1]

		prDetail, err := pr.fetchPullRequestDetail(repoName, pullRequest.GetNumber())
		if err != nil {
			fmt.Printf("failed to fetch pull request detail! for %d,  err:%s", pullRequest.GetID(), err)
			continue
		}

		firstCommitTime, err := pr.fetchFirstCommitTime(repoName, pullRequest.GetNumber())
		if err != nil {
			fmt.Printf("failed to fetch pull request first commit time! for %d,  err:%s", pullRequest.GetID(), err)
			continue
		}

		prEntity := &entities.PullRequest{
			PrID:            int(pullRequest.GetID()),
			PrNumber:        pullRequest.GetNumber(),
			Title:           pullRequest.GetTitle(),
			BranchName:      prDetail.GetHead().GetRef(),
			Body:            pullRequest.GetBody(),
			RepositoryName:  repoName,
			Comments:        pullRequest.GetComments(),
			ReviewComments:  prDetail.GetReviewComments(),
			Labels:          pr.getLabels(prDetail.Labels),
			Commits:         prDetail.GetCommits(),
			Additions:       prDetail.GetAdditions(),
			Deletions:       prDetail.GetDeletions(),
			ChangedFiles:    prDetail.GetChangedFiles(),
			FirstCommitAt:   firstCommitTime,
			PrCreatedAt:     pullRequest.GetCreatedAt(),
			PrUpdatedAt:     pullRequest.GetUpdatedAt(),
			PrClosedAt:      pullRequest.GetClosedAt(),
			PrMergedAt:      prDetail.GetMergedAt(),
			OwnerLogin:      pullRequest.User.GetLogin(),
			OwnerID:         int(pullRequest.User.GetID()),
			TransformStatus: constants.TransformationStatusPending,
			TransformAt:     time.Time{},
		}

		err = pr.Repository.Save(prEntity)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (pr *PullRequest) fetchPullRequestDetail(repoName string, pullRequestNumber int) (*github.PullRequest, error) {
	prDetail, res, err := pr.PullRequestClient.Get(context.Background(), pr.GithubOrganizationName, repoName, pullRequestNumber)
	if err != nil {
		return nil, err
	}

	helpers.RateLimitCoolDown(&res.Rate)

	return prDetail, nil
}

func (pr *PullRequest) getLabels(labels []*github.Label) string {
	l := []string{}

	for _, label := range labels {
		l = append(l, label.GetName())
	}

	jsonLabels, err := json.Marshal(l)
	if err != nil {
		jsonLabels = []byte{}
	}

	return string(jsonLabels)
}

func (pr *PullRequest) fetchFirstCommitTime(repoName string, pullRequestNumber int) (time.Time, error) {
	commitList, res, err := pr.PullRequestClient.ListCommits(context.Background(), pr.GithubOrganizationName, repoName, pullRequestNumber, nil)
	if err != nil {
		return time.Time{}, err
	}

	helpers.RateLimitCoolDown(&res.Rate)

	if len(commitList) == 0 {
		return time.Time{}, nil
	}
	return commitList[0].GetCommit().GetAuthor().GetDate(), nil
}
