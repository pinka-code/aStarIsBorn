package model

type Repository struct {
	ID          int64
	Name        string
	FullName    string
	Description string
	Stars       int
	Language    string
	Forks       int
	URL         string
}