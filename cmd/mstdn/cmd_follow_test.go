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
				q := r.URL.Query().Get("q")
				if q == "mattn" {
					fmt.Fprintln(w, `[{"id": 123}]`)
					return
				} else if q == "different_id" {
					fmt.Fprintln(w, `[{"id": 1234567}]`)
					return
				} else if q == "empty" {
					fmt.Fprintln(w, `[]`)
					return
				}
			case "/api/v1/accounts/123/follow":
				fmt.Fprintln(w, `{"id": 123}`)
				ok = true
				return
			}
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		},
		func(app *cli.App) {
			err := app.Run([]string{"mstdn", "follow", "mattn"})
			if err != nil {
				t.Fatalf("should not be fail: %v", err)
			}
		},
		func(app *cli.App) {
			err := app.Run([]string{"mstdn", "follow"})
			if err == nil {
				t.Fatalf("should be fail: %v", err)
			}
		},
		func(app *cli.App) {
			err := app.Run([]string{"mstdn", "follow", "fail"})
			if err == nil {
				t.Fatalf("should be fail: %v", err)
			}
		},
		func(app *cli.App) {
			err := app.Run([]string{"mstdn", "follow", "empty"})
			if err != nil {
				t.Fatalf("should not be fail: %v", err)
			}
		},
		func(app *cli.App) {
			err := app.Run([]string{"mstdn", "follow", "different_id"})
			if err == nil {
				t.Fatalf("should be fail: %v", err)
			}
		},
	)
	if !ok {
		t.Fatal("something wrong to sequence to follow account")
	}
}
