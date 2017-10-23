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
		fmt.Fprintln(w, `{"username": "zzz"}`)
		return
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	_, err := client.GetAccount(context.Background(), "1")
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	a, err := client.GetAccount(context.Background(), "1234567")
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
		fmt.Fprintln(w, `{"username": "zzz"}`)
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
		fmt.Fprintln(w, `{"username": "zzz"}`)
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
		fmt.Fprintln(w, `[{"content": "foo"}, {"content": "bar"}]`)
		return
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	_, err := client.GetAccountStatuses(context.Background(), "123", nil)
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	ss, err := client.GetAccountStatuses(context.Background(), "1234567", nil)
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

func TestGetAccountFollowers(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/accounts/1234567/followers" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		fmt.Fprintln(w, `[{"username": "foo"}, {"username": "bar"}]`)
		return
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	_, err := client.GetAccountFollowers(context.Background(), "123", nil)
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	fl, err := client.GetAccountFollowers(context.Background(), "1234567", nil)
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if len(fl) != 2 {
		t.Fatalf("result should be two: %d", len(fl))
	}
	if fl[0].Username != "foo" {
		t.Fatalf("want %q but %q", "foo", fl[0].Username)
	}
	if fl[1].Username != "bar" {
		t.Fatalf("want %q but %q", "bar", fl[1].Username)
	}
}

func TestGetAccountFollowing(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/accounts/1234567/following" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		fmt.Fprintln(w, `[{"username": "foo"}, {"username": "bar"}]`)
		return
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	_, err := client.GetAccountFollowing(context.Background(), "123", nil)
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	fl, err := client.GetAccountFollowing(context.Background(), "1234567", nil)
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if len(fl) != 2 {
		t.Fatalf("result should be two: %d", len(fl))
	}
	if fl[0].Username != "foo" {
		t.Fatalf("want %q but %q", "foo", fl[0].Username)
	}
	if fl[1].Username != "bar" {
		t.Fatalf("want %q but %q", "bar", fl[1].Username)
	}
}

func TestGetBlocks(t *testing.T) {
	canErr := true
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if canErr {
			canErr = false
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		fmt.Fprintln(w, `[{"username": "foo"}, {"username": "bar"}]`)
		return
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	_, err := client.GetBlocks(context.Background(), nil)
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	bl, err := client.GetBlocks(context.Background(), nil)
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
	rel, err := client.AccountFollow(context.Background(), "123")
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	rel, err = client.AccountFollow(context.Background(), "1234567")
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if rel.ID != "1234567" {
		t.Fatalf("want %q but %q", "1234567", rel.ID)
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
	rel, err := client.AccountUnfollow(context.Background(), "123")
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	rel, err = client.AccountUnfollow(context.Background(), "1234567")
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if rel.ID != "1234567" {
		t.Fatalf("want %q but %q", "1234567", rel.ID)
	}
	if rel.Following {
		t.Fatalf("want %t but %t", false, rel.Following)
	}
}

func TestAccountBlock(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/accounts/1234567/block" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		fmt.Fprintln(w, `{"id":1234567,"blocking":true}`)
		return
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	_, err := client.AccountBlock(context.Background(), "123")
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	rel, err := client.AccountBlock(context.Background(), "1234567")
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if rel.ID != "1234567" {
		t.Fatalf("want %q but %q", "1234567", rel.ID)
	}
	if !rel.Blocking {
		t.Fatalf("want %t but %t", true, rel.Blocking)
	}
}

func TestAccountUnblock(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/accounts/1234567/unblock" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		fmt.Fprintln(w, `{"id":1234567,"blocking":false}`)
		return
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	_, err := client.AccountUnblock(context.Background(), "123")
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	rel, err := client.AccountUnblock(context.Background(), "1234567")
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if rel.ID != "1234567" {
		t.Fatalf("want %q but %q", "1234567", rel.ID)
	}
	if rel.Blocking {
		t.Fatalf("want %t but %t", false, rel.Blocking)
	}
}

func TestAccountMute(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/accounts/1234567/mute" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		fmt.Fprintln(w, `{"id":1234567,"muting":true}`)
		return
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	_, err := client.AccountMute(context.Background(), "123")
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	rel, err := client.AccountMute(context.Background(), "1234567")
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if rel.ID != "1234567" {
		t.Fatalf("want %q but %q", "1234567", rel.ID)
	}
	if !rel.Muting {
		t.Fatalf("want %t but %t", true, rel.Muting)
	}
}

func TestAccountUnmute(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/accounts/1234567/unmute" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		fmt.Fprintln(w, `{"id":1234567,"muting":false}`)
		return
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	_, err := client.AccountUnmute(context.Background(), "123")
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	rel, err := client.AccountUnmute(context.Background(), "1234567")
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if rel.ID != "1234567" {
		t.Fatalf("want %q but %q", "1234567", rel.ID)
	}
	if rel.Muting {
		t.Fatalf("want %t but %t", false, rel.Muting)
	}
}

func TestGetAccountRelationship(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ids := r.URL.Query()["id[]"]
		if ids[0] == "1234567" && ids[1] == "8901234" {
			fmt.Fprintln(w, `[{"id":1234567},{"id":8901234}]`)
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	_, err := client.GetAccountRelationships(context.Background(), []string{"123", "456"})
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	rels, err := client.GetAccountRelationships(context.Background(), []string{"1234567", "8901234"})
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if rels[0].ID != "1234567" {
		t.Fatalf("want %q but %q", "1234567", rels[0].ID)
	}
	if rels[1].ID != "8901234" {
		t.Fatalf("want %q but %q", "8901234", rels[1].ID)
	}
}

func TestAccountsSearch(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query()["q"][0] != "foo" {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		fmt.Fprintln(w, `[{"username": "foobar"}, {"username": "barfoo"}]`)
		return
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	_, err := client.AccountsSearch(context.Background(), "zzz", 2)
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	res, err := client.AccountsSearch(context.Background(), "foo", 2)
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if len(res) != 2 {
		t.Fatalf("result should be two: %d", len(res))
	}
	if res[0].Username != "foobar" {
		t.Fatalf("want %q but %q", "foobar", res[0].Username)
	}
	if res[1].Username != "barfoo" {
		t.Fatalf("want %q but %q", "barfoo", res[1].Username)
	}
}

func TestFollowRemoteUser(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.PostFormValue("uri") != "foo@success.social" {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		fmt.Fprintln(w, `{"username": "zzz"}`)
		return
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	_, err := client.FollowRemoteUser(context.Background(), "foo@fail.social")
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	ru, err := client.FollowRemoteUser(context.Background(), "foo@success.social")
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if ru.Username != "zzz" {
		t.Fatalf("want %q but %q", "zzz", ru.Username)
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
		fmt.Fprintln(w, `[{"username": "foo"}, {"username": "bar"}]`)
		return
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	_, err := client.GetFollowRequests(context.Background(), nil)
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	fReqs, err := client.GetFollowRequests(context.Background(), nil)
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
	err := client.FollowRequestAuthorize(context.Background(), "123")
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	err = client.FollowRequestAuthorize(context.Background(), "1234567")
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
	err := client.FollowRequestReject(context.Background(), "123")
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	err = client.FollowRequestReject(context.Background(), "1234567")
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
		fmt.Fprintln(w, `[{"username": "foo"}, {"username": "bar"}]`)
		return
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	_, err := client.GetMutes(context.Background(), nil)
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	mutes, err := client.GetMutes(context.Background(), nil)
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
