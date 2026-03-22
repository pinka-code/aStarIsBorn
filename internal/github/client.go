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

func (c *Client) GetCommitsWithFilesSince(owner, repo string, since time.Time) ([]model.CommitWithFiles, error) {
	var all []model.CommitWithFiles

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

		for _, commit := range commits {
			detailURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits/%s", owner, repo, commit.SHA)

			var detail struct {
				SHA   string `json:"sha"`
				Files []struct {
					Filename string `json:"filename"`
				} `json:"files"`
				Commit struct {
					Author struct {
						Date time.Time `json:"date"`
					} `json:"author"`
				} `json:"commit"`
				Author *model.Contributor `json:"author"`
			}

			if err := doJSONInto(detailURL, &detail); err != nil {
				return nil, err
			}

			var files []model.File
			for _, f := range detail.Files {
				files = append(files, model.File{Filename: f.Filename})
			}

			all = append(all, model.CommitWithFiles{
				SHA:    detail.SHA,
				Author: detail.Author,
				Date:   detail.Commit.Author.Date,
				Files:  files,
			})
		}

		if len(commits) < 100 {
			break
		}
	}

	return all, nil
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

func doJSONInto(url string, v interface{}) error {
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(v)
}
