package main

import (
	"fmt"
	"os"
	"repo-selector/internal/daily"
	"repo-selector/internal/github"
	"time"
)

func main() {
	date := time.Now()

	client := github.NewClient(os.Getenv("GITHUB_TOKEN"))

	result, err := daily.SelectDailyRepository(
		client,
		"data/history.json",
		date,
	)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	err = daily.BuildDailySnapshot(
		client,
		result.Repo,
		"data/today.json",
		date,
	)
	if err != nil {
		fmt.Println("Snapshot error:", err)
		return
	}

	fmt.Println("Date:", result.Date)
	fmt.Println("Criteria:", result.Criteria)
	fmt.Println("Selected repo:", result.Repo.FullName)
}
