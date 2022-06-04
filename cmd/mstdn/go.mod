module github.com/mattn/go-mastodon/cmd/mstdn

go 1.16

replace github.com/mattn/go-mastodon => ../..

require (
	github.com/PuerkitoBio/goquery v1.8.0
	github.com/fatih/color v1.13.0
	github.com/mattn/go-mastodon v0.0.4
	github.com/mattn/go-tty v0.0.4
	github.com/urfave/cli v1.22.9
	golang.org/x/net v0.0.0-20220531201128-c960675eff93
)
