package mastodon

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAccount(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/accounts/1234567" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		fmt.Fprintln(w, `{"Username": "zzz"}`)
		return
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	a, err := client.GetAccount(context.Background(), 1)
	if err == nil {
		t.Fatalf("should not be fail: %v", err)
	}
	a, err = client.GetAccount(context.Background(), 1234567)
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if a.Username != "zzz" {
		t.Fatalf("want %q but %q", "zzz", a.Username)
	}
}

func TestGetAccountCurrentUser(t *testing.T) {
	canErr := true
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if canErr {
			canErr = false
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		fmt.Fprintln(w, `{"Username": "zzz"}`)
		return
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	_, err := client.GetAccountCurrentUser(context.Background())
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	a, err := client.GetAccountCurrentUser(context.Background())
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if a.Username != "zzz" {
		t.Fatalf("want %q but %q", "zzz", a.Username)
	}
}

func TestAccountUpdate(t *testing.T) {
	canErr := true
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if canErr {
			canErr = false
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		fmt.Fprintln(w, `{"Username": "zzz"}`)
		return
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	_, err := client.AccountUpdate(context.Background(), &Profile{})
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	a, err := client.AccountUpdate(context.Background(), &Profile{
		DisplayName: String("display_name"),
		Note:        String("note"),
		Avatar:      "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAUoAAADrCAYAAAA...",
		Header:      "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAUoAAADrCAYAAAA...",
	})
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if a.Username != "zzz" {
		t.Fatalf("want %q but %q", "zzz", a.Username)
	}
}

func TestGetAccountStatuses(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/accounts/1234567/statuses" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		fmt.Fprintln(w, `[{"Content": "foo"}, {"Content": "bar"}]`)
		return
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	_, err := client.GetAccountStatuses(context.Background(), 123)
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	ss, err := client.GetAccountStatuses(context.Background(), 1234567)
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if ss[0].Content != "foo" {
		t.Fatalf("want %q but %q", "foo", ss[0].Content)
	}
	if ss[1].Content != "bar" {
		t.Fatalf("want %q but %q", "bar", ss[1].Content)
	}
}

func TestGetBlocks(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `[{"Username": "foo"}, {"Username": "bar"}]`)
		return
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	bl, err := client.GetBlocks(context.Background())
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if len(bl) != 2 {
		t.Fatalf("result should be two: %d", len(bl))
	}
	if bl[0].Username != "foo" {
		t.Fatalf("want %q but %q", "foo", bl[0].Username)
	}
	if bl[1].Username != "bar" {
		t.Fatalf("want %q but %q", "bar", bl[1].Username)
	}
}

func TestAccountFollow(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/accounts/1234567/follow" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		fmt.Fprintln(w, `{"id":1234567,"following":true}`)
		return
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	rel, err := client.AccountFollow(context.Background(), 123)
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	rel, err = client.AccountFollow(context.Background(), 1234567)
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if rel.ID != 1234567 {
		t.Fatalf("want %d but %d", 1234567, rel.ID)
	}
	if !rel.Following {
		t.Fatalf("want %t but %t", true, rel.Following)
	}
}

func TestAccountUnfollow(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/accounts/1234567/unfollow" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		fmt.Fprintln(w, `{"id":1234567,"following":false}`)
		return
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	rel, err := client.AccountUnfollow(context.Background(), 123)
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	rel, err = client.AccountUnfollow(context.Background(), 1234567)
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if rel.ID != 1234567 {
		t.Fatalf("want %d but %d", 1234567, rel.ID)
	}
	if rel.Following {
		t.Fatalf("want %t but %t", false, rel.Following)
	}
}

func TestGetFollowRequests(t *testing.T) {
	canErr := true
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if canErr {
			canErr = false
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		fmt.Fprintln(w, `[{"Username": "foo"}, {"Username": "bar"}]`)
		return
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	_, err := client.GetFollowRequests(context.Background())
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	fReqs, err := client.GetFollowRequests(context.Background())
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if len(fReqs) != 2 {
		t.Fatalf("result should be two: %d", len(fReqs))
	}
	if fReqs[0].Username != "foo" {
		t.Fatalf("want %q but %q", "foo", fReqs[0].Username)
	}
	if fReqs[1].Username != "bar" {
		t.Fatalf("want %q but %q", "bar", fReqs[1].Username)
	}
}

func TestFollowRequestAuthorize(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/follow_requests/1234567/authorize" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		}
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	err := client.FollowRequestAuthorize(context.Background(), 123)
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	err = client.FollowRequestAuthorize(context.Background(), 1234567)
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
}

func TestFollowRequestReject(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/follow_requests/1234567/reject" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		}
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	err := client.FollowRequestReject(context.Background(), 123)
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	err = client.FollowRequestReject(context.Background(), 1234567)
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
}

func TestGetMutes(t *testing.T) {
	canErr := true
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if canErr {
			canErr = false
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		fmt.Fprintln(w, `[{"Username": "foo"}, {"Username": "bar"}]`)
		return
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	_, err := client.GetMutes(context.Background())
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	mutes, err := client.GetMutes(context.Background())
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if len(mutes) != 2 {
		t.Fatalf("result should be two: %d", len(mutes))
	}
	if mutes[0].Username != "foo" {
		t.Fatalf("want %q but %q", "foo", mutes[0].Username)
	}
	if mutes[1].Username != "bar" {
		t.Fatalf("want %q but %q", "bar", mutes[1].Username)
	}
}
