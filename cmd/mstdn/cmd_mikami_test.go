package main

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/urfave/cli"
)

func TestCmdMikami(t *testing.T) {
	ok := false
	buf := bytes.NewBuffer(nil)
	testWithServer(
		func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Query().Get("q") == "三上" {
				ok = true
				fmt.Fprintln(w, `<div class="post"><div class="mst_content"><a href="http://example.com/@test/1"><p>三上</p></a></div></div>`)
			}
		},
		func(app *cli.App) {
			app.Writer = buf
			err := app.Run([]string{"mstdn", "mikami"})
			if err != nil {
				t.Fatalf("should not be fail: %v", err)
			}
		},
	)
	if !ok {
		t.Fatal("should be search Mikami")
	}
	result := buf.String()
	if !strings.Contains(result, "http://example.com/@test/1") {
		t.Fatalf("%q should be contained in output of search: %s", "http://example.com/@test/1", result)
	}
	if !strings.Contains(result, "三上") {
		t.Fatalf("%q should be contained in output of search: %s", "三上", result)
	}
}
