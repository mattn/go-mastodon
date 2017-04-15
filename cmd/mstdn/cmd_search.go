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
		fmt.Println(result)
	}
	for _, result := range results.Statuses {
		fmt.Println(result)
	}
	for _, result := range results.Hashtags {
		fmt.Println(result)
	}
	return nil
}
