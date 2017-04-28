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
data: {"content": "foo"}
event: notification
data: {"type": "mention"}
event: delete
data: 1234567
:thump
	`)
	go func() {
		defer close(q)
		err := handleReader(context.Background(), q, r)
		if err != nil {
			t.Fatalf("should not be fail: %v", err)
		}
	}()
	var passUpdate, passNotification, passDelete bool
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
		}
	}
	if !passUpdate || !passNotification || !passDelete {
		t.Fatalf("have not passed through somewhere: update %t, notification %t, delete %t",
			passUpdate, passNotification, passDelete)
	}
}

func TestStreamingPublic(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/streaming/public" {
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
	q, err := client.StreamingPublic(ctx, false)
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	time.AfterFunc(3*time.Second, func() {
		cancel()
		close(q)
	})
	events := []Event{}
	for e := range q {
		events = append(events, e)
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
