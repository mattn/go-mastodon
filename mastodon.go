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

// Account hold information for mastodon account.
type Account struct {
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

// GetAccount return Account.
func (c *Client) GetAccount(id int) (*Account, error) {
	var account Account
	err := c.doAPI("GET", fmt.Sprintf("/api/v1/accounts/%d", id), nil, &account)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

// GetAccountFollowers return followers list.
func (c *Client) GetAccountFollowers(id int64) ([]*Account, error) {
	var accounts []*Account
	err := c.doAPI("GET", fmt.Sprintf("/api/v1/accounts/%d/followers", id), nil, &accounts)
	if err != nil {
		return nil, err
	}
	return accounts, nil
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

// UpdateEvent is struct for passing status event to app.
type UpdateEvent struct{ Status *Status }

func (e *UpdateEvent) event() {}

// NotificationEvent is struct for passing notification event to app.
type NotificationEvent struct{}

func (e *NotificationEvent) event() {}

// DeleteEvent is struct for passing deletion event to app.
type DeleteEvent struct{ ID int64 }

func (e *DeleteEvent) event() {}

// ErrorEvent is struct for passing errors to app.
type ErrorEvent struct{ err error }

func (e *ErrorEvent) event()        {}
func (e *ErrorEvent) Error() string { return e.err.Error() }

// Event is interface passing events to app.
type Event interface {
	event()
}

// StreamingPublic return channel to read events.
func (c *Client) StreamingPublic(ctx context.Context) (chan Event, error) {
	url, err := url.Parse(c.config.Server)
	if err != nil {
		return nil, err
	}
	url.Path = path.Join(url.Path, "/api/v1/streaming/public")

	var resp *http.Response

	q := make(chan Event, 10)
	go func() {
		defer ctx.Done()

		for {
			req, err := http.NewRequest("GET", url.String(), nil)
			if err == nil {
				req.Header.Set("Authorization", "Bearer "+c.config.AccessToken)
				resp, err = c.Do(req)
			}
			if err == nil {
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
							err = json.Unmarshal([]byte(token[1]), &status)
							if err == nil {
								q <- &UpdateEvent{&status}
							}
						case "notification":
						case "delete":
						}
					default:
					}
				}
				resp.Body.Close()
				err = ctx.Err()
				if err == nil {
					break
				}
			} else {
				q <- &ErrorEvent{err}
			}
			time.Sleep(3 * time.Second)
		}
	}()
	go func() {
		<-ctx.Done()
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()
	return q, nil
}

// GetAccount return Account.
func (c *Client) Follow(uri string) (*Account, error) {
	params := url.Values{}
	params.Set("uri", uri)

	var account Account
	err := c.doAPI("POST", "/api/v1/follows", params, &account)
	if err != nil {
		return nil, err
	}
	return &account, nil
}
