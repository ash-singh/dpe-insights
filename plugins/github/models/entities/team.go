package entities

// Team model, table: extracted_teams.
type Team struct {
	ID          int    `json:"id,omitempty"`
	Name        string `json:"name" db:"name"`
	Slug        string `json:"slug" db:"slug"`
	Description string `json:"description" db:"description"`
}

// TeamRepository model, table: extracted_team_repositories.
type TeamRepository struct {
	GithubTeamID       int `json:"github_team_id" db:"github_team_id"`
	GithubRepositoryID int `json:"github_repository_id" db:"github_repository_id"`
}

// TeamUser model, table: extracted_team_users.
type TeamUser struct {
	GithubTeamID int `json:"github_team_id" db:"github_team_id"`
	GithubUserID int `json:"github_user_id" db:"github_user_id"`
}
