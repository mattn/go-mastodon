package main

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"os/signal"
	"strings"

	"github.com/mattn/go-mastodon"
	"github.com/urfave/cli"
)

type SimpleJSON struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Acct     string `json:"acct"`
	Avatar   string `json:"avatar"`
	Content  string `json:"content"`
}

func cmdStream(c *cli.Context) error {
	asJSON := c.Bool("json")
	asSimpleJSON := c.Bool("simplejson")

	client := c.App.Metadata["client"].(*mastodon.Client)
	config := c.App.Metadata["config"].(*mastodon.Config)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)

	var q chan mastodon.Event
	var err error

	t := c.String("type")
	if t == "public" {
		q, err = client.StreamingPublic(ctx)
	} else if t == "" || t == "public/local" {
		q, err = client.StreamingPublicLocal(ctx)
	} else if strings.HasPrefix(t, "user:") {
		q, err = client.StreamingUser(ctx, t[5:])
	} else if strings.HasPrefix(t, "hashtag:") {
		q, err = client.StreamingHashtag(ctx, t[8:])
	} else {
		return errors.New("invalid type")
	}
	if err != nil {
		return err
	}
	go func() {
		<-sc
		cancel()
		close(q)
	}()

	s := newScreen(config)
	for e := range q {
		if asJSON {
			json.NewEncoder(c.App.Writer).Encode(e)
		} else if asSimpleJSON {
			if t, ok := e.(*mastodon.UpdateEvent); ok {
				json.NewEncoder(c.App.Writer).Encode(&SimpleJSON{
					ID:       t.Status.ID,
					Username: t.Status.Account.Username,
					Acct:     t.Status.Account.Acct,
					Avatar:   t.Status.Account.AvatarStatic,
					Content:  textContent(t.Status.Content),
				})
			}
		} else {
			switch t := e.(type) {
			case *mastodon.UpdateEvent:
				s.displayStatus(c.App.Writer, t.Status)
			case *mastodon.NotificationEvent:
				// TODO s.displayStatus(c.App.Writer, t.Notification.Status)
			case *mastodon.ErrorEvent:
				s.displayError(c.App.Writer, t)
			}
		}
	}
	return nil
}
