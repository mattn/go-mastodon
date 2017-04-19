package main

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/fatih/color"
	"github.com/mattn/go-mastodon"
	"github.com/urfave/cli"
)

func acct(acct, host string) string {
	if !strings.Contains(acct, "@") {
		acct += "@" + host
	}
	return acct
}

func cmdTimeline(c *cli.Context) error {
	client := c.App.Metadata["client"].(*mastodon.Client)
	config := c.App.Metadata["config"].(*mastodon.Config)
	timeline, err := client.GetTimelineHome(context.Background())
	if err != nil {
		return err
	}
	u, err := url.Parse(config.Server)
	if err != nil {
		return err
	}
	for i := len(timeline) - 1; i >= 0; i-- {
		t := timeline[i]
		if t.Reblog != nil {
			color.Set(color.FgHiRed)
			fmt.Fprint(c.App.Writer, acct(t.Account.Acct, u.Host))
			color.Set(color.Reset)
			fmt.Fprint(c.App.Writer, " rebloged ")
			color.Set(color.FgHiBlue)
			fmt.Fprintln(c.App.Writer, acct(t.Reblog.Account.Acct, u.Host))
			fmt.Fprintln(c.App.Writer, textContent(t.Reblog.Content))
			color.Set(color.Reset)
		} else {
			color.Set(color.FgHiRed)
			fmt.Fprintln(c.App.Writer, acct(t.Account.Acct, u.Host))
			color.Set(color.Reset)
			fmt.Fprintln(c.App.Writer, textContent(t.Content))
		}
	}
	return nil
}
