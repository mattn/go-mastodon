package main

import (
	"context"
	"fmt"

	"github.com/mattn/go-mastodon"
	"github.com/urfave/cli"
)

func cmdFollowers(c *cli.Context) error {
	client := c.App.Metadata["client"].(*mastodon.Client)
	account, err := client.GetAccountCurrentUser(context.Background())
	if err != nil {
		return err
	}
	followers, err := client.GetAccountFollowers(context.Background(), account.ID)
	if err != nil {
		return err
	}
	for _, follower := range followers {
		fmt.Fprintf(c.App.Writer, "%v,%v\n", follower.ID, follower.Username)
	}
	return nil
}
