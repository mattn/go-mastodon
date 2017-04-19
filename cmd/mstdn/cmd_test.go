package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"

	"github.com/mattn/go-mastodon"
	"github.com/urfave/cli"
)

func testWithServer(h http.HandlerFunc, testFunc func(*cli.App)) string {
	ts := httptest.NewServer(h)
	defer ts.Close()

	cli.OsExiter = func(n int) {}

	client := mastodon.NewClient(&mastodon.Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})

	var buf bytes.Buffer
	app := makeApp()
	app.Writer = &buf
	app.Metadata = map[string]interface{}{
		"client": client,
		"config": &mastodon.Config{
			Server: "example.com",
		},
	}
	testFunc(app)
	return buf.String()
}
