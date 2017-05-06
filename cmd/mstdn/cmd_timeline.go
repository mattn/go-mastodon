package main

import (
	"context"

	"github.com/mattn/go-mastodon"
	"github.com/urfave/cli"
)

func cmdTimeline(c *cli.Context) error {
	client := c.App.Metadata["client"].(*mastodon.Client)
	config := c.App.Metadata["config"].(*mastodon.Config)
	timeline, err := client.GetTimelineHome(context.Background(), nil)
	if err != nil {
		return err
	}
	s := newScreen(config)
	for i := len(timeline) - 1; i >= 0; i-- {
		s.displayStatus(c.App.Writer, timeline[i])
	}
	return nil
}
