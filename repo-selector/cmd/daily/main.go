package main

import (
	"fmt"
	"math/rand"
	"os"
	"repo-selector/internal/config"
	"repo-selector/internal/github"
	"repo-selector/internal/selector"
	"repo-selector/internal/storage"
	"time"
)

func main() {
	date := time.Now()
	seed := selector.SeedFromDate(date)
	rng := rand.New(rand.NewSource(seed))

	criteria := resolveCriteria(rng)
	query := selector.BuildQuery(criteria)

	client := github.NewClient(os.Getenv("GITHUB_TOKEN"))

	items, total, err := client.SearchRepositories(query, 1)
	if err != nil || len(items) == 0 {
		fmt.Println("No repositories found:", err)
		return
	}

	if total > 1000 {
		total = 1000
	}

	index := selector.PickIndex(seed, total)
	pos := index % len(items)

	repo := items[pos]

	idFloat, ok := repo["id"].(float64)
	if !ok {
		fmt.Println("Invalid repo ID")
		return
	}
	repoID := int64(idFloat)

	store, _ := storage.Load("data/history.json")

	if store.Exists(repoID) {
		fmt.Println("Repo already selected previously, skipping (rare case)")
		return
	}

	store.Add(date.Format("2006-01-02"), repoID)
	store.Save()

	fmt.Println("Date:", date.Format("2006-01-02"))
	fmt.Println("Criteria:", criteria.Pretty())
	fmt.Println("Selected repo:", repo["full_name"])
}

func resolveCriteria(rng *rand.Rand) config.Criteria {
	forcedLang := os.Getenv("LANGUAGE")
	if forcedLang != "" {
		return config.Criteria{
			MinContributors: true,
			Stars:           config.StarsMedium,
			Language:        forcedLang,
			Size:            config.SizeMedium,
		}
	}

	starsOptions := []config.StarsRange{
		config.StarsLow,
		config.StarsMedium,
		config.StarsHigh,
	}

	sizeOptions := []config.SizeRange{
		config.SizeSmall,
		config.SizeMedium,
		config.SizeLarge,
	}

	languages := []string{
		"go", "python", "typescript", "rust", "java", "c++",
	}

	return config.Criteria{
		MinContributors: rng.Intn(2) == 0,
		Stars:           starsOptions[rng.Intn(len(starsOptions))],
		Language:        languages[rng.Intn(len(languages))],
		Size:            sizeOptions[rng.Intn(len(sizeOptions))],
	}
}
