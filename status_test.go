package mastodon

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestGetFavourites(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `[{"content": "foo"}, {"content": "bar"}]`)
		return
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	favs, err := client.GetFavourites(context.Background(), nil)
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

func TestGetBookmarks(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `[{"content": "foo"}, {"content": "bar"}]`)
		return
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	books, err := client.GetBookmarks(context.Background(), nil)
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if len(books) != 2 {
		t.Fatalf("result should be two: %d", len(books))
	}
	if books[0].Content != "foo" {
		t.Fatalf("want %q but %q", "foo", books[0].Content)
	}
	if books[1].Content != "bar" {
		t.Fatalf("want %q but %q", "bar", books[1].Content)
	}
}

func TestGetStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/statuses/1234567" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		fmt.Fprintln(w, `{"content": "zzz", "emojis":[{"shortcode":"💩", "url":"http://example.com", "static_url": "http://example.com/static"}]}`)
		return
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	_, err := client.GetStatus(context.Background(), "123")
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	status, err := client.GetStatus(context.Background(), "1234567")
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if status.Content != "zzz" {
		t.Fatalf("want %q but %q", "zzz", status.Content)
	}
	if len(status.Emojis) != 1 {
		t.Fatal("should have emojis")
	}
	if status.Emojis[0].ShortCode != "💩" {
		t.Fatalf("want %q but %q", "💩", status.Emojis[0].ShortCode)
	}
	if status.Emojis[0].URL != "http://example.com" {
		t.Fatalf("want %q but %q", "https://example.com", status.Emojis[0].URL)
	}
	if status.Emojis[0].StaticURL != "http://example.com/static" {
		t.Fatalf("want %q but %q", "https://example.com/static", status.Emojis[0].StaticURL)
	}
}

func TestGetStatusCard(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/statuses/1234567/card" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		fmt.Fprintln(w, `{"title": "zzz"}`)
		return
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	_, err := client.GetStatusCard(context.Background(), "123")
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	card, err := client.GetStatusCard(context.Background(), "1234567")
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if card.Title != "zzz" {
		t.Fatalf("want %q but %q", "zzz", card.Title)
	}
}

func TestGetStatusContext(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/statuses/1234567/context" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		fmt.Fprintln(w, `{"ancestors": [{"content": "zzz"},{"content": "bbb"}]}`)
		return
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	_, err := client.GetStatusContext(context.Background(), "123")
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	context, err := client.GetStatusContext(context.Background(), "1234567")
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if len(context.Ancestors) != 2 {
		t.Fatalf("Ancestors should have 2 entries but %q", len(context.Ancestors))
	}
	if context.Ancestors[0].Content != "zzz" {
		t.Fatalf("want %q but %q", "zzz", context.Ancestors[0].Content)
	}
	if context.Ancestors[1].Content != "bbb" {
		t.Fatalf("want %q but %q", "bbb", context.Ancestors[1].Content)
	}
	if len(context.Descendants) > 0 {
		t.Fatalf("Descendants should not be included")
	}
}

func TestGetRebloggedBy(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/statuses/1234567/reblogged_by" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		fmt.Fprintln(w, `[{"username": "foo"}, {"username": "bar"}]`)
		return
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	_, err := client.GetRebloggedBy(context.Background(), "123", nil)
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	rbs, err := client.GetRebloggedBy(context.Background(), "1234567", nil)
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
		fmt.Fprintln(w, `[{"username": "foo"}, {"username": "bar"}]`)
		return
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	_, err := client.GetFavouritedBy(context.Background(), "123", nil)
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	fbs, err := client.GetFavouritedBy(context.Background(), "1234567", nil)
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
		fmt.Fprintln(w, `{"content": "zzz"}`)
		return
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	_, err := client.Reblog(context.Background(), "123")
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	status, err := client.Reblog(context.Background(), "1234567")
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
		fmt.Fprintln(w, `{"content": "zzz"}`)
		return
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	_, err := client.Unreblog(context.Background(), "123")
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	status, err := client.Unreblog(context.Background(), "1234567")
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
		fmt.Fprintln(w, `{"content": "zzz"}`)
		return
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	_, err := client.Favourite(context.Background(), "123")
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	status, err := client.Favourite(context.Background(), "1234567")
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
		fmt.Fprintln(w, `{"content": "zzz"}`)
		return
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	_, err := client.Unfavourite(context.Background(), "123")
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	status, err := client.Unfavourite(context.Background(), "1234567")
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if status.Content != "zzz" {
		t.Fatalf("want %q but %q", "zzz", status.Content)
	}
}

