package transformer

import (
	"fmt"
	"os"
	"time"

	"github.com/sendinblue/dpe-insights/plugins/github/models/entities"
)

// PullRequestRepository interface.
type PullRequestRepository interface {
	Fetch() ([]entities.PullRequest, error)
	MarkTransformed(prIDs []int)
}

// TransformedPullRequestRepository interface.
type TransformedPullRequestRepository interface {
	FetchTeam(ownerID int, repositoryName string) string
	Save(e *entities.TransformedPullRequest) error
}

// PullRequest trans extracted pull request data.
type PullRequest struct {
	Repository                       PullRequestRepository
	TransformedPullRequestRepository TransformedPullRequestRepository
}

// fetchRecordsToProcess fetch record for processing.
func (pr *PullRequest) fetchRecordsToProcess() ([]entities.PullRequest, error) {
	result, err := pr.Repository.Fetch()
	if err != nil {
		return nil, err
	}
	return result, nil
}

// transform exported PullRequest data.
func (pr *PullRequest) transform(prEntity entities.PullRequest) *entities.TransformedPullRequest {
	duration := 0
	status := "closed"
	var syncAt time.Time
	var timeDiff time.Duration

	if prEntity.PrClosedAt.IsZero() {
		status = "open"
		timeDiff = time.Since(prEntity.PrCreatedAt)
	} else {
		timeDiff = prEntity.PrClosedAt.Sub(prEntity.PrCreatedAt)
	}

	duration = int(timeDiff.Hours())

	teamSlug := pr.TransformedPullRequestRepository.FetchTeam(prEntity.OwnerID, prEntity.RepositoryName)

	if teamSlug == "global" {
		fmt.Println(prEntity.OwnerID , " has no team")
		fmt.Println(prEntity)
		os.Exit(1)
	}
	transformedPR := entities.TransformedPullRequest{
		PrID:           prEntity.PrID,
		PrNumber:       prEntity.PrNumber,
		Duration:       duration,
		Status:         status,
		Title:          prEntity.Title,
		BranchName:     prEntity.BranchName,
		Body:           prEntity.Body,
		RepositoryName: prEntity.RepositoryName,
		OwnerID:        prEntity.OwnerID,
		OwnerLoginName: prEntity.OwnerLogin,
		TeamSlug:       teamSlug,
		Commits:        prEntity.Commits,
		Labels:         prEntity.Labels,
		ReviewComments: prEntity.ReviewComments,
		Additions:      prEntity.Additions,
		Deletions:      prEntity.Deletions,
		ChangedFiles:   prEntity.ChangedFiles,
		FirstCommitAt:  prEntity.FirstCommitAt,
		PrCreatedAt:    prEntity.PrCreatedAt,
		SyncAt:         syncAt,
		PrClosedAt:     prEntity.PrClosedAt,
		PrUpdatedAt:    prEntity.PrUpdatedAt,
		PrMergedAt:     prEntity.PrMergedAt,
	}

	return &transformedPR
}

// markAsTransformed mark pull request as transformed.
func (pr *PullRequest) markAsTransformed(prIds []int) {
	pr.Repository.MarkTransformed(prIds)
}

// TransformRecords transform pull request data.
func (pr *PullRequest) TransformRecords() (int, error) {
	var processedPRs []int
	count := 0
	for {
		records, err := pr.fetchRecordsToProcess()
		if err != nil {
			return 0, err
		}

		if len(records) == 0 {
			break
		}

		for _, r := range records {
			// transform new data
			e := pr.transform(r)

			err := pr.TransformedPullRequestRepository.Save(e)
			if err != nil {
				fmt.Println(fmt.Errorf("%q: %w", "Transformed pull request save err", err))
				continue
			}

			processedPRs = append(processedPRs, e.PrID)
		}

		pr.markAsTransformed(processedPRs)

		count += len(processedPRs)

		processedPRs = []int{}
	}

	return count, nil
}
