package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/fatih/color"
	"github.com/mattn/go-mastodon"
	"github.com/urfave/cli"
)

func cmdStream(c *cli.Context) error {
	client := c.App.Metadata["client"].(*mastodon.Client)
	ctx, cancel := context.WithCancel(context.Background())
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	q, err := client.StreamingPublic(ctx)
	if err != nil {
		return err
	}
	go func() {
		<-sc
		cancel()
		close(q)
	}()
	for e := range q {
		switch t := e.(type) {
		case *mastodon.UpdateEvent:
			color.Set(color.FgHiRed)
			fmt.Fprintln(c.App.Writer, t.Status.Account.Username)
			color.Set(color.Reset)
			fmt.Fprintln(c.App.Writer, textContent(t.Status.Content))
		case *mastodon.ErrorEvent:
			color.Set(color.FgYellow)
			fmt.Fprintln(c.App.Writer, t.Error())
			color.Set(color.Reset)
		}
	}
	return nil
}
