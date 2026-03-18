package selector

import (
	"fmt"
	"repo-selector/internal/config"
)

func BuildQuery(c config.Criteria) string {
	query := "is:public archived:false"

	// stars
	switch c.Stars {
	case config.StarsLow:
		query += " stars:1..50"
	case config.StarsMedium:
		query += " stars:50..200"
	case config.StarsHigh:
		query += " stars:>200"
	}

	// language
	if c.Language != "" {
		query += fmt.Sprintf(" language:%s", c.Language)
	}

	// activity
	query += " pushed:>2023-01-01"

	// size (approx via GitHub size field)
	switch c.Size {
	case config.SizeSmall:
		query += " size:<1000"
	case config.SizeMedium:
		query += " size:1000..10000"
	case config.SizeLarge:
		query += " size:>10000"
	}

	// proxy contributors
	if c.MinContributors {
		query += " forks:>10"
	}

	return query
}