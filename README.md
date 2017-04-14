# go-mastodon

[![Build Status](https://travis-ci.org/mattn/go-mastodon.png?branch=master)](https://travis-ci.org/mattn/go-mastodon)
[![Coverage Status](https://coveralls.io/repos/mattn/go-mastodon/badge.png?branch=HEAD)](https://coveralls.io/r/mattn/go-mastodon?branch=HEAD)
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

## Installation

```
$ go get github.com/mattn/go-mastodon
```

## License

MIT

## Author

Yasuhiro Matsumoto (a.k.a. mattn)
