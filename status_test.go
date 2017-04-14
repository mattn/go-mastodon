package mastodon

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetFavourites(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `[{"Content": "foo"}, {"Content": "bar"}]`)
		return
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	favs, err := client.GetFavourites()
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if len(favs) != 2 {
		t.Fatalf("result should be two: %d", len(favs))
	}
	if favs[0].Content != "foo" {
		t.Fatalf("want %q but %q", "foo", favs[0].Content)
	}
	if favs[1].Content != "bar" {
		t.Fatalf("want %q but %q", "bar", favs[1].Content)
	}
}
