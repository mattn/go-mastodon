package main

import (
	"fmt"

	"github.com/mattn/go-mastodon"
	"github.com/urfave/cli"
)

func cmdInstance(c *cli.Context) error {
	client := c.App.Metadata["client"].(*mastodon.Client)
	instance, err := client.GetInstance()
	if err != nil {
		return err
	}
	fmt.Printf("URI        : %s\n", instance.URI)
	fmt.Printf("Title      : %s\n", instance.Title)
	fmt.Printf("Description: %s\n", instance.Description)
	fmt.Printf("EMail      : %s\n", instance.EMail)
	return nil
}
