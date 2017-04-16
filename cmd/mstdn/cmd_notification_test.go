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
			if r.URL.Path != "/api/v1/notifications" {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			}
			fmt.Fprintln(w, `[{"type": "rebloged", "status": {"content": "foo"}}]`)
			return
		},
		func(app *cli.App) {
			app.Run([]string{"mstdn", "notification"})
		},
	)
	if !strings.Contains(out, "rebloged") {
		t.Fatalf("%q should be contained in output of instance command: %v", "rebloged", out)
	}
}
