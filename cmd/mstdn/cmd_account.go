package main

import (
	"fmt"

	"github.com/mattn/go-mastodon"
	"github.com/urfave/cli"
)

func cmdAccount(c *cli.Context) error {
	client := c.App.Metadata["client"].(*mastodon.Client)
	account, err := client.GetAccountCurrentUser()
	if err != nil {
		return err
	}
	fmt.Printf("URI           : %v\n", account.Acct)
	fmt.Printf("ID            : %v\n", account.ID)
	fmt.Printf("Username      : %v\n", account.Username)
	fmt.Printf("Acct          : %v\n", account.Acct)
	fmt.Printf("DisplayName   : %v\n", account.DisplayName)
	fmt.Printf("Locked        : %v\n", account.Locked)
	fmt.Printf("CreatedAt     : %v\n", account.CreatedAt.Local())
	fmt.Printf("FollowersCount: %v\n", account.FollowersCount)
	fmt.Printf("FollowingCount: %v\n", account.FollowingCount)
	fmt.Printf("StatusesCount : %v\n", account.StatusesCount)
	fmt.Printf("Note          : %v\n", account.Note)
	fmt.Printf("URL           : %v\n", account.URL)
	return nil
}
