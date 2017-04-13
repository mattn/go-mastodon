# go-mastodon

***EXPERIMENTAL***

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
