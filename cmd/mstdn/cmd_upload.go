package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/mattn/go-mastodon"
	"github.com/urfave/cli"
)

func cmdUpload(c *cli.Context) error {
	if !c.Args().Present() {
		return errors.New("arguments required")
	}
	client := c.App.Metadata["client"].(*mastodon.Client)
	for i := 0; i < c.NArg(); i++ {
		attachment, err := client.UploadMedia(context.Background(), c.Args().Get(i))
		if err != nil {
			return err
		}
		if i > 0 {
			fmt.Fprintln(c.App.Writer)
		}
		fmt.Fprintf(c.App.Writer, "ID        : %v\n", attachment.ID)
		fmt.Fprintf(c.App.Writer, "Type      : %v\n", attachment.Type)
		fmt.Fprintf(c.App.Writer, "URL       : %v\n", attachment.URL)
		fmt.Fprintf(c.App.Writer, "RemoteURL : %v\n", attachment.RemoteURL)
		fmt.Fprintf(c.App.Writer, "PreviewURL: %v\n", attachment.PreviewURL)
		fmt.Fprintf(c.App.Writer, "TextURL   : %v\n", attachment.TextURL)
	}
	return nil
}
