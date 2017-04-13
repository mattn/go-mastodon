package mastodon

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type Config struct {
	Server       string
	ClientID     string
	ClientSecret string
	AccessToken  string
}

type client struct {
	http.Client
	config *Config
}

func NewClient(config *Config) *client {
	return &client{
		Client: *http.DefaultClient,
		config: config,
	}
}

func (c *client) Authenticate(username, password string) error {
	params := url.Values{}
	params.Set("client_id", c.config.ClientID)
	params.Set("client_secret", c.config.ClientSecret)
	params.Set("grant_type", "password")
	params.Set("username", username)
	params.Set("password", password)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", c.config.Server, "/oauth/token"), strings.NewReader(params.Encode()))
	if err != nil {
		return err
	}
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

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

type Timeline struct {
	ID                 int         `json:"id"`
	CreatedAt          time.Time   `json:"created_at"`
	InReplyToID        interface{} `json:"in_reply_to_id"`
	InReplyToAccountID interface{} `json:"in_reply_to_account_id"`
	Sensitive          bool        `json:"sensitive"`
	SpoilerText        string      `json:"spoiler_text"`
	Visibility         string      `json:"visibility"`
	Application        interface{} `json:"application"`
	Account            struct {
		ID             int       `json:"id"`
		Username       string    `json:"username"`
		Acct           string    `json:"acct"`
		DisplayName    string    `json:"display_name"`
		Locked         bool      `json:"locked"`
		CreatedAt      time.Time `json:"created_at"`
		FollowersCount int       `json:"followers_count"`
		FollowingCount int       `json:"following_count"`
		StatusesCount  int       `json:"statuses_count"`
		Note           string    `json:"note"`
		URL            string    `json:"url"`
		Avatar         string    `json:"avatar"`
		AvatarStatic   string    `json:"avatar_static"`
		Header         string    `json:"header"`
		HeaderStatic   string    `json:"header_static"`
	} `json:"account"`
	MediaAttachments []interface{} `json:"media_attachments"`
	Mentions         []interface{} `json:"mentions"`
	Tags             []interface{} `json:"tags"`
	URI              string        `json:"uri"`
	Content          string        `json:"content"`
	URL              string        `json:"url"`
	ReblogsCount     int           `json:"reblogs_count"`
	FavouritesCount  int           `json:"favourites_count"`
	Reblog           interface{}   `json:"reblog"`
	Favourited       interface{}   `json:"favourited"`
	Reblogged        interface{}   `json:"reblogged"`
}

func (c *client) GetTimeline(path string) ([]Timeline, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", c.config.Server, "/api/v1/timelines/home"), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.config.AccessToken)
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var timeline []Timeline
	err = json.NewDecoder(io.TeeReader(resp.Body, os.Stdout)).Decode(&timeline)
	if err != nil {
		return nil, err
	}
	return timeline, nil
}
