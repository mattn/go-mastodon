package main

import (
	"context"
	"fmt"

	"github.com/fatih/color"
	"github.com/mattn/go-mastodon"
	"github.com/urfave/cli"
)

func cmdNotification(c *cli.Context) error {
	client := c.App.Metadata["client"].(*mastodon.Client)
	notifications, err := client.GetNotifications(context.Background(), nil)
	if err != nil {
		return err
	}
	for _, n := range notifications {
		if n.Status != nil {
			color.Set(color.FgHiRed)
			fmt.Fprint(c.App.Writer, n.Account.Acct)
			color.Set(color.Reset)
			fmt.Fprintln(c.App.Writer, " "+n.Type)
			s := n.Status
			fmt.Fprintln(c.App.Writer, textContent(s.Content))
		}
	}
	return nil
}
