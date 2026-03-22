package storage

import (
	"a-star-is-born/internal/model"
	"bufio"
	"encoding/json"
	"os"
)

type HistoryStore struct {
	Path string
	Data []*model.Repository
}

func LoadHistory(path string) (*HistoryStore, error) {
	file, err := os.Open(path)
	if err != nil {
		return &HistoryStore{
			Path: path,
			Data: []*model.Repository{},
		}, nil
	}
	defer file.Close()

	var repos []*model.Repository

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var repo model.Repository
		if err := json.Unmarshal(scanner.Bytes(), &repo); err != nil {
			continue // skip ligne invalide
		}
		repos = append(repos, &repo)
	}

	return &HistoryStore{
		Path: path,
		Data: repos,
	}, nil
}

func (s *HistoryStore) Append(r *model.Repository) error {
	file, err := os.OpenFile(s.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(r)
}

func (s *HistoryStore) GetRepository(repoID int64) (*model.Repository, bool) {
	for _, r := range s.Data {
		if r.ID == repoID {
			return r, true
		}
	}
	return nil, false
}

func (s *HistoryStore) UpdateRepository(r *model.Repository) {
	for i, existing := range s.Data {
		if existing.ID == r.ID {
			s.Data[i] = r
			return
		}
	}
	s.Data = append(s.Data, r)
}

func NewRepository(r model.Repository) *model.Repository {
	return &model.Repository{
		ID:              r.ID,
		Name:            r.Name,
		FullName:        r.FullName,
		Owner:           r.Owner,
		HTMLURL:         r.HTMLURL,
		Description:     r.Description,
		Language:        r.Language,
		StargazersCount: r.StargazersCount,
		ForksCount:      r.ForksCount,
		OpenIssuesCount: r.OpenIssuesCount,
		WatchersCount:   r.WatchersCount,
		Size:            r.Size,
		DefaultBranch:   r.DefaultBranch,
		CreatedAt:       r.CreatedAt,
		UpdatedAt:       r.UpdatedAt,
		PushedAt:        r.PushedAt,
	}
}
