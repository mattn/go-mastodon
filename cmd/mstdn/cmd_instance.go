package main

import (
	"fmt"

	"github.com/mattn/go-mastodon"
	"github.com/urfave/cli"
)

func cmdInstance(c *cli.Context) error {
	client := c.App.Metadata["client"].(*mastodon.Client)
	instance, err := client.GetInstance()
	if err != nil {
		return err
	}
	fmt.Fprintf(c.App.Writer, "URI        : %s\n", instance.URI)
	fmt.Fprintf(c.App.Writer, "Title      : %s\n", instance.Title)
	fmt.Fprintf(c.App.Writer, "Description: %s\n", instance.Description)
	fmt.Fprintf(c.App.Writer, "EMail      : %s\n", instance.EMail)
	return nil
}
