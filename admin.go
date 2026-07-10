package mastodon

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// AdminAccount holds the admin-level view of an account.
type AdminAccount struct {
	ID            ID        `json:"id"`
	Username      string    `json:"username"`
	Domain        string    `json:"domain"`
	CreatedAt     time.Time `json:"created_at"`
	Email         string    `json:"email"`
	IP            string    `json:"ip"`
	Locale        string    `json:"locale"`
	InviteRequest string    `json:"invite_request"`
	Confirmed     bool      `json:"confirmed"`
	Approved      bool      `json:"approved"`
	Disabled      bool      `json:"disabled"`
	Silenced      bool      `json:"silenced"`
	Suspended     bool      `json:"suspended"`
	Account       *Account  `json:"account"`
}

// AdminReport holds the admin-level view of a report.
type AdminReport struct {
	ID                   ID            `json:"id"`
	ActionTaken          bool          `json:"action_taken"`
	ActionTakenAt        *time.Time    `json:"action_taken_at"`
	Category             string        `json:"category"`
	Comment              string        `json:"comment"`
	Forwarded            bool          `json:"forwarded"`
	CreatedAt            time.Time     `json:"created_at"`
	UpdatedAt            time.Time     `json:"updated_at"`
	Account              *AdminAccount `json:"account"`
	TargetAccount        *AdminAccount `json:"target_account"`
	AssignedAccount      *AdminAccount `json:"assigned_account"`
	ActionTakenByAccount *AdminAccount `json:"action_taken_by_account"`
	Statuses             []*Status     `json:"statuses"`
}

// AdminAccountAction specifies a moderation action against an account.
type AdminAccountAction struct {
	// Type is one of none, sensitive, disable, silence and suspend.
	Type                  string
	Text                  string
	ReportID              ID
	WarningPresetID       ID
	SendEmailNotification bool
}

// GetAdminAccounts returns accounts from the admin view.
func (c *Client) GetAdminAccounts(ctx context.Context, pg *Pagination) ([]*AdminAccount, error) {
	var accounts []*AdminAccount
	err := c.doAPI(ctx, http.MethodGet, "/api/v1/admin/accounts", nil, &accounts, pg)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

// GetAdminAccount returns an account from the admin view.
func (c *Client) GetAdminAccount(ctx context.Context, id ID) (*AdminAccount, error) {
	var account AdminAccount
	err := c.doAPI(ctx, http.MethodGet, fmt.Sprintf("/api/v1/admin/accounts/%s", url.PathEscape(string(id))), nil, &account, nil)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

// AdminAccountPerformAction performs a moderation action against the account.
func (c *Client) AdminAccountPerformAction(ctx context.Context, id ID, action *AdminAccountAction) error {
	params := url.Values{}
	params.Set("type", action.Type)
	if action.Text != "" {
		params.Set("text", action.Text)
	}
	if action.ReportID != "" {
		params.Set("report_id", string(action.ReportID))
	}
	if action.WarningPresetID != "" {
		params.Set("warning_preset_id", string(action.WarningPresetID))
	}
	if action.SendEmailNotification {
		params.Set("send_email_notification", strconv.FormatBool(action.SendEmailNotification))
	}
	return c.doAPI(ctx, http.MethodPost, fmt.Sprintf("/api/v1/admin/accounts/%s/action", url.PathEscape(string(id))), params, nil, nil)
}

// AdminAccountApprove approves the pending account.
func (c *Client) AdminAccountApprove(ctx context.Context, id ID) (*AdminAccount, error) {
	return c.adminAccountPost(ctx, id, "approve")
}

// AdminAccountReject rejects the pending account.
func (c *Client) AdminAccountReject(ctx context.Context, id ID) (*AdminAccount, error) {
	return c.adminAccountPost(ctx, id, "reject")
}

// AdminAccountEnable re-enables the disabled account.
func (c *Client) AdminAccountEnable(ctx context.Context, id ID) (*AdminAccount, error) {
	return c.adminAccountPost(ctx, id, "enable")
}

// AdminAccountUnsensitive removes the sensitive flag from the account.
func (c *Client) AdminAccountUnsensitive(ctx context.Context, id ID) (*AdminAccount, error) {
	return c.adminAccountPost(ctx, id, "unsensitive")
}

// AdminAccountUnsilence unsilences the account.
func (c *Client) AdminAccountUnsilence(ctx context.Context, id ID) (*AdminAccount, error) {
	return c.adminAccountPost(ctx, id, "unsilence")
}

// AdminAccountUnsuspend unsuspends the account.
func (c *Client) AdminAccountUnsuspend(ctx context.Context, id ID) (*AdminAccount, error) {
	return c.adminAccountPost(ctx, id, "unsuspend")
}

func (c *Client) adminAccountPost(ctx context.Context, id ID, action string) (*AdminAccount, error) {
	var account AdminAccount
	err := c.doAPI(ctx, http.MethodPost, fmt.Sprintf("/api/v1/admin/accounts/%s/%s", url.PathEscape(string(id)), action), nil, &account, nil)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

// GetAdminReports returns reports from the admin view.
func (c *Client) GetAdminReports(ctx context.Context, pg *Pagination) ([]*AdminReport, error) {
	var reports []*AdminReport
	err := c.doAPI(ctx, http.MethodGet, "/api/v1/admin/reports", nil, &reports, pg)
	if err != nil {
		return nil, err
	}
	return reports, nil
}

// GetAdminReport returns a report from the admin view.
func (c *Client) GetAdminReport(ctx context.Context, id ID) (*AdminReport, error) {
	var report AdminReport
	err := c.doAPI(ctx, http.MethodGet, fmt.Sprintf("/api/v1/admin/reports/%s", url.PathEscape(string(id))), nil, &report, nil)
	if err != nil {
		return nil, err
	}
	return &report, nil
}

// AdminReportAssignToSelf assigns the report to the current user.
func (c *Client) AdminReportAssignToSelf(ctx context.Context, id ID) (*AdminReport, error) {
	return c.adminReportPost(ctx, id, "assign_to_self")
}

// AdminReportUnassign unassigns the report.
func (c *Client) AdminReportUnassign(ctx context.Context, id ID) (*AdminReport, error) {
	return c.adminReportPost(ctx, id, "unassign")
}

// AdminReportResolve marks the report as resolved.
func (c *Client) AdminReportResolve(ctx context.Context, id ID) (*AdminReport, error) {
	return c.adminReportPost(ctx, id, "resolve")
}

// AdminReportReopen reopens the resolved report.
func (c *Client) AdminReportReopen(ctx context.Context, id ID) (*AdminReport, error) {
	return c.adminReportPost(ctx, id, "reopen")
}

func (c *Client) adminReportPost(ctx context.Context, id ID, action string) (*AdminReport, error) {
	var report AdminReport
	err := c.doAPI(ctx, http.MethodPost, fmt.Sprintf("/api/v1/admin/reports/%s/%s", url.PathEscape(string(id)), action), nil, &report, nil)
	if err != nil {
		return nil, err
	}
	return &report, nil
}
