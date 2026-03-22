package storage

import (
	"encoding/json"
	"os"
)

type Store struct {
	Path string
	Data map[string]int64
}

func Load(path string) (*Store, error) {
	file, err := os.Open(path)
	if err != nil {
		return &Store{Path: path, Data: map[string]int64{}}, nil
	}
	defer file.Close()

	data := map[string]int64{}
	json.NewDecoder(file).Decode(&data)

	return &Store{Path: path, Data: data}, nil
}

func (s *Store) Save() error {
	file, _ := os.Create(s.Path)
	defer file.Close()
	return json.NewEncoder(file).Encode(s.Data)
}

func (s *Store) Exists(repoID int64) bool {
	for _, v := range s.Data {
		if v == repoID {
			return true
		}
	}
	return false
}

func (s *Store) Add(date string, repoID int64) {
	s.Data[date] = repoID
}