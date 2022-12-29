package mastodon

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

const (
	redirectURI = "urn:ietf:wg:oauth:2.0:oob"
)

func TestDoAPI(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("max_id") == "999" {
			w.Header().Set("Link", `<:>; rel="next"`)
		} else {
			w.Header().Set("Link", `<http://example.com?max_id=234>; rel="next", <http://example.com?since_id=890>; rel="prev"`)
		}
		fmt.Fprintln(w, `[{"username": "foo"}, {"username": "bar"}]`)
	}))
	defer ts.Close()

	c := NewClient(&Config{Server: ts.URL})
	var accounts []Account
	err := c.doAPI(context.Background(), http.MethodGet, "/", nil, &accounts, &Pagination{
		MaxID: "999",
	})
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}

	pg := &Pagination{
		MaxID:   "123",
		SinceID: "789",
		Limit:   10,
	}
	err = c.doAPI(context.Background(), http.MethodGet, "/", url.Values{}, &accounts, pg)
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if pg.MaxID != "234" {
		t.Fatalf("want %q but %q", "234", pg.MaxID)
	}
	if pg.SinceID != "890" {
		t.Fatalf("want %q but %q", "890", pg.SinceID)
	}
	if accounts[0].Username != "foo" {
		t.Fatalf("want %q but %q", "foo", accounts[0].Username)
	}
	if accounts[1].Username != "bar" {
		t.Fatalf("want %q but %q", "bar", accounts[1].Username)
	}

	pg = &Pagination{
		MaxID:   "123",
		SinceID: "789",
		Limit:   10,
	}
	err = c.doAPI(context.Background(), http.MethodGet, "/", nil, &accounts, pg)
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if pg.MaxID != "234" {
		t.Fatalf("want %q but %q", "234", pg.MaxID)
	}
	if pg.SinceID != "890" {
		t.Fatalf("want %q but %q", "890", pg.SinceID)
	}
	if accounts[0].Username != "foo" {
		t.Fatalf("want %q but %q", "foo", accounts[0].Username)
	}
	if accounts[1].Username != "bar" {
		t.Fatalf("want %q but %q", "bar", accounts[1].Username)
	}

	// *Pagination is nil
	err = c.doAPI(context.Background(), http.MethodGet, "/", nil, &accounts, nil)
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if accounts[0].Username != "foo" {
		t.Fatalf("want %q but %q", "foo", accounts[0].Username)
	}
	if accounts[1].Username != "bar" {
		t.Fatalf("want %q but %q", "bar", accounts[1].Username)
	}
}

func TestAuthenticate(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.FormValue("username") != "valid" || r.FormValue("password") != "user" {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		fmt.Fprintln(w, `{"access_token": "zoo"}`)
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
	})
	err := client.Authenticate(context.Background(), "invalid", "user")
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}

	client = NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
	})
	err = client.Authenticate(context.Background(), "valid", "user")
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
}

func TestAuthenticateWithCancel(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(3 * time.Second)
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
	})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := client.Authenticate(ctx, "invalid", "user")
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	if want := fmt.Sprintf("Post %q: context canceled", ts.URL+"/oauth/token"); want != err.Error() {
		t.Fatalf("want %q but %q", want, err.Error())
	}
}

// DEPRECATED: AuthenticateApp is deprecated and replaced by GetAppAccessToken
func TestAuthenticateApp(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.FormValue("client_id") != "foo" || r.FormValue("client_secret") != "bar" {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		fmt.Fprintln(w, `{"name":"zzz","website":"yyy","vapid_key":"xxx"}`)
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bat",
	})
	err := client.AuthenticateApp(context.Background())
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}

	client = NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
	})
	err = client.AuthenticateApp(context.Background())
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
}

func TestGetAppAccessToken(t *testing.T) {
	wantAccessToken := "applicationAccessToken"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.FormValue("client_id") != "foo" || r.FormValue("client_secret") != "bar" {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, `{"accesS_token":"%s"}\n`, wantAccessToken)
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bat",
	})

	err := client.GetAppAccessToken(context.Background(), redirectURI)
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}

	client = NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
	})

	err = client.GetAppAccessToken(context.Background(), redirectURI)
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	gotAccessToken := client.Config.AccessToken

	if wantAccessToken != gotAccessToken {
		t.Fatalf("want %s but got %s", wantAccessToken, gotAccessToken)
	}
}

