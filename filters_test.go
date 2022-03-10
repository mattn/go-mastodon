package mastodon

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"
	"time"
)

func TestGetFilters(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `[{"id": "6191", "phrase": "rust", "context": ["home"], "whole_word": true, "expires_at": "2019-05-21T13:47:31.333Z", "irreversible": false}, {"id": "5580", "phrase": "@twitter.com", "context": ["home", "notifications", "public", "thread"], "whole_word": false, "expires_at": null, "irreversible": true}]`)
		return
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	d, err := time.Parse(time.RFC3339Nano, "2019-05-21T13:47:31.333Z")
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	tf := []Filter{
		{
			ID:           ID("6191"),
			Phrase:       "rust",
			Context:      []string{"home"},
			WholeWord:    true,
			ExpiresAt:    d,
			Irreversible: false,
		},
		{
			ID:           ID("5580"),
			Phrase:       "@twitter.com",
			Context:      []string{"notifications", "home", "thread", "public"},
			WholeWord:    false,
			ExpiresAt:    time.Time{},
			Irreversible: true,
		},
	}

	filters, err := client.GetFilters(context.Background())
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if len(filters) != 2 {
		t.Fatalf("result should be two: %d", len(filters))
	}
	for i, f := range tf {
		if filters[i].ID != f.ID {
			t.Fatalf("want %q but %q", string(f.ID), filters[i].ID)
		}
		if filters[i].Phrase != f.Phrase {
			t.Fatalf("want %q but %q", f.Phrase, filters[i].Phrase)
		}
		sort.Strings(filters[i].Context)
		sort.Strings(f.Context)
		if strings.Join(filters[i].Context, ", ") != strings.Join(f.Context, ", ") {
			t.Fatalf("want %q but %q", f.Context, filters[i].Context)
		}
		if filters[i].ExpiresAt != f.ExpiresAt {
			t.Fatalf("want %q but %q", f.ExpiresAt, filters[i].ExpiresAt)
		}
		if filters[i].WholeWord != f.WholeWord {
			t.Fatalf("want %t but %t", f.WholeWord, filters[i].WholeWord)
		}
		if filters[i].Irreversible != f.Irreversible {
			t.Fatalf("want %t but %t", f.Irreversible, filters[i].Irreversible)
		}
	}
}

func TestGetFilter(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/filters/1" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		fmt.Fprintln(w, `{"id": "1", "phrase": "rust", "context": ["home"], "whole_word": true, "expires_at": "2019-05-21T13:47:31.333Z", "irreversible": false}`)
		return
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	_, err := client.GetFilter(context.Background(), "2")
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	d, err := time.Parse(time.RFC3339Nano, "2019-05-21T13:47:31.333Z")
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	tf := Filter{
		ID:           ID("1"),
		Phrase:       "rust",
		Context:      []string{"home"},
		WholeWord:    true,
		ExpiresAt:    d,
		Irreversible: false,
	}
	filter, err := client.GetFilter(context.Background(), "1")
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if filter.ID != tf.ID {
		t.Fatalf("want %q but %q", string(tf.ID), filter.ID)
	}
	if filter.Phrase != tf.Phrase {
		t.Fatalf("want %q but %q", tf.Phrase, filter.Phrase)
	}
	sort.Strings(filter.Context)
	sort.Strings(tf.Context)
	if strings.Join(filter.Context, ", ") != strings.Join(tf.Context, ", ") {
		t.Fatalf("want %q but %q", tf.Context, filter.Context)
	}
	if filter.ExpiresAt != tf.ExpiresAt {
		t.Fatalf("want %q but %q", tf.ExpiresAt, filter.ExpiresAt)
	}
	if filter.WholeWord != tf.WholeWord {
		t.Fatalf("want %t but %t", tf.WholeWord, filter.WholeWord)
	}
	if filter.Irreversible != tf.Irreversible {
		t.Fatalf("want %t but %t", tf.Irreversible, filter.Irreversible)
	}
}

