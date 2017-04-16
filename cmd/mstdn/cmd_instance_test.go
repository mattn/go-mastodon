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
			if r.URL.Path != "/api/v1/instance" {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			}
			fmt.Fprintln(w, `{"Title": "zzz"}`)
			return
		},
		func(app *cli.App) {
			app.Run([]string{"mstdn", "instance"})
		},
	)
	if !strings.Contains(out, "zzz") {
		t.Fatalf("%q should be contained in output of instance command: %v", "zzz", out)
	}
}
