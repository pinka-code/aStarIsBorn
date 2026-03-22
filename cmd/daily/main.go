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

	repo, err := daily.SelectDailyRepository(client, "data/history.jsonl", date)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("📅 Date:", date.Format("2006-01-02"))
	fmt.Println("📦 Repo:", repo.FullName)
	fmt.Println("🔗 URL:", repo.HTMLURL)
	fmt.Println("⭐ Stars:", repo.StargazersCount)
	fmt.Println("💻 Language:", repo.Language)

	if repo.DeepWikiURL != "" {
		fmt.Println("🧠 DeepWiki:", repo.DeepWikiURL)
	}

	if repo.Description != "" {
		fmt.Println("📝 Description:", repo.Description)
	} else {
		fmt.Println("📝 Description: (none)")
	}
}
