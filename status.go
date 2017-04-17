package mastodon

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// Status is struct to hold status.
type Status struct {
	ID                 int64        `json:"id"`
	CreatedAt          time.Time    `json:"created_at"`
	InReplyToID        interface{}  `json:"in_reply_to_id"`
	InReplyToAccountID interface{}  `json:"in_reply_to_account_id"`
	Sensitive          bool         `json:"sensitive"`
	SpoilerText        string       `json:"spoiler_text"`
	Visibility         string       `json:"visibility"`
	Application        Application  `json:"application"`
	Account            Account      `json:"account"`
	MediaAttachments   []Attachment `json:"media_attachments"`
	Mentions           []Mention    `json:"mentions"`
	Tags               []Tag        `json:"tags"`
	URI                string       `json:"uri"`
	Content            string       `json:"content"`
	URL                string       `json:"url"`
	ReblogsCount       int64        `json:"reblogs_count"`
	FavouritesCount    int64        `json:"favourites_count"`
	Reblog             *Status      `json:"reblog"`
	Favourited         interface{}  `json:"favourited"`
	Reblogged          interface{}  `json:"reblogged"`
}

// Context hold information for mastodon context.
type Context struct {
	Ancestors   []*Status `json:"ancestors"`
	Descendants []*Status `json:"descendants"`
}

// Card hold information for mastodon card.
type Card struct {
	URL         string `json:"url"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       string `json:"image"`
}

// GetFavourites return the favorite list of the current user.
func (c *Client) GetFavourites(ctx context.Context) ([]*Status, error) {
	var statuses []*Status
	err := c.doAPI(ctx, http.MethodGet, "/api/v1/favourites", nil, &statuses)
	if err != nil {
		return nil, err
	}
	return statuses, nil
}

// GetStatus return status specified by id.
func (c *Client) GetStatus(ctx context.Context, id int64) (*Status, error) {
	var status Status
	err := c.doAPI(ctx, http.MethodGet, fmt.Sprintf("/api/v1/statuses/%d", id), nil, &status)
	if err != nil {
		return nil, err
	}
	return &status, nil
}

// GetStatusContext return status specified by id.
func (c *Client) GetStatusContext(ctx context.Context, id int64) (*Context, error) {
	var context Context
	err := c.doAPI(ctx, http.MethodGet, fmt.Sprintf("/api/v1/statuses/%d/context", id), nil, &context)
	if err != nil {
		return nil, err
	}
	return &context, nil
}

// GetStatusCard return status specified by id.
func (c *Client) GetStatusCard(ctx context.Context, id int64) (*Card, error) {
	var card Card
	err := c.doAPI(ctx, http.MethodGet, fmt.Sprintf("/api/v1/statuses/%d/card", id), nil, &card)
	if err != nil {
		return nil, err
	}
	return &card, nil
}

// GetRebloggedBy returns the account list of the user who reblogged the toot of id.
func (c *Client) GetRebloggedBy(ctx context.Context, id int64) ([]*Account, error) {
	var accounts []*Account
	err := c.doAPI(ctx, http.MethodGet, fmt.Sprintf("/api/v1/statuses/%d/reblogged_by", id), nil, &accounts)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

// GetFavouritedBy returns the account list of the user who liked the toot of id.
func (c *Client) GetFavouritedBy(ctx context.Context, id int64) ([]*Account, error) {
	var accounts []*Account
	err := c.doAPI(ctx, http.MethodGet, fmt.Sprintf("/api/v1/statuses/%d/favourited_by", id), nil, &accounts)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

// Reblog is reblog the toot of id and return status of reblog.
func (c *Client) Reblog(ctx context.Context, id int64) (*Status, error) {
	var status Status
	err := c.doAPI(ctx, http.MethodPost, fmt.Sprintf("/api/v1/statuses/%d/reblog", id), nil, &status)
	if err != nil {
		return nil, err
	}
	return &status, nil
}

// Unreblog is unreblog the toot of id and return status of the original toot.
func (c *Client) Unreblog(ctx context.Context, id int64) (*Status, error) {
	var status Status
	err := c.doAPI(ctx, http.MethodPost, fmt.Sprintf("/api/v1/statuses/%d/unreblog", id), nil, &status)
	if err != nil {
		return nil, err
	}
	return &status, nil
}

// Favourite is favourite the toot of id and return status of the favourite toot.
func (c *Client) Favourite(ctx context.Context, id int64) (*Status, error) {
	var status Status
	err := c.doAPI(ctx, http.MethodPost, fmt.Sprintf("/api/v1/statuses/%d/favourite", id), nil, &status)
	if err != nil {
		return nil, err
	}
	return &status, nil
}

// Unfavourite is unfavourite the toot of id and return status of the unfavourite toot.
func (c *Client) Unfavourite(ctx context.Context, id int64) (*Status, error) {
	var status Status
	err := c.doAPI(ctx, http.MethodPost, fmt.Sprintf("/api/v1/statuses/%d/unfavourite", id), nil, &status)
	if err != nil {
		return nil, err
	}
	return &status, nil
}

// GetTimelineHome return statuses from home timeline.
func (c *Client) GetTimelineHome(ctx context.Context) ([]*Status, error) {
	var statuses []*Status
	err := c.doAPI(ctx, http.MethodGet, "/api/v1/timelines/home", nil, &statuses)
	if err != nil {
		return nil, err
	}
	return statuses, nil
}

// GetTimelineHashtag return statuses from tagged timeline.
func (c *Client) GetTimelineHashtag(ctx context.Context, tag string) ([]*Status, error) {
	var statuses []*Status
	err := c.doAPI(ctx, http.MethodGet, fmt.Sprintf("/api/v1/timelines/tag/%s", (&url.URL{Path: tag}).EscapedPath()), nil, &statuses)
	if err != nil {
		return nil, err
	}
	return statuses, nil
}

// PostStatus post the toot.
func (c *Client) PostStatus(ctx context.Context, toot *Toot) (*Status, error) {
	params := url.Values{}
	params.Set("status", toot.Status)
	if toot.InReplyToID > 0 {
		params.Set("in_reply_to_id", fmt.Sprint(toot.InReplyToID))
	}
	// TODO: media_ids, senstitive, spoiler_text, visibility
	//params.Set("visibility", "public")

	var status Status
	err := c.doAPI(ctx, http.MethodPost, "/api/v1/statuses", params, &status)
	if err != nil {
		return nil, err
	}
	return &status, nil
}

// DeleteStatus delete the toot.
func (c *Client) DeleteStatus(ctx context.Context, id int64) error {
	return c.doAPI(ctx, http.MethodDelete, fmt.Sprintf("/api/v1/statuses/%d", id), nil, nil)
}

// Search search content with query.
func (c *Client) Search(ctx context.Context, q string, resolve bool) (*Results, error) {
	params := url.Values{}
	params.Set("q", q)
	params.Set("resolve", fmt.Sprint(resolve))
	var results Results
	err := c.doAPI(ctx, http.MethodGet, "/api/v1/search", params, &results)
	if err != nil {
		return nil, err
	}
	return &results, nil
}

// PostMedia upload a media attachment.
func (c *Client) UploadMedia(ctx context.Context, file string) (*Attachment, error) {
	var attachment Attachment
	err := c.doAPI(ctx, http.MethodPost, "/api/v1/media", file, &attachment)
	if err != nil {
		return nil, err
	}
	return &attachment, nil
}
