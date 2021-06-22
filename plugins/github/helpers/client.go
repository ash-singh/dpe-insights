package helpers

import (
	"context"

	"github.com/sendinblue/dpe-insights/core/config"
	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
)

// NewClient new github client.
func NewClient() *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.NewConfig().PluginGithubOauth2Token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	return client
}
