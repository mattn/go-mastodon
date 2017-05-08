package mastodon

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetNotifications(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/notifications":
			fmt.Fprintln(w, `[{"id": 122, "action_taken": false}, {"id": 123, "action_taken": true}]`)
			return
		case "/api/v1/notifications/123":
			fmt.Fprintln(w, `{"id": 123, "action_taken": true}`)
			return
		case "/api/v1/notifications/clear":
			fmt.Fprintln(w, `{}`)
			return
		}
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
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
	if ns[0].ID != 122 {
		t.Fatalf("want %v but %v", 122, ns[0].ID)
	}
	if ns[1].ID != 123 {
		t.Fatalf("want %v but %v", 123, ns[1].ID)
	}
	n, err := client.GetNotification(context.Background(), 123)
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if n.ID != 123 {
		t.Fatalf("want %v but %v", 123, n.ID)
	}
	err = client.ClearNotifications(context.Background())
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
}
