package mastodon

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
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

// AppConfig is a setting for registering applications.
type AppConfig struct {
	http.Client
	Server     string
	ClientName string

	// Where the user should be redirected after authorization (for no redirect, use urn:ietf:wg:oauth:2.0:oob)
	RedirectURIs string

	// This can be a space-separated list of the following items: "read", "write" and "follow".
	Scopes string

	// Optional.
	Website string
}

// Application is mastodon application.
type Application struct {
	ID           int64  `json:"id"`
	RedirectURI  string `json:"redirect_uri"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

// RegisterApp returns the mastodon application.
func RegisterApp(appConfig *AppConfig) (*Application, error) {
	params := url.Values{}
	params.Set("client_name", appConfig.ClientName)
	params.Set("redirect_uris", appConfig.RedirectURIs)
	params.Set("scopes", appConfig.Scopes)
	params.Set("website", appConfig.Website)

	url, err := url.Parse(appConfig.Server)
	if err != nil {
		return nil, err
	}
	url.Path = path.Join(url.Path, "/api/v1/apps")

	req, err := http.NewRequest("POST", url.String(), strings.NewReader(params.Encode()))
	if err != nil {
		return nil, err
	}
	resp, err := appConfig.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	app := &Application{}
	err = json.NewDecoder(resp.Body).Decode(app)
	if err != nil {
		return nil, err
	}

	return app, nil
}

type Visibility int64

type Toot struct {
	Status      string  `json:"status"`
	InReplyToID int64   `json:"in_reply_to_id"`
	MediaIDs    []int64 `json:"in_reply_to_id"`
	Sensitive   bool    `json:"sensitive"`
	SpoilerText string  `json:"spoiler_text"`
	Visibility  string  `json:"visibility"`
}

type Status struct {
	ID                 int64       `json:"id"`
	CreatedAt          time.Time   `json:"created_at"`
	InReplyToID        interface{} `json:"in_reply_to_id"`
	InReplyToAccountID interface{} `json:"in_reply_to_account_id"`
	Sensitive          bool        `json:"sensitive"`
	SpoilerText        string      `json:"spoiler_text"`
	Visibility         string      `json:"visibility"`
	Application        interface{} `json:"application"`
	Account            struct {
		ID             int64     `json:"id"`
		Username       string    `json:"username"`
		Acct           string    `json:"acct"`
		DisplayName    string    `json:"display_name"`
		Locked         bool      `json:"locked"`
		CreatedAt      time.Time `json:"created_at"`
		FollowersCount int64     `json:"followers_count"`
		FollowingCount int64     `json:"following_count"`
		StatusesCount  int64     `json:"statuses_count"`
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
	ReblogsCount     int64         `json:"reblogs_count"`
	FavouritesCount  int64         `json:"favourites_count"`
	Reblog           interface{}   `json:"reblog"`
	Favourited       interface{}   `json:"favourited"`
	Reblogged        interface{}   `json:"reblogged"`
}

func (c *client) GetTimelineHome() ([]*Status, error) {
	url, err := url.Parse(c.config.Server)
	if err != nil {
		return nil, err
	}
	url.Path = path.Join(url.Path, "/api/v1/timelines/home")

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.config.AccessToken)
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var statuses []*Status
	err = json.NewDecoder(resp.Body).Decode(&statuses)
	if err != nil {
		return nil, err
	}
	return statuses, nil
}

func (c *client) PostStatus(toot *Toot) (*Status, error) {
	params := url.Values{}
	params.Set("status", toot.Status)
	//params.Set("in_reply_to_id", fmt.Sprint(toot.InReplyToID))
	// TODO: media_ids, senstitive, spoiler_text, visibility
	//params.Set("visibility", "public")

	url, err := url.Parse(c.config.Server)
	if err != nil {
		return nil, err
	}
	url.Path = path.Join(url.Path, "/api/v1/statuses")

	req, err := http.NewRequest("POST", url.String(), strings.NewReader(params.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.config.AccessToken)
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var status Status
	err = json.NewDecoder(resp.Body).Decode(&status)
	if err != nil {
		return nil, err
	}
	return &status, nil
}

type UpdateEvent struct {
	Status *Status
}

func (e *UpdateEvent) event() {}

type NotificationEvent struct {
}

func (e *NotificationEvent) event() {}

type DeleteEvent struct {
	ID int64
}

func (e *DeleteEvent) event() {}

type Event interface {
	event()
}

func (c *client) StreamingPublic(ctx context.Context) (chan Event, error) {
	url, err := url.Parse(c.config.Server)
	if err != nil {
		return nil, err
	}
	url.Path = path.Join(url.Path, "/api/v1/streaming/public")

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.config.AccessToken)
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	q := make(chan Event)
	go func() {
		defer ctx.Done()
		name := ""
		s := bufio.NewScanner(resp.Body)
		for s.Scan() {
			line := s.Text()
			token := strings.SplitN(line, ":", 2)
			if len(token) != 2 {
				continue
			}
			switch strings.TrimSpace(token[0]) {
			case "event":
				name = strings.TrimSpace(token[1])
			case "data":
				switch name {
				case "update":
					var status Status
					json.Unmarshal([]byte(token[1]), &status)
					q <- &UpdateEvent{&status}
				case "notification":
				case "delete":
				}
			}
		}
		fmt.Println(s.Err())
	}()
	go func() {
		<-ctx.Done()
		resp.Body.Close()
	}()
	return q, nil
}
