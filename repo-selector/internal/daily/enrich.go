package daily

import (
	"fmt"
	"sort"
	"time"

	"repo-selector/internal/github"
	"repo-selector/internal/model"
	"repo-selector/internal/storage"
)

func BuildDailySnapshot(
	client *github.Client,
	repo model.Repository,
	outputPath string,
	date time.Time,
) error {
	owner := repo.Owner.Login
	name := repo.Name

	contributors, err := client.GetAllContributors(owner, name)
	if err != nil {
		return fmt.Errorf("contributors fetch failed: %w", err)
	}

	since := date.AddDate(0, 0, -30)

	commits, err := client.GetCommitsSince(owner, name, since)
	if err != nil {
		return fmt.Errorf("commits fetch failed: %w", err)
	}

	stats := buildContributorStats(contributors, commits)

	snapshot := model.Snapshot{
		Date:         date.Format("2006-01-02"),
		Repository:   repo,
		Contributors: stats,
	}

	if err := storage.SaveJSON(outputPath, snapshot); err != nil {
		return fmt.Errorf("save snapshot failed: %w", err)
	}

	return nil
}

func buildContributorStats(
	contributors []model.Contributor,
	commits []model.Commit,
) []model.ContributorStats {
	statsMap := make(map[string]*model.ContributorStats)

	for _, c := range contributors {
		statsMap[c.Login] = &model.ContributorStats{
			Login:         c.Login,
			AvatarURL:     c.AvatarURL,
			HTMLURL:       c.HTMLURL,
			Contributions: c.Contributions,
			CommitCount:   0,
		}
	}

	for _, commit := range commits {
		if commit.Author == nil {
			continue
		}

		login := commit.Author.Login

		stat, ok := statsMap[login]
		if !ok {
			stat = &model.ContributorStats{
				Login: login,
			}
			statsMap[login] = stat
		}

		stat.CommitCount++

		commitDate := commit.Commit.Author.Date
		if commitDate.After(stat.LastCommitAt) {
			stat.LastCommitAt = commitDate
		}
	}

	result := make([]model.ContributorStats, 0, len(statsMap))
	for _, v := range statsMap {
		result = append(result, *v)
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].CommitCount == result[j].CommitCount {
			return result[i].Contributions > result[j].Contributions
		}
		return result[i].CommitCount > result[j].CommitCount
	})

	return result
}
