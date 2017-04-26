package main

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
	"github.com/urfave/cli"
)

func cmdMikami(c *cli.Context) error {
	doc, err := goquery.NewDocument("http://mastodonsearch.jp/cross/?q=三上")
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
