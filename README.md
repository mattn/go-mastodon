# go-mastodon

[![Build Status](https://travis-ci.org/mattn/go-mastodon.svg?branch=master)](https://travis-ci.org/mattn/go-mastodon)
[![Coverage Status](https://coveralls.io/repos/github/mattn/go-mastodon/badge.svg?branch=master)](https://coveralls.io/github/mattn/go-mastodon?branch=master)
[![GoDoc](https://godoc.org/github.com/mattn/go-mastodon?status.svg)](http://godoc.org/github.com/mattn/go-mastodon)
[![Go Report Card](https://goreportcard.com/badge/github.com/mattn/go-mastodon)](https://goreportcard.com/report/github.com/mattn/go-mastodon)

## Usage

```go
c := mastodon.NewClient(&mastodon.Config{
	Server:       "https://mstdn.jp",
	ClientID:     "client-id",
	ClientSecret: "client-secret",
})
err := c.Authenticate("your-username", "your-password")
if err != nil {
	log.Fatal(err)
}
timeline, err := c.GetTimeline("/api/v1/timelines/home")
if err != nil {
	log.Fatal(err)
}
```
## Status of implementations

* [x] GET /api/v1/accounts/:id
* [x] GET /api/v1/accounts/verify_credentials
* [ ] PATCH /api/v1/accounts/update_credentials
* [x] GET /api/v1/accounts/:id/followers
* [x] GET /api/v1/accounts/:id/following
* [ ] GET /api/v1/accounts/:id/statuses
* [x] POST /api/v1/accounts/:id/follow
* [x] POST /api/v1/accounts/:id/unfollow
* [x] GET /api/v1/accounts/:id/block
* [x] GET /api/v1/accounts/:id/unblock
* [x] GET /api/v1/accounts/:id/mute
* [x] GET /api/v1/accounts/:id/unmute
* [x] GET /api/v1/accounts/relationships
* [x] GET /api/v1/accounts/search
* [x] POST /api/v1/apps
* [x] GET /api/v1/blocks
* [x] GET /api/v1/favourites
* [x] GET /api/v1/follow_requests
* [ ] POST /api/v1/follow_requests/authorize
* [ ] POST /api/v1/follow_requests/reject
* [x] POST /api/v1/follows
* [x] GET /api/v1/instance
* [ ] POST /api/v1/media
* [ ] GET /api/v1/mutes
* [x] GET /api/v1/notifications
* [x] GET /api/v1/notifications/:id
* [x] POST /api/v1/notifications/clear
* [ ] GET /api/v1/reports
* [ ] POST /api/v1/reports
* [ ] GET /api/v1/search
* [x] GET /api/v1/statuses/:id
* [x] GET /api/v1/statuses/:id/context
* [x] GET /api/v1/statuses/:id/card
* [ ] GET /api/v1/statuses/:id/reblogged_by
* [ ] GET /api/v1/statuses/:id/favourited_by
* [ ] POST /api/v1/statuses
* [x] DELETE /api/v1/statuses/:id
* [ ] POST /api/v1/statuses/:id/reblog
* [ ] POST /api/v1/statuses/:id/unreblog
* [ ] POST /api/v1/statuses/:id/favourite
* [ ] POST /api/v1/statuses/:id/unfavourite
* [x] GET /api/v1/timelines/home
* [x] GET /api/v1/timelines/public
* [x] GET /api/v1/timelines/tag/:hashtag

## Installation

```
$ go get github.com/mattn/go-mastodon
```

## License

MIT

## Author

Yasuhiro Matsumoto (a.k.a. mattn)
