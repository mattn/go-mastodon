package main

import (
	"context"
	"errors"
	"strconv"

	"github.com/mattn/go-mastodon"
	"github.com/urfave/cli"
)

func cmdDelete(c *cli.Context) error {
	client := c.App.Metadata["client"].(*mastodon.Client)
	if !c.Args().Present() {
		return errors.New("arguments required")
	}
	for i := 0; i < c.NArg(); i++ {
		id, err := strconv.ParseInt(c.Args().Get(i), 10, 64)
		if err != nil {
			return err
		}
		err = client.DeleteStatus(context.Background(), id)
		if err != nil {
			return err
		}
	}
	return nil
}
