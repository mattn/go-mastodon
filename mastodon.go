// Package mastodon provides functions and structs for accessing the mastodon API.
package mastodon

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/tomnomnom/linkheader"
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

func (c *Client) doAPI(ctx context.Context, method string, uri string, params interface{}, res interface{}, pg *Pagination) error {
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
			if pg != nil {
				values = pg.setValues(values)
			}
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
		if method == http.MethodGet && pg != nil {
			u.RawQuery = pg.toValues().Encode()
		}
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

	if resp.StatusCode != http.StatusOK {
		return parseAPIError("bad request", resp)
	} else if res == nil {
		return nil
	} else if pg != nil {
		if lh := resp.Header.Get("Link"); lh != "" {
			pg2, err := newPagination(lh)
			if err != nil {
				return err
			}
			*pg = *pg2
		}
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
func (c *Client) Authenticate(ctx context.Context, username, password string) error {
	params := url.Values{
		"client_id":     {c.config.ClientID},
		"client_secret": {c.config.ClientSecret},
		"grant_type":    {"password"},
		"username":      {username},
		"password":      {password},
		"scope":         {"read write follow"},
	}

	return c.authenticate(ctx, params)
}

// AuthenticateToken logs in using a grant token returned by Application.AuthURI.
//
// redirectURI should be the same as Application.RedirectURI.
func (c *Client) AuthenticateToken(ctx context.Context, authCode, redirectURI string) error {
	params := url.Values{
		"client_id":     {c.config.ClientID},
		"client_secret": {c.config.ClientSecret},
		"grant_type":    {"authorization_code"},
		"code":          {authCode},
		"redirect_uri":  {redirectURI},
	}

	return c.authenticate(ctx, params)
}

func (c *Client) authenticate(ctx context.Context, params url.Values) error {
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

	var res struct {
		AccessToken string `json:"access_token"`
	}
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return err
	}
	c.config.AccessToken = res.AccessToken
	return nil
}

// Toot is struct to post status.
type Toot struct {
	Status      string `json:"status"`
	InReplyToID ID     `json:"in_reply_to_id"`
	MediaIDs    []ID   `json:"media_ids"`
	Sensitive   bool   `json:"sensitive"`
	SpoilerText string `json:"spoiler_text"`
	Visibility  string `json:"visibility"`
}

// Mention hold information for mention.
type Mention struct {
	URL      string `json:"url"`
	Username string `json:"username"`
	Acct     string `json:"acct"`
	ID       ID     `json:"id"`
}

// Tag hold information for tag.
type Tag struct {
	Name    string    `json:"name"`
	URL     string    `json:"url"`
	History []History `json:"history"`
}

// History hold information for history.
type History struct {
	Day      string `json:"day"`
	Uses     int64  `json:"uses"`
	Accounts int64  `json:"accounts"`
}

// Attachment hold information for attachment.
type Attachment struct {
	ID          ID             `json:"id"`
	Type        string         `json:"type"`
	URL         string         `json:"url"`
	RemoteURL   string         `json:"remote_url"`
	PreviewURL  string         `json:"preview_url"`
	TextURL     string         `json:"text_url"`
	Description string         `json:"description"`
	Meta        AttachmentMeta `json:"meta"`
}

// AttachmentMeta holds information for attachment metadata.
type AttachmentMeta struct {
	Original Size `json:"original"`
	Small    Size `json:"small"`
}

// Size holds information for attatchment size.
type Size struct {
	Width  int64   `json:"width"`
	Height int64   `json:"height"`
	Size   string  `json:"size"`
	Aspect float64 `json:"aspect"`
}

// Emoji hold information for CustomEmoji.
type Emoji struct {
	ShortCode       string `json:"shortcode"`
	StaticURL       string `json:"static_url"`
	URL             string `json:"url"`
	VisibleInPicker bool   `json:"visible_in_picker"`
}

// Results hold information for search result.
type Results struct {
	Accounts []*Account `json:"accounts"`
	Statuses []*Status  `json:"statuses"`
	Hashtags []string   `json:"hashtags"`
}

// Pagination is a struct for specifying the get range.
type Pagination struct {
	MaxID   ID
	SinceID ID
	Limit   int64
}

func newPagination(rawlink string) (*Pagination, error) {
	if rawlink == "" {
		return nil, errors.New("empty link header")
	}

	p := &Pagination{}
	for _, link := range linkheader.Parse(rawlink) {
		switch link.Rel {
		case "next":
			maxID, err := getPaginationID(link.URL, "max_id")
			if err != nil {
				return nil, err
			}
			p.MaxID = maxID
		case "prev":
			sinceID, err := getPaginationID(link.URL, "since_id")
			if err != nil {
				return nil, err
			}
			p.SinceID = sinceID
		}
	}

	return p, nil
}

func getPaginationID(rawurl, key string) (ID, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return "", err
	}

	val := u.Query().Get(key)
	if val == "" {
		return "", nil
	}

	id, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return "", err
	}

	return ID(fmt.Sprint(id)), nil
}

func (p *Pagination) toValues() url.Values {
	return p.setValues(url.Values{})
}

func (p *Pagination) setValues(params url.Values) url.Values {
	if p.MaxID != "" {
		params.Set("max_id", string(p.MaxID))
	} else if p.SinceID != "" {
		params.Set("since_id", string(p.SinceID))
	}
	if p.Limit > 0 {
		params.Set("limit", fmt.Sprint(p.Limit))
	}

	return params
}
