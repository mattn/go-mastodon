package mastodon

import "time"

// Poll hold information for mastodon polls.
type Poll struct {
	ID         ID           `json:"id"`
	ExpiresAt  time.Time    `json:"expires_at"`
	Expired    bool         `json:"expired"`
	Multiple   bool         `json:"multiple"`
	VotesCount int64        `json:"votes_count"`
	Options    []PollOption `json:"options"`
	Voted      bool         `json:"voted"`
}

// Poll hold information for a mastodon poll option.
type PollOption struct {
	Title      string `json:"title"`
	VotesCount int64  `json:"votes_count"`
}
