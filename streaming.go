package mastodon

import (
	"bufio"
	"context"
	"encoding/json"
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

// StreamingPublic return channel to read events.
func (c *Client) StreamingPublic(ctx context.Context) (chan Event, error) {
	u, err := url.Parse(c.config.Server)
	if err != nil {
		return nil, err
	}
	u.Path = path.Join(u.Path, "/api/v1/streaming/public")

	var resp *http.Response

	q := make(chan Event, 10)
	go func() {
		defer ctx.Done()

		for {
			req, err := http.NewRequest(http.MethodGet, u.String(), nil)
			if err == nil {
				req.Header.Set("Authorization", "Bearer "+c.config.AccessToken)
				resp, err = c.Do(req)
			}
			if err == nil {
				err = handleReader(ctx, q, resp.Body)
				resp.Body.Close()
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
