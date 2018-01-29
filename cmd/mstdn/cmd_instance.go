package main

import (
	"context"
	"fmt"
	"sort"

	"github.com/mattn/go-mastodon"
	"github.com/urfave/cli"
)

func cmdInstance(c *cli.Context) error {
	client := c.App.Metadata["client"].(*mastodon.Client)
	instance, err := client.GetInstance(context.Background())
	if err != nil {
		return err
	}
	fmt.Fprintf(c.App.Writer, "URI        : %s\n", instance.URI)
	fmt.Fprintf(c.App.Writer, "Title      : %s\n", instance.Title)
	fmt.Fprintf(c.App.Writer, "Description: %s\n", instance.Description)
	fmt.Fprintf(c.App.Writer, "EMail      : %s\n", instance.EMail)
	if instance.Version != "" {
		fmt.Fprintf(c.App.Writer, "Version    : %s\n", instance.Version)
	}
	if instance.Thumbnail != "" {
		fmt.Fprintf(c.App.Writer, "Thumbnail  : %s\n", instance.Thumbnail)
	}
	if instance.URLs != nil {
		var keys []string
		for _, k := range instance.URLs {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fmt.Fprintf(c.App.Writer, "%s: %s\n", k, instance.URLs[k])
		}
	}
	if instance.Stats != nil {
		fmt.Fprintf(c.App.Writer, "User Count   : %v\n", instance.Stats.UserCount)
		fmt.Fprintf(c.App.Writer, "Status Count : %v\n", instance.Stats.StatusCount)
		fmt.Fprintf(c.App.Writer, "Domain Count : %v\n", instance.Stats.DomainCount)
	}
	return nil
}
