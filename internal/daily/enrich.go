package daily

import (
	"a-star-is-born/internal/github"
	"a-star-is-born/internal/model"
	"a-star-is-born/internal/storage"
	"fmt"
	"time"
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

	commits, err := client.GetCommitsWithFilesSince(owner, name, since)
	if err != nil {
		return fmt.Errorf("commits fetch failed: %w", err)
	}

	snapshot := model.Snapshot{
		Date:         date.Format("2006-01-02"),
		Repository:   repo,
		Contributors: contributors,
		Commits:      commits,
	}

	if err := storage.SaveJSON(outputPath, snapshot); err != nil {
		return fmt.Errorf("save snapshot failed: %w", err)
	}

	return nil
}
