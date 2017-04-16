package main

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/urfave/cli"
)

func TestCmdToot(t *testing.T) {
	toot := ""
	testWithServer(
		func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/api/v1/statuses" {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			}
			toot = r.FormValue("status")
			fmt.Fprintln(w, `{"ID": 2345}`)
			return
		},
		func(app *cli.App) {
			app.Run([]string{"mstdn", "toot", "foo"})
		},
	)
	if toot != "foo" {
		t.Fatalf("want %q, got %q", "foo", toot)
	}
}
