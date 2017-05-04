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
func (c *Client) GetNotifications(ctx context.Context, pg *Pagination) ([]*Notification, *Pagination, error) {
	var notifications []*Notification
	retPG, err := c.doAPI(ctx, http.MethodGet, "/api/v1/notifications", nil, &notifications, pg)
	if err != nil {
		return nil, nil, err
	}
	return notifications, retPG, nil
}

// GetNotification return notification.
func (c *Client) GetNotification(ctx context.Context, id int64) (*Notification, error) {
	var notification Notification
	_, err := c.doAPI(ctx, http.MethodGet, fmt.Sprintf("/api/v1/notifications/%d", id), nil, &notification, nil)
	if err != nil {
		return nil, err
	}
	return &notification, nil
}

// ClearNotifications clear notifications.
func (c *Client) ClearNotifications(ctx context.Context) error {
	_, err := c.doAPI(ctx, http.MethodPost, "/api/v1/notifications/clear", nil, nil, nil)
	return err
}
