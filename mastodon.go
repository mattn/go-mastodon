package mastodon

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

// Config is a setting for access mastodon APIs.
type Config struct {
	Server       string
	ClientID     string
	ClientSecret string
	AccessToken  string
}

// Client is a API client for mastodon.
type Client struct {
	http.Client
	config *Config
}

func (c *Client) doAPI(method string, uri string, params url.Values, res interface{}) error {
	url, err := url.Parse(c.config.Server)
	if err != nil {
		return err
	}
	url.Path = path.Join(url.Path, uri)

	var resp *http.Response
	req, err := http.NewRequest(method, url.String(), strings.NewReader(params.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.config.AccessToken)
	resp, err = c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if res == nil {
		return nil
	}

	if method == "GET" && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad request: %v", resp.Status)
	}

	return json.NewDecoder(resp.Body).Decode(&res)
}

// NewClient return new mastodon API client.
func NewClient(config *Config) *Client {
	return &Client{
		Client: *http.DefaultClient,
		config: config,
	}
}

// Authenticate get access-token to the API.
func (c *Client) Authenticate(username, password string) error {
	params := url.Values{}
	params.Set("client_id", c.config.ClientID)
	params.Set("client_secret", c.config.ClientSecret)
	params.Set("grant_type", "password")
	params.Set("username", username)
	params.Set("password", password)
	params.Set("scope", "read write follow")

	url, err := url.Parse(c.config.Server)
	if err != nil {
		return err
	}
	url.Path = path.Join(url.Path, "/oauth/token")

	req, err := http.NewRequest("POST", url.String(), strings.NewReader(params.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad authorization: %v", resp.Status)
	}

	res := struct {
		AccessToken string `json:"access_token"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return err
	}
	c.config.AccessToken = res.AccessToken
	return nil
}

// Toot is struct to post status.
type Toot struct {
	Status      string  `json:"status"`
	InReplyToID int64   `json:"in_reply_to_id"`
	MediaIDs    []int64 `json:"media_ids"`
	Sensitive   bool    `json:"sensitive"`
	SpoilerText string  `json:"spoiler_text"`
	Visibility  string  `json:"visibility"`
}

// Status is struct to hold status.
type Status struct {
	ID                 int64         `json:"id"`
	CreatedAt          time.Time     `json:"created_at"`
	InReplyToID        interface{}   `json:"in_reply_to_id"`
	InReplyToAccountID interface{}   `json:"in_reply_to_account_id"`
	Sensitive          bool          `json:"sensitive"`
	SpoilerText        string        `json:"spoiler_text"`
	Visibility         string        `json:"visibility"`
	Application        interface{}   `json:"application"`
	Account            Account       `json:"account"`
	MediaAttachments   []interface{} `json:"media_attachments"`
	Mentions           []interface{} `json:"mentions"`
	Tags               []interface{} `json:"tags"`
	URI                string        `json:"uri"`
	Content            string        `json:"content"`
	URL                string        `json:"url"`
	ReblogsCount       int64         `json:"reblogs_count"`
	FavouritesCount    int64         `json:"favourites_count"`
	Reblog             interface{}   `json:"reblog"`
	Favourited         interface{}   `json:"favourited"`
	Reblogged          interface{}   `json:"reblogged"`
}

// GetTimelineHome return statuses from home timeline.
func (c *Client) GetTimelineHome() ([]*Status, error) {
	var statuses []*Status
	err := c.doAPI("GET", "/api/v1/timelines/home", nil, &statuses)
	if err != nil {
		return nil, err
	}
	return statuses, nil
}

// PostStatus post the toot.
func (c *Client) PostStatus(toot *Toot) (*Status, error) {
	params := url.Values{}
	params.Set("status", toot.Status)
	if toot.InReplyToID > 0 {
		params.Set("in_reply_to_id", fmt.Sprint(toot.InReplyToID))
	}
	// TODO: media_ids, senstitive, spoiler_text, visibility
	//params.Set("visibility", "public")

	var status Status
	err := c.doAPI("POST", "/api/v1/statuses", params, &status)
	if err != nil {
		return nil, err
	}
	return &status, nil
}
