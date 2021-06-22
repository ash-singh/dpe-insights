package importer

import (
	"context"
	"testing"

	"github.com/sendinblue/dpe-insights/plugins/github/models/entities"
	"github.com/go-sql-driver/mysql"
	"github.com/google/go-github/v32/github"
	"github.com/stretchr/testify/mock"
)

type mockUserRepository struct {
	mock.Mock
}

func (m *mockUserRepository) Save(e *entities.User) error {
	args := m.Called(e)
	return args.Error(0)
}

type mockUserClient struct {
	mock.Mock
}

func (m *mockUserClient) ListMembers(ctx context.Context, org string, opts *github.ListMembersOptions) ([]*github.User, *github.Response, error) {
	args := m.Called()
	return args.Get(0).([]*github.User), args.Get(1).(*github.Response), args.Error(2)
}

func TestUserImporter_ImportUsers(t *testing.T) {
	mockRepository := new(mockUserRepository)
	mockClient := new(mockUserClient)

	var recordList []*github.User

	user := &github.User{
		ID:        github.Int64(int64(1)),
		Login:     github.String("test-user"),
		Type:      github.String("user"),
		SiteAdmin: github.Bool(false),
		Name:      github.String("Test User"),
	}

	recordList = append(recordList, user)

	e := &entities.User{
		GithubID:        1,
		GithubLoginName: "test-user",
		GithubUserType:  "user",
		IsSiteAdmin:     false,
	}

	mGitResponse := &github.Response{
		Rate: github.Rate{
			Limit:     5000,
			Remaining: 1000,
			Reset:     github.Timestamp{},
		},
	}

	// setup expectations
	mockClient.
		On("ListMembers").Return(recordList, mGitResponse, nil)

	mockRepository.
		On("Save", e).Return(nil)

	importer := &UserImporter{
		Repository:             mockRepository,
		UserClient:             mockClient,
		GithubOrganizationName: "TesOrg",
	}

	err := importer.ImportUsers()
	if err != nil {
		t.Fatal("User import failed")
	}

	mockClient.AssertExpectations(t)
	mockRepository.AssertExpectations(t)
}

func TestUserImporter_ImportUsers_WithFailedSave(t *testing.T) {
	mockRepository := new(mockUserRepository)
	mockClient := new(mockUserClient)

	var recordList []*github.User

	user := &github.User{
		ID:        github.Int64(int64(1)),
		Login:     github.String("test-user"),
		Type:      github.String("user"),
		SiteAdmin: github.Bool(false),
		Name:      github.String("Test User"),
	}

	recordList = append(recordList, user)

	e := &entities.User{
		GithubID:        1,
		GithubLoginName: "test-user",
		GithubUserType:  "user",
		IsSiteAdmin:     false,
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
	mockClient.
		On("ListMembers").Return(recordList, mGitResponse, nil)

	mockRepository.
		On("Save", e).Return(insertErr)

	importer := &UserImporter{
		Repository:             mockRepository,
		UserClient:             mockClient,
		GithubOrganizationName: "TesOrg",
	}

	err := importer.ImportUsers()
	if err != nil {
		t.Fatal("User import failed")
	}

	mockClient.AssertExpectations(t)
	mockRepository.AssertExpectations(t)
}

func TestUserImporter_ImportUsers_WithGithubFetchError(t *testing.T) {
	mockRepository := new(mockUserRepository)
	mockClient := new(mockUserClient)

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
	mockClient.
		On("ListMembers").Return([]*github.User{}, mGitResponse, githubError)

	importer := &UserImporter{
		Repository:             mockRepository,
		UserClient:             mockClient,
		GithubOrganizationName: "TesOrg",
	}

	err := importer.ImportUsers()

	if err == nil {
		t.Fatal("Github fetch error is not managed")
	}

	mockClient.AssertExpectations(t)
}
