package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/mattn/go-mastodon"
)

func TestCmdStream(t *testing.T) {
	var ts *httptest.Server
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/streaming/public/local" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		f, _ := w.(http.Flusher)
		fmt.Fprintln(w, `
event: update
data: {"content": "foo", "account":{"acct":"FOO"}}
		`)
		f.Flush()

		fmt.Fprintln(w, `
event: update
data: {"content": "bar", "account":{"acct":"BAR"}}
		`)
		f.Flush()
		return
	}))
	defer ts.Close()

	config := &mastodon.Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	}
	client := mastodon.NewClient(config)

	var buf bytes.Buffer
	app := makeApp()
	app.Writer = &buf
	app.Metadata = map[string]interface{}{
		"client": client,
		"config": config,
	}

	stop := func() {
		time.Sleep(5 * time.Second)
		if sig, ok := app.Metadata["signal"]; ok {
			sig.(chan os.Signal) <- os.Interrupt
			return
		}
		panic("timeout")
	}

	var out string

	go stop()
	app.Run([]string{"mstdn", "stream"})
	out = buf.String()
	if !strings.Contains(out, "FOO@") {
		t.Fatalf("%q should be contained in output of command: %v", "FOO@", out)
	}
	if !strings.Contains(out, "foo") {
		t.Fatalf("%q should be contained in output of command: %v", "foo", out)
	}

	go stop()
	app.Run([]string{"mstdn", "stream", "--simplejson"})
	out = buf.String()
	if !strings.Contains(out, "FOO@") {
		t.Fatalf("%q should be contained in output of command: %v", "FOO@", out)
	}
	if !strings.Contains(out, "foo") {
		t.Fatalf("%q should be contained in output of command: %v", "foo", out)
	}
}