func TestGetUserAccessToken(t *testing.T) {
	wantAccessToken := "userAccessToken"
	authCode := "AuthorizationCode"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.FormValue("client_id") != "foo" || r.FormValue("client_secret") != "bar" || r.FormValue("code") != authCode {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, `{"accesS_token":"%s"}\n`, wantAccessToken)
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bat",
	})

	err := client.GetUserAccessToken(context.Background(), authCode, redirectURI)
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}

	client = NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
	})

	err = client.GetUserAccessToken(context.Background(), authCode, redirectURI)
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	gotAccessToken := client.Config.AccessToken

	if wantAccessToken != gotAccessToken {
		t.Fatalf("want %s but got %s", wantAccessToken, gotAccessToken)
	}
}

func TestPostStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer zoo" {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		fmt.Fprintln(w, `{"access_token": "zoo"}`)
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
	})
	_, err := client.PostStatus(context.Background(), &Toot{
		Status: "foobar",
	})
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}

	client = NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	_, err = client.PostStatus(context.Background(), &Toot{
		Status: "foobar",
	})
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
}

func TestPostStatusWithCancel(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(3 * time.Second)
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
	})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := client.PostStatus(ctx, &Toot{
		Status: "foobar",
	})
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	if want := fmt.Sprintf("Post %q: context canceled", ts.URL+"/api/v1/statuses"); want != err.Error() {
		t.Fatalf("want %q but %q", want, err.Error())
	}
}
func TestPostStatusParams(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/statuses" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		r.ParseForm()
		if r.FormValue("media_ids[]") != "" && r.FormValue("poll[options][]") != "" {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		}
		s := Status{
			ID:      ID("1"),
			Content: fmt.Sprintf("<p>%s</p>", r.FormValue("status")),
		}
		if r.FormValue("in_reply_to_id") != "" {
			s.InReplyToID = ID(r.FormValue("in_reply_to_id"))
		}
		if r.FormValue("visibility") != "" {
			s.Visibility = (r.FormValue("visibility"))
		}
		if r.FormValue("language") != "" {
			s.Language = (r.FormValue("language"))
		}
		if r.FormValue("sensitive") == "true" {
			s.Sensitive = true
			s.SpoilerText = fmt.Sprintf("<p>%s</p>", r.FormValue("spoiler_text"))
		}
		if r.FormValue("media_ids[]") != "" {
			for _, id := range r.Form["media_ids[]"] {
				s.MediaAttachments = append(s.MediaAttachments,
					Attachment{ID: ID(id)})
			}
		}
		if r.FormValue("poll[options][]") != "" {
			p := Poll{}
			for _, opt := range r.Form["poll[options][]"] {
				p.Options = append(p.Options, PollOption{
					Title:      opt,
					VotesCount: 0,
				})
			}
			if r.FormValue("poll[multiple]") == "true" {
				p.Multiple = true
			}
			s.Poll = &p
		}
		json.NewEncoder(w).Encode(s)
	}))
	defer ts.Close()
	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	s, err := client.PostStatus(context.Background(), &Toot{
		Status:      "foobar",
		InReplyToID: ID("2"),
		Visibility:  "unlisted",
		Language:    "sv",
		Sensitive:   true,
		SpoilerText: "bar",
		MediaIDs:    []ID{"1", "2"},
		Poll: &TootPoll{
			Options: []string{"A", "B"},
		},
	})
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if len(s.MediaAttachments) > 0 && s.Poll != nil {
		t.Fatal("should not fail, can't have both Media and Poll")
	}
	if s.Content != "<p>foobar</p>" {
		t.Fatalf("want %q but %q", "<p>foobar</p>", s.Content)
	}
	if s.InReplyToID != "2" {
		t.Fatalf("want %q but %q", "2", s.InReplyToID)
	}
	if s.Visibility != "unlisted" {
		t.Fatalf("want %q but %q", "unlisted", s.Visibility)
	}
	if s.Language != "sv" {
		t.Fatalf("want %q but %q", "sv", s.Language)
	}
	if s.Sensitive != true {
		t.Fatalf("want %t but %t", true, s.Sensitive)
	}
	if s.SpoilerText != "<p>bar</p>" {
		t.Fatalf("want %q but %q", "<p>bar</p>", s.SpoilerText)
	}
	s, err = client.PostStatus(context.Background(), &Toot{
		Status: "foobar",
		Poll: &TootPoll{
			Multiple:   true,
			Options:    []string{"A", "B"},
			HideTotals: true,
		},
	})
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if s.Poll == nil {
		t.Fatalf("poll should not be %v", s.Poll)
	}
	if len(s.Poll.Options) != 2 {
		t.Fatalf("want %q but %q", 2, len(s.Poll.Options))
	}
	if s.Poll.Options[0].Title != "A" {
		t.Fatalf("want %q but %q", "A", s.Poll.Options[0].Title)
	}
	if s.Poll.Options[1].Title != "B" {
		t.Fatalf("want %q but %q", "B", s.Poll.Options[1].Title)
	}
	if s.Poll.Multiple != true {
		t.Fatalf("want %t but %t", true, s.Poll.Multiple)
	}
}

func TestUpdateStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer zoo" {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		fmt.Fprintln(w, `{"access_token": "zoo"}`)
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
	})
	_, err := client.UpdateStatus(context.Background(), &Toot{
		Status: "foobar",
	}, ID("1"))
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}

	client = NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	_, err = client.UpdateStatus(context.Background(), &Toot{
		Status: "foobar",
	}, ID("1"))
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
}

func TestUpdateStatusWithCancel(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(3 * time.Second)
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
	})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := client.UpdateStatus(ctx, &Toot{
		Status: "foobar",
	}, ID("1"))
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	if want := fmt.Sprintf("Put %q: context canceled", ts.URL+"/api/v1/statuses/1"); want != err.Error() {
		t.Fatalf("want %q but %q", want, err.Error())
	}
}
func TestUpdateStatusParams(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/statuses/1" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		r.ParseForm()
		if r.FormValue("media_ids[]") != "" && r.FormValue("poll[options][]") != "" {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		}
		s := Status{
			ID:      ID("1"),
			Content: fmt.Sprintf("<p>%s</p>", r.FormValue("status")),
		}
		if r.FormValue("in_reply_to_id") != "" {
			s.InReplyToID = ID(r.FormValue("in_reply_to_id"))
		}
		if r.FormValue("visibility") != "" {
			s.Visibility = (r.FormValue("visibility"))
		}
		if r.FormValue("language") != "" {
			s.Language = (r.FormValue("language"))
		}
		if r.FormValue("sensitive") == "true" {
			s.Sensitive = true
			s.SpoilerText = fmt.Sprintf("<p>%s</p>", r.FormValue("spoiler_text"))
		}
		if r.FormValue("media_ids[]") != "" {
			for _, id := range r.Form["media_ids[]"] {
				s.MediaAttachments = append(s.MediaAttachments,
					Attachment{ID: ID(id)})
			}
		}
		if r.FormValue("poll[options][]") != "" {
			p := Poll{}
			for _, opt := range r.Form["poll[options][]"] {
				p.Options = append(p.Options, PollOption{
					Title:      opt,
					VotesCount: 0,
				})
			}
			if r.FormValue("poll[multiple]") == "true" {
				p.Multiple = true
			}
			s.Poll = &p
		}
		json.NewEncoder(w).Encode(s)
	}))
	defer ts.Close()
	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	s, err := client.UpdateStatus(context.Background(), &Toot{
		Status:      "foobar",
		InReplyToID: ID("2"),
		Visibility:  "unlisted",
		Language:    "sv",
		Sensitive:   true,
		SpoilerText: "bar",
		MediaIDs:    []ID{"1", "2"},
		Poll: &TootPoll{
			Options: []string{"A", "B"},
		},
	}, ID("1"))
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if len(s.MediaAttachments) > 0 && s.Poll != nil {
		t.Fatal("should not fail, can't have both Media and Poll")
	}
	if s.Content != "<p>foobar</p>" {
		t.Fatalf("want %q but %q", "<p>foobar</p>", s.Content)
	}
	if s.InReplyToID != "2" {
		t.Fatalf("want %q but %q", "2", s.InReplyToID)
	}
	if s.Visibility != "unlisted" {
		t.Fatalf("want %q but %q", "unlisted", s.Visibility)
	}
	if s.Language != "sv" {
		t.Fatalf("want %q but %q", "sv", s.Language)
	}
	if s.Sensitive != true {
		t.Fatalf("want %t but %t", true, s.Sensitive)
	}
	if s.SpoilerText != "<p>bar</p>" {
		t.Fatalf("want %q but %q", "<p>bar</p>", s.SpoilerText)
	}
	s, err = client.UpdateStatus(context.Background(), &Toot{
		Status: "foobar",
		Poll: &TootPoll{
			Multiple: true,
			Options:  []string{"A", "B"},
		},
	}, ID("1"))
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if s.Poll == nil {
		t.Fatalf("poll should not be %v", s.Poll)
	}
	if len(s.Poll.Options) != 2 {
		t.Fatalf("want %q but %q", 2, len(s.Poll.Options))
	}
	if s.Poll.Options[0].Title != "A" {
		t.Fatalf("want %q but %q", "A", s.Poll.Options[0].Title)
	}
	if s.Poll.Options[1].Title != "B" {
		t.Fatalf("want %q but %q", "B", s.Poll.Options[1].Title)
	}
	if s.Poll.Multiple != true {
		t.Fatalf("want %t but %t", true, s.Poll.Multiple)
	}
}

