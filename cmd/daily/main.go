package main

import (
	"a-star-is-born/internal/daily"
	"a-star-is-born/internal/github"
	"fmt"
	"os"
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
