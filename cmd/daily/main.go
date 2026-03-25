package main

import (
	"a-star-is-born/internal/daily"
	"a-star-is-born/internal/github"
	"fmt"
	"net/smtp"
	"os"
	"strings"
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

	description := repo.Description
	if description == "" {
		description = "(none)"
	}

	deepWiki := ""
	if repo.DeepWikiURL != "" {
		deepWiki = fmt.Sprintf("🧠 DeepWiki: %s\n", repo.DeepWikiURL)
	}

	msg := fmt.Sprintf(
		"Subject: A Star is born Newsletter\n\n"+
			"📅 Date: %s\n"+
			"📦 Repo: %s\n"+
			"🔗 URL: %s\n"+
			"⭐ Stars: %d\n"+
			"💻 Language: %s\n"+
			"%s"+
			"📝 Description: %s\n",
		date.Format("2006-01-02"),
		repo.FullName,
		repo.HTMLURL,
		repo.StargazersCount,
		repo.Language,
		deepWiki,
		description,
	)

	// Sending email
	from := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASS")
	subscribers := os.Getenv("SUBSCRIBERS")
	toList := strings.Split(subscribers, ",")

	for _, to := range toList {
		err := smtp.SendMail(
			"smtp.gmail.com:587",
			smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
			from,
			[]string{to},
			[]byte(msg),
		)

		if err != nil {
			fmt.Println("❌ Error", err)
			continue
		}
	}
	fmt.Println("✅ Email sent")
}
