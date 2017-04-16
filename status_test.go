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

func TestGetStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/statuses/1234567" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		fmt.Fprintln(w, `{"Content": "zzz"}`)
		return
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	_, err := client.GetStatus(123)
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	status, err := client.GetStatus(1234567)
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if status.Content != "zzz" {
		t.Fatalf("want %q but %q", "zzz", status.Content)
	}
}

func TestGetRebloggedBy(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/statuses/1234567/reblogged_by" {
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
	_, err := client.GetRebloggedBy(123)
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	rbs, err := client.GetRebloggedBy(1234567)
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if len(rbs) != 2 {
		t.Fatalf("result should be two: %d", len(rbs))
	}
	if rbs[0].Username != "foo" {
		t.Fatalf("want %q but %q", "foo", rbs[0].Username)
	}
	if rbs[1].Username != "bar" {
		t.Fatalf("want %q but %q", "bar", rbs[0].Username)
	}
}

func TestGetFavouritedBy(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/statuses/1234567/favourited_by" {
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
	_, err := client.GetFavouritedBy(123)
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	fbs, err := client.GetFavouritedBy(1234567)
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if len(fbs) != 2 {
		t.Fatalf("result should be two: %d", len(fbs))
	}
	if fbs[0].Username != "foo" {
		t.Fatalf("want %q but %q", "foo", fbs[0].Username)
	}
	if fbs[1].Username != "bar" {
		t.Fatalf("want %q but %q", "bar", fbs[0].Username)
	}
}

func TestReblog(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/statuses/1234567/reblog" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		fmt.Fprintln(w, `{"Content": "zzz"}`)
		return
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	_, err := client.Reblog(123)
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	status, err := client.Reblog(1234567)
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if status.Content != "zzz" {
		t.Fatalf("want %q but %q", "zzz", status.Content)
	}
}

func TestUnreblog(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/statuses/1234567/unreblog" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		fmt.Fprintln(w, `{"Content": "zzz"}`)
		return
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	_, err := client.Unreblog(123)
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	status, err := client.Unreblog(1234567)
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if status.Content != "zzz" {
		t.Fatalf("want %q but %q", "zzz", status.Content)
	}
}

func TestFavourite(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/statuses/1234567/favourite" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		fmt.Fprintln(w, `{"Content": "zzz"}`)
		return
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	_, err := client.Favourite(123)
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	status, err := client.Favourite(1234567)
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if status.Content != "zzz" {
		t.Fatalf("want %q but %q", "zzz", status.Content)
	}
}

func TestUnfavourite(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/statuses/1234567/unfavourite" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		fmt.Fprintln(w, `{"Content": "zzz"}`)
		return
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	_, err := client.Unfavourite(123)
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	status, err := client.Unfavourite(1234567)
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if status.Content != "zzz" {
		t.Fatalf("want %q but %q", "zzz", status.Content)
	}
}