func TestCreateFilter(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.PostFormValue("phrase") != "rust" && r.PostFormValue("phrase") != "@twitter.com" {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if r.PostFormValue("phrase") == "rust" {
			fmt.Fprintln(w, `{"id": "1", "phrase": "rust", "context": ["home"], "whole_word": true, "expires_at": "2019-05-21T13:47:31.333Z", "irreversible": true}`)
			return
		} else {
			fmt.Fprintln(w, `{"id": "2", "phrase": "@twitter.com", "context": ["home", "notifications", "public", "thread"], "whole_word": false, "expires_at": null, "irreversible": false}`)
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
	_, err := client.CreateFilter(context.Background(), nil)
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	_, err = client.CreateFilter(context.Background(), &Filter{Context: []string{"home"}})
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	_, err = client.CreateFilter(context.Background(), &Filter{Phrase: "rust"})
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	_, err = client.CreateFilter(context.Background(), &Filter{Phrase: "Test", Context: []string{"home"}})
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}

	d, err := time.Parse(time.RFC3339Nano, "2019-05-21T13:47:31.333Z")
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	tf := []Filter{
		{
			ID:           ID("1"),
			Phrase:       "rust",
			Context:      []string{"home"},
			WholeWord:    true,
			ExpiresAt:    d,
			Irreversible: true,
		},
		{
			ID:           ID("2"),
			Phrase:       "@twitter.com",
			Context:      []string{"notifications", "home", "thread", "public"},
			WholeWord:    false,
			ExpiresAt:    time.Time{},
			Irreversible: false,
		},
	}
	for _, f := range tf {
		filter, err := client.CreateFilter(context.Background(), &f)
		if err != nil {
			t.Fatalf("should not be fail: %v", err)
		}
		if filter.ID != f.ID {
			t.Fatalf("want %q but %q", string(f.ID), filter.ID)
		}
		if filter.Phrase != f.Phrase {
			t.Fatalf("want %q but %q", f.Phrase, filter.Phrase)
		}
		sort.Strings(filter.Context)
		sort.Strings(f.Context)
		if strings.Join(filter.Context, ", ") != strings.Join(f.Context, ", ") {
			t.Fatalf("want %q but %q", f.Context, filter.Context)
		}
		if filter.ExpiresAt != f.ExpiresAt {
			t.Fatalf("want %q but %q", f.ExpiresAt, filter.ExpiresAt)
		}
		if filter.WholeWord != f.WholeWord {
			t.Fatalf("want %t but %t", f.WholeWord, filter.WholeWord)
		}
		if filter.Irreversible != f.Irreversible {
			t.Fatalf("want %t but %t", f.Irreversible, filter.Irreversible)
		}
	}
}

func TestUpdateFilter(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/filters/1" {
			fmt.Fprintln(w, `{"id": "1", "phrase": "rust", "context": ["home"], "whole_word": true, "expires_at": "2019-05-21T13:47:31.333Z", "irreversible": true}`)
			return
		} else if r.URL.Path == "/api/v1/filters/2" {
			fmt.Fprintln(w, `{"id": "2", "phrase": "@twitter.com", "context": ["home", "notifications", "public", "thread"], "whole_word": false, "expires_at": null, "irreversible": false}`)
			return
		} else {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
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
	_, err := client.UpdateFilter(context.Background(), ID("1"), nil)
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	_, err = client.UpdateFilter(context.Background(), ID(""), &Filter{Phrase: ""})
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	_, err = client.UpdateFilter(context.Background(), ID("2"), &Filter{Phrase: ""})
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	_, err = client.UpdateFilter(context.Background(), ID("2"), &Filter{Phrase: "rust"})
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	_, err = client.UpdateFilter(context.Background(), ID("3"), &Filter{Phrase: "rust", Context: []string{"home"}})
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}

	d, err := time.Parse(time.RFC3339Nano, "2019-05-21T13:47:31.333Z")
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	tf := []Filter{
		{
			ID:           ID("1"),
			Phrase:       "rust",
			Context:      []string{"home"},
			WholeWord:    true,
			ExpiresAt:    d,
			Irreversible: true,
		},
		{
			ID:           ID("2"),
			Phrase:       "@twitter.com",
			Context:      []string{"notifications", "home", "thread", "public"},
			WholeWord:    false,
			ExpiresAt:    time.Time{},
			Irreversible: false,
		},
	}
	for _, f := range tf {
		filter, err := client.UpdateFilter(context.Background(), f.ID, &f)
		if err != nil {
			t.Fatalf("should not be fail: %v", err)
		}
		if filter.ID != f.ID {
			t.Fatalf("want %q but %q", string(f.ID), filter.ID)
		}
		if filter.Phrase != f.Phrase {
			t.Fatalf("want %q but %q", f.Phrase, filter.Phrase)
		}
		sort.Strings(filter.Context)
		sort.Strings(f.Context)
		if strings.Join(filter.Context, ", ") != strings.Join(f.Context, ", ") {
			t.Fatalf("want %q but %q", f.Context, filter.Context)
		}
		if filter.ExpiresAt != f.ExpiresAt {
			t.Fatalf("want %q but %q", f.ExpiresAt, filter.ExpiresAt)
		}
		if filter.WholeWord != f.WholeWord {
			t.Fatalf("want %t but %t", f.WholeWord, filter.WholeWord)
		}
		if filter.Irreversible != f.Irreversible {
			t.Fatalf("want %t but %t", f.Irreversible, filter.Irreversible)
		}
	}
}

func TestDeleteFilter(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/filters/1" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		if r.Method != "DELETE" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusMethodNotAllowed)
			return
		}
		return
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	err := client.DeleteFilter(context.Background(), "2")
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	err = client.DeleteFilter(context.Background(), "1")
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
}
