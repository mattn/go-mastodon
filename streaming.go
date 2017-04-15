package mastodon

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

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

func handleReader(ctx context.Context, q chan Event, r io.Reader) error {
	name := ""
	s := bufio.NewScanner(r)
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
				err := json.Unmarshal([]byte(token[1]), &status)
				if err == nil {
					q <- &UpdateEvent{&status}
				}
			case "notification":
			case "delete":
			}
		default:
		}
	}
	return ctx.Err()
}

func (c *Client) streaming(ctx context.Context, p string, tag string) (chan Event, error) {
	u, err := url.Parse(c.config.Server)
	if err != nil {
		return nil, err
	}
	u.Path = path.Join(u.Path, "/api/v1/streaming/"+p)

	params := url.Values{}
	params.Set("tag", tag)
	var resp *http.Response

	q := make(chan Event, 10)
	go func() {
		defer ctx.Done()

		for {
			var in io.Reader
			if tag != "" {
				in = strings.NewReader(params.Encode())
			}
			req, err := http.NewRequest(http.MethodGet, u.String(), in)
			if err == nil {
				req.Header.Set("Authorization", "Bearer "+c.config.AccessToken)
				resp, err = c.Do(req)
				if resp.StatusCode != 200 {
					err = fmt.Errorf("bad request: %v", resp.Status)
				}
			}
			if err == nil {
				err = handleReader(ctx, q, resp.Body)
				if err == nil {
					break
				}
			} else {
				q <- &ErrorEvent{err}
			}
			resp.Body.Close()
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

// StreamingPublic return channel to read events on public.
func (c *Client) StreamingPublic(ctx context.Context) (chan Event, error) {
	return c.streaming(ctx, "public", "")
}

// StreamingHome return channel to read events on home.
func (c *Client) StreamingHome(ctx context.Context) (chan Event, error) {
	return c.streaming(ctx, "home", "")
}

// StreamingHashtag return channel to read events on tagged timeline.
func (c *Client) StreamingHashtag(ctx context.Context, tag string) (chan Event, error) {
	return c.streaming(ctx, "hashtag", tag)
}
