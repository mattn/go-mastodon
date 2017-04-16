package main

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/urfave/cli"
)

func TestCmdInstance(t *testing.T) {
	out := testWithServer(
		func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/api/v1/instance":
				fmt.Fprintln(w, `{"title": "zzz"}`)
				return
			}
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		},
		func(app *cli.App) {
			app.Run([]string{"mstdn", "instance"})
		},
	)
	if !strings.Contains(out, "zzz") {
		t.Fatalf("%q should be contained in output of command: %v", "zzz", out)
	}
}
