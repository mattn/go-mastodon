package mastodon

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestHandleReader(t *testing.T) {
	q := make(chan Event)
	r := strings.NewReader(`
event: update
data: {content: error}
event: update
data: {"content": "foo"}
event: notification
data: {"type": "mention"}
event: delete
data: 1234567
:thump
	`)
	go func() {
		defer close(q)
		err := handleReader(q, r)
		if err != nil {
			t.Fatalf("should not be fail: %v", err)
		}
	}()
	var passUpdate, passNotification, passDelete, passError bool
	for e := range q {
		switch event := e.(type) {
		case *UpdateEvent:
			passUpdate = true
			if event.Status.Content != "foo" {
				t.Fatalf("want %q but %q", "foo", event.Status.Content)
			}
		case *NotificationEvent:
			passNotification = true
			if event.Notification.Type != "mention" {
				t.Fatalf("want %q but %q", "mention", event.Notification.Type)
			}
		case *DeleteEvent:
			passDelete = true
			if event.ID != 1234567 {
				t.Fatalf("want %d but %d", 1234567, event.ID)
			}
		case *ErrorEvent:
			passError = true
			if event.err == nil {
				t.Fatalf("should be fail: %v", event.err)
			}
		}
	}
	if !passUpdate || !passNotification || !passDelete || !passError {
		t.Fatalf("have not passed through somewhere: "+
			"update %t, notification %t, delete %t, error %t",
			passUpdate, passNotification, passDelete, passError)
	}
}

func TestStreaming(t *testing.T) {
	var isEnd bool
	canErr := true
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isEnd {
			return
		} else if canErr {
			canErr = false
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		f := w.(http.Flusher)
		fmt.Fprintln(w, `
event: update
data: {"content": "foo"}
		`)
		f.Flush()
		isEnd = true
	}))
	defer ts.Close()

	c := NewClient(&Config{Server: ":"})
	_, err := c.streaming(context.Background(), "", nil)
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}

	c = NewClient(&Config{Server: ts.URL})
	ctx, cancel := context.WithCancel(context.Background())
	time.AfterFunc(time.Second, cancel)
	q, err := c.streaming(ctx, "", nil)
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	var cnt int
	var passError, passUpdate bool
	for e := range q {
		switch event := e.(type) {
		case *ErrorEvent:
			passError = true
			if event.err == nil {
				t.Fatalf("should be fail: %v", event.err)
			}
		case *UpdateEvent:
			cnt++
			passUpdate = true
			if event.Status.Content != "foo" {
				t.Fatalf("want %q but %q", "foo", event.Status.Content)
			}
		}
	}
	if cnt != 1 {
		t.Fatalf("result should be one: %d", cnt)
	}
	if !passError || !passUpdate {
		t.Fatalf("have not passed through somewhere: error %t, update %t", passError, passUpdate)
	}
}

func TestDoStreaming(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.(http.Flusher).Flush()
		time.Sleep(time.Second)
	}))
	defer ts.Close()

	c := NewClient(&Config{Server: ts.URL})

	req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	time.AfterFunc(time.Millisecond, cancel)
	req = req.WithContext(ctx)

	q := make(chan Event)
	go func() {
		defer close(q)
		c.doStreaming(req, q)
		if err != nil {
			t.Fatalf("should not be fail: %v", err)
		}
	}()
	var passError bool
	for e := range q {
		if event, ok := e.(*ErrorEvent); ok {
			passError = true
			if event.err == nil {
				t.Fatalf("should be fail: %v", event.err)
			}
		}
	}
	if !passError {
		t.Fatalf("have not passed through: error %t", passError)
	}
}

func TestStreamingUser(t *testing.T) {
	var isEnd bool
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isEnd {
			return
		} else if r.URL.Path != "/api/v1/streaming/user" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		f, _ := w.(http.Flusher)
		fmt.Fprintln(w, `
event: update
data: {"content": "foo"}
		`)
		f.Flush()
		isEnd = true
	}))
	defer ts.Close()

	c := NewClient(&Config{Server: ts.URL})
	ctx, cancel := context.WithCancel(context.Background())
	time.AfterFunc(time.Second, cancel)
	q, err := c.StreamingUser(ctx)
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	events := []Event{}
	for e := range q {
		if _, ok := e.(*ErrorEvent); !ok {
			events = append(events, e)
		}
	}
	if len(events) != 1 {
		t.Fatalf("result should be one: %d", len(events))
	}
	if events[0].(*UpdateEvent).Status.Content != "foo" {
		t.Fatalf("want %q but %q", "foo", events[0].(*UpdateEvent).Status.Content)
	}
}

func TestStreamingPublic(t *testing.T) {
	var isEnd bool
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isEnd {
			return
		} else if r.URL.Path != "/api/v1/streaming/public/local" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		f, _ := w.(http.Flusher)
		fmt.Fprintln(w, `
event: update
data: {"content": "foo"}
		`)
		f.Flush()

		fmt.Fprintln(w, `
event: update
data: {"content": "bar"}
		`)
		f.Flush()
		isEnd = true
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	ctx, cancel := context.WithCancel(context.Background())
	q, err := client.StreamingPublic(ctx, true)
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	time.AfterFunc(time.Second, cancel)
	events := []Event{}
	for e := range q {
		if _, ok := e.(*ErrorEvent); !ok {
			events = append(events, e)
		}
	}
	if len(events) != 2 {
		t.Fatalf("result should be two: %d", len(events))
	}
	if events[0].(*UpdateEvent).Status.Content != "foo" {
		t.Fatalf("want %q but %q", "foo", events[0].(*UpdateEvent).Status.Content)
	}
	if events[1].(*UpdateEvent).Status.Content != "bar" {
		t.Fatalf("want %q but %q", "bar", events[1].(*UpdateEvent).Status.Content)
	}
}

func TestStreamingHashtag(t *testing.T) {
	var isEnd bool
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isEnd {
			return
		} else if r.URL.Path != "/api/v1/streaming/hashtag/local" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		f, _ := w.(http.Flusher)
		fmt.Fprintln(w, `
event: update
data: {"content": "foo"}
		`)
		f.Flush()
		isEnd = true
	}))
	defer ts.Close()

	client := NewClient(&Config{Server: ts.URL})
	ctx, cancel := context.WithCancel(context.Background())
	time.AfterFunc(time.Second, cancel)
	q, err := client.StreamingHashtag(ctx, "hashtag", true)
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	events := []Event{}
	for e := range q {
		if _, ok := e.(*ErrorEvent); !ok {
			events = append(events, e)
		}
	}
	if len(events) != 1 {
		t.Fatalf("result should be one: %d", len(events))
	}
	if events[0].(*UpdateEvent).Status.Content != "foo" {
		t.Fatalf("want %q but %q", "foo", events[0].(*UpdateEvent).Status.Content)
	}
}
