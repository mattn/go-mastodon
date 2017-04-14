package mastodon

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// Account hold information for mastodon account.
type Account struct {
	ID             int64     `json:"id"`
	Username       string    `json:"username"`
	Acct           string    `json:"acct"`
	DisplayName    string    `json:"display_name"`
	Locked         bool      `json:"locked"`
	CreatedAt      time.Time `json:"created_at"`
	FollowersCount int64     `json:"followers_count"`
	FollowingCount int64     `json:"following_count"`
	StatusesCount  int64     `json:"statuses_count"`
	Note           string    `json:"note"`
	URL            string    `json:"url"`
	Avatar         string    `json:"avatar"`
	AvatarStatic   string    `json:"avatar_static"`
	Header         string    `json:"header"`
	HeaderStatic   string    `json:"header_static"`
}

// GetAccount return Account.
func (c *Client) GetAccount(id int) (*Account, error) {
	var account Account
	err := c.doAPI(http.MethodGet, fmt.Sprintf("/api/v1/accounts/%d", id), nil, &account)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

// GetAccountCurrentUser return Account of current user.
func (c *Client) GetAccountCurrentUser() (*Account, error) {
	var account Account
	err := c.doAPI(http.MethodGet, "/api/v1/accounts/verify_credentials", nil, &account)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

// GetAccountFollowers return followers list.
func (c *Client) GetAccountFollowers(id int64) ([]*Account, error) {
	var accounts []*Account
	err := c.doAPI(http.MethodGet, fmt.Sprintf("/api/v1/accounts/%d/followers", id), nil, &accounts)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

// GetAccountFollowing return following list.
func (c *Client) GetAccountFollowing(id int64) ([]*Account, error) {
	var accounts []*Account
	err := c.doAPI(http.MethodGet, fmt.Sprintf("/api/v1/accounts/%d/following", id), nil, &accounts)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

// GetBlocks return block list.
func (c *Client) GetBlocks() ([]*Account, error) {
	var accounts []*Account
	err := c.doAPI(http.MethodGet, "/api/v1/blocks", nil, &accounts)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

// Relationship hold information for relation-ship to the account.
type Relationship struct {
	ID         int64 `json:"id"`
	Following  bool  `json:"following"`
	FollowedBy bool  `json:"followed_by"`
	Blocking   bool  `json:"blocking"`
	Muting     bool  `json:"muting"`
	Requested  bool  `json:"requested"`
}

// AccountFollow follow the account.
func (c *Client) AccountFollow(id int64) (*Relationship, error) {
	var relationship Relationship
	err := c.doAPI(http.MethodPost, fmt.Sprintf("/api/v1/accounts/%d/follow", id), nil, &relationship)
	if err != nil {
		return nil, err
	}
	return &relationship, nil
}

// AccountUnfollow unfollow the account.
func (c *Client) AccountUnfollow(id int64) (*Relationship, error) {
	var relationship Relationship
	err := c.doAPI(http.MethodPost, fmt.Sprintf("/api/v1/accounts/%d/unfollow", id), nil, &relationship)
	if err != nil {
		return nil, err
	}
	return &relationship, nil
}

// AccountBlock block the account.
func (c *Client) AccountBlock(id int64) (*Relationship, error) {
	var relationship Relationship
	err := c.doAPI(http.MethodPost, fmt.Sprintf("/api/v1/accounts/%d/block", id), nil, &relationship)
	if err != nil {
		return nil, err
	}
	return &relationship, nil
}

// AccountUnblock unblock the account.
func (c *Client) AccountUnblock(id int64) (*Relationship, error) {
	var relationship Relationship
	err := c.doAPI(http.MethodPost, fmt.Sprintf("/api/v1/accounts/%d/unblock", id), nil, &relationship)
	if err != nil {
		return nil, err
	}
	return &relationship, nil
}

// AccountMute mute the account.
func (c *Client) AccountMute(id int64) (*Relationship, error) {
	var relationship Relationship
	err := c.doAPI(http.MethodPost, fmt.Sprintf("/api/v1/accounts/%d/mute", id), nil, &relationship)
	if err != nil {
		return nil, err
	}
	return &relationship, nil
}

// AccountUnmute unmute the account.
func (c *Client) AccountUnmute(id int64) (*Relationship, error) {
	var relationship Relationship
	err := c.doAPI(http.MethodPost, fmt.Sprintf("/api/v1/accounts/%d/unmute", id), nil, &relationship)
	if err != nil {
		return nil, err
	}
	return &relationship, nil
}

// GetAccountRelationship return relationship for the account.
func (c *Client) GetAccountRelationship(id int64) ([]*Relationship, error) {
	params := url.Values{}
	params.Set("id", fmt.Sprint(id))

	var relationships []*Relationship
	err := c.doAPI(http.MethodGet, "/api/v1/accounts/relationship", params, &relationships)
	if err != nil {
		return nil, err
	}
	return relationships, nil
}

// AccountsSearch search accounts by query.
func (c *Client) AccountsSearch(q string, limit int64) ([]*Account, error) {
	params := url.Values{}
	params.Set("q", q)
	params.Set("limit", fmt.Sprint(limit))

	var accounts []*Account
	err := c.doAPI(http.MethodGet, "/api/v1/accounts/search", params, &accounts)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

// Follow send follow-request.
func (c *Client) FollowRemoteUser(uri string) (*Account, error) {
	params := url.Values{}
	params.Set("uri", uri)

	var account Account
	err := c.doAPI(http.MethodPost, "/api/v1/follows", params, &account)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

// GetFollowRequests return follow-requests.
func (c *Client) GetFollowRequests() ([]*Account, error) {
	var accounts []*Account
	err := c.doAPI(http.MethodGet, "/api/v1/follow_requests", nil, &accounts)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}
