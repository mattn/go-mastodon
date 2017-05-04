package main

import (
	"context"
	"fmt"
	"time"

	"github.com/mattn/go-mastodon"
	"github.com/urfave/cli"
)

func cmdFollowers(c *cli.Context) error {
	client := c.App.Metadata["client"].(*mastodon.Client)
	config := c.App.Metadata["config"].(*mastodon.Config)

	account, err := client.GetAccountCurrentUser(context.Background())
	if err != nil {
		return err
	}
	var maxID *int64
	var followers []*mastodon.Account
	for {
		fs, pg, err := client.GetAccountFollowers(
			context.Background(), account.ID, &mastodon.Pagination{MaxID: maxID})
		if err != nil {
			return err
		}
		followers = append(followers, fs...)
		if pg.MaxID == nil {
			break
		}
		maxID = pg.MaxID
		time.Sleep(10 * time.Second)
	}
	s := newScreen(config)
	for _, follower := range followers {
		fmt.Fprintf(c.App.Writer, "%v,%v\n", follower.ID, s.acct(follower.Acct))
	}
	return nil
}
