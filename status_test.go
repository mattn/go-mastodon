package mastodon

import (
	"context"
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
	favs, err := client.GetFavourites(context.Background())
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
	_, err := client.GetStatus(context.Background(), 123)
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	status, err := client.GetStatus(context.Background(), 1234567)
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
	_, err := client.GetRebloggedBy(context.Background(), 123)
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	rbs, err := client.GetRebloggedBy(context.Background(), 1234567)
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
		t.Fatalf("want %q but %q", "bar", rbs[1].Username)
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
	_, err := client.GetFavouritedBy(context.Background(), 123)
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	fbs, err := client.GetFavouritedBy(context.Background(), 1234567)
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
		t.Fatalf("want %q but %q", "bar", fbs[1].Username)
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
	_, err := client.Reblog(context.Background(), 123)
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	status, err := client.Reblog(context.Background(), 1234567)
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
	_, err := client.Unreblog(context.Background(), 123)
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	status, err := client.Unreblog(context.Background(), 1234567)
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
	_, err := client.Favourite(context.Background(), 123)
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	status, err := client.Favourite(context.Background(), 1234567)
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
	_, err := client.Unfavourite(context.Background(), 123)
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	status, err := client.Unfavourite(context.Background(), 1234567)
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if status.Content != "zzz" {
		t.Fatalf("want %q but %q", "zzz", status.Content)
	}
}

func TestUploadMedia(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/media" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		if r.Method != "POST" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		fmt.Fprintln(w, `{"ID": 123}`)
		return
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	attachment, err := client.UploadMedia(context.Background(), "testdata/logo.png")
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if attachment.ID != 123 {
		t.Fatalf("want %q but %q", 123, attachment.ID)
	}
}
