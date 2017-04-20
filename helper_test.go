package mastodon

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
)

const wantBase64 = "data:image/png;base64,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"

func TestBase64EncodeFileName(t *testing.T) {
	// Error in os.Open.
	uri, err := Base64EncodeFileName("fail")
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}

	// Success.
	uri, err = Base64EncodeFileName("testdata/logo.png")
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if uri != wantBase64 {
		t.Fatalf("want %q but %q", wantBase64, uri)
	}
}

func TestBase64Encode(t *testing.T) {
	// Error in file.Stat.
	uri, err := Base64Encode(nil)
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}

	// Error in file.Read.
	logo, err := os.Open("testdata/logo.png")
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	_, err = ioutil.ReadAll(logo)
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	uri, err = Base64Encode(logo)
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}

	// Success.
	logo, err = os.Open("testdata/logo.png")
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	uri, err = Base64Encode(logo)
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if uri != wantBase64 {
		t.Fatalf("want %q but %q", wantBase64, uri)
	}
}

func TestString(t *testing.T) {
	s := "test"
	sp := String(s)
	if *sp != s {
		t.Fatalf("want %q but %q", s, *sp)
	}
}

func TestParseAPIError(t *testing.T) {
	// No api error.
	r := ioutil.NopCloser(strings.NewReader(`<html><head><title>404</title></head></html>`))
	err := parseAPIError("bad request", &http.Response{Status: "404 Not Found", Body: r})
	want := "bad request: 404 Not Found"
	if err.Error() != want {
		t.Fatalf("want %q but %q", want, err.Error())
	}

	// With api error.
	r = ioutil.NopCloser(strings.NewReader(`{"error":"Record not found"}`))
	err = parseAPIError("bad request", &http.Response{Status: "404 Not Found", Body: r})
	want = "bad request: 404 Not Found: Record not found"
	if err.Error() != want {
		t.Fatalf("want %q but %q", want, err.Error())
	}
}
