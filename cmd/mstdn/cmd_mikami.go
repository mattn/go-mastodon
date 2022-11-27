package main

import (
	"github.com/urfave/cli/v2"
)

func cmdMikami(c *cli.Context) error {
	return xSearch(c.App.Metadata["xsearch_url"].(string), "三上", c.App.Writer)
}
