package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"

	"github.com/fatih/color"
	"github.com/mattn/go-mastodon"
	"github.com/urfave/cli"
)

type SimpleJSON struct {
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Content  string `json:"content"`
}

func cmdStream(c *cli.Context) error {
	asJSON := c.Bool("json")
	asSimpleJSON := c.Bool("simplejson")
	client := c.App.Metadata["client"].(*mastodon.Client)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	q, err := client.StreamingPublic(ctx)
	if err != nil {
		return err
	}
	go func() {
		<-sc
		cancel()
		close(q)
	}()
	for e := range q {
		if asJSON {
			json.NewEncoder(c.App.Writer).Encode(e)
		} else if asSimpleJSON {
			if t, ok := e.(*mastodon.UpdateEvent); ok {
				json.NewEncoder(c.App.Writer).Encode(&SimpleJSON{
					Username: t.Status.Account.Username,
					Avatar:   t.Status.Account.AvatarStatic,
					Content:  textContent(t.Status.Content),
				})
			}
		} else {
			switch t := e.(type) {
			case *mastodon.UpdateEvent:
				color.Set(color.FgHiRed)
				fmt.Fprintln(c.App.Writer, t.Status.Account.Username)
				color.Set(color.Reset)
				fmt.Fprintln(c.App.Writer, textContent(t.Status.Content))
			case *mastodon.ErrorEvent:
				color.Set(color.FgYellow)
				fmt.Fprintln(c.App.Writer, t.Error())
				color.Set(color.Reset)
			}
		}
	}
	return nil
}
