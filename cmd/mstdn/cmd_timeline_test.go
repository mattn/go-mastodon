package main

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/urfave/cli"
)

func TestCmdTimeline(t *testing.T) {
	out := testWithServer(
		func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/api/v1/timelines/home":
				fmt.Fprintln(w, `[{"content": "home"}]`)
				return
			case "/api/v1/timelines/public":
				fmt.Fprintln(w, `[{"content": "public"}]`)
				return
			case "/api/v1/conversations":
				fmt.Fprintln(w, `[{"id": "4", "unread":false, "last_status" : {"content": "direct"}}]`)
				return
			}
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		},
		func(app *cli.App) {
			app.Run([]string{"mstdn", "timeline"})
			app.Run([]string{"mstdn", "timeline-home"})
			app.Run([]string{"mstdn", "timeline-public"})
			app.Run([]string{"mstdn", "timeline-local"})
			app.Run([]string{"mstdn", "timeline-direct"})
		},
	)
	want := strings.Join([]string{
		"@example.com",
		"home",
		"@example.com",
		"home",
		"@example.com",
		"public",
		"@example.com",
		"public",
		"@example.com",
		"direct",
	}, "\n") + "\n"
	if !strings.Contains(out, want) {
		t.Fatalf("%q should be contained in output of command: %v", want, out)
	}
}
