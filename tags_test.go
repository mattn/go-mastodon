package mastodon

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTagUnfollow(t *testing.T) {
	t.Parallel()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{
			"name": "Test",
			"url": "http://mastodon.example/tags/test",
			"history": [
				{
				"day": "1668556800",
				"accounts": "0",
				"uses": "0"
				},
				{
				"day": "1668470400",
				"accounts": "0",
				"uses": "0"
				},
				{
				"day": "1668384000",
				"accounts": "0",
				"uses": "0"
				},
				{
				"day": "1668297600",
				"accounts": "1",
				"uses": "1"
				},
				{
				"day": "1668211200",
				"accounts": "0",
				"uses": "0"
				},
				{
				"day": "1668124800",
				"accounts": "0",
				"uses": "0"
				},
				{
				"day": "1668038400",
				"accounts": "0",
				"uses": "0"
				}
			],
			"following": false
		}`)
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	tag, err := client.TagUnfollow(context.Background(), "Test")
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if tag.Name != "Test" {
		t.Fatalf("want %q but %q", "Test", tag.Name)
	}
	if tag.Following {
		t.Fatalf("want %t but %t", false, tag.Following)
	}
}
