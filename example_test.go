package mastodon_test

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mattn/go-mastodon"
)

func ExampleRegisterApp() {
	app, err := mastodon.RegisterApp(context.Background(), &mastodon.AppConfig{
		Server:     "https://mstdn.jp",
		ClientName: "client-name",
		Scopes:     "read write follow",
		Website:    "https://github.com/mattn/go-mastodon",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("client-id    : %s\n", app.ClientID)
	fmt.Printf("client-secret: %s\n", app.ClientSecret)
}

func ExampleClient() {
	c := mastodon.NewClient(&mastodon.Config{
		Server:       "https://mstdn.jp",
		ClientID:     "client-id",
		ClientSecret: "client-secret",
	})
	err := c.Authenticate(context.Background(), "your-email", "your-password")
	if err != nil {
		log.Fatal(err)
	}
	timeline, err := c.GetTimelineHome(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	for i := len(timeline) - 1; i >= 0; i-- {
		fmt.Println(timeline[i])
	}
}

func ExamplePagination() {
	c := mastodon.NewClient(&mastodon.Config{
		Server:       "https://mstdn.jp",
		ClientID:     "client-id",
		ClientSecret: "client-secret",
	})
	var followers []*mastodon.Account
	var pg mastodon.Pagination
	for {
		fs, err := c.GetAccountFollowers(context.Background(), "1", &pg)
		if err != nil {
			log.Fatal(err)
		}
		followers = append(followers, fs...)
		if pg.MaxID == "" {
			break
		}
		time.Sleep(10 * time.Second)
	}
	for _, f := range followers {
		fmt.Println(f.Acct)
	}
}
