package importer

import (
	"context"
	"testing"

	"github.com/sendinblue/dpe-insights/plugins/github/models/entities"
	"github.com/google/go-github/v32/github"
	"github.com/stretchr/testify/mock"
)

type mockRepository struct {
	mock.Mock
}

func (m *mockRepository) Insert(e entities.Repository) error {
	return nil
}

func (m *mockRepository) Update(e entities.Repository) error {
	return nil
}

func (m *mockRepository) Fetch() ([]entities.Repository, error) {
	args := m.Called()
	return args.Get(0).([]entities.Repository), args.Error(1)
}

type mockRepositoryClient struct{}

func (m *mockRepositoryClient) ListByOrg(ctx context.Context, org string, opts *github.RepositoryListByOrgOptions) ([]*github.Repository, *github.Response, error) {
	var recordList []*github.Repository

	repository := &github.Repository{
		ID:   github.Int64(int64(1)),
		Name: github.String("Test Repository"),
	}

	recordList = append(recordList, repository)

	return recordList, &github.Response{Rate: github.Rate{
		Limit:     5000,
		Remaining: 1000,
		Reset:     github.Timestamp{},
	}}, nil
}

func TestRepositoryImporter_ImportRepositories(t *testing.T) {
	mRepository := &mockRepository{}
	mClient := &mockRepositoryClient{}

	importer := &RepositoryImporter{
		Repository:             mRepository,
		RepositoryClient:       mClient,
		GithubOrganizationName: "TestRep",
	}

	err := importer.ImportRepositories()
	if err != nil {
		t.Fatal("Repository import failed")
	}
}
