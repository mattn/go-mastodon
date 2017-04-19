package main

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/urfave/cli"
)

func TestCmdFollow(t *testing.T) {
	ok := false
	testWithServer(
		func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/api/v1/accounts/search":
				fmt.Fprintln(w, `[{"id": 123}]`)
				return
			case "/api/v1/accounts/123/follow":
				fmt.Fprintln(w, `{"id": 123}`)
				ok = true
				return
			}
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		},
		func(app *cli.App) {
			app.Run([]string{"mstdn", "follow", "mattn"})
		},
	)
	if !ok {
		t.Fatal("something wrong to sequence to follow account")
	}
}
