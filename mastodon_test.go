package mastodon

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
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
		return
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
		return
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
	if want := "Post " + ts.URL + "/oauth/token: context canceled"; want != err.Error() {
		t.Fatalf("want %q but %q", want, err.Error())
	}
}

func TestPostStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer zoo" {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		fmt.Fprintln(w, `{"access_token": "zoo"}`)
		return
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
		return
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
	if want := "Post " + ts.URL + "/api/v1/statuses: context canceled"; want != err.Error() {
		t.Fatalf("want %q but %q", want, err.Error())
	}
}

func TestGetTimelineHome(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `[{"content": "foo"}, {"content": "bar"}]`)
		return
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
		return
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
	if want := "Get " + ts.URL + "/api/v1/timelines/home: context canceled"; want != err.Error() {
		t.Fatalf("want %q but %q", want, err.Error())
	}
}

func TestForTheCoverages(t *testing.T) {
	(*UpdateEvent)(nil).event()
	(*NotificationEvent)(nil).event()
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
	if err == nil {
		t.Fatalf("should be fail: %v", err)
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
		SinceID: "789",
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
	if after.Get("since_id") != "" {
		t.Fatalf("result should be empty string: %q", after.Get("since_id"))
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
}
