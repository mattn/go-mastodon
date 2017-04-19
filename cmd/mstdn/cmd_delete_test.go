package main

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/urfave/cli"
)

func TestCmdDelete(t *testing.T) {
	ok := false
	f := func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/statuses/123":
			fmt.Fprintln(w, `{}`)
			ok = true
			return
		}
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	testWithServer(
		f, func(app *cli.App) {
			app.Run([]string{"mstdn", "delete", "122"})
		},
	)
	if ok {
		t.Fatal("something wrong to sequence to follow account")
	}

	ok = false
	testWithServer(
		f, func(app *cli.App) {
			app.Run([]string{"mstdn", "delete", "123"})
		},
	)
	if !ok {
		t.Fatal("something wrong to sequence to follow account")
	}
}
