package mastodon

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Notification hold information for mastodon notification.
type Notification struct {
	ID        int64     `json:"id"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
	Account   Account   `json:"account"`
	Status    *Status   `json:"status"`
}

// GetNotifications return notifications.
func (c *Client) GetNotifications(ctx context.Context) ([]*Notification, error) {
	var notifications []*Notification
	err := c.doAPI(ctx, http.MethodGet, "/api/v1/notifications", nil, &notifications)
	if err != nil {
		return nil, err
	}
	return notifications, nil
}

// GetNotifications return notification.
func (c *Client) GetNotification(ctx context.Context, id int64) (*Notification, error) {
	var notification Notification
	err := c.doAPI(ctx, http.MethodGet, fmt.Sprintf("/api/v1/notifications/%d", id), nil, &notification)
	if err != nil {
		return nil, err
	}
	return &notification, nil
}

// ClearNotifications clear notifications.
func (c *Client) ClearNotifications(ctx context.Context) error {
	return c.doAPI(ctx, http.MethodPost, "/api/v1/notifications/clear", nil, nil)
}
