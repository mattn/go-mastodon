package main

import (
	"context"
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
	config := c.App.Metadata["config"].(*mastodon.Config)

	results, err := client.Search(context.Background(), argstr(c), false)
	if err != nil {
		return err
	}
	s := newScreen(config)
	if len(results.Accounts) > 0 {
		fmt.Fprintln(c.App.Writer, "===ACCOUNT===")
		for _, result := range results.Accounts {
			fmt.Fprintf(c.App.Writer, "%v,%v\n", result.ID, s.acct(result.Acct))
		}
		fmt.Fprintln(c.App.Writer)
	}
	if len(results.Statuses) > 0 {
		fmt.Fprintln(c.App.Writer, "===STATUS===")
		for _, result := range results.Statuses {
			s.displayStatus(c.App.Writer, result)
		}
		fmt.Fprintln(c.App.Writer)
	}
	if len(results.Hashtags) > 0 {
		fmt.Fprintln(c.App.Writer, "===HASHTAG===")
		for _, result := range results.Hashtags {
			fmt.Fprintf(c.App.Writer, "#%v\n", result)
		}
		fmt.Fprintln(c.App.Writer)
	}
	return nil
}
