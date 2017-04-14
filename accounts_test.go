package mastodon

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetBlocks(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/blocks" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		fmt.Fprintln(w, `[{"Username": "foo"}, {"Username": "bar"}]`)
		return
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	bl, err := client.GetBlocks()
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if len(bl) != 2 {
		t.Fatalf("result should be two: %d", len(bl))
	}
	if bl[0].Username != "foo" {
		t.Fatalf("want %q but %q", "foo", bl[0].Username)
	}
	if bl[1].Username != "bar" {
		t.Fatalf("want %q but %q", "bar", bl[0].Username)
	}
}

func TestAccountFollow(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/accounts/1234567/follow" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		fmt.Fprintln(w, `{"id":1234567,"following":true}`)
		return
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	rel, err := client.AccountFollow(123)
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	rel, err = client.AccountFollow(1234567)
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if rel.ID != 1234567 {
		t.Fatalf("want %d but %d", 1234567, rel.ID)
	}
	if !rel.Following {
		t.Fatalf("want %t but %t", true, rel.Following)
	}
}

func TestAccountUnfollow(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/accounts/1234567/unfollow" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		fmt.Fprintln(w, `{"id":1234567,"following":false}`)
		return
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	rel, err := client.AccountUnfollow(123)
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	rel, err = client.AccountUnfollow(1234567)
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if rel.ID != 1234567 {
		t.Fatalf("want %d but %d", 1234567, rel.ID)
	}
	if rel.Following {
		t.Fatalf("want %t but %t", false, rel.Following)
	}
}
