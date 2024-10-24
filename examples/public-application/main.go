package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mattn/go-mastodon"
)

func main() {
	// Register the application
	appConfig := &mastodon.AppConfig{
		Server:       "https://mastodon.social",
		ClientName:   "publicApp",
		Scopes:       "read write push",
		Website:      "https://github.com/mattn/go-mastodon",
		RedirectURIs: "urn:ietf:wg:oauth:2.0:oob",
	}

	app, err := mastodon.RegisterApp(context.Background(), appConfig)
	if err != nil {
		log.Fatal(err)
	}

	config := &mastodon.Config{
		Server:       "https://mastodon.social",
		ClientID:     app.ClientID,
		ClientSecret: app.ClientSecret,
	}

	// Create the client
	c := mastodon.NewClient(config)

	// Get an Access Token & Sets it in the client config
	err = c.GetAppAccessToken(context.Background(), app.RedirectURI)
	if err != nil {
		log.Fatal(err)
	}

	// Save credentials for later usage if you wish to do so, config file, database, etc...
	fmt.Println("ClientID:", c.Config.ClientID)
	fmt.Println("ClientSecret:", c.Config.ClientSecret)
	fmt.Println("Access Token:", c.Config.AccessToken)

	// Lookup and account id
	acc, err := c.AccountLookup(context.Background(), "coolapso")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(acc)

	pager := mastodon.Pagination{
		Limit: 10,
	}

	// Get the the usernames of users following the account ID
	followers, err := c.GetAccountFollowers(context.Background(), acc.ID, &pager)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range followers {
		fmt.Println(f.Username)
	}
}
