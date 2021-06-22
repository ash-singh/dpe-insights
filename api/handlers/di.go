package handlers

import (
	"sync"

	"github.com/sendinblue/dpe-insights/core/config"
	"github.com/sendinblue/dpe-insights/core/databases/mysql"
	"github.com/sendinblue/dpe-insights/plugins/github/helpers"
	"github.com/sendinblue/dpe-insights/plugins/github/models/repositories"
	pullRequest "github.com/sendinblue/dpe-insights/plugins/github/models/repositories/pull_request"
	"github.com/sendinblue/dpe-insights/plugins/github/services/importer"
	"github.com/sendinblue/dpe-insights/plugins/github/services/syncdata"
	"github.com/sendinblue/dpe-insights/plugins/github/services/transformer"
)

type diContainer struct {
	teamImporter           *importer.TeamImporter
	repositoryImporter     *importer.RepositoryImporter
	userImporter           *importer.UserImporter
	syncData               *syncdata.PullRequest
	prImporter             *importer.PullRequest
	prTransformer          *transformer.PullRequest
	prCountRepository      *pullRequest.CountRepository
	tprRepository          *pullRequest.TransformedPullRequest
	tprCountRepository     *pullRequest.TransformedPRCountRepository
	tprTeamCountRepository *pullRequest.TransformedTeamPRCountRepository
}

var (
	container *diContainer
	once      sync.Once
)

func (di *diContainer) GetPullRequestImporter() *importer.PullRequest {
	return di.prImporter
}

func (di *diContainer) GetTeamImporter() *importer.TeamImporter {
	return di.teamImporter
}

func (di *diContainer) GetRepositoryImporter() *importer.RepositoryImporter {
	return di.repositoryImporter
}

func (di *diContainer) GetUserImporter() *importer.UserImporter {
	return di.userImporter
}

func (di *diContainer) GetPullRequestTransformer() *transformer.PullRequest {
	return di.prTransformer
}

func (di *diContainer) GetSyncDataService() *syncdata.PullRequest {
	return di.syncData
}

func (di *diContainer) GetPullRequestCountRepository() *pullRequest.CountRepository {
	return di.prCountRepository
}

func (di *diContainer) GetTransformedPullRequestRepository() *pullRequest.TransformedPullRequest {
	return di.tprRepository
}

func (di *diContainer) GetTransformedPullRequestCountRepository() *pullRequest.TransformedPRCountRepository {
	return di.tprCountRepository
}

func (di *diContainer) GetTransformedTeamPullRequestCountRepository() *pullRequest.TransformedTeamPRCountRepository {
	return di.tprTeamCountRepository
}

func newDIContainer() *diContainer {
	di := &diContainer{}

	conf := config.NewConfig()

	db, _ := mysql.NewDB(conf)

	gitClient := helpers.NewClient()

	di.prCountRepository = &pullRequest.CountRepository{Db: *db}

	di.tprRepository = &pullRequest.TransformedPullRequest{Db: *db}

	di.tprCountRepository = &pullRequest.TransformedPRCountRepository{Db: *db}

	di.tprTeamCountRepository = &pullRequest.TransformedTeamPRCountRepository{Db: *db}

	di.prImporter = &importer.PullRequest{
		Repository:             &pullRequest.PullRequest{Db: *db},
		SearchClient:           gitClient.Search,
		PullRequestClient:      gitClient.PullRequests,
		GithubOrganizationName: conf.PluginGithubOrganizationName,
	}

	di.teamImporter = &importer.TeamImporter{
		TeamClient:             gitClient.Teams,
		Repository:             &repositories.Team{Db: *db},
		GithubOrganizationName: conf.PluginGithubOrganizationName,
	}

	di.repositoryImporter = &importer.RepositoryImporter{
		RepositoryClient:       gitClient.Repositories,
		Repository:             &repositories.Repository{Db: *db},
		GithubOrganizationName: conf.PluginGithubOrganizationName,
	}

	di.userImporter = &importer.UserImporter{
		UserClient:             gitClient.Organizations,
		Repository:             &repositories.User{Db: *db},
		GithubOrganizationName: conf.PluginGithubOrganizationName,
	}

	di.syncData = &syncdata.PullRequest{
		GitClient:       gitClient,
		Repository:      pullRequest.PullRequest{Db: *db},
		GithubOwnerName: conf.PluginGithubOwnerName,
	}

	di.prTransformer = &transformer.PullRequest{
		Repository:                       &pullRequest.PullRequest{Db: *db},
		TransformedPullRequestRepository: di.GetTransformedPullRequestRepository(),
	}

	return di
}

func getDIContainer() *diContainer {
	once.Do(func() {
		container = newDIContainer()
	})
	return container
}
