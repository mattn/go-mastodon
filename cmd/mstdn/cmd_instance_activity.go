package main

import (
	"context"
	"fmt"

	"github.com/mattn/go-mastodon"
	"github.com/urfave/cli"
)

func cmdInstanceActivity(c *cli.Context) error {
	client := c.App.Metadata["client"].(*mastodon.Client)
	activities, err := client.GetInstanceActivity(context.Background())
	if err != nil {
		return err
	}
	for _, activity := range activities {
		fmt.Fprintf(c.App.Writer, "Logins        : %v\n", activity.Logins)
		fmt.Fprintf(c.App.Writer, "Registrations : %v\n", activity.Registrations)
		fmt.Fprintf(c.App.Writer, "Statuses      : %v\n", activity.Statuses)
		fmt.Fprintf(c.App.Writer, "Week          : %v\n", activity.Week)
	}
	return nil
}
