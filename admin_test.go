package mastodon

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAdminAccounts(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/admin/accounts" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		fmt.Fprintln(w, `[{"id": "1", "username": "foo", "email": "foo@example.com", "account": {"acct": "foo"}}, {"id": "2", "username": "bar", "suspended": true}]`)
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:      ts.URL,
		AccessToken: "zoo",
	})
	accounts, err := client.GetAdminAccounts(context.Background(), nil)
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if len(accounts) != 2 {
		t.Fatalf("result should be two: %d", len(accounts))
	}
	if accounts[0].Email != "foo@example.com" {
		t.Fatalf("want %q but %q", "foo@example.com", accounts[0].Email)
	}
	if accounts[0].Account.Acct != "foo" {
		t.Fatalf("want %q but %q", "foo", accounts[0].Account.Acct)
	}
	if !accounts[1].Suspended {
		t.Fatalf("want %v but %v", true, accounts[1].Suspended)
	}
}

func TestGetAdminAccount(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/admin/accounts/1234567" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		fmt.Fprintln(w, `{"id": "1234567", "username": "foo", "confirmed": true}`)
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:      ts.URL,
		AccessToken: "zoo",
	})
	_, err := client.GetAdminAccount(context.Background(), "123")
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	account, err := client.GetAdminAccount(context.Background(), "1234567")
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if account.Username != "foo" {
		t.Fatalf("want %q but %q", "foo", account.Username)
	}
	if !account.Confirmed {
		t.Fatalf("want %v but %v", true, account.Confirmed)
	}
}

func TestAdminAccountPerformAction(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/admin/accounts/1234567/action" || r.Method != http.MethodPost {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if r.PostForm.Get("type") != "silence" || r.PostForm.Get("text") != "spam" {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		fmt.Fprintln(w, `{}`)
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:      ts.URL,
		AccessToken: "zoo",
	})
	err := client.AdminAccountPerformAction(context.Background(), "1234567", &AdminAccountAction{
		Type: "silence",
		Text: "spam",
	})
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
}

func TestAdminAccountActions(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		switch r.URL.Path {
		case "/api/v1/admin/accounts/123/approve",
			"/api/v1/admin/accounts/123/reject",
			"/api/v1/admin/accounts/123/enable",
			"/api/v1/admin/accounts/123/unsensitive",
			"/api/v1/admin/accounts/123/unsilence",
			"/api/v1/admin/accounts/123/unsuspend":
			fmt.Fprintln(w, `{"id": "123", "username": "foo"}`)
			return
		}
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:      ts.URL,
		AccessToken: "zoo",
	})
	for _, f := range []func(context.Context, ID) (*AdminAccount, error){
		client.AdminAccountApprove,
		client.AdminAccountReject,
		client.AdminAccountEnable,
		client.AdminAccountUnsensitive,
		client.AdminAccountUnsilence,
		client.AdminAccountUnsuspend,
	} {
		account, err := f(context.Background(), "123")
		if err != nil {
			t.Fatalf("should not be fail: %v", err)
		}
		if account.ID != "123" {
			t.Fatalf("want %q but %q", "123", account.ID)
		}
	}
}

func TestGetAdminReports(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/admin/reports" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		fmt.Fprintln(w, `[{"id": "1", "comment": "spam", "target_account": {"id": "2", "username": "bar"}, "statuses": [{"content": "zzz"}]}]`)
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:      ts.URL,
		AccessToken: "zoo",
	})
	reports, err := client.GetAdminReports(context.Background(), nil)
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if len(reports) != 1 {
		t.Fatalf("result should be one: %d", len(reports))
	}
	if reports[0].Comment != "spam" {
		t.Fatalf("want %q but %q", "spam", reports[0].Comment)
	}
	if reports[0].TargetAccount.Username != "bar" {
		t.Fatalf("want %q but %q", "bar", reports[0].TargetAccount.Username)
	}
	if reports[0].Statuses[0].Content != "zzz" {
		t.Fatalf("want %q but %q", "zzz", reports[0].Statuses[0].Content)
	}
}

func TestGetAdminReport(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/admin/reports/1234567" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		fmt.Fprintln(w, `{"id": "1234567", "action_taken": true}`)
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:      ts.URL,
		AccessToken: "zoo",
	})
	report, err := client.GetAdminReport(context.Background(), "1234567")
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if !report.ActionTaken {
		t.Fatalf("want %v but %v", true, report.ActionTaken)
	}
}

func TestAdminReportActions(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		switch r.URL.Path {
		case "/api/v1/admin/reports/123/assign_to_self",
			"/api/v1/admin/reports/123/unassign",
			"/api/v1/admin/reports/123/resolve",
			"/api/v1/admin/reports/123/reopen":
			fmt.Fprintln(w, `{"id": "123"}`)
			return
		}
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:      ts.URL,
		AccessToken: "zoo",
	})
	for _, f := range []func(context.Context, ID) (*AdminReport, error){
		client.AdminReportAssignToSelf,
		client.AdminReportUnassign,
		client.AdminReportResolve,
		client.AdminReportReopen,
	} {
		report, err := f(context.Background(), "123")
		if err != nil {
			t.Fatalf("should not be fail: %v", err)
		}
		if report.ID != "123" {
			t.Fatalf("want %q but %q", "123", report.ID)
		}
	}
}

func TestRevokeToken(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/oauth/revoke" || r.Method != http.MethodPost {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if r.PostForm.Get("token") != "zoo" || r.PostForm.Get("client_id") != "foo" || r.PostForm.Get("client_secret") != "bar" {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		fmt.Fprintln(w, `{}`)
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	if err := client.RevokeToken(context.Background()); err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if client.Config.AccessToken != "" {
		t.Fatalf("access token should be cleared: %q", client.Config.AccessToken)
	}
}
