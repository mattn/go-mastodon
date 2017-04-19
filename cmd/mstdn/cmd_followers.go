package main

import (
	"context"
	"fmt"

	"github.com/mattn/go-mastodon"
	"github.com/urfave/cli"
)

func cmdFollowers(c *cli.Context) error {
	client := c.App.Metadata["client"].(*mastodon.Client)
	config := c.App.Metadata["config"].(*mastodon.Config)

	account, err := client.GetAccountCurrentUser(context.Background())
	if err != nil {
		return err
	}
	followers, err := client.GetAccountFollowers(context.Background(), account.ID)
	if err != nil {
		return err
	}
	s := newScreen(config)
	for _, follower := range followers {
		fmt.Fprintf(c.App.Writer, "%v,%v\n", follower.ID, s.acct(follower.Acct))
	}
	return nil
}
