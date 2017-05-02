package main

import (
	"fmt"
	"io"
	"net/url"

	"github.com/PuerkitoBio/goquery"
	"github.com/urfave/cli"
)

func cmdXSearch(c *cli.Context) error {
	return xSearch(c.App.Metadata["xsearch_url"].(string), c.Args().First(), c.App.Writer)
}

func xSearch(xsearchRawurl, query string, w io.Writer) error {
	u, err := url.Parse(xsearchRawurl)
	if err != nil {
		return err
	}
	params := url.Values{}
	params.Set("q", query)
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
		fmt.Fprintf(w, "%s\n", href)
		fmt.Fprintf(w, "%s\n\n", text)
	})
	return nil
}
