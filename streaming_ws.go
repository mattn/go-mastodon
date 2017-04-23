package mastodon

import (
	"context"
	"encoding/json"
	"net/url"
	"path"

	"github.com/gorilla/websocket"
)

// Stream is a struct of data that flows in streaming.
type Stream struct {
	Event   string      `json:"event"`
	Payload interface{} `json:"payload"`
}

// StreamingWSPublic return channel to read events on public using WebSocket.
func (c *Client) StreamingWSPublic(ctx context.Context) (chan Event, error) {
	return c.streamingWS(ctx, "public", "")
}

// StreamingWSPublicLocal return channel to read events on public local using WebSocket.
func (c *Client) StreamingWSPublicLocal(ctx context.Context) (chan Event, error) {
	return c.streamingWS(ctx, "public:local", "")
}

// StreamingWSUser return channel to read events on home using WebSocket.
func (c *Client) StreamingWSUser(ctx context.Context, user string) (chan Event, error) {
	return c.streamingWS(ctx, "user", "")
}

// StreamingWSHashtag return channel to read events on tagged timeline using WebSocket.
func (c *Client) StreamingWSHashtag(ctx context.Context, tag string) (chan Event, error) {
	return c.streamingWS(ctx, "hashtag", tag)
}

// StreamingWSHashtagLocal return channel to read events on tagged local timeline using WebSocket.
func (c *Client) StreamingWSHashtagLocal(ctx context.Context, tag string) (chan Event, error) {
	return c.streamingWS(ctx, "hashtag:local", tag)
}

func (c *Client) streamingWS(ctx context.Context, stream, tag string) (chan Event, error) {
	params := url.Values{}
	params.Set("access_token", c.config.AccessToken)
	params.Set("stream", stream)
	if tag != "" {
		params.Set("tag", tag)
	}

	u, err := changeWebSocketScheme(c.config.Server)
	if err != nil {
		return nil, err
	}
	u.Path = path.Join(u.Path, "/api/v1/streaming")
	u.RawQuery = params.Encode()

	q := make(chan Event)
	go func() {
		for {
			err := c.handleWS(ctx, u.String(), q)
			if err != nil {
				return
			}
		}
	}()

	return q, nil
}

func (c *Client) handleWS(ctx context.Context, rawurl string, q chan Event) error {
	conn, err := c.dialRedirect(rawurl)
	if err != nil {
		q <- &ErrorEvent{err: err}

		// End.
		return err
	}
	defer conn.Close()

	for {
		select {
		case <-ctx.Done():
			q <- &ErrorEvent{err: ctx.Err()}

			// End.
			return ctx.Err()
		default:
		}

		var s Stream
		err := conn.ReadJSON(&s)
		if err != nil {
			q <- &ErrorEvent{err: err}

			// Reconnect.
			break
		}

		err = nil
		switch s.Event {
		case "update":
			var status Status
			err = json.Unmarshal([]byte(s.Payload.(string)), &status)
			if err == nil {
				q <- &UpdateEvent{Status: &status}
			}
		case "notification":
			var notification Notification
			err = json.Unmarshal([]byte(s.Payload.(string)), &notification)
			if err == nil {
				q <- &NotificationEvent{Notification: &notification}
			}
		case "delete":
			if err == nil {
				q <- &DeleteEvent{ID: int64(s.Payload.(float64))}
			}
		}
		if err != nil {
			q <- &ErrorEvent{err}
		}
	}

	return nil
}

func (c *Client) dialRedirect(rawurl string) (conn *websocket.Conn, err error) {
	for {
		conn, rawurl, err = c.dial(rawurl)
		if err != nil {
			return nil, err
		} else if conn != nil {
			return conn, nil
		}
	}
}

func (c *Client) dial(rawurl string) (*websocket.Conn, string, error) {
	conn, resp, err := c.Dial(rawurl, nil)
	if err != nil && err != websocket.ErrBadHandshake {
		return nil, "", err
	}
	defer resp.Body.Close()

	if loc := resp.Header.Get("Location"); loc != "" {
		u, err := changeWebSocketScheme(loc)
		if err != nil {
			return nil, "", err
		}

		return nil, u.String(), nil
	}

	return conn, "", err
}

func changeWebSocketScheme(rawurl string) (*url.URL, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}

	switch u.Scheme {
	case "http":
		u.Scheme = "ws"
	case "https":
		u.Scheme = "wss"
	}

	return u, nil
}
