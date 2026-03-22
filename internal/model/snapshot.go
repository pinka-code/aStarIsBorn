package model

type Snapshot struct {
	Date         string            `json:"date"`
	Repository   Repository        `json:"repository"`
	Contributors []Contributor     `json:"contributors"`
	Commits      []CommitWithFiles `json:"commits"`
}
