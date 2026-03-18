package github

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	Token string
}

func (c *Client) SearchRepositories(query string, page int) ([]map[string]interface{}, int, error) {
	url := fmt.Sprintf("https://api.github.com/search/repositories?q=%s&per_page=30&page=%d", query, page)

	req, _ := http.NewRequest("GET", url, nil)
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	var result struct {
		TotalCount int                      `json:"total_count"`
		Items      []map[string]interface{} `json:"items"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	return result.Items, result.TotalCount, err
}