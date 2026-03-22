package model

import "time"

type Snapshot struct {
	Date         string             `json:"date"`
	Repository   Repository         `json:"repository"`
	Contributors []ContributorStats `json:"contributors"`
}

type ContributorStats struct {
	Login         string    `json:"login"`
	AvatarURL     string    `json:"avatar_url"`
	HTMLURL       string    `json:"html_url"`
	Contributions int       `json:"contributions"`
	CommitCount   int       `json:"commit_count"`
	LastCommitAt  time.Time `json:"last_commit_at"`
}
