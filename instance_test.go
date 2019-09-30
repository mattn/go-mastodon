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
		fmt.Fprintln(w, `{"title": "mastodon", "uri": "http://mstdn.example.com", "description": "test mastodon", "email": "mstdn@mstdn.example.com", "contact_account": {"username": "mattn"}}`)
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
	if ins.URI != "http://mstdn.example.com" {
		t.Fatalf("want %q but %q", "http://mstdn.example.com", ins.URI)
	}
	if ins.Description != "test mastodon" {
		t.Fatalf("want %q but %q", "test mastodon", ins.Description)
	}
	if ins.EMail != "mstdn@mstdn.example.com" {
		t.Fatalf("want %q but %q", "mstdn@mstdn.example.com", ins.EMail)
	}
	if ins.ContactAccount.Username != "mattn" {
		t.Fatalf("want %q but %q", "mattn", ins.ContactAccount.Username)
	}
}

func TestGetInstanceMore(t *testing.T) {
	canErr := true
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if canErr {
			canErr = false
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		fmt.Fprintln(w, `{"title": "mastodon", "uri": "http://mstdn.example.com", "description": "test mastodon", "email": "mstdn@mstdn.example.com", "version": "0.0.1", "urls":{"foo":"http://stream1.example.com", "bar": "http://stream2.example.com"}, "thumbnail": "http://mstdn.example.com/logo.png", "stats":{"user_count":1, "status_count":2, "domain_count":3}}}`)
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
	if ins.URI != "http://mstdn.example.com" {
		t.Fatalf("want %q but %q", "mastodon", ins.URI)
	}
	if ins.Description != "test mastodon" {
		t.Fatalf("want %q but %q", "test mastodon", ins.Description)
	}
	if ins.EMail != "mstdn@mstdn.example.com" {
		t.Fatalf("want %q but %q", "mstdn@mstdn.example.com", ins.EMail)
	}
	if ins.Version != "0.0.1" {
		t.Fatalf("want %q but %q", "0.0.1", ins.Version)
	}
	if ins.URLs["foo"] != "http://stream1.example.com" {
		t.Fatalf("want %q but %q", "http://stream1.example.com", ins.Version)
	}
	if ins.URLs["bar"] != "http://stream2.example.com" {
		t.Fatalf("want %q but %q", "http://stream2.example.com", ins.Version)
	}
	if ins.Thumbnail != "http://mstdn.example.com/logo.png" {
		t.Fatalf("want %q but %q", "http://mstdn.example.com/logo.png", ins.Thumbnail)
	}
	if ins.Stats == nil {
		t.Fatal("status should be nil")
	}
	if ins.Stats.UserCount != 1 {
		t.Fatalf("want %v but %v", 1, ins.Stats.UserCount)
	}
	if ins.Stats.StatusCount != 2 {
		t.Fatalf("want %v but %v", 2, ins.Stats.StatusCount)
	}
	if ins.Stats.DomainCount != 3 {
		t.Fatalf("want %v but %v", 3, ins.Stats.DomainCount)
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
