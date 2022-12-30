package mastodon

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetNotifications(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/notifications":
			if r.URL.Query().Get("exclude_types[]") == "follow" {
				fmt.Fprintln(w, `[{"id": 321, "action_taken": true}]`)
			} else {
				fmt.Fprintln(w, `[{"id": 122, "action_taken": false}, {"id": 123, "action_taken": true}]`)
			}
			return
		case "/api/v1/notifications/123":
			fmt.Fprintln(w, `{"id": 123, "action_taken": true}`)
			return
		case "/api/v1/notifications/clear":
			fmt.Fprintln(w, `{}`)
			return
		case "/api/v1/notifications/123/dismiss":
			fmt.Fprintln(w, `{}`)
			return
		}
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	ns, err := client.GetNotifications(context.Background(), nil)
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if len(ns) != 2 {
		t.Fatalf("result should be two: %d", len(ns))
	}
	if ns[0].ID != "122" {
		t.Fatalf("want %v but %v", "122", ns[0].ID)
	}
	if ns[1].ID != "123" {
		t.Fatalf("want %v but %v", "123", ns[1].ID)
	}
	nse, err := client.GetNotificationsExclude(context.Background(), &[]string{"follow"}, nil)
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if len(nse) != 1 {
		t.Fatalf("result should be one: %d", len(nse))
	}
	if nse[0].ID != "321" {
		t.Fatalf("want %v but %v", "321", nse[0].ID)
	}
	n, err := client.GetNotification(context.Background(), "123")
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if n.ID != "123" {
		t.Fatalf("want %v but %v", "123", n.ID)
	}
	err = client.ClearNotifications(context.Background())
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	err = client.DismissNotification(context.Background(), "123")
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
}

func TestPushSubscription(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/push/subscription":
			fmt.Fprintln(w, ` {"id":1,"endpoint":"https://example.org","alerts":{"follow":true,"favourite":"true","reblog":"true","mention":"true"},"server_key":"foobar"}`)
			return
		}
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})

	enabled := new(Sbool)
	*enabled = true
	alerts := PushAlerts{Follow: enabled, Favourite: enabled, Reblog: enabled, Mention: enabled}

	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	shared := make([]byte, 16)
	_, err = rand.Read(shared)
	if err != nil {
		t.Fatal(err)
	}

	testSub := func(sub *PushSubscription, err error) {
		if err != nil {
			t.Fatalf("should not be fail: %v", err)
		}
		if sub.ID != "1" {
			t.Fatalf("want %v but %v", "1", sub.ID)
		}
		if sub.Endpoint != "https://example.org" {
			t.Fatalf("want %v but %v", "https://example.org", sub.Endpoint)
		}
		if sub.ServerKey != "foobar" {
			t.Fatalf("want %v but %v", "foobar", sub.ServerKey)
		}
		if *sub.Alerts.Favourite != true {
			t.Fatalf("want %v but %v", true, *sub.Alerts.Favourite)
		}
		if *sub.Alerts.Mention != true {
			t.Fatalf("want %v but %v", true, *sub.Alerts.Mention)
		}
		if *sub.Alerts.Reblog != true {
			t.Fatalf("want %v but %v", true, *sub.Alerts.Reblog)
		}
		if *sub.Alerts.Follow != true {
			t.Fatalf("want %v but %v", true, *sub.Alerts.Follow)
		}
	}

	sub, err := client.AddPushSubscription(context.Background(), "http://example.org", priv.PublicKey, shared, alerts)
	testSub(sub, err)

	sub, err = client.GetPushSubscription(context.Background())
	testSub(sub, err)

	sub, err = client.UpdatePushSubscription(context.Background(), &alerts)
	testSub(sub, err)

	err = client.RemovePushSubscription(context.Background())
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
}
