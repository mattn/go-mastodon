package mastodon

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestHandleReader(t *testing.T) {
	large := "large"
	largeContent := strings.Repeat(large, 2*(bufio.MaxScanTokenSize/len(large)))

	q := make(chan Event)
	r := strings.NewReader(fmt.Sprintf(`
event: update
data: {content: error}
event: update
data: {"content": "foo"}
event: update
data: {"content": "%s"}
event: notification
data: {"type": "mention"}
event: delete
data: 1234567
event: status.update
data: {"content": "foo"}
event: conversation
data: {"id":"819516","unread":true,"accounts":[{"id":"108892712797543112","username":"a","acct":"a@pl.nulled.red","display_name":"a","locked":false,"bot":true,"discoverable":false,"group":false,"created_at":"2022-08-27T00:00:00.000Z","note":"a (pleroma edition)","url":"https://pl.nulled.red/users/a","avatar":"https://files.mastodon.social/cache/accounts/avatars/108/892/712/797/543/112/original/975674b2caa61034.png","avatar_static":"https://files.mastodon.social/cache/accounts/avatars/108/892/712/797/543/112/original/975674b2caa61034.png","header":"https://files.mastodon.social/cache/accounts/headers/108/892/712/797/543/112/original/f61d0382356caa0e.png","header_static":"https://files.mastodon.social/cache/accounts/headers/108/892/712/797/543/112/original/f61d0382356caa0e.png","followers_count":0,"following_count":0,"statuses_count":362,"last_status_at":"2022-11-13","emojis":[],"fields":[]}],"last_status":{"id":"109346889330629417","created_at":"2022-11-15T08:31:57.476Z","in_reply_to_id":null,"in_reply_to_account_id":null,"sensitive":false,"spoiler_text":"","visibility":"direct","language":null,"uri":"https://pl.nulled.red/objects/c869c5be-c184-4706-a45d-3459d9aa711c","url":"https://pl.nulled.red/objects/c869c5be-c184-4706-a45d-3459d9aa711c","replies_count":0,"reblogs_count":0,"favourites_count":0,"edited_at":null,"favourited":false,"reblogged":false,"muted":false,"bookmarked":false,"content":"test <span class=\"h-card\"><a class=\"u-url mention\" href=\"https://mastodon.social/@trwnh\" rel=\"nofollow noopener noreferrer\" target=\"_blank\">@<span>trwnh</span></a></span>","filtered":[],"reblog":null,"account":{"id":"108892712797543112","username":"a","acct":"a@pl.nulled.red","display_name":"a","locked":false,"bot":true,"discoverable":false,"group":false,"created_at":"2022-08-27T00:00:00.000Z","note":"a (pleroma edition)","url":"https://pl.nulled.red/users/a","avatar":"https://files.mastodon.social/cache/accounts/avatars/108/892/712/797/543/112/original/975674b2caa61034.png","avatar_static":"https://files.mastodon.social/cache/accounts/avatars/108/892/712/797/543/112/original/975674b2caa61034.png","header":"https://files.mastodon.social/cache/accounts/headers/108/892/712/797/543/112/original/f61d0382356caa0e.png","header_static":"https://files.mastodon.social/cache/accounts/headers/108/892/712/797/543/112/original/f61d0382356caa0e.png","followers_count":0,"following_count":0,"statuses_count":362,"last_status_at":"2022-11-13","emojis":[],"fields":[]},"media_attachments":[],"mentions":[{"id":"14715","username":"trwnh","url":"https://mastodon.social/@trwnh","acct":"trwnh"}],"tags":[],"emojis":[],"card":null,"poll":null}}
:thump
	`, largeContent))
	var wg sync.WaitGroup
	wg.Add(1)
	errs := make(chan error, 1)
	go func() {
		defer wg.Done()
		defer close(q)
		err := handleReader(q, r)
		if err != nil {
			t.Errorf("should not be fail: %v", err)
		}
		errs <- err
	}()
	var passUpdate, passUpdateLarge, passNotification, passDelete, passError bool
	for e := range q {
		switch event := e.(type) {
		case *UpdateEvent:
			if event.Status.Content == "foo" {
				passUpdate = true
			} else if event.Status.Content == largeContent {
				passUpdateLarge = true
			} else {
				t.Fatalf("bad update content: %q", event.Status.Content)
			}
		case *UpdateEditEvent:
			if event.Status.Content == "foo" {
				passUpdate = true
			} else if event.Status.Content == largeContent {
				passUpdateLarge = true
			} else {
				t.Fatalf("bad update content: %q", event.Status.Content)
			}
		case *ConversationEvent:
			passNotification = true
			if event.Conversation.ID != "819516" {
				t.Fatalf("want %q but %q", "819516", event.Conversation.ID)
			}
		case *NotificationEvent:
			passNotification = true
			if event.Notification.Type != "mention" {
				t.Fatalf("want %q but %q", "mention", event.Notification.Type)
			}
		case *DeleteEvent:
			passDelete = true
			if event.ID != "1234567" {
				t.Fatalf("want %q but %q", "1234567", event.ID)
			}
		case *ErrorEvent:
			passError = true
			if event.Err == nil {
				t.Fatalf("should be fail: %v", event.Err)
			}
		}
	}
	if !passUpdate || !passUpdateLarge || !passNotification || !passDelete || !passError {
		t.Fatalf("have not passed through somewhere: "+
			"update: %t, update (large): %t, notification: %t, delete: %t, error: %t",
			passUpdate, passUpdateLarge, passNotification, passDelete, passError)
	}
	wg.Wait()
	err := <-errs
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
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
	var passUpdate bool
	for e := range q {
		switch event := e.(type) {
		case *ErrorEvent:
			if event.Err != nil && !errors.Is(event.Err, context.Canceled) {
				t.Fatalf("should be fail: %v", event.Err)
			}
		case *UpdateEvent:
			cnt++
			passUpdate = true
			if event.Status.Content != "foo" {
				t.Fatalf("want %q but %q", "foo", event.Status.Content)
			}
		case *UpdateEditEvent:
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
	if !passUpdate {
		t.Fatalf("have not passed through somewhere: update %t", passUpdate)
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
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(q)
		c.doStreaming(req, q)
		if err != nil {
			t.Errorf("should not be fail: %v", err)
		}
	}()
	var passError bool
	for e := range q {
		if event, ok := e.(*ErrorEvent); ok {
			passError = true
			if event.Err == nil {
				t.Fatalf("should be fail: %v", event.Err)
			}
		}
	}
	if !passError {
		t.Fatalf("have not passed through: error %t", passError)
	}
	wg.Wait()
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

func TestStreamingList(t *testing.T) {
	var isEnd bool
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isEnd {
			return
		} else if r.URL.Path != "/api/v1/streaming/list" {
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
	q, err := client.StreamingList(ctx, "1")
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

func TestStreamingDirect(t *testing.T) {
	var isEnd bool
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isEnd {
			return
		} else if r.URL.Path != "/api/v1/streaming/direct" {
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
	q, err := client.StreamingDirect(ctx)
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
