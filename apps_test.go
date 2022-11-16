package mastodon

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRegisterApp(t *testing.T) {
	isNotJSON := true
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		} else if r.URL.Path != "/api/v1/apps" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		} else if r.FormValue("redirect_uris") != "urn:ietf:wg:oauth:2.0:oob" {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		} else if isNotJSON {
			isNotJSON = false
			fmt.Fprintln(w, `<html><head><title>Apps</title></head></html>`)
			return
		}
		fmt.Fprintln(w, `{"id": 123, "client_id": "foo", "client_secret": "bar"}`)
	}))
	defer ts.Close()

	// Status not ok.
	_, err := RegisterApp(context.Background(), &AppConfig{
		Server:       ts.URL,
		RedirectURIs: "/",
	})
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}

	// Error in url.Parse
	_, err = RegisterApp(context.Background(), &AppConfig{
		Server: ":",
	})
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}

	// Error in json.NewDecoder
	_, err = RegisterApp(context.Background(), &AppConfig{
		Server: ts.URL,
	})
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}

	// Success.
	app, err := RegisterApp(context.Background(), &AppConfig{
		Server: ts.URL,
		Scopes: "read write follow",
	})
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if string(app.ID) != "123" {
		t.Fatalf("want %q but %q", "bar", app.ClientSecret)
	}
	if app.ClientID != "foo" {
		t.Fatalf("want %q but %q", "foo", app.ClientID)
	}
	if app.ClientSecret != "bar" {
		t.Fatalf("want %q but %q", "bar", app.ClientSecret)
	}
}

func TestRegisterAppWithCancel(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(3 * time.Second)
		fmt.Fprintln(w, `{"client_id": "foo", "client_secret": "bar"}`)
	}))
	defer ts.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := RegisterApp(ctx, &AppConfig{
		Server: ts.URL,
		Scopes: "read write follow",
	})
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	if want := fmt.Sprintf("Post %q: context canceled", ts.URL+"/api/v1/apps"); want != err.Error() {
		t.Fatalf("want %q but %q", want, err.Error())
	}
}

func TestVerifyAppCredentials(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer zoo" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		if r.URL.Path != "/api/v1/apps/verify_credentials" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		fmt.Fprintln(w, `{"name":"zzz","website":"yyy","vapid_key":"xxx"}`)
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zip",
	})
	_, err := client.VerifyAppCredentials(context.Background())
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}

	client = NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	a, err := client.VerifyAppCredentials(context.Background())
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if a.Name != "zzz" {
		t.Fatalf("want %q but %q", "zzz", a.Name)
	}
	if a.Website != "yyy" {
		t.Fatalf("want %q but %q", "yyy", a.Name)
	}
	if a.VapidKey != "xxx" {
		t.Fatalf("want %q but %q", "xxx", a.Name)
	}
}
