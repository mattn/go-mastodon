package main

import (
	"errors"
	"log"

	"github.com/mattn/go-mastodon"
	"github.com/urfave/cli"
)

func cmdToot(c *cli.Context) error {
	if !c.Args().Present() {
		return errors.New("arguments required")
	}

	var toot string
	ff := c.String("ff")
	if ff != "" {
		text, err := readFile(ff)
		if err != nil {
			log.Fatal(err)
		}
		toot = string(text)
	} else {
		toot = argstr(c)
	}
	client := c.App.Metadata["client"].(*mastodon.Client)
	_, err := client.PostStatus(&mastodon.Toot{
		Status: toot,
	})
	return err
}
