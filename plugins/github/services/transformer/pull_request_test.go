package transformer

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/sendinblue/dpe-insights/core/constants"
	"github.com/sendinblue/dpe-insights/plugins/github/models/entities"
	"github.com/stretchr/testify/mock"
)

var errTransformer = errors.New("repositoryError")

func TransformerError(msg string) error {
	return fmt.Errorf("%w: %s", errTransformer, msg)
}

type mockPullRequestRepository struct {
	mock.Mock
}

type mockTransformedPullRequest struct {
	mock.Mock
}

func (m *mockPullRequestRepository) Fetch() ([]entities.PullRequest, error) {
	args := m.Called()
	return args.Get(0).([]entities.PullRequest), args.Error(1)
}

func (m *mockPullRequestRepository) MarkTransformed(prIDs []int) {
	m.Called(prIDs)
}

func (mtpr *mockTransformedPullRequest) Save(transformedPR *entities.TransformedPullRequest) error {
	return nil
}

func (mtpr *mockTransformedPullRequest) FetchTeam(ownerID int, repositoryName string) string {
	return "test-team"
}

// TestTransformRecords_FetchError test pull request transformer failure.
func TestTransformRecords_FetchError(t *testing.T) {
	mockTransformedPullRequest := &mockTransformedPullRequest{}

	mockRepository := new(mockPullRequestRepository)
	// setup expectations
	mockRepository.On("Fetch").Return([]entities.PullRequest{}, TransformerError("transformer error"))

	transformer := PullRequest{
		Repository:                       mockRepository,
		TransformedPullRequestRepository: mockTransformedPullRequest,
	}

	r, err := transformer.TransformRecords()

	if err == nil {
		t.Fatalf("expected transformation err, got %d", r)
	}

	// assert that the expectations were met
	mockRepository.AssertExpectations(t)
}

// TestTransformRecords_Success test pull request transformer success.
func TestTransformRecords_Success(t *testing.T) {
	mockTransformedPullRequest := &mockTransformedPullRequest{}

	records := []entities.PullRequest{
		{
			ID:              1111,
			PrID:            2222,
			PrNumber:        3333,
			Title:           "DPE-21 Test pull request title",
			Body:            "Test body",
			RepositoryName:  "dpe-insights",
			Comments:        0,
			ReviewComments:  0,
			Commits:         0,
			Additions:       0,
			Deletions:       0,
			ChangedFiles:    0,
			TransformStatus: constants.TransformationStatusPending,
			OwnerLogin:      "",
			OwnerID:         0,
			FirstCommitAt:   time.Time{},
			PrCreatedAt:     time.Time{},
			PrUpdatedAt:     time.Time{},
			PrClosedAt:      time.Time{},
			PrMergedAt:      time.Time{},
			TransformAt:     time.Time{},
		},
	}
	mockRepository := new(mockPullRequestRepository)

	// setup expectations
	mockRepository.
		On("Fetch").Return(records, nil).Once().
		// pull request is marked transformed
		On("MarkTransformed", []int{2222}).Once().
		On("Fetch").Return([]entities.PullRequest{}, nil)

	transformer := PullRequest{
		Repository:                       mockRepository,
		TransformedPullRequestRepository: mockTransformedPullRequest,
	}

	r, _ := transformer.TransformRecords()

	if r != 1 {
		t.Fatalf("expected transformation record count to be 1, got %d", r)
	}
	// assert that the expectations were met
	mockRepository.AssertExpectations(t)
}

// TestTransformRecords_NoRecords test pull request transformer no records.
func TestTransformRecords_NoRecords(t *testing.T) {
	mockTransformedPullRequest := &mockTransformedPullRequest{}

	mockRepository := new(mockPullRequestRepository)

	// setup expectations
	mockRepository.
		On("Fetch").Return([]entities.PullRequest{}, nil).Once()

	transformer := PullRequest{
		Repository:                       mockRepository,
		TransformedPullRequestRepository: mockTransformedPullRequest,
	}

	r, err := transformer.TransformRecords()
	if err != nil {
		t.Fatalf("expected no transformation err, got %d", err)
	}

	if r != 0 {
		t.Fatalf("expected transformation record count to be 0, got %d", r)
	}

	// assert that the expectations were met
	mockRepository.AssertExpectations(t)
	mockRepository.AssertNotCalled(t, "MarkTransformed")
}
