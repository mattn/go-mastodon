package main

import (
	"fmt"
	"net/url"

	"github.com/PuerkitoBio/goquery"
	"github.com/urfave/cli"
)

func cmdXSearch(c *cli.Context) error {
	u, err := url.Parse("http://mastodonsearch.jp/cross/")
	if err != nil {
		return err
	}
	params := url.Values{}
	params.Set("q", c.Args().First())
	u.RawQuery = params.Encode()
	doc, err := goquery.NewDocument(u.String())
	if err != nil {
		return err
	}
	doc.Find(".post").Each(func(n int, elem *goquery.Selection) {
		href, ok := elem.Find(".mst_content a").Attr("href")
		if !ok {
			return
		}
		text := elem.Find(".mst_content p").Text()
		fmt.Println(href)
		fmt.Println(text)
		fmt.Println()
	})
	return nil
}
