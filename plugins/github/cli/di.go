package github

import (
	"sync"

	"github.com/sendinblue/dpe-insights/core/config"
	"github.com/sendinblue/dpe-insights/core/databases/mysql"
	"github.com/sendinblue/dpe-insights/plugins/github/helpers"
	"github.com/sendinblue/dpe-insights/plugins/github/models/repositories"
	pullRequestRepository "github.com/sendinblue/dpe-insights/plugins/github/models/repositories/pull_request"
	"github.com/sendinblue/dpe-insights/plugins/github/services/importer"
	"github.com/sendinblue/dpe-insights/plugins/github/services/syncdata"
	"github.com/sendinblue/dpe-insights/plugins/github/services/transformer"
)

type diContainer struct {
	config              *config.Config
	teamImporter        *importer.TeamImporter
	repositoryImporter  *importer.RepositoryImporter
	userImporter        *importer.UserImporter
	releaseImporter     *importer.ReleaseImporter
	pullRequestImporter *importer.PullRequest

	syncData               *syncdata.PullRequest
	pullRequestTransformer *transformer.PullRequest

	pullRequestCountRepository *pullRequestRepository.CountRepository

	tprRepository          *pullRequestRepository.TransformedPullRequest
	tprCountRepository     *pullRequestRepository.TransformedPRCountRepository
	tprTeamCountRepository *pullRequestRepository.TransformedTeamPRCountRepository
}

var (
	container *diContainer
	once      sync.Once
)

// GetConfig returns Config.
func (di *diContainer) GetConfig() *config.Config {
	return di.config
}

// GetTeamImporter returns TeamImporter.
func (di *diContainer) GetTeamImporter() *importer.TeamImporter {
	return di.teamImporter
}

// GetRepositoryImporter returns RepositoryImporter.
func (di *diContainer) GetRepositoryImporter() *importer.RepositoryImporter {
	return di.repositoryImporter
}

// GetUserImporter returns UserImporter.
func (di *diContainer) GetUserImporter() *importer.UserImporter {
	return di.userImporter
}

// GetSyncDataService returns sync data.
func (di *diContainer) GetSyncDataService() *syncdata.PullRequest {
	return di.syncData
}

// GetPullRequestTransformer returns pull request transformer.
func (di *diContainer) GetPullRequestTransformer() *transformer.PullRequest {
	return di.pullRequestTransformer
}

// GetReleaseImporter returns ReleaseImporter.
func (di *diContainer) GetReleaseImporter() *importer.ReleaseImporter {
	return di.releaseImporter
}

// GetPullRequestImporter returns PullRequest importer.
func (di *diContainer) GetPullRequestImporter() *importer.PullRequest {
	return di.pullRequestImporter
}

// GetPullRequestCountRepository returns Pull request count repository.
func (di *diContainer) GetPullRequestCountRepository() *pullRequestRepository.CountRepository {
	return di.pullRequestCountRepository
}

// GetTransformedPullRequestRepository return transformed pull request repository.
func (di *diContainer) GetTransformedPullRequestRepository() *pullRequestRepository.TransformedPullRequest {
	return di.tprRepository
}

// GetTransformedPullRequestCountRepository returns transformed pull request count repository.
func (di *diContainer) GetTransformedPullRequestCountRepository() *pullRequestRepository.TransformedPRCountRepository {
	return di.tprCountRepository
}

// GetTransformedTeamPullRequestCountRepository returns transformed teams pull request count repository.
func (di *diContainer) GetTransformedTeamPullRequestCountRepository() *pullRequestRepository.TransformedTeamPRCountRepository {
	return di.tprTeamCountRepository
}

func newDIContainer() *diContainer {
	dic := &diContainer{}

	dic.config = config.NewConfig()

	db, _ := mysql.NewDB(dic.config)

	gitClient := helpers.NewClient()

	dic.teamImporter = &importer.TeamImporter{
		Repository:             &repositories.Team{Db: *db},
		TeamClient:             gitClient.Teams,
		GithubOrganizationName: dic.config.PluginGithubOrganizationName,
	}

	dic.repositoryImporter = &importer.RepositoryImporter{
		Repository:             &repositories.Repository{Db: *db},
		RepositoryClient:       gitClient.Repositories,
		GithubOrganizationName: dic.config.PluginGithubOrganizationName,
	}

	dic.userImporter = &importer.UserImporter{
		Repository:             &repositories.User{Db: *db},
		UserClient:             gitClient.Organizations,
		GithubOrganizationName: dic.config.PluginGithubOrganizationName,
	}

	dic.syncData = &syncdata.PullRequest{
		Repository:      pullRequestRepository.PullRequest{Db: *db},
		GitClient:       helpers.NewClient(),
		GithubOwnerName: dic.config.PluginGithubOwnerName,
	}

	dic.pullRequestTransformer = &transformer.PullRequest{
		Repository:                       &pullRequestRepository.PullRequest{Db: *db},
		TransformedPullRequestRepository: &pullRequestRepository.TransformedPullRequest{Db: *db},
	}

	dic.releaseImporter = &importer.ReleaseImporter{
		ReleaseClient:     gitClient.Repositories,
		Repository:        &repositories.Repository{Db: *db},
		ReleaseRepository: &repositories.Release{Db: *db},
		OrganizationName:  dic.config.PluginGithubOrganizationName,
	}

	dic.pullRequestImporter = &importer.PullRequest{
		Repository:             &pullRequestRepository.PullRequest{Db: *db},
		SearchClient:           gitClient.Search,
		PullRequestClient:      gitClient.PullRequests,
		GithubOrganizationName: dic.config.PluginGithubOrganizationName,
	}

	dic.pullRequestCountRepository = &pullRequestRepository.CountRepository{Db: *db}
	dic.tprRepository = &pullRequestRepository.TransformedPullRequest{Db: *db}
	dic.tprCountRepository = &pullRequestRepository.TransformedPRCountRepository{Db: *db}
	dic.tprTeamCountRepository = &pullRequestRepository.TransformedTeamPRCountRepository{Db: *db}

	return dic
}

func getDIContainer() *diContainer {
	once.Do(func() {
		container = newDIContainer()
	})
	return container
}
