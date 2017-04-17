package mastodon

import (
	"context"
	"net/http"
)

// Report hold information for mastodon report.
type Report struct {
	ID          int64 `json:"id"`
	ActionTaken bool  `json:"action_taken"`
}

// GetReports return report of the current user.
func (c *Client) GetReports(ctx context.Context) ([]*Report, error) {
	var reports []*Report
	err := c.doAPI(ctx, http.MethodGet, "/api/v1/reports", nil, &reports)
	if err != nil {
		return nil, err
	}
	return reports, nil
}

//  Report reports the report
func (c *Client) Report(ctx context.Context, id int64) (*Report, error) {
	var report Report
	err := c.doAPI(ctx, http.MethodPost, "/api/v1/reports", nil, &report)
	if err != nil {
		return nil, err
	}
	return &report, nil
}
