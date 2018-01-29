package mastodon

import (
	"context"
	"net/http"
)

// Instance hold information for mastodon instance.
type Instance struct {
	URI         string            `json:"uri"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	EMail       string            `json:"email"`
	Version     string            `json:"version,omitempty"`
	URLs        map[string]string `json:"urls,omitempty"`
	Stats       *InstanceStats    `json:"stats,omitempty"`
	Thumbnail   string            `json:"thumbnail,omitempty"`
}

// InstanceStats hold information for mastodon instance stats.
type InstanceStats struct {
	UserCount   int64 `json:"user_count"`
	StatusCount int64 `json:"status_count"`
	DomainCount int64 `json:"domain_count"`
}

// GetInstance return Instance.
func (c *Client) GetInstance(ctx context.Context) (*Instance, error) {
	var instance Instance
	err := c.doAPI(ctx, http.MethodGet, "/api/v1/instance", nil, &instance, nil)
	if err != nil {
		return nil, err
	}
	return &instance, nil
}
