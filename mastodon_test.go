package mastodon

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestAuthenticate(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.FormValue("username") != "valid" || r.FormValue("password") != "user" {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		fmt.Fprintln(w, `{"AccessToken": "zoo"}`)
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
		fmt.Fprintln(w, `{"AccessToken": "zoo"}`)
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
		fmt.Fprintln(w, `[{"Content": "foo"}, {"Content": "bar"}]`)
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
	tl, err := client.GetTimelineHome(context.Background())
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
	_, err := client.GetTimelineHome(ctx)
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

func TestLinkHeader(t *testing.T) {
	tests := []struct {
		header []string
		rel    string
		want   []string
	}{
		{
			header: []string{`<http://example.com/?max_id=3>; rel="foo"`},
			rel:    "boo",
			want:   nil,
		},
		{
			header: []string{`<http://example.com/?max_id=3>; rel="foo"`},
			rel:    "foo",
			want:   []string{"http://example.com/?max_id=3"},
		},
		{
			header: []string{`<http://example.com/?max_id=3>; rel="foo1"`},
			rel:    "foo",
			want:   nil,
		},
		{
			header: []string{`<http://example.com/?max_id=3>; rel="foo", <http://example.com/?max_id=4>; rel="bar"`},
			rel:    "foo",
			want:   []string{"http://example.com/?max_id=3"},
		},
		{
			header: []string{`<http://example.com/?max_id=3>; rel="foo", <http://example.com/?max_id=4>; rel="bar"`},
			rel:    "bar",
			want:   []string{"http://example.com/?max_id=4"},
		},
	}

	for _, test := range tests {
		h := make(http.Header)
		for _, he := range test.header {
			h.Add("Link", he)
		}
		got := linkHeader(h, test.rel)
		if !reflect.DeepEqual(got, test.want) {
			t.Fatalf("want %v but %v", test.want, got)
		}
	}
}
