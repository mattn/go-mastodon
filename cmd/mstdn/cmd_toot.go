package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/mattn/go-mastodon"
	"github.com/urfave/cli"
)

func cmdToot(c *cli.Context) error {
	var toot string
	ff := c.String("ff")
	if ff != "" {
		text, err := readFile(ff)
		if err != nil {
			return err
		}
		toot = string(text)
	} else {
		if !c.Args().Present() {
			return errors.New("arguments required")
		}
		toot = argstr(c)
	}
	client := c.App.Metadata["client"].(*mastodon.Client)
	_, err := client.PostStatus(context.Background(), &mastodon.Toot{
		Status:      toot,
		InReplyToID: mastodon.ID(fmt.Sprint(c.Int64("i"))),
	})
	return err
}
