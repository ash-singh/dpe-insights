package entities

import "time"

// PullRequestCount model, table: extracted_pull_request_count .
type PullRequestCount struct {
	ID     int    `json:"id,omitempty"`
	Total  int    `json:"total,omitempty"`
	Closed int    `json:"closed"`
	Open   int    `json:"open"`
	Date   string `json:"date"`
}

// PullRequest model, table: extracted_pull_requests .
type PullRequest struct {
	ID              int       `json:"id,omitempty"`
	PrID            int       `json:"pr_id" db:"pr_id"`
	PrNumber        int       `json:"pr_number" db:"pr_number"`
	Title           string    `json:"title"`
	BranchName      string    `json:"branch_name" db:"branch_name"`
	Body            string    `json:"body"`
	RepositoryName  string    `json:"repository_name" db:"repository_name"`
	Comments        int       `json:"comments" db:"comments"`
	Labels          string    `json:"labels" db:"labels"`
	ReviewComments  int       `json:"review_comments" db:"review_comments"`
	Commits         int       `json:"commits" db:"commits"`
	Additions       int       `json:"additions" db:"additions"`
	Deletions       int       `json:"deletions" db:"deletions"`
	ChangedFiles    int       `json:"changed_files" db:"changed_files"`
	TransformStatus string    `json:"transform_status" db:"transform_status"`
	OwnerLogin      string    `json:"owner_login" db:"owner_login"`
	OwnerID         int       `json:"owner_id" db:"owner_id"`
	FirstCommitAt   time.Time `json:"first_commit_at" db:"first_commit_at"`
	PrCreatedAt     time.Time `json:"pr_created_at" db:"pr_created_at"`
	PrUpdatedAt     time.Time `json:"pr_updated_at" db:"pr_updated_at"`
	PrClosedAt      time.Time `json:"pr_closed_at" db:"pr_closed_at"`
	PrMergedAt      time.Time `json:"pr_merged_at" db:"pr_merged_at"`
	TransformAt     time.Time `json:"transform_at" db:"transform_at"`
}

// TransformedPullRequest model, table: transformed_pull_request_data .
type TransformedPullRequest struct {
	ID             int       `json:"id,omitempty"`
	PrID           int       `json:"pr_id" db:"pr_id"`
	PrNumber       int       `json:"pr_number" db:"pr_number"`
	Duration       int       `json:"duration" db:"duration"`
	Status         string    `json:"status" db:"status"`
	Title          string    `json:"title"`
	BranchName     string    `json:"branch_name" db:"branch_name"`
	Body           string    `json:"body"`
	RepositoryName string    `json:"repository_name" db:"repository_name"`
	OwnerLoginName string    `json:"owner_login_name" db:"owner_login_name"`
	OwnerID        int       `json:"owner_id" db:"owner_id"`
	TeamSlug       string    `json:"team_slug" db:"team_slug"`
	Comments       int       `json:"comments" db:"comments"`
	Labels         string    `json:"labels" db:"labels"`
	ReviewComments int       `json:"review_comments" db:"review_comments"`
	Commits        int       `json:"commits" db:"commits"`
	Additions      int       `json:"additions" db:"additions"`
	Deletions      int       `json:"deletions" db:"deletions"`
	ChangedFiles   int       `json:"changed_files" db:"changed_files"`
	FirstCommitAt  time.Time `json:"first_commit_at" db:"first_commit_at"`
	PrCreatedAt    time.Time `json:"created_at" db:"pr_created_at"`
	PrClosedAt     time.Time `json:"closed_at" db:"pr_closed_at"`
	PrUpdatedAt    time.Time `json:"pr_updated_at" db:"pr_updated_at"`
	PrMergedAt     time.Time `json:"pr_merged_at" db:"pr_merged_at"`
	SyncAt         time.Time `json:"sync_at" db:"sync_at"`
}

// TransformedPullRequestCount model, table: transformed_pull_request_count .
type TransformedPullRequestCount struct {
	ID     int    `json:"id,omitempty"`
	Total  int    `json:"total,omitempty"`
	Closed int    `json:"closed"`
	Open   int    `json:"open"`
	Date   string `json:"date"`
}

// TransformedTeamPullRequestCount model, table: transformed_team_pull_request_count .
type TransformedTeamPullRequestCount struct {
	ID              int     `json:"id,omitempty"`
	Total           int     `json:"total,omitempty"`
	Closed          int     `json:"closed"`
	Open            int     `json:"open"`
	CloseTotalRatio float64 `json:"close_total_ratio" db:"close_total_ratio"`
	TeamSlug        string  `json:"team_slug" db:"team_slug"`
	Date            string  `json:"date" db:"date"`
}

// SyncRecord used for sync of pull request data.
type SyncRecord struct {
	PrID           int    `json:"pr_id" db:"pr_id"`
	PrNumber       int    `json:"pr_number" db:"pr_number"`
	RepositoryName string `json:"repository_name" db:"repository_name"`
}
