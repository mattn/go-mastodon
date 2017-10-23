package main

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"os/signal"
	"strings"
	"text/template"

	"github.com/mattn/go-mastodon"
	"github.com/urfave/cli"
)

// SimpleJSON is a struct for output JSON for data to be simple used
type SimpleJSON struct {
	ID       mastodon.ID `json:"id"`
	Username string      `json:"username"`
	Acct     string      `json:"acct"`
	Avatar   string      `json:"avatar"`
	Content  string      `json:"content"`
}

func checkFlag(f ...bool) bool {
	n := 0
	for _, on := range f {
		if on {
			n++
		}
	}
	return n > 1
}

func cmdStream(c *cli.Context) error {
	asJSON := c.Bool("json")
	asSimpleJSON := c.Bool("simplejson")
	asFormat := c.String("template")

	if checkFlag(asJSON, asSimpleJSON, asFormat != "") {
		return errors.New("cannot speicify two or three options in --json/--simplejson/--template")
	}
	tx, err := template.New("mstdn").Funcs(template.FuncMap{
		"nl": func(s string) string {
			return s + "\n"
		},
		"text": func(s string) string {
			return textContent(s)
		},
	}).Parse(asFormat)
	if err != nil {
		return err
	}

	client := c.App.Metadata["client"].(*mastodon.Client)
	config := c.App.Metadata["config"].(*mastodon.Config)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)

	var q chan mastodon.Event

	t := c.String("type")
	if t == "public" {
		q, err = client.StreamingPublic(ctx, false)
	} else if t == "" || t == "public/local" {
		q, err = client.StreamingPublic(ctx, true)
	} else if strings.HasPrefix(t, "user:") {
		q, err = client.StreamingUser(ctx)
	} else if strings.HasPrefix(t, "hashtag:") {
		q, err = client.StreamingHashtag(ctx, t[8:], false)
	} else {
		return errors.New("invalid type")
	}
	if err != nil {
		return err
	}
	go func() {
		<-sc
		cancel()
	}()

	c.App.Metadata["signal"] = sc

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
		} else if asFormat != "" {
			tx.ExecuteTemplate(c.App.Writer, "mstdn", e)
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
