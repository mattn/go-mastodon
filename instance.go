package mastodon

import (
	"context"
	"net/http"
)

// Instance hold information for mastodon instance.
type Instance struct {
	URI         string `json:"uri"`
	Title       string `json:"title"`
	Description string `json:"description"`
	EMail       string `json:"email"`
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
