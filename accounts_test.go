package mastodon

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetAccount(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/accounts/1234567" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		fmt.Fprintln(w, `{"username": "zzz"}`)
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

func TestAccountLookup(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/accounts/lookup" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		acct := r.URL.Query().Get("acct")
		if acct != "foo@bar" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		fmt.Fprintln(w, `{"username": "foo@bar"}`)
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	_, err := client.AccountLookup(context.Background(), "123")
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	a, err := client.AccountLookup(context.Background(), "foo@bar")
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if a.Username != "foo@bar" {
		t.Fatalf("want %q but %q", "foo@bar", a.Username)
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
	tbool := true
	fields := []Field{{"foo", "bar", time.Time{}}, {"dum", "baz", time.Time{}}}
	source := AccountSource{Language: String("de"), Privacy: String("public"), Sensitive: &tbool}
	a, err := client.AccountUpdate(context.Background(), &Profile{
		DisplayName: String("display_name"),
		Note:        String("note"),
		Locked:      &tbool,
		Fields:      &fields,
		Source:      &source,
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

func TestGetAccountPinnedStatuses(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/accounts/1234567/statuses" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		pinned := r.URL.Query().Get("pinned")
		if pinned != "true" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		fmt.Fprintln(w, `[{"content": "foo"}, {"content": "bar"}]`)
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	_, err := client.GetAccountPinnedStatuses(context.Background(), "123")
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	ss, err := client.GetAccountPinnedStatuses(context.Background(), "1234567")
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

func TestGetEndorsements(t *testing.T) {
	canErr := true
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if canErr {
			canErr = false
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		fmt.Fprintln(w, `[
  {
    "id": "952529",
    "username": "foo",
    "acct": "alayna@desvox.es",
    "display_name": "Alayna Desirae",
    "locked": true,
    "bot": false,
    "created_at": "2019-10-26T23:12:06.570Z",
    "note": "experiencing ________ difficulties<br>22y/o INFP in Oklahoma",
    "url": "https://desvox.es/users/alayna",
    "avatar": "https://files.mastodon.social/accounts/avatars/000/952/529/original/6534122046d050d5.png",
    "avatar_static": "https://files.mastodon.social/accounts/avatars/000/952/529/original/6534122046d050d5.png",
    "header": "https://files.mastodon.social/accounts/headers/000/952/529/original/496f1f817e042ade.png",
    "header_static": "https://files.mastodon.social/accounts/headers/000/952/529/original/496f1f817e042ade.png",
    "followers_count": 0,
    "following_count": 0,
    "statuses_count": 955,
    "last_status_at": "2019-11-23T07:05:50.682Z",
    "emojis": [],
    "fields": []
  },
  {
    "id": "832844",
    "username": "bar",
    "acct": "a9@broadcast.wolfgirl.engineering",
    "display_name": "vivienne :collar: ",
    "locked": true,
    "bot": false,
    "created_at": "2019-06-12T18:55:12.053Z",
    "note": "borderline nsfw, considered a schedule I drug by nixon<br>waiting for the year of the illumos desktop",
    "url": "https://broadcast.wolfgirl.engineering/users/a9",
    "avatar": "https://files.mastodon.social/accounts/avatars/000/832/844/original/ae1de0b8fb63d1c6.png",
    "avatar_static": "https://files.mastodon.social/accounts/avatars/000/832/844/original/ae1de0b8fb63d1c6.png",
    "header": "https://files.mastodon.social/accounts/headers/000/832/844/original/5088e4a16e6d8736.png",
    "header_static": "https://files.mastodon.social/accounts/headers/000/832/844/original/5088e4a16e6d8736.png",
    "followers_count": 43,
    "following_count": 67,
    "statuses_count": 5906,
    "last_status_at": "2019-11-23T05:23:47.911Z",
    "emojis": [
      {
        "shortcode": "collar",
        "url": "https://files.mastodon.social/custom_emojis/images/000/106/920/original/80953b9cd96ec4dc.png",
        "static_url": "https://files.mastodon.social/custom_emojis/images/000/106/920/static/80953b9cd96ec4dc.png",
        "visible_in_picker": true
      }
    ],
    "fields": []
  }
]`)
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	_, err := client.GetEndorsements(context.Background(), nil)
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	endorsements, err := client.GetEndorsements(context.Background(), nil)
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if len(endorsements) != 2 {
		t.Fatalf("result should be two: %d", len(endorsements))
	}
	if endorsements[0].Username != "foo" {
		t.Fatalf("want %q but %q", "foo", endorsements[0].Username)
	}
	if endorsements[1].Username != "bar" {
		t.Fatalf("want %q but %q", "bar", endorsements[1].Username)
	}
}

func TestAccountFollow(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/accounts/1234567/follow" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		fmt.Fprintln(w, `{"id":1234567,"following":true}`)
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	_, err := client.AccountFollow(context.Background(), "123")
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	rel, err := client.AccountFollow(context.Background(), "1234567")
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
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	_, err := client.AccountUnfollow(context.Background(), "123")
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	rel, err := client.AccountUnfollow(context.Background(), "1234567")
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
func TestAccountsSearchResolve(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query()["q"][0] != "foo" {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if r.FormValue("resolve") != "true" {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		fmt.Fprintln(w, `[{"username": "foobar"}, {"username": "barfoo"}]`)
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	_, err := client.AccountsSearchResolve(context.Background(), "zzz", 2, false)
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	res, err := client.AccountsSearchResolve(context.Background(), "foo", 2, true)
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
func TestGetFollowedTags(t *testing.T) {
	t.Parallel()
	canErr := true
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if canErr {
			canErr = false
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		fmt.Fprintln(w, `[
  {
    "name": "Test1",
    "url": "http://mastodon.example/tags/test1",
    "history": [
      {
        "day": "1668211200",
        "accounts": "0",
        "uses": "0"
      },
      {
        "day": "1668124800",
        "accounts": "0",
        "uses": "0"
      },
      {
        "day": "1668038400",
        "accounts": "0",
        "uses": "0"
      }
    ],
    "following": true
  },
  {
    "name": "Test2",
    "url": "http://mastodon.example/tags/test2",
    "history": [
      {
        "day": "1668211200",
        "accounts": "0",
        "uses": "0"
      }
    ],
    "following": true
  }
]`)
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	_, err := client.GetFollowedTags(context.Background(), nil)
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	followedTags, err := client.GetFollowedTags(context.Background(), nil)
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if len(followedTags) != 2 {
		t.Fatalf("result should be two: %d", len(followedTags))
	}
	if followedTags[0].Name != "Test1" {
		t.Fatalf("want %q but %q", "Test1", followedTags[0].Name)
	}
	if followedTags[0].URL != "http://mastodon.example/tags/test1" {
		t.Fatalf("want %q but got %q", "http://mastodon.example/tags/test1", followedTags[0].URL)
	}
	if !followedTags[0].Following {
		t.Fatalf("want following, but got false")
	}
	if len(followedTags[0].History) != 3 {
		t.Fatalf("expecting first tag history length to be %d but got %d", 3, len(followedTags[0].History))
	}
	if followedTags[1].Name != "Test2" {
		t.Fatalf("want %q but %q", "Test2", followedTags[1].Name)
	}
	if followedTags[1].URL != "http://mastodon.example/tags/test2" {
		t.Fatalf("want %q but got %q", "http://mastodon.example/tags/test2", followedTags[1].URL)
	}
	if !followedTags[1].Following {
		t.Fatalf("want following, but got false")
	}
	if len(followedTags[1].History) != 1 {
		t.Fatalf("expecting first tag history length to be %d but got %d", 1, len(followedTags[1].History))
	}
}