func TestGetTimelineHome(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `[{"content": "foo"}, {"content": "bar"}]`)
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
	})
	_, err := client.PostStatus(context.Background(), &Toot{
		Status: "foobar",
	})
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}

	client = NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	tl, err := client.GetTimelineHome(context.Background(), nil)
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

func TestGetTimelineHomeWithCancel(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(3 * time.Second)
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := client.GetTimelineHome(ctx, nil)
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	if want := fmt.Sprintf("Get %q: context canceled", ts.URL+"/api/v1/timelines/home"); want != err.Error() {
		t.Fatalf("want %q but %q", want, err.Error())
	}
}

func TestForTheCoverages(t *testing.T) {
	(*UpdateEvent)(nil).event()
	(*UpdateEditEvent)(nil).event()
	(*NotificationEvent)(nil).event()
	(*ConversationEvent)(nil).event()
	(*DeleteEvent)(nil).event()
	(*ErrorEvent)(nil).event()
	_ = (&ErrorEvent{io.EOF}).Error()
}

func TestNewPagination(t *testing.T) {
	_, err := newPagination("")
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}

	_, err = newPagination(`<:>; rel="next"`)
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}

	_, err = newPagination(`<:>; rel="prev"`)
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}

	_, err = newPagination(`<http://example.com?min_id=abc>; rel="prev"`)
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}

	pg, err := newPagination(`<http://example.com?max_id=123>; rel="next", <http://example.com?since_id=789>; rel="prev"`)
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if pg.MaxID != "123" {
		t.Fatalf("want %q but %q", "123", pg.MaxID)
	}
	if pg.SinceID != "789" {
		t.Fatalf("want %q but %q", "789", pg.SinceID)
	}
}

func TestGetPaginationID(t *testing.T) {
	_, err := getPaginationID(":", "max_id")
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}

	_, err = getPaginationID("http://example.com?max_id=abc", "max_id")
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}

	id, err := getPaginationID("http://example.com?max_id=123", "max_id")
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if id != "123" {
		t.Fatalf("want %q but %q", "123", id)
	}
}

func TestPaginationSetValues(t *testing.T) {
	p := &Pagination{
		MaxID:   "123",
		SinceID: "456",
		MinID:   "789",
		Limit:   10,
	}
	before := url.Values{"key": {"value"}}
	after := p.setValues(before)
	if after.Get("key") != "value" {
		t.Fatalf("want %q but %q", "value", after.Get("key"))
	}
	if after.Get("max_id") != "123" {
		t.Fatalf("want %q but %q", "123", after.Get("max_id"))
	}
	if after.Get("since_id") != "456" {
		t.Fatalf("want %q but %q", "456", after.Get("since_id"))
	}
	if after.Get("min_id") != "789" {
		t.Fatalf("want %q but %q", "789", after.Get("min_id"))
	}
	if after.Get("limit") != "10" {
		t.Fatalf("want %q but %q", "10", after.Get("limit"))
	}

	p = &Pagination{
		MaxID:   "",
		SinceID: "789",
	}
	before = url.Values{}
	after = p.setValues(before)
	if after.Get("max_id") != "" {
		t.Fatalf("result should be empty string: %q", after.Get("max_id"))
	}
	if after.Get("since_id") != "789" {
		t.Fatalf("want %q but %q", "789", after.Get("since_id"))
	}
	if after.Get("min_id") != "" {
		t.Fatalf("result should be empty string: %q", after.Get("min_id"))
	}
}
