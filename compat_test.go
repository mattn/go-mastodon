package mastodon_test

import (
	"encoding/json"
	"reflect"
	"slices"
	"testing"

	"github.com/mattn/go-mastodon"
)

func TestIDUnmarshalJSON(t *testing.T) {
	tests := []struct {
		in   string
		want mastodon.ID
	}{
		{`"123"`, "123"},
		{`123`, "123"},
		{`null`, ""},
	}
	for _, test := range tests {
		var id mastodon.ID
		if err := json.Unmarshal([]byte(test.in), &id); err != nil {
			t.Fatalf("should not be fail: %v", err)
		}
		if id != test.want {
			t.Fatalf("want %q but %q for input %s", test.want, id, test.in)
		}
	}
}

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

func TestIDCompareNonNumeric(t *testing.T) {
	// Non-numeric IDs (e.g. ULIDs used by GoToSocial) must not panic.
	ids := []mastodon.ID{
		"01F8MH5ZYAS9XKD4NK1FSD5J1Z",
		"01F8MH0BBE73B7VKDGZRE2M3XM",
		"123",
	}

	slices.SortFunc(ids, mastodon.ID.Compare)
	want := []mastodon.ID{
		"123",
		"01F8MH0BBE73B7VKDGZRE2M3XM",
		"01F8MH5ZYAS9XKD4NK1FSD5J1Z",
	}

	if got, want := ids, want; !reflect.DeepEqual(got, want) {
		t.Fatalf("invalid sorted slices:\ngot= %q\nwant=%q", got, want)
	}
}
