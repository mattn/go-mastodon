package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/urfave/cli"
)

func TestCmdXSearch(t *testing.T) {
	testWithServer(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, `<div class="post"><div class="mst_content"><a href="http://example.com/@test/1"><p>test status</p></a></div></div>`)
		},
		func(app *cli.App) {
			err := app.Run([]string{"mstdn", "xsearch", "test"})
			if err != nil {
				t.Fatalf("should not be fail: %v", err)
			}
		},
	)
}

func TestXSearch(t *testing.T) {
	canErr := true
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if canErr {
			canErr = false
			http.Error(w, http.StatusText(http.StatusInternalServerError), 9999)
			return
		} else if r.URL.Query().Get("q") == "empty" {
			fmt.Fprintln(w, `<div class="post"><div class="mst_content"><a><p>test status</p></a></div></div>`)
			return
		}

		fmt.Fprintln(w, `<div class="post"><div class="mst_content"><a href="http://example.com/@test/1"><p>test status</p></a></div></div>`)
	}))
	defer ts.Close()

	err := xSearch(":", "", nil)
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}

	err = xSearch(ts.URL, "", nil)
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}

	buf := bytes.NewBuffer(nil)
	err = xSearch(ts.URL, "empty", buf)
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	result := buf.String()
	if result != "" {
		t.Fatalf("the search result should be empty: %s", result)
	}

	buf = bytes.NewBuffer(nil)
	err = xSearch(ts.URL, "test", buf)
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	result = buf.String()
	if !strings.Contains(result, "http://example.com/@test/1") {
		t.Fatalf("%q should be contained in output of search: %s", "http://example.com/@test/1", result)
	}
	if !strings.Contains(result, "test status") {
		t.Fatalf("%q should be contained in output of search: %s", "test status", result)
	}
}
