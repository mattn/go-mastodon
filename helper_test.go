package mastodon

import (
	"testing"
)

func TestString(t *testing.T) {
	s := "test"
	sp := String(s)
	if *sp != s {
		t.Fatalf("want %q but %q", s, *sp)
	}
}
