package mastodon

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestStreamingWSPublic(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(wsMock))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	ctx, cancel := context.WithCancel(context.Background())
	q, err := client.StreamingWSPublic(ctx)
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}

	wsTest(t, q, cancel)
}

func TestStreamingWSPublicLocal(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(wsMock))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	ctx, cancel := context.WithCancel(context.Background())
	q, err := client.StreamingWSPublicLocal(ctx)
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}

	wsTest(t, q, cancel)
}

func TestStreamingWSUser(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(wsMock))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	ctx, cancel := context.WithCancel(context.Background())
	q, err := client.StreamingWSUser(ctx)
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}

	wsTest(t, q, cancel)
}

func TestStreamingWSHashtag(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(wsMock))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	ctx, cancel := context.WithCancel(context.Background())
	q, err := client.StreamingWSHashtag(ctx, "zzz")
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}

	wsTest(t, q, cancel)
}

func TestStreamingWSHashtagLocal(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(wsMock))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	ctx, cancel := context.WithCancel(context.Background())
	q, err := client.StreamingWSHashtagLocal(ctx, "zzz")
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}

	wsTest(t, q, cancel)
}

func wsMock(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/v1/streaming" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	u := websocket.Upgrader{}
	conn, err := u.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	err = conn.WriteMessage(websocket.TextMessage,
		[]byte(`{"event":"update","payload":"{\"content\":\"foo\"}"}`))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = conn.WriteMessage(websocket.TextMessage,
		[]byte(`{"event":"update","payload":"{\"content\":\"bar\"}"}`))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	time.Sleep(10 * time.Second)
}

func wsTest(t *testing.T, q chan Event, cancel func()) {
	time.AfterFunc(time.Second, func() {
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

func TestDial(t *testing.T) {
	canErr := true
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if canErr {
			canErr = false
			http.Redirect(w, r, ":", http.StatusMovedPermanently)
			return
		}

		http.Redirect(w, r, "http://www.example.com/", http.StatusMovedPermanently)
	}))
	defer ts.Close()

	client := NewClient(&Config{})
	_, _, err := client.dial(":")
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}

	_, rawurl, err := client.dial("ws://" + ts.Listener.Addr().String())
	if err == nil {
		t.Fatalf("should not be fail: %v", err)
	}

	_, rawurl, err = client.dial("ws://" + ts.Listener.Addr().String())
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if rawurl != "ws://www.example.com/" {
		t.Fatalf("want %q but %q", "ws://www.example.com/", rawurl)
	}
}

func TestChangeWebSocketScheme(t *testing.T) {
	_, err := changeWebSocketScheme(":")
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}

	u, err := changeWebSocketScheme("http://example.com/")
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if u.Scheme != "ws" {
		t.Fatalf("want %q but %q", "ws", u.Scheme)
	}

	u, err = changeWebSocketScheme("https://example.com/")
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if u.Scheme != "wss" {
		t.Fatalf("want %q but %q", "wss", u.Scheme)
	}
}
