package mastodon_test

import (
	"reflect"
	"slices"
	"testing"

	"github.com/mattn/go-mastodon"
)

func TestIDCompare(t *testing.T) {
	ids := []mastodon.ID{
		"123",
		"103",
		"",
		"0",
		"103",
		"122",
	}

	slices.SortFunc(ids, mastodon.ID.Compare)
	want := []mastodon.ID{
		"",
		"0",
		"103",
		"103",
		"122",
		"123",
	}

	if got, want := ids, want; !reflect.DeepEqual(got, want) {
		t.Fatalf("invalid sorted slices:\ngot= %q\nwant=%q", got, want)
	}
}
