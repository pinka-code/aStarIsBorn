package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Client struct {
	Token string
	http  *http.Client
}

func NewClient(token string) *Client {
	return &Client{
		Token: token,
		http: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout: 5 * time.Second,
				}).DialContext,
				TLSHandshakeTimeout: 5 * time.Second,
				MaxIdleConns:        10,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     30 * time.Second,
			},
		},
	}
}

func (c *Client) SearchRepositories(query string, page int) ([]map[string]interface{}, int, error) {
	params := url.Values{}
	params.Set("q", query)
	params.Set("per_page", "30")
	params.Set("page", strconv.Itoa(page))

	url := "https://api.github.com/search/repositories?" + params.Encode()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, 0, err
	}

	req.Header.Set("User-Agent", "my-github-client")
	req.Header.Set("Accept", "application/vnd.github+json")

	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, 0, fmt.Errorf("github API error: %s - %s", resp.Status, string(body))
	}

	var result struct {
		TotalCount int                      `json:"total_count"`
		Items      []map[string]interface{} `json:"items"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, 0, err
	}

	return result.Items, result.TotalCount, nil
}
