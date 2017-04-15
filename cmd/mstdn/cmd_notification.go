package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/mattn/go-mastodon"
	"github.com/urfave/cli"
)

func cmdNotification(c *cli.Context) error {
	client := c.App.Metadata["client"].(*mastodon.Client)
	notifications, err := client.GetNotifications()
	if err != nil {
		return err
	}
	for _, n := range notifications {
		if n.Status != nil {
			color.Set(color.FgHiRed)
			fmt.Print(n.Account.Username)
			color.Set(color.Reset)
			fmt.Println(" " + n.Type)
			s := n.Status
			fmt.Println(textContent(s.Content))
		}
	}
	return nil
}