func TestBookmark(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/statuses/1234567/bookmark" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		fmt.Fprintln(w, `{"content": "zzz"}`)
		return
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	_, err := client.Bookmark(context.Background(), "123")
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	status, err := client.Bookmark(context.Background(), "1234567")
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if status.Content != "zzz" {
		t.Fatalf("want %q but %q", "zzz", status.Content)
	}
}

func TestUnbookmark(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/statuses/1234567/unbookmark" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		fmt.Fprintln(w, `{"content": "zzz"}`)
		return
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	_, err := client.Unbookmark(context.Background(), "123")
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	status, err := client.Unbookmark(context.Background(), "1234567")
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if status.Content != "zzz" {
		t.Fatalf("want %q but %q", "zzz", status.Content)
	}
}

func TestGetTimelinePublic(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("local") == "" {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		fmt.Fprintln(w, `[{"content": "foo"}, {"content": "bar"}]`)
	}))
	defer ts.Close()

	client := NewClient(&Config{Server: ts.URL})
	_, err := client.GetTimelinePublic(context.Background(), false, nil)
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	tl, err := client.GetTimelinePublic(context.Background(), true, nil)
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if len(tl) != 2 {
		t.Fatalf("result should be two: %d", len(tl))
	}
	if tl[0].Content != "foo" {
		t.Fatalf("want %q but %q", "foo", tl[0].Content)
	}
	if tl[1].Content != "bar" {
		t.Fatalf("want %q but %q", "bar", tl[1].Content)
	}
}

func TestGetTimelineDirect(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `[{"id": "4", "unread":false, "last_status" : {"content": "zzz"}}, {"id": "3", "unread":true, "last_status" : {"content": "bar"}}]`)
	}))
	defer ts.Close()

	client := NewClient(&Config{Server: ts.URL})
	tl, err := client.GetTimelineDirect(context.Background(), nil)
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if len(tl) != 2 {
		t.Fatalf("result should be two: %d", len(tl))
	}
	if tl[0].Content != "zzz" {
		t.Fatalf("want %q but %q", "foo", tl[0].Content)
	}
	if tl[1].Content != "bar" {
		t.Fatalf("want %q but %q", "bar", tl[1].Content)
	}
}

func TestGetTimelineHashtag(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/timelines/tag/zzz" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		fmt.Fprintln(w, `[{"content": "zzz"},{"content": "yyy"}]`)
		return
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	_, err := client.GetTimelineHashtag(context.Background(), "notfound", false, nil)
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	tags, err := client.GetTimelineHashtag(context.Background(), "zzz", true, nil)
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if len(tags) != 2 {
		t.Fatalf("should have %q entries but %q", "2", len(tags))
	}
	if tags[0].Content != "zzz" {
		t.Fatalf("want %q but %q", "zzz", tags[0].Content)
	}
	if tags[1].Content != "yyy" {
		t.Fatalf("want %q but %q", "zzz", tags[1].Content)
	}
}

func TestGetTimelineList(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/timelines/list/1" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		fmt.Fprintln(w, `[{"content": "zzz"},{"content": "yyy"}]`)
		return
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	_, err := client.GetTimelineList(context.Background(), "notfound", nil)
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	tags, err := client.GetTimelineList(context.Background(), "1", nil)
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if len(tags) != 2 {
		t.Fatalf("should have %q entries but %q", "2", len(tags))
	}
	if tags[0].Content != "zzz" {
		t.Fatalf("want %q but %q", "zzz", tags[0].Content)
	}
	if tags[1].Content != "yyy" {
		t.Fatalf("want %q but %q", "zzz", tags[1].Content)
	}
}

