package main

import (
	"context"
	"errors"

	"github.com/mattn/go-mastodon"
	"github.com/urfave/cli"
)

func cmdDelete(c *cli.Context) error {
	client := c.App.Metadata["client"].(*mastodon.Client)
	if !c.Args().Present() {
		return errors.New("arguments required")
	}
	for i := 0; i < c.NArg(); i++ {
		err := client.DeleteStatus(context.Background(), mastodon.ID(c.Args().Get(i)))
		if err != nil {
			return err
		}
	}
	return nil
}
