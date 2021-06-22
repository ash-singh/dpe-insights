package syncdata

import (
	"context"
	"fmt"

	"github.com/sendinblue/dpe-insights/plugins/github/models/entities"
	pullrequest "github.com/sendinblue/dpe-insights/plugins/github/models/repositories/pull_request"
	"github.com/google/go-github/v32/github"
)

// PullRequest Sync service for pull request data.
type PullRequest struct {
	Repository      pullrequest.PullRequest // PullRequest Repository
	GitClient       *github.Client          // GitClient for calling github api
	GithubOwnerName string
}

// syncOpenPullRequest pull request date with github.
func (pr *PullRequest) syncOpenPullRequest() (int, error) {
	pc := 0
	records, err := pr.Repository.FetchRecordsForSync()
	if err != nil {
		return 0, err
	}

	for _, record := range records {
		updatedPrRecord, _, err := pr.GitClient.PullRequests.Get(
			context.Background(),
			pr.GithubOwnerName,
			record.RepositoryName,
			record.PrNumber)
		if err != nil {
			fmt.Println(err)
			continue
		}
		prEntity := &entities.PullRequest{
			PrID:           record.PrID,
			PrNumber:       updatedPrRecord.GetNumber(),
			Title:          updatedPrRecord.GetTitle(),
			Body:           updatedPrRecord.GetBody(),
			Comments:       updatedPrRecord.GetComments(),
			Commits:        updatedPrRecord.GetCommits(),
			ReviewComments: updatedPrRecord.GetReviewComments(),
			Additions:      updatedPrRecord.GetAdditions(),
			Deletions:      updatedPrRecord.GetDeletions(),
			ChangedFiles:   updatedPrRecord.GetChangedFiles(),
			PrCreatedAt:    updatedPrRecord.GetCreatedAt(),
			PrUpdatedAt:    updatedPrRecord.GetUpdatedAt(),
			PrMergedAt:     updatedPrRecord.GetMergedAt(),
			PrClosedAt:     updatedPrRecord.GetClosedAt(),
		}

		err = pr.Repository.Update(prEntity)

		if err != nil {
			fmt.Println("Update pull request failed: ", err)
			continue
		}

		pc++
	}
	return pc, nil
}

// Sync pull request data.
func (pr *PullRequest) Sync() int {
	total := 0
	for {
		recordCount, err := pr.syncOpenPullRequest()
		if err != nil {
			fmt.Println(err)
			return 0
		}

		if recordCount == 0 {
			break
		}

		fmt.Println("Processed : ", recordCount)

		total += recordCount
	}

	return total
}
