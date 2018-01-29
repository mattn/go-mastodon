package mastodon

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
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

func TestGetInstanceActivity(t *testing.T) {
	canErr := true
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if canErr {
			canErr = false
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		fmt.Fprintln(w, `[{"week":"1516579200","statuses":"1","logins":"1","registrations":"0"}]`)
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server: ts.URL,
	})
	_, err := client.GetInstanceActivity(context.Background())
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	activity, err := client.GetInstanceActivity(context.Background())
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if activity[0].Week != Unixtime(time.Unix(1516579200, 0)) {
		t.Fatalf("want %v but %v", Unixtime(time.Unix(1516579200, 0)), activity[0].Week)
	}
	if activity[0].Logins != 1 {
		t.Fatalf("want %q but %q", 1, activity[0].Logins)
	}
}

func TestGetInstancePeers(t *testing.T) {
	canErr := true
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if canErr {
			canErr = false
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		fmt.Fprintln(w, `["mastodon.social","mstdn.jp"]`)
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server: ts.URL,
	})
	_, err := client.GetInstancePeers(context.Background())
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	peers, err := client.GetInstancePeers(context.Background())
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if peers[0] != "mastodon.social" {
		t.Fatalf("want %q but %q", "mastodon.social", peers[0])
	}
	if peers[1] != "mstdn.jp" {
		t.Fatalf("want %q but %q", "mstdn.jp", peers[1])
	}
}
