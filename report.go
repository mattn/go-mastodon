package mastodon

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// Report hold information for mastodon report.
type Report struct {
	ID          int64 `json:"id"`
	ActionTaken bool  `json:"action_taken"`
}

// GetReports return report of the current user.
func (c *Client) GetReports(ctx context.Context) ([]*Report, error) {
	var reports []*Report
	err := c.doAPI(ctx, http.MethodGet, "api/v1/reports", nil, &reports, nil)
	if err != nil {
		return nil, err
	}
	return reports, nil
}

// Report reports the report
func (c *Client) Report(ctx context.Context, accountID int64, ids []int64, comment string) (*Report, error) {
	params := url.Values{}
	params.Set("account_id", fmt.Sprint(accountID))
	for _, id := range ids {
		params.Add("status_ids[]", fmt.Sprint(id))
	}
	params.Set("comment", comment)
	var report Report
	err := c.doAPI(ctx, http.MethodPost, "api/v1/reports", params, &report, nil)
	if err != nil {
		return nil, err
	}
	return &report, nil
}
