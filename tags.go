package mastodon

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// TagInfo gets statistics and information about a tag
func (c *Client) TagInfo(ctx context.Context, tag string) (*FollowedTag, error) {
	var hashtag FollowedTag
	err := c.doAPI(ctx, http.MethodGet, fmt.Sprintf("/api/v1/tags/%s", url.PathEscape(string(tag))), nil, &hashtag, nil)
	if err != nil {
		return nil, err
	}
	return &hashtag, nil
}

// TagFollow lets you follow a hashtag
func (c *Client) TagFollow(ctx context.Context, tag string) (*FollowedTag, error) {
	var hashtag FollowedTag
	err := c.doAPI(ctx, http.MethodPost, fmt.Sprintf("/api/v1/tags/%s/follow", url.PathEscape(string(tag))), nil, &hashtag, nil)
	if err != nil {
		return nil, err
	}
	return &hashtag, nil
}

// TagUnfollow unfollows a hashtag.
func (c *Client) TagUnfollow(ctx context.Context, ID string) (*FollowedTag, error) {
	var tag FollowedTag
	err := c.doAPI(ctx, http.MethodPost, fmt.Sprintf("/api/v1/tags/%s/unfollow", ID), nil, &tag, nil)
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

// TagsFollowed returns a list of hashtags you follow.
func (c *Client) TagsFollowed(ctx context.Context, pg *Pagination) ([]*FollowedTag, error) {
	var hashtags []*FollowedTag
	err := c.doAPI(ctx, http.MethodGet, "/api/v1/followed_tags", nil, &hashtags, pg)
	if err != nil {
		return nil, err
	}
	return hashtags, nil
}
