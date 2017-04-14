package mastodon

import (
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
func (c *Client) GetNotifications() ([]*Notification, error) {
	var notifications []*Notification
	err := c.doAPI(http.MethodGet, "/api/v1/notifications", nil, &notifications)
	if err != nil {
		return nil, err
	}
	return notifications, nil
}

// GetNotifications return notifications.
func (c *Client) GetNotification(id int64) (*Notification, error) {
	var notification Notification
	err := c.doAPI(http.MethodGet, fmt.Sprintf("/api/v1/notifications/%d", id), nil, &notification)
	if err != nil {
		return nil, err
	}
	return &notification, nil
}

// ClearNotifications clear notifications.
func (c *Client) ClearNotifications() error {
	return c.doAPI(http.MethodPost, "/api/v1/notifications/clear", nil, nil)
}
