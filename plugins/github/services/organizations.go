package services

import (
	"context"

	"github.com/sendinblue/dpe-insights/plugins/github/helpers"
	"github.com/google/go-github/v32/github"
)

// GetOrganization Get organization details from github for a given orgID.
func GetOrganization(orgID int64) (*github.Organization, error) {
	client := helpers.NewClient()

	// get organization detail
	organization, _, err := client.Organizations.GetByID(context.Background(), orgID)

	return organization, err
}

// GetOrganizations Get list of all the organization from github.
func GetOrganizations() ([]*github.Organization, error) {
	client := helpers.NewClient()

	// list all organization for the authenticated user
	organizations, _, err := client.Organizations.List(context.Background(), "", nil)

	return organizations, err
}
