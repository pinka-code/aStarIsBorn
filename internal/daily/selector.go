package daily

import (
	"a-star-is-born/internal/config"
	"a-star-is-born/internal/github"
	"a-star-is-born/internal/model"
	"a-star-is-born/internal/selector"
	"a-star-is-born/internal/storage"
	"fmt"
	"math/rand"
	"time"
)

func SelectDailyRepository(client *github.Client, storePath string, date time.Time) (*model.Repository, error) {
	repo, err := selectDailyRepo(client, date)
	if err != nil {
		return nil, err
	}

	enriched := enrichWithDeepWiki(*repo)

	if err := storeRepository(storePath, enriched); err != nil {
		return nil, err
	}

	return &enriched, nil
}

func selectDailyRepo(client *github.Client, date time.Time) (*model.Repository, error) {
	seed := selector.SeedFromDate(date)
	rng := rand.New(rand.NewSource(seed))

	criteria := config.ResolveCriteria(rng)
	query := selector.BuildQuery(criteria)

	items, total, err := client.SearchRepositories(query, 1)
	if err != nil {
		return &model.Repository{}, fmt.Errorf("search failed: %w", err)
	}

	if len(items) == 0 {
		return &model.Repository{}, fmt.Errorf("no repositories found")
	}

	if total > 1000 {
		total = 1000
	}

	index := selector.PickIndex(seed, total)
	pos := index % len(items)

	return &items[pos], nil
}

func storeRepository(storePath string, repo model.Repository) error {
	store, err := storage.LoadHistory(storePath)
	if err != nil {
		return err
	}

	if _, exists := store.GetRepository(repo.ID); !exists {
		return store.Append(&repo)
	}

	return nil
}

func enrichWithDeepWiki(repo model.Repository) model.Repository {
	if repo.DeepWikiURL == "" {
		repo.DeepWikiURL = fmt.Sprintf(
			"https://deepwiki.com/%s/%s",
			repo.Owner.Login,
			repo.Name,
		)
	}
	return repo
}
