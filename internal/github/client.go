package github

import (
	"a-star-is-born/internal/model"
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

func (c *Client) newRequest(method, rawURL string, queryParams url.Values) (*http.Request, error) {
	if queryParams != nil {
		rawURL = rawURL + "?" + queryParams.Encode()
	}

	req, err := http.NewRequest(method, rawURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "my-github-client")
	req.Header.Set("Accept", "application/vnd.github+json")

	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	return req, nil
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github API error: %s - %s", resp.Status, string(body))
	}

	return body, nil
}

func parseJSON[T any](data []byte, target *T) error {
	return json.Unmarshal(data, target)
}

func doJSON[T any](c *Client, method, rawURL string, params url.Values) (T, error) {
	var result T

	req, err := c.newRequest(method, rawURL, params)
	if err != nil {
		return result, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return result, err
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return result, err
	}

	return result, nil
}

func (c *Client) SearchRepositories(query string, page int) ([]model.Repository, int, error) {
	params := url.Values{}
	params.Set("q", query)
	params.Set("per_page", "30")
	params.Set("page", strconv.Itoa(page))

	type response struct {
		TotalCount int                `json:"total_count"`
		Items      []model.Repository `json:"items"`
	}

	res, err := doJSON[response](
		c,
		"GET",
		"https://api.github.com/search/repositories",
		params,
	)
	if err != nil {
		return nil, 0, err
	}

	return res.Items, res.TotalCount, nil
}

func (c *Client) GetAllContributors(owner, repo string) ([]model.Contributor, error) {
	var all []model.Contributor

	for page := 1; ; page++ {
		items, err := c.GetContributors(owner, repo, page)
		if err != nil {
			return nil, err
		}

		if len(items) == 0 {
			break
		}

		all = append(all, items...)
	}

	return all, nil
}

func (c *Client) GetContributors(owner, repo string, page int) ([]model.Contributor, error) {
	params := url.Values{}
	params.Set("per_page", "100")
	params.Set("page", strconv.Itoa(page))

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/contributors", owner, repo)

	return doJSON[[]model.Contributor](c, "GET", url, params)
}

func (c *Client) GetCommitsSince(owner, repo string, since time.Time) ([]model.Commit, error) {
	var all []model.Commit

	for page := 1; ; page++ {
		params := url.Values{}
		params.Set("per_page", "100")
		params.Set("page", strconv.Itoa(page))
		params.Set("since", since.Format(time.RFC3339))

		url := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits", owner, repo)

		commits, err := doJSON[[]model.Commit](c, "GET", url, params)
		if err != nil {
			return nil, err
		}

		if len(commits) == 0 {
			break
		}

		all = append(all, commits...)

		if len(commits) < 100 {
			break
		}
	}

	return all, nil
}
