package importer

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/sendinblue/dpe-insights/plugins/github/models/entities"
	"github.com/go-sql-driver/mysql"
	"github.com/google/go-github/v32/github"
	"github.com/stretchr/testify/mock"
)

type mockReleaseClient struct {
	mock.Mock
}

func (m *mockReleaseClient) ListReleases(ctx context.Context, org string, repo string, opts *github.ListOptions) ([]*github.RepositoryRelease, *github.Response, error) {
	args := m.Called(ctx, org, repo, opts)
	return args.Get(0).([]*github.RepositoryRelease), args.Get(1).(*github.Response), args.Error(2)
}

type mockReleaseRepository struct {
	mock.Mock
}

func (r *mockReleaseRepository) Save(e *entities.Release) error {
	args := r.Called(e)
	return args.Error(0)
}

type mkRepository struct {
	mock.Mock
}

func (r *mkRepository) Fetch() ([]entities.Repository, error) {
	args := r.Called()
	return args.Get(0).([]entities.Repository), args.Error(1)
}

var errNew = errors.New("error")

func showError(msg string) error {
	return fmt.Errorf("%q: %w", msg, errNew)
}

func TestReleaseImporter_ListReleasesError(t *testing.T) {
	mRepository := new(mkRepository)
	mReleaseClient := new(mockReleaseClient)

	repoRecords := []entities.Repository{
		{
			ID:   1,
			Name: "Test Repo",
		},
	}

	ctx := context.Background()
	org := "TestOrg"
	var releaseData []*github.RepositoryRelease
	opts := &github.ListOptions{
		PerPage: recordsPerPage,
	}

	mGitResponse := &github.Response{
		Rate: github.Rate{
			Limit:     5000,
			Remaining: 1000,
			Reset:     github.Timestamp{},
		},
	}

	mRepository.On("Fetch").Return(repoRecords, nil)
	mReleaseClient.On("ListReleases", ctx, org, repoRecords[0].Name, opts).Return(releaseData, mGitResponse, showError("List Release Error"))

	releaseImporter := &ReleaseImporter{
		Repository:       mRepository,
		ReleaseClient:    mReleaseClient,
		OrganizationName: org,
	}

	err := releaseImporter.ImportReleases(true)
	if err != nil {
		t.Fatalf("expected repository err, got %s", err)
	}

	mRepository.AssertExpectations(t)
	mReleaseClient.AssertExpectations(t)
}

func TestReleaseImporter_SaveError(t *testing.T) {
	timeNow := time.Now()
	tEntity := &entities.Release{
		ReleaseID:    1,
		Title:        "Test Release",
		TagName:      "v0.0.1",
		Body:         "Release description",
		RepositoryID: 1,
		CreatedAt:    timeNow,
		PublishedAt:  timeNow,
	}
	insertErr := &mysql.MySQLError{
		Number:  11212,
		Message: "Save Error",
	}
	mReleaseRepository := new(mockReleaseRepository)
	mRepository := new(mkRepository)
	mReleaseClient := new(mockReleaseClient)

	repoRecords := []entities.Repository{
		{
			ID:   1,
			Name: "Test Repo",
		},
	}

	ctx := context.Background()
	org := "TestOrg"
	currTime := github.Timestamp{Time: timeNow}
	var releaseData []*github.RepositoryRelease
	releaseEntity := &github.RepositoryRelease{
		ID:          github.Int64(int64(1)),
		Name:        github.String("Test Release"),
		TagName:     github.String("v0.0.1"),
		Body:        github.String("Release description"),
		CreatedAt:   &currTime,
		PublishedAt: &currTime,
	}
	releaseData = append(releaseData, releaseEntity)
	opts := &github.ListOptions{
		PerPage: recordsPerPage,
	}

	mGitResponse := &github.Response{
		Rate: github.Rate{
			Limit:     5000,
			Remaining: 1000,
			Reset:     github.Timestamp{},
		},
	}

	mRepository.On("Fetch").Return(repoRecords, nil)
	mReleaseRepository.On("Save", tEntity).Return(insertErr)
	mReleaseClient.On("ListReleases", ctx, org, repoRecords[0].Name, opts).Return(releaseData, mGitResponse, nil)

	releaseImporter := &ReleaseImporter{
		ReleaseRepository: mReleaseRepository,
		Repository:        mRepository,
		ReleaseClient:     mReleaseClient,
		OrganizationName:  org,
	}

	err := releaseImporter.ImportReleases(true)
	if err != nil {
		t.Fatalf("expected repository err, got %s", err)
	}

	mRepository.AssertExpectations(t)
	mReleaseRepository.AssertExpectations(t)
	mReleaseClient.AssertExpectations(t)
}

func TestReleaseImporter_FetchError(t *testing.T) {
	mRepository := new(mkRepository)
	mRepository.On("Fetch").Return([]entities.Repository{}, showError("Fetch Repository Error"))

	releaseImporter := &ReleaseImporter{
		Repository:       mRepository,
		OrganizationName: "TestOrg",
	}
	err := releaseImporter.ImportReleases(true)
	if err == nil {
		t.Fatalf("expected repository err, got %s", err)
	}

	mRepository.AssertExpectations(t)
}

func TestReleaseImporter_ImportReleases(t *testing.T) {
	mReleaseClient := new(mockReleaseClient)
	mReleaseRepository := new(mockReleaseRepository)
	mRepository := new(mkRepository)

	ctx := context.Background()
	org := "TestOrg"
	timeNow := time.Now()
	currTime := github.Timestamp{Time: timeNow}
	var releaseData []*github.RepositoryRelease

	releaseEntity := &github.RepositoryRelease{
		ID:          github.Int64(int64(1)),
		Name:        github.String("Test Release"),
		TagName:     github.String("v0.0.1"),
		Body:        github.String("Release description"),
		CreatedAt:   &currTime,
		PublishedAt: &currTime,
	}
	releaseData = append(releaseData, releaseEntity)

	tEntity := &entities.Release{
		ReleaseID:    1,
		Title:        "Test Release",
		TagName:      "v0.0.1",
		Body:         "Release description",
		RepositoryID: 1,
		CreatedAt:    timeNow,
		PublishedAt:  timeNow,
	}

	opts := &github.ListOptions{
		PerPage: recordsPerPage,
	}
	repoRecords := []entities.Repository{
		{
			ID:   1,
			Name: "Test Repo",
		},
	}

	mGitResponse := &github.Response{
		Rate: github.Rate{
			Limit:     5000,
			Remaining: 1000,
			Reset:     github.Timestamp{},
		},
	}

	mRepository.On("Fetch").Return(repoRecords, nil)
	mReleaseClient.On("ListReleases", ctx, org, repoRecords[0].Name, opts).Return(releaseData, mGitResponse, nil)
	mReleaseRepository.On("Save", tEntity).Return(nil)

	releaseImporter := &ReleaseImporter{
		ReleaseClient:     mReleaseClient,
		Repository:        mRepository,
		ReleaseRepository: mReleaseRepository,
		OrganizationName:  org,
	}

	err := releaseImporter.ImportReleases(true)
	if err != nil {
		t.Fatal("Release import failed")
	}

	mRepository.AssertExpectations(t)
	mReleaseClient.AssertExpectations(t)
	mReleaseRepository.AssertExpectations(t)
}
