package main

import (
	"context"
	"fmt"

	"github.com/fatih/color"
	"github.com/mattn/go-mastodon"
	"github.com/urfave/cli"
)

func cmdTimeline(c *cli.Context) error {
	client := c.App.Metadata["client"].(*mastodon.Client)
	timeline, err := client.GetTimelineHome(context.Background())
	if err != nil {
		return err
	}
	for i := len(timeline) - 1; i >= 0; i-- {
		t := timeline[i]
		color.Set(color.FgHiRed)
		fmt.Fprintln(c.App.Writer, t.Account.Username)
		color.Set(color.Reset)
		fmt.Fprintln(c.App.Writer, textContent(t.Content))
	}
	return nil
}
