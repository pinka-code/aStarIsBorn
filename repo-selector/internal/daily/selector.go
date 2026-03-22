package daily

import (
	"fmt"
	"math/rand"
	"time"

	"repo-selector/internal/config"
	"repo-selector/internal/github"
	"repo-selector/internal/model"
	"repo-selector/internal/selector"
	"repo-selector/internal/storage"
)

type Result struct {
	Date     string
	Criteria string
	Repo     model.Repository
}

func SelectDailyRepository(client *github.Client, storePath string, date time.Time) (*Result, error) {
	seed := selector.SeedFromDate(date)
	rng := rand.New(rand.NewSource(seed))

	criteria := config.ResolveCriteria(rng)
	query := selector.BuildQuery(criteria)

	items, total, err := client.SearchRepositories(query, 1)
	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}

	if len(items) == 0 {
		return nil, fmt.Errorf("no repositories found")
	}

	if total > 1000 {
		total = 1000
	}

	index := selector.PickIndex(seed, total)
	pos := index % len(items)

	repo := items[pos]

	store, err := storage.Load(storePath)
	if err != nil {
		return nil, err
	}

	if !store.Exists(repo.ID) {
		store.Add(date.Format("2006-01-02"), repo.ID)

		if err := store.Save(); err != nil {
			return nil, err
		}
	}

	return &Result{
		Date:     date.Format("2006-01-02"),
		Criteria: criteria.Pretty(),
		Repo:     repo,
	}, nil
}
