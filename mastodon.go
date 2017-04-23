package mastodon

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/gorilla/websocket"
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
	websocket.Dialer
	config   *Config
	interval time.Duration
}

type page struct {
	next string
}

func linkHeader(h http.Header, rel string) []string {
	var links []string
	for _, v := range h["Link"] {
		var p string
		for len(v) > 0 {
			i := strings.Index(v, ";")
			if i < 0 {
				break
			}
			e := i
			i++
			for i < len(v) {
				if v[i] != ' ' {
					break
				}
				i++
			}
			p = strings.TrimSpace(v[i:])
			if !strings.HasPrefix(p, "rel=") {
				break
			}
			i += 4
			pos := strings.Index(p[4:], `,`)
			if pos > 0 {
				p = p[4 : 4+pos]
				i += pos
			} else {
				p = p[4:]
				i = len(v) - 1
			}
			if k := strings.Trim(p, `"`); k == rel {
				links = append(links, strings.Trim(v[:e], "<>"))
			}
			i++
			for i < len(v) {
				if v[i] != ' ' {
					break
				}
				i++
			}
			v = v[i:]
		}
	}
	return links
}

func (c *Client) doAPI(ctx context.Context, method string, uri string, params interface{}, res interface{}, next *bool) error {
	u, err := url.Parse(c.config.Server)
	if err != nil {
		return err
	}
	u.Path = path.Join(u.Path, uri)

	var req *http.Request
	ct := "application/x-www-form-urlencoded"
	if values, ok := params.(url.Values); ok {
		var body io.Reader
		if method == http.MethodGet {
			u.RawQuery = values.Encode()
		} else {
			body = strings.NewReader(values.Encode())
		}
		req, err = http.NewRequest(method, u.String(), body)
		if err != nil {
			return err
		}
	} else if file, ok := params.(string); ok {
		f, err := os.Open(file)
		if err != nil {
			return err
		}
		defer f.Close()

		var buf bytes.Buffer
		mw := multipart.NewWriter(&buf)
		part, err := mw.CreateFormFile("file", filepath.Base(file))
		if err != nil {
			return err
		}
		_, err = io.Copy(part, f)
		if err != nil {
			return err
		}
		err = mw.Close()
		if err != nil {
			return err
		}
		req, err = http.NewRequest(method, u.String(), &buf)
		if err != nil {
			return err
		}
		ct = mw.FormDataContentType()
	} else {
		req, err = http.NewRequest(method, u.String(), nil)
		if err != nil {
			return err
		}
	}
	req = req.WithContext(ctx)
	req.Header.Set("Authorization", "Bearer "+c.config.AccessToken)
	if params != nil {
		req.Header.Set("Content-Type", ct)
	}

	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if next != nil && params != nil {
		nl := linkHeader(resp.Header, "next")
		*next = false
		if len(nl) > 0 {
			u, err = url.Parse(nl[0])
			if err == nil {
				for k, v := range u.Query() {
					params.(url.Values)[k] = v
				}
			}
			*next = true
		}
	}
	if resp.StatusCode != http.StatusOK {
		return parseAPIError("bad request", resp)
	} else if res == nil {
		return nil
	}
	return json.NewDecoder(resp.Body).Decode(&res)
}

// NewClient return new mastodon API client.
func NewClient(config *Config) *Client {
	return &Client{
		Client:   *http.DefaultClient,
		config:   config,
		interval: 10 * time.Second,
	}
}

// Authenticate get access-token to the API.
func (c *Client) Authenticate(ctx context.Context, username, password string) error {
	params := url.Values{}
	params.Set("client_id", c.config.ClientID)
	params.Set("client_secret", c.config.ClientSecret)
	params.Set("grant_type", "password")
	params.Set("username", username)
	params.Set("password", password)
	params.Set("scope", "read write follow")

	u, err := url.Parse(c.config.Server)
	if err != nil {
		return err
	}
	u.Path = path.Join(u.Path, "/oauth/token")

	req, err := http.NewRequest(http.MethodPost, u.String(), strings.NewReader(params.Encode()))
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return parseAPIError("bad authorization", resp)
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

// Mention hold information for mention.
type Mention struct {
	URL      string `json:"url"`
	Username string `json:"username"`
	Acct     string `json:"acct"`
	ID       int64  `json:"id"`
}

// Tag hold information for tag.
type Tag struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// Attachment hold information for attachment.
type Attachment struct {
	ID         int64  `json:"id"`
	Type       string `json:"type"`
	URL        string `json:"url"`
	RemoteURL  string `json:"remote_url"`
	PreviewURL string `json:"preview_url"`
	TextURL    string `json:"text_url"`
}

// Results hold information for search result.
type Results struct {
	Accounts []*Account `json:"accounts"`
	Statuses []*Status  `json:"statuses"`
	Hashtags []string   `json:"hashtags"`
}
