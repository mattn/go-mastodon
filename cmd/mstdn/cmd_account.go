package main

import (
	"context"
	"fmt"

	"github.com/mattn/go-mastodon"
	"github.com/urfave/cli"
)

func cmdAccount(c *cli.Context) error {
	client := c.App.Metadata["client"].(*mastodon.Client)
	account, err := client.GetAccountCurrentUser(context.Background())
	if err != nil {
		return err
	}
	fmt.Fprintf(c.App.Writer, "URI           : %v\n", account.Acct)
	fmt.Fprintf(c.App.Writer, "ID            : %v\n", account.ID)
	fmt.Fprintf(c.App.Writer, "Username      : %v\n", account.Username)
	fmt.Fprintf(c.App.Writer, "Acct          : %v\n", account.Acct)
	fmt.Fprintf(c.App.Writer, "DisplayName   : %v\n", account.DisplayName)
	fmt.Fprintf(c.App.Writer, "Locked        : %v\n", account.Locked)
	fmt.Fprintf(c.App.Writer, "CreatedAt     : %v\n", account.CreatedAt.Local())
	fmt.Fprintf(c.App.Writer, "FollowersCount: %v\n", account.FollowersCount)
	fmt.Fprintf(c.App.Writer, "FollowingCount: %v\n", account.FollowingCount)
	fmt.Fprintf(c.App.Writer, "StatusesCount : %v\n", account.StatusesCount)
	fmt.Fprintf(c.App.Writer, "Note          : %v\n", textContent(account.Note))
	fmt.Fprintf(c.App.Writer, "URL           : %v\n", account.URL)
	return nil
}
