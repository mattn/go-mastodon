package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mattn/go-mastodon"
)

// Create client with credentials from user generated application
func main() {
	config := &mastodon.Config{
		Server:       "https://mastodon.social",
		ClientID:     "ClientKey",
		ClientSecret: "ClientSecret",
		AccessToken:  "AccessToken",
	}

	// Create the client
	c := mastodon.NewClient(config)

	// Post a toot
	finalText := "this is the content of my new post!"
	visibility := "public"

	toot := mastodon.Toot{
		Status:     finalText,
		Visibility: visibility,
	}

	post, err := c.PostStatus(context.Background(), &toot)
	if err != nil {
		log.Fatalf("%#v\n", err)
	}

	fmt.Println("My new post is:", post)
}
