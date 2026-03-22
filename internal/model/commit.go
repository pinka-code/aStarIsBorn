package model

import "time"

type Commit struct {
	SHA    string `json:"sha"`
	Commit struct {
		Author struct {
			Name  string    `json:"name"`
			Email string    `json:"email"`
			Date  time.Time `json:"date"`
		} `json:"author"`
		Message string `json:"message"`
	} `json:"commit"`
	Author *Contributor `json:"author"`
}

type CommitWithFiles struct {
	Author *Contributor `json:"author"`
	SHA    string       `json:"sha"`
	Date   time.Time    `json:"date"`
	Files  []File       `json:"files"`
}

type File struct {
	Filename string `json:"filename"`
}
