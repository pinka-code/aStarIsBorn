package model

import "time"

type Node struct {
	Name          string         `json:"name"`
	Path          string         `json:"path,omitempty"`
	Children      []*Node        `json:"children,omitempty"`
	Contributions map[string]int `json:"contributions,omitempty"`
	Dominant      string         `json:"dominant,omitempty"`
	Dominance     float64        `json:"dominance,omitempty"`
	LastCommitAt  *time.Time     `json:"last_commit_at,omitempty"`
}
