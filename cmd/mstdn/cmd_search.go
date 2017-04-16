package main

import (
	"errors"
	"fmt"

	"github.com/mattn/go-mastodon"
	"github.com/urfave/cli"
)

func cmdSearch(c *cli.Context) error {
	if !c.Args().Present() {
		return errors.New("arguments required")
	}

	client := c.App.Metadata["client"].(*mastodon.Client)
	results, err := client.Search(argstr(c), false)
	if err != nil {
		return err
	}
	for _, result := range results.Accounts {
		fmt.Fprintln(c.App.Writer, result)
	}
	for _, result := range results.Statuses {
		fmt.Fprintln(c.App.Writer, result)
	}
	for _, result := range results.Hashtags {
		fmt.Fprintln(c.App.Writer, result)
	}
	return nil
}
