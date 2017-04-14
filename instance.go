package mastodon

import "net/http"

// Instance hold information for mastodon instance.
type Instance struct {
	URI         string `json:"uri"`
	Title       string `json:"title"`
	Description string `json:"description"`
	EMail       string `json:"email"`
}

// GetInstance return Instance.
func (c *Client) GetInstance() (*Instance, error) {
	var instance Instance
	err := c.doAPI(http.MethodGet, "/api/v1/instance", nil, &instance)
	if err != nil {
		return nil, err
	}
	return &instance, nil
}
