package main

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/urfave/cli"
)

func TestCmdFollowers(t *testing.T) {
	out := testWithServer(
		func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/api/v1/accounts/verify_credentials":
				fmt.Fprintln(w, `{"id": 123}`)
				return
			case "/api/v1/accounts/123/followers":
				w.Header().Set("Link", `<http://example.com?since_id=890>; rel="prev"`)
				fmt.Fprintln(w, `[{"id": 234, "username": "ZZZ", "acct": "zzz"}]`)
				return
			}
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		},
		func(app *cli.App) {
			app.Run([]string{"mstdn", "followers"})
		},
	)
	if !strings.Contains(out, "zzz") {
		t.Fatalf("%q should be contained in output of command: %v", "zzz", out)
	}
}
