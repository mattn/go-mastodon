package main

import (
	"context"
	"errors"

	"github.com/mattn/go-mastodon"
	"github.com/urfave/cli"
)

func cmdFollow(c *cli.Context) error {
	client := c.App.Metadata["client"].(*mastodon.Client)
	if !c.Args().Present() {
		return errors.New("arguments required")
	}
	for i := 0; i < c.NArg(); i++ {
		account, err := client.AccountsSearch(context.Background(), c.Args().Get(i), 1)
		if err != nil {
			return err
		}
		if len(account) == 0 {
			continue
		}
		_, err = client.AccountFollow(context.Background(), account[0].ID)
		if err != nil {
			return err
		}
	}
	return nil
}
