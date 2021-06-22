package entities

// User model, table: extracted_users.
type User struct {
	ID              int    `json:"id,omitempty"`
	GithubID        int    `json:"github_id" db:"github_id"`
	GithubLoginName string `json:"github_login_name" db:"github_login_name"`
	GithubUserType  string `json:"github_user_type" db:"github_user_type"`
	IsSiteAdmin     bool   `json:"is_site_admin" db:"is_site_admin"`
}
