package main

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/urfave/cli"
)

func TestCmdSearch(t *testing.T) {
	out := testWithServer(
		func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/api/v1/search":
				fmt.Fprintln(w, `{"accounts": [{"id": 234, "acct": "zzz"}], "statuses":[{"id": 345, "content": "yyy"}], "hashtags": ["www", "わろす"]}`)
				return
			}
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		},
		func(app *cli.App) {
			app.Run([]string{"mstdn", "search", "zzz"})
		},
	)
	for _, s := range []string{"zzz", "yyy", "www", "わろす"} {
		if !strings.Contains(out, s) {
			t.Fatalf("%q should be contained in output of command: %v", s, out)
		}
	}
}
