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

func cmdTimelineHome(c *cli.Context) error {
	return cmdTimeline(c)
}

func cmdTimelinePublic(c *cli.Context) error {
	client := c.App.Metadata["client"].(*mastodon.Client)
	config := c.App.Metadata["config"].(*mastodon.Config)
	timeline, err := client.GetTimelinePublic(context.Background(), false, nil)
	if err != nil {
		return err
	}
	s := newScreen(config)
	for i := len(timeline) - 1; i >= 0; i-- {
		s.displayStatus(c.App.Writer, timeline[i])
	}
	return nil
}

func cmdTimelineLocal(c *cli.Context) error {
	client := c.App.Metadata["client"].(*mastodon.Client)
	config := c.App.Metadata["config"].(*mastodon.Config)
	timeline, err := client.GetTimelinePublic(context.Background(), true, nil)
	if err != nil {
		return err
	}
	s := newScreen(config)
	for i := len(timeline) - 1; i >= 0; i-- {
		s.displayStatus(c.App.Writer, timeline[i])
	}
	return nil
}

func cmdTimelineDirect(c *cli.Context) error {
	client := c.App.Metadata["client"].(*mastodon.Client)
	config := c.App.Metadata["config"].(*mastodon.Config)
	timeline, err := client.GetTimelineDirect(context.Background(), nil)
	if err != nil {
		return err
	}
	s := newScreen(config)
	for i := len(timeline) - 1; i >= 0; i-- {
		s.displayStatus(c.App.Writer, timeline[i])
	}
	return nil
}
