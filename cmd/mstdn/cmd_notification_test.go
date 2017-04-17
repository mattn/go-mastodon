package main

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/urfave/cli"
)

func TestCmdNotification(t *testing.T) {
	out := testWithServer(
		func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/api/v1/notifications":
				fmt.Fprintln(w, `[{"type": "rebloged", "status": {"content": "foo"}}]`)
				return
			}
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		},
		func(app *cli.App) {
			app.Run([]string{"mstdn", "notification"})
		},
	)
	if !strings.Contains(out, "rebloged") {
		t.Fatalf("%q should be contained in output of command: %v", "rebloged", out)
	}
}
