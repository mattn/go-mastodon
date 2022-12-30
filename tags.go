package mastodon

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// Tag hold information for tag.
type Tag struct {
	Name      string      `json:"name"`
	URL       string      `json:"url"`
	History   []History   `json:"history"`
	Following interface{} `json:"following"`
}

// TagInfo gets statistics and information about a tag
func (c *Client) TagInfo(ctx context.Context, tag string) (*Tag, error) {
	var hashtag Tag
	err := c.doAPI(ctx, http.MethodGet, fmt.Sprintf("/api/v1/tags/%s", url.PathEscape(string(tag))), nil, &hashtag, nil)
	if err != nil {
		return nil, err
	}
	return &hashtag, nil
}

// TagFollow lets you follow a hashtag
func (c *Client) TagFollow(ctx context.Context, tag string) (*Tag, error) {
	var hashtag Tag
	err := c.doAPI(ctx, http.MethodPost, fmt.Sprintf("/api/v1/tags/%s/follow", url.PathEscape(string(tag))), nil, &hashtag, nil)
	if err != nil {
		return nil, err
	}
	return &hashtag, nil
}

// TagUnfollow lets you unfollow a hashtag
func (c *Client) TagUnfollow(ctx context.Context, tag string) (*Tag, error) {
	var hashtag Tag
	err := c.doAPI(ctx, http.MethodPost, fmt.Sprintf("/api/v1/tags/%s/unfollow", url.PathEscape(string(tag))), nil, &hashtag, nil)
	if err != nil {
		return nil, err
	}
	return &hashtag, nil
}

func (c *Client) TagsFollowed(ctx context.Context, pg *Pagination) ([]*Tag, error) {
	var hashtags []*Tag
	err := c.doAPI(ctx, http.MethodGet, "/api/v1/followed_tags", nil, &hashtags, pg)
	if err != nil {
		return nil, err
	}
	return hashtags, nil
}
