package mastodon

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetReports(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/reports" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		fmt.Fprintln(w, `[{"id": 122, "action_taken": false}, {"id": 123, "action_taken": true}]`)
		return
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	rs, err := client.GetReports(context.Background())
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if len(rs) != 2 {
		t.Fatalf("result should be two: %d", len(rs))
	}
	if rs[0].ID != 122 {
		t.Fatalf("want %v but %v", 122, rs[0].ID)
	}
	if rs[1].ID != 123 {
		t.Fatalf("want %v but %v", 123, rs[1].ID)
	}
}

func TestReport(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/reports" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		if r.FormValue("account_id") != "122" && r.FormValue("account_id") != "123" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		if r.FormValue("account_id") == "122" {
			fmt.Fprintln(w, `{"id": 1234, "action_taken": false}`)
		} else {
			fmt.Fprintln(w, `{"id": 1234, "action_taken": true}`)
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
	rp, err := client.Report(context.Background(), 121, nil, "")
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	rp, err = client.Report(context.Background(), 122, nil, "")
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if rp.ID != 1234 {
		t.Fatalf("want %v but %v", 1234, rp.ID)
	}
	if rp.ActionTaken {
		t.Fatalf("want %v but %v", true, rp.ActionTaken)
	}
	rp, err = client.Report(context.Background(), 123, []int64{567}, "")
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if rp.ID != 1234 {
		t.Fatalf("want %v but %v", 1234, rp.ID)
	}
	if !rp.ActionTaken {
		t.Fatalf("want %v but %v", false, rp.ActionTaken)
	}
}
