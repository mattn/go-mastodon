package mastodon

import (
	"encoding/json"
	"testing"
	"time"
)

func TestUnixtimeRoundTrip(t *testing.T) {
	var ut Unixtime
	if err := json.Unmarshal([]byte(`1546300800`), &ut); err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if got, want := time.Time(ut).Unix(), int64(1546300800); got != want {
		t.Fatalf("want %v but %v", want, got)
	}

	b, err := json.Marshal(ut)
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if string(b) != "1546300800" {
		t.Fatalf("want %q but %q", "1546300800", string(b))
	}
}

func TestUnixTimeStringRoundTrip(t *testing.T) {
	var ut UnixTimeString
	if err := json.Unmarshal([]byte(`"1546300800"`), &ut); err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if got, want := ut.Unix(), int64(1546300800); got != want {
		t.Fatalf("want %v but %v", want, got)
	}

	b, err := json.Marshal(ut)
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if string(b) != `"1546300800"` {
		t.Fatalf("want %q but %q", `"1546300800"`, string(b))
	}

	// The marshaled form must unmarshal back to the same time.
	var ut2 UnixTimeString
	if err := json.Unmarshal(b, &ut2); err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if !ut2.Equal(ut.Time) {
		t.Fatalf("want %v but %v", ut.Time, ut2.Time)
	}
}
