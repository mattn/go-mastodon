package mastodon

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetLists(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/lists" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		fmt.Fprintln(w, `[{"id": "1", "title": "foo"}, {"id": "2", "title": "bar"}]`)
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	lists, err := client.GetLists(context.Background())
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if len(lists) != 2 {
		t.Fatalf("result should be two: %d", len(lists))
	}
	if lists[0].Title != "foo" {
		t.Fatalf("want %q but %q", "foo", lists[0].Title)
	}
	if lists[1].Title != "bar" {
		t.Fatalf("want %q but %q", "bar", lists[1].Title)
	}
}

func TestGetAccountLists(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/accounts/1/lists" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		fmt.Fprintln(w, `[{"id": "1", "title": "foo"}, {"id": "2", "title": "bar"}]`)
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	_, err := client.GetAccountLists(context.Background(), "2")
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	lists, err := client.GetAccountLists(context.Background(), "1")
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if len(lists) != 2 {
		t.Fatalf("result should be two: %d", len(lists))
	}
	if lists[0].Title != "foo" {
		t.Fatalf("want %q but %q", "foo", lists[0].Title)
	}
	if lists[1].Title != "bar" {
		t.Fatalf("want %q but %q", "bar", lists[1].Title)
	}
}

func TestGetListAccounts(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/lists/1/accounts" {
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
	_, err := client.GetListAccounts(context.Background(), "2")
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	accounts, err := client.GetListAccounts(context.Background(), "1")
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if len(accounts) != 2 {
		t.Fatalf("result should be two: %d", len(accounts))
	}
	if accounts[0].Username != "foo" {
		t.Fatalf("want %q but %q", "foo", accounts[0].Username)
	}
	if accounts[1].Username != "bar" {
		t.Fatalf("want %q but %q", "bar", accounts[1].Username)
	}
}

func TestGetList(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/lists/1" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		fmt.Fprintln(w, `{"id": "1", "title": "foo"}`)
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	_, err := client.GetList(context.Background(), "2")
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	list, err := client.GetList(context.Background(), "1")
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if list.Title != "foo" {
		t.Fatalf("want %q but %q", "foo", list.Title)
	}
}

func TestCreateList(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.PostFormValue("title") != "foo" {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		fmt.Fprintln(w, `{"id": "1", "title": "foo"}`)
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	_, err := client.CreateList(context.Background(), "")
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	list, err := client.CreateList(context.Background(), "foo")
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if list.Title != "foo" {
		t.Fatalf("want %q but %q", "foo", list.Title)
	}
}

func TestRenameList(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/lists/1" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		if r.PostFormValue("title") != "bar" {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		fmt.Fprintln(w, `{"id": "1", "title": "bar"}`)
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	_, err := client.RenameList(context.Background(), "2", "bar")
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	list, err := client.RenameList(context.Background(), "1", "bar")
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if list.Title != "bar" {
		t.Fatalf("want %q but %q", "bar", list.Title)
	}
}

func TestDeleteList(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/lists/1" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		if r.Method != "DELETE" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusMethodNotAllowed)
			return
		}
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	err := client.DeleteList(context.Background(), "2")
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	err = client.DeleteList(context.Background(), "1")
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
}

func TestAddToList(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/lists/1/accounts" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		if r.PostFormValue("account_ids[]") != "1" {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	err := client.AddToList(context.Background(), "1", "1")
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
}

func TestRemoveFromList(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/lists/1/accounts" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		if r.Method != "DELETE" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusMethodNotAllowed)
			return
		}
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	err := client.RemoveFromList(context.Background(), "1", "1")
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
}
