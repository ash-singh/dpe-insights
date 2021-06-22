package importer

import (
	"context"
	"fmt"

	"github.com/sendinblue/dpe-insights/plugins/github/helpers"
	"github.com/sendinblue/dpe-insights/plugins/github/models/entities"
	"github.com/google/go-github/v32/github"
)

// ReleaseClient interface.
type ReleaseClient interface {
	ListReleases(ctx context.Context, owner, repo string, opts *github.ListOptions) ([]*github.RepositoryRelease, *github.Response, error)
}

// ReleaseRepository interface.
type ReleaseRepository interface {
	Save(releaseEntity *entities.Release) error
}

// OrgRepository interface.
type OrgRepository interface {
	Fetch() ([]entities.Repository, error)
}

// ReleaseImporter import release data from git.
type ReleaseImporter struct {
	ReleaseClient     ReleaseClient
	ReleaseRepository ReleaseRepository
	Repository        OrgRepository
	OrganizationName  string
}

// ImportReleases Import all releases.
func (r *ReleaseImporter) ImportReleases(isDaily bool) error {
	reposArr, err := r.Repository.Fetch()
	if err != nil {
		return err
	}

	releasesMaxPage := getMaxPage(isDaily)
	for _, repo := range reposArr {
		releases := r.getReleases(repo.Name, releasesMaxPage)
		r.saveReleaseData(releases, repo.ID)
	}

	return nil
}

// getReleases Get maxReleaseCount records.
func (r *ReleaseImporter) getReleases(repoName string, maxPage int) []*github.RepositoryRelease {
	opts := &github.ListOptions{
		PerPage: recordsPerPage,
	}
	// get all repository releases
	var allReleases []*github.RepositoryRelease
	for {
		releases, resp, err := r.ReleaseClient.ListReleases(context.Background(), r.OrganizationName, repoName, opts)
		if err != nil {
			fmt.Println(err)
			break
		}

		helpers.RateLimitCoolDown(&resp.Rate)

		allReleases = append(allReleases, releases...)

		if resp.NextPage == 0 || opts.Page >= maxPage {
			break
		}

		opts.Page = resp.NextPage
	}

	return allReleases
}

// saveReleaseData Save ReleaseImporter data.
func (r *ReleaseImporter) saveReleaseData(rData []*github.RepositoryRelease, repoID int) {
	for _, release := range rData {
		releaseEntity := &entities.Release{
			ReleaseID:    int(release.GetID()),
			Title:        release.GetName(),
			TagName:      release.GetTagName(),
			Body:         release.GetBody(),
			RepositoryID: repoID,
			AuthorLogin:  release.GetAuthor().GetLogin(),
			CreatedAt:    release.GetCreatedAt().Time,
			PublishedAt:  release.GetPublishedAt().Time,
		}

		err := r.ReleaseRepository.Save(releaseEntity)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func getMaxPage(isDaily bool) int {
	if isDaily {
		return 0
	}

	return 5
}
