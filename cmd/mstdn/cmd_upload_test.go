package main

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/urfave/cli"
)

func TestCmdUpload(t *testing.T) {
	out := testWithServer(
		func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/api/v1/media":
				fmt.Fprintln(w, `{"id": 123}`)
				return
			}
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		},
		func(app *cli.App) {
			app.Run([]string{"mstdn", "upload", "../../testdata/logo.png"})
		},
	)
	if !strings.Contains(out, "123") {
		t.Fatalf("%q should be contained in output of command: %v", "123", out)
	}
}