func TestGetTimelineMedia(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("local") == "" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		fmt.Fprintln(w, `[{"content": "zzz"},{"content": "yyy"}]`)
		return
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	_, err := client.GetTimelineMedia(context.Background(), false, nil)
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	tags, err := client.GetTimelineMedia(context.Background(), true, nil)
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if len(tags) != 2 {
		t.Fatalf("should have %q entries but %q", "2", len(tags))
	}
	if tags[0].Content != "zzz" {
		t.Fatalf("want %q but %q", "zzz", tags[0].Content)
	}
	if tags[1].Content != "yyy" {
		t.Fatalf("want %q but %q", "zzz", tags[1].Content)
	}
}

func TestDeleteStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/statuses/1234567" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		if r.Method != "DELETE" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusMethodNotAllowed)
			return
		}
		return
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	err := client.DeleteStatus(context.Background(), "123")
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	err = client.DeleteStatus(context.Background(), "1234567")
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
}

func TestSearch(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v2/search" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		if r.RequestURI != "/api/v2/search?q=q&resolve=false" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusBadRequest)
			return
		}

		fmt.Fprintln(w, `
			{"accounts":[{"username": "zzz"},{"username": "yyy"}],
			"statuses":[{"content": "aaa"}],
			"hashtags":[{"name": "tag"},{"name": "tag2"},{"name": "tag3"}]
		}`)
		return
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	ret, err := client.Search(context.Background(), "q", false)
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if len(ret.Accounts) != 2 {
		t.Fatalf("Accounts have %q entries, but %q", "2", len(ret.Accounts))
	}
	if ret.Accounts[0].Username != "zzz" {
		t.Fatalf("Accounts Username should %q , but %q", "zzz", ret.Accounts[0].Username)
	}
	if len(ret.Statuses) != 1 {
		t.Fatalf("Statuses have %q entries, but %q", "1", len(ret.Statuses))
	}
	if ret.Statuses[0].Content != "aaa" {
		t.Fatalf("Statuses Content should %q , but %q", "aaa", ret.Statuses[0].Content)
	}
	if len(ret.Hashtags) != 3 {
		t.Fatalf("Hashtags have %q entries, but %q", "3", len(ret.Hashtags))
	}
	if ret.Hashtags[2].Name != "tag3" {
		t.Fatalf("Hashtags[2] should %q , but %q", "tag3", ret.Hashtags[2])
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
		fmt.Fprintln(w, `{"id": 123}`)
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
	if attachment.ID != "123" {
		t.Fatalf("want %q but %q", "123", attachment.ID)
	}
	file, err := os.Open("testdata/logo.png")
	if err != nil {
		t.Fatalf("could not open file: %v", err)
	}
	defer file.Close()
	writerAttachment, err := client.UploadMediaFromReader(context.Background(), file)
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if writerAttachment.ID != "123" {
		t.Fatalf("want %q but %q", "123", attachment.ID)
	}
}

func TestGetConversations(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/conversations" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		}
		fmt.Fprintln(w, `[{"id": "4", "unread":false, "last_status" : {"content": "zzz"}}, {"id": "3", "unread":true, "last_status" : {"content": "bar"}}]`)
		return
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	convs, err := client.GetConversations(context.Background(), nil)
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if len(convs) != 2 {
		t.Fatalf("result should be 2: %d", len(convs))
	}
	if convs[0].ID != "4" {
		t.Fatalf("want %q but %q", "4", convs[0].ID)
	}
	if convs[0].LastStatus.Content != "zzz" {
		t.Fatalf("want %q but %q", "zzz", convs[0].LastStatus.Content)
	}
	if convs[1].Unread != true {
		t.Fatalf("unread should be true: %t", convs[1].Unread)
	}
}

func TestDeleteConversation(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/conversations/12345678" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		if r.Method != "DELETE" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusMethodNotAllowed)
			return
		}
		return
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "hoge",
	})
	err := client.DeleteConversation(context.Background(), "12345678")
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
}

func TestMarkConversationsAsRead(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/conversations/111111/read" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		if r.Method != "POST" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		return
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	err := client.MarkConversationAsRead(context.Background(), "111111")
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
}
