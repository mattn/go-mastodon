package mastodon

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestTagInfo(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/tags/test" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		fmt.Fprintln(w, `
		{
			"name": "test",
			"url": "http://mastodon.example/tags/test",
			"history": [
			  {
				"day": "1668124800",
				"accounts": "1",
				"uses": "2"
			  },
			  {
				"day": "1668038400",
				"accounts": "0",
				"uses": "0"
			  }
			]
		}`)
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	_, err := client.TagInfo(context.Background(), "foo")
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	tag, err := client.TagInfo(context.Background(), "test")
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if tag.Name != "test" {
		t.Fatalf("want %q but %q", "test", tag.Name)
	}
	if tag.URL != "http://mastodon.example/tags/test" {
		t.Fatalf("want %q but %q", "http://mastodon.example/tags/test", tag.URL)
	}
	if len(tag.History) != 2 {
		t.Fatalf("result should be two: %d", len(tag.History))
	}
	uts := UnixTimeString{ time.Unix(1668124800, 0) }
	if tag.History[0].Day != uts {
		t.Fatalf("want %q but %q", uts, tag.History[0].Day)
	}
	if tag.History[0].Accounts != 1 {
		t.Fatalf("want %q but %q", 1, tag.History[0].Accounts)
	}
	if tag.History[0].Uses != 2 {
		t.Fatalf("want %q but %q", 2, tag.History[0].Uses)
	}
	if tag.Following != false {
		t.Fatalf("want %v but %v", nil, tag.Following)
	}
}

func TestTagFollow(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/tags/test/follow" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		fmt.Fprintln(w, `
		{
			"name": "test",
			"url": "http://mastodon.example/tags/test",
			"history": [
			  {
				"day": "1668124800",
				"accounts": "1",
				"uses": "2"
			  },
			  {
				"day": "1668038400",
				"accounts": "0",
				"uses": "0"
			  }
			],
			"following": true
		}`)
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	_, err := client.TagFollow(context.Background(), "foo")
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	tag, err := client.TagFollow(context.Background(), "test")
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if tag.Name != "test" {
		t.Fatalf("want %q but %q", "test", tag.Name)
	}
	if tag.URL != "http://mastodon.example/tags/test" {
		t.Fatalf("want %q but %q", "http://mastodon.example/tags/test", tag.URL)
	}
	if len(tag.History) != 2 {
		t.Fatalf("result should be two: %d", len(tag.History))
	}
	uts := UnixTimeString{ time.Unix(1668124800, 0) }
	if tag.History[0].Day != uts {
		t.Fatalf("want %q but %q", uts, tag.History[0].Day)
	}
	if tag.History[0].Accounts != 1 {
		t.Fatalf("want %q but %q", 1, tag.History[0].Accounts)
	}
	if tag.History[0].Uses != 2 {
		t.Fatalf("want %q but %q", 2, tag.History[0].Uses)
	}
	if tag.Following != true {
		t.Fatalf("want %v but %v", true, tag.Following)
	}
}

func TestTagUnfollow(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/tags/test/unfollow" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		fmt.Fprintln(w, `
		{
			"name": "test",
			"url": "http://mastodon.example/tags/test",
			"history": [
			  {
				"day": "1668124800",
				"accounts": "1",
				"uses": "2"
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

	_, err := client.TagUnfollow(context.Background(), "foo")
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	tag, err := client.TagUnfollow(context.Background(), "test")
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if tag.Name != "test" {
		t.Fatalf("want %q but %q", "test", tag.Name)
	}
	if tag.URL != "http://mastodon.example/tags/test" {
		t.Fatalf("want %q but %q", "http://mastodon.example/tags/test", tag.URL)
	}
	if len(tag.History) != 2 {
		t.Fatalf("result should be two: %d", len(tag.History))
	}
	uts := UnixTimeString{ time.Unix(1668124800, 0) }
	if tag.History[0].Day != uts {
		t.Fatalf("want %q but %q", uts, tag.History[0].Day)
	}
	if tag.History[0].Accounts != 1 {
		t.Fatalf("want %q but %q", 1, tag.History[0].Accounts)
	}
	if tag.History[0].Uses != 2 {
		t.Fatalf("want %q but %q", 2, tag.History[0].Uses)
	}
	if tag.Following != false {
		t.Fatalf("want %v but %v", false, tag.Following)
	}
}

func TestTagsFollowed(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/followed_tags" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		if r.FormValue("limit") == "1" {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, `
		[{
			"name": "test",
			"url": "http://mastodon.example/tags/test",
			"history": [
			  {
				"day": "1668124800",
				"accounts": "1",
				"uses": "2"
			  },
			  {
				"day": "1668038400",
				"accounts": "0",
				"uses": "0"
			  }
			],
			"following": true
		},
		{
			"name": "foo",
			"url": "http://mastodon.example/tags/foo",
			"history": [
			  {
				"day": "1668124800",
				"accounts": "1",
				"uses": "2"
			  },
			  {
				"day": "1668038400",
				"accounts": "0",
				"uses": "0"
			  }
			],
			"following": true
		}]`)
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	_, err := client.TagsFollowed(context.Background(), &Pagination{Limit: 1})
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	tags, err := client.TagsFollowed(context.Background(), nil)
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if len(tags) != 2 {
		t.Fatalf("want %q but %q", 2, len(tags))
	}
	if tags[0].Name != "test" {
		t.Fatalf("want %q but %q", "test", tags[0].Name)
	}
	if tags[0].URL != "http://mastodon.example/tags/test" {
		t.Fatalf("want %q but %q", "http://mastodon.example/tags/test", tags[0].URL)
	}
	if len(tags[0].History) != 2 {
		t.Fatalf("result should be two: %d", len(tags[0].History))
	}
	uts := UnixTimeString{ time.Unix(1668124800, 0) }
	if tags[0].History[0].Day != uts {
		t.Fatalf("want %q but %q", uts, tags[0].History[0].Day)
	}
	if tags[0].History[0].Accounts != 1 {
		t.Fatalf("want %q but %q", 1, tags[0].History[0].Accounts)
	}
	if tags[0].History[0].Uses != 2 {
		t.Fatalf("want %q but %q", 2, tags[0].History[0].Uses)
	}
	if tags[0].Following != true {
		t.Fatalf("want %v but %v", true, tags[0].Following)
	}
}
