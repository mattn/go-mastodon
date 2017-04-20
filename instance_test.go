package mastodon

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetInstance(t *testing.T) {
	canErr := true
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if canErr {
			canErr = false
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		fmt.Fprintln(w, `{"title": "mastodon"}`)
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	_, err := client.GetInstance(context.Background())
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	ins, err := client.GetInstance(context.Background())
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if ins.Title != "mastodon" {
		t.Fatalf("want %q but %q", "mastodon", ins.Title)
	}
}
