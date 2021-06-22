package pullrequest

import (
	"fmt"
	"strconv"
	"time"

	"github.com/sendinblue/dpe-insights/core/constants"
	mysqldb "github.com/sendinblue/dpe-insights/core/databases/mysql"
	"github.com/sendinblue/dpe-insights/core/helpers"
	"github.com/sendinblue/dpe-insights/plugins/github/models/entities"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// TransformedPullRequest repository for pull requests.
type TransformedPullRequest struct {
	Db sqlx.DB
}

const (
	transformedPRTbl = "transformed_pull_request_data"
)

// Insert open pull request data.
func (tpr *TransformedPullRequest) Insert(e *entities.TransformedPullRequest) error {
	query := "INSERT INTO " + transformedPRTbl +
		"(pr_id, pr_number, duration, status, title, branch_name, body, " +
		" repository_name, owner_id, owner_login_name, team_slug," +
		" comments, labels, review_comments, commits, additions, deletions, changed_files, " +
		" first_commit_at, pr_updated_at, pr_created_at, pr_closed_at, pr_merged_at) " +
		"VALUES (:pr_id, :pr_number, :duration,:status, :title, :branch_name, :body, " +
		" :repository_name, :owner_id, :owner_login_name, :team_slug, " +
		" :comments, :labels, :review_comments, :commits, :additions, :deletions, :changed_files, " +
		" :first_commit_at, :pr_updated_at, :pr_created_at, :pr_closed_at, :pr_merged_at)"
	var err error

	if _, err = tpr.Db.NamedExec(query, *e); err != nil {
		return err
	}

	return nil
}

// Update pull request data.
func (tpr *TransformedPullRequest) Update(e *entities.TransformedPullRequest) error {
	var err error

	query := "UPDATE " + transformedPRTbl +
		" SET duration = ?, title = ?, branch_name =?, body = ?, team_slug = ?, status = ?," +
		" comments = ?, labels = ?, review_comments = ?,commits = ?, additions = ?,deletions = ?, changed_files = ?, " +
		" first_commit_at = ?, pr_updated_at = ?, pr_closed_at = ?, pr_merged_at = ?" +
		" WHERE pr_id = ?"

	qry, args, err := sqlx.In(
		query,
		e.Duration,
		e.Title,
		e.BranchName,
		e.Body,
		e.TeamSlug,
		e.Status,
		e.Comments,
		e.Labels,
		e.ReviewComments,
		e.Commits,
		e.Additions,
		e.Deletions,
		e.ChangedFiles,
		e.FirstCommitAt,
		e.PrUpdatedAt,
		e.PrClosedAt,
		e.PrMergedAt,
		e.PrID)
	if err != nil {
		return err
	}

	if _, err := tpr.Db.Exec(qry, args...); err != nil {
		return err
	}

	return nil
}

// FetchCountStats Fetch fetch records to transform.
func (tpr *TransformedPullRequest) FetchCountStats(dateFilter time.Time) (*entities.TransformedPullRequestCount, error) {
	var transformedPrCount entities.TransformedPullRequestCount
	var err error
	var total int
	var closed int

	totalQuery := "SELECT count(*) as total FROM " + transformedPRTbl + " WHERE pr_created_at < '" + dateFilter.Format(constants.DateFormatISO) + "'"

	if err = tpr.Db.QueryRowx(totalQuery).Scan(&total); err != nil {
		return nil, fmt.Errorf("%q: %w", "FetchCountStats total count fetch err", err)
	}

	closedQuery := "SELECT count(*) as closed FROM " + transformedPRTbl + " WHERE status = 'closed' AND pr_created_at < '" + dateFilter.Format(constants.DateFormatISO) + "'"

	if err = tpr.Db.QueryRowx(closedQuery).Scan(&closed); err != nil {
		return nil, fmt.Errorf("%q: %w", "FetchCountStats closed count fetch err", err)
	}

	transformedPrCount.Open = total - closed
	transformedPrCount.Total = total
	transformedPrCount.Closed = closed
	transformedPrCount.Date = dateFilter.Format(constants.DateFormatISO)

	return &transformedPrCount, nil
}

// FetchTeam for a pull request.
func (tpr *TransformedPullRequest) FetchTeam(ownerID int, repositoryName string) string {
	rteams := tpr.fetchRepositoryTeams(repositoryName)

	uteams := tpr.fetchUserTeams(ownerID)

	team := tpr.fetchMatchingTeam(rteams, uteams)

	return team
}

func (tpr *TransformedPullRequest) fetchMatchingTeam(repoTeams []string, userTeams []string) string {
	var teams []string
	filterTeams := []string{"sendinblue", "php-review", "go-review", "node-js-review"}

	userTeams = helpers.Difference(userTeams, filterTeams)
	repoTeams = helpers.Difference(repoTeams, filterTeams)

	hash := make(map[string]bool)
	for _, e := range userTeams {
		hash[e] = true
	}

	for _, e := range repoTeams {
		// If elements present in the hashmap then append intersection list.
		if hash[e] {
			teams = append(teams, e)
		}
	}

	if len(teams) > 0 {
		return teams[0]
	}

	if len(userTeams) > 0 {
		return userTeams[0]
	}

	if len(repoTeams) > 0 {
		return repoTeams[0]
	}

	return constants.NoTeam
}

func (tpr *TransformedPullRequest) fetchUserTeams(userID int) []string {
	var teams []string
	var err error

	query := "SELECT t2.slug" +
		" FROM extracted_team_users as t1 " +
		" INNER JOIN extracted_teams as t2 " +
		" ON t1.github_team_id = t2.id " +
		" WHERE t1.github_user_id= '" + strconv.Itoa(userID) + "'"

	if err = tpr.Db.Select(&teams, query); err != nil {
		fmt.Println(err)
	}

	return teams
}

func (tpr *TransformedPullRequest) fetchRepositoryTeams(repositoryName string) []string {
	var teams []string
	var err error

	query := "SELECT t3.slug" +
		" FROM extracted_team_repositories as t1 " +
		" INNER JOIN extracted_repositories as t2 " +
		" ON t1.github_repository_id = t2.id " +
		" INNER JOIN dpe_insights.extracted_teams as t3" +
		" ON t1.github_team_id = t3.id" +
		" WHERE t2.name= '" + repositoryName + "'"

	if err = tpr.Db.Select(&teams, query); err != nil {
		fmt.Println(err)
	}

	return teams
}

// FetchTeamsPullRequestCount Fetch teams - pull request count.
func (tpr *TransformedPullRequest) FetchTeamsPullRequestCount(fromDate, toDate string) (*[]entities.TransformedTeamPullRequestCount, error) {
	result := make(map[string]entities.TransformedTeamPullRequestCount)

	globalPrTotal, err := tpr.fetchTotalTeamPullRequestCount(&result, fromDate, toDate)
	if err != nil {
		return nil, err
	}

	globalPrClosed, err := tpr.fetchClosedTeamPullRequestCount(&result, fromDate, toDate)
	if err != nil {
		return nil, err
	}

	result[constants.NoTeam] = entities.TransformedTeamPullRequestCount{
		Total:    globalPrTotal,
		Closed:   globalPrClosed,
		Open:     globalPrTotal - globalPrClosed,
		TeamSlug: constants.NoTeam,
		Date:     fromDate,
	}

	var records []entities.TransformedTeamPullRequestCount

	for _, v := range result {
		records = append(records, v)
	}

	return &records, nil
}

// InsertOrUpdate Insert or update pull request data.
func (tpr *TransformedPullRequest) InsertOrUpdate(e *entities.TransformedPullRequest) error {
	err := tpr.Insert(e)
	if err != nil {
		if driverErr, ok := err.(*mysql.MySQLError); ok {
			if driverErr.Number == mysqldb.ERDupEntry {
				return tpr.Update(e)
			}
		}
		return err
	}
	return nil
}

func (tpr *TransformedPullRequest) fetchTotalTeamPullRequestCount(result *map[string]entities.TransformedTeamPullRequestCount, fromDate, toDate string) (int, error) {
	teamMap := *result

	query := "SELECT count(id) AS 'count', team_slug as 'team'" +
		" FROM " + transformedPRTbl +
		" WHERE pr_created_at >= '" + fromDate + "'" +
		" AND pr_created_at < '" + toDate + "'" +
		" GROUP BY team_slug"

	rows, err := tpr.Db.Query(query)
	if err != nil {
		return 0, err
	}

	globalPrTotal := 0

	for rows.Next() {
		var team string
		var count int
		err = rows.Scan(&count, &team)
		if err != nil {
			continue
		}

		teamMap[team] = entities.TransformedTeamPullRequestCount{
			Total:    count,
			Closed:   0,
			Open:     0,
			TeamSlug: team,
			Date:     fromDate,
		}
		globalPrTotal += count
	}

	return globalPrTotal, nil
}

func (tpr *TransformedPullRequest) fetchClosedTeamPullRequestCount(result *map[string]entities.TransformedTeamPullRequestCount, fromDate, toDate string) (int, error) {
	teamMap := *result

	query := "SELECT count(id) AS 'count', team_slug as 'team'" +
		" FROM " + transformedPRTbl +
		" WHERE pr_created_at >= '" + fromDate + "'" +
		" AND pr_created_at < '" + toDate + "'" +
		" AND status='closed' GROUP BY team_slug"

	rows, err := tpr.Db.Query(query)
	if err != nil {
		return 0, err
	}

	globalPrClosed := 0

	for rows.Next() {
		var team string
		var count int
		err = rows.Scan(&count, &team)

		if err != nil {
			continue
		}

		r, ok := teamMap[team]

		if !ok {
			teamMap[team] = entities.TransformedTeamPullRequestCount{
				Total:    count,
				Closed:   0,
				Open:     0,
				TeamSlug: team,
				Date:     fromDate,
			}
			continue
		}

		globalPrClosed += count

		r.Closed = count
		r.Open = r.Total - r.Closed

		teamMap[team] = r
	}

	return globalPrClosed, nil
}

// Save insert or update pull request data.
func (tpr *TransformedPullRequest) Save(e *entities.TransformedPullRequest) error {
	err := tpr.Insert(e)
	if err != nil {
		if driverErr, ok := err.(*mysql.MySQLError); ok {
			if driverErr.Number == mysqldb.ERDupEntry {
				return tpr.Update(e)
			}
		}
		return err
	}
	return nil
}
