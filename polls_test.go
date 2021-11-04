package mastodon

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetPoll(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/polls/1234567" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		fmt.Fprintln(w, `{"id": "1234567", "expires_at": "2019-12-05T04:05:08.302Z", "expired": true, "multiple": false, "votes_count": 10, "voters_count": null, "voted": true, "own_votes": [1], "options": [{"title": "accept", "votes_count": 6}, {"title": "deny", "votes_count": 4}], "emojis":[{"shortcode":"💩", "url":"http://example.com", "static_url": "http://example.com/static"}]}`)
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	_, err := client.GetPoll(context.Background(), "123")
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	poll, err := client.GetPoll(context.Background(), "1234567")
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if poll.Expired != true {
		t.Fatalf("want %t but %t", true, poll.Expired)
	}
	if poll.Multiple != false {
		t.Fatalf("want %t but %t", true, poll.Multiple)
	}
	if poll.VotesCount != 10 {
		t.Fatalf("want %d but %d", 10, poll.VotesCount)
	}
	if poll.VotersCount != 0 {
		t.Fatalf("want %d but %d", 0, poll.VotersCount)
	}
	if poll.Voted != true {
		t.Fatalf("want %t but %t", true, poll.Voted)
	}
	if len(poll.OwnVotes) != 1 {
		t.Fatalf("should have own votes")
	}
	if poll.OwnVotes[0] != 1 {
		t.Fatalf("want %d but %d", 1, poll.OwnVotes[0])
	}
	if len(poll.Options) != 2 {
		t.Fatalf("should have 2 options")
	}
	if poll.Options[0].Title != "accept" {
		t.Fatalf("want %q but %q", "accept", poll.Options[0].Title)
	}
	if poll.Options[0].VotesCount != 6 {
		t.Fatalf("want %q but %q", 6, poll.Options[0].VotesCount)
	}
	if poll.Options[1].Title != "deny" {
		t.Fatalf("want %q but %q", "deny", poll.Options[1].Title)
	}
	if poll.Options[1].VotesCount != 4 {
		t.Fatalf("want %q but %q", 4, poll.Options[1].VotesCount)
	}
	if len(poll.Emojis) != 1 {
		t.Fatal("should have emojis")
	}
	if poll.Emojis[0].ShortCode != "💩" {
		t.Fatalf("want %q but %q", "💩", poll.Emojis[0].ShortCode)
	}
	if poll.Emojis[0].URL != "http://example.com" {
		t.Fatalf("want %q but %q", "https://example.com", poll.Emojis[0].URL)
	}
	if poll.Emojis[0].StaticURL != "http://example.com/static" {
		t.Fatalf("want %q but %q", "https://example.com/static", poll.Emojis[0].StaticURL)
	}
}

func TestPollVote(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/polls/1234567/votes" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		if r.Method != "POST" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusMethodNotAllowed)
			return
		}
		fmt.Fprintln(w, `{"id": "1234567", "expires_at": "2019-12-05T04:05:08.302Z", "expired": false, "multiple": false, "votes_count": 10, "voters_count": null, "voted": true, "own_votes": [1], "options": [{"title": "accept", "votes_count": 6}, {"title": "deny", "votes_count": 4}], "emojis":[]}`)
	}))
	defer ts.Close()

	client := NewClient(&Config{
		Server:       ts.URL,
		ClientID:     "foo",
		ClientSecret: "bar",
		AccessToken:  "zoo",
	})
	poll, err := client.PollVote(context.Background(), ID("1234567"), 1)
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if poll.Expired != false {
		t.Fatalf("want %t but %t", false, poll.Expired)
	}
	if poll.Multiple != false {
		t.Fatalf("want %t but %t", true, poll.Multiple)
	}
	if poll.VotesCount != 10 {
		t.Fatalf("want %d but %d", 10, poll.VotesCount)
	}
	if poll.VotersCount != 0 {
		t.Fatalf("want %d but %d", 0, poll.VotersCount)
	}
	if poll.Voted != true {
		t.Fatalf("want %t but %t", true, poll.Voted)
	}
	if len(poll.OwnVotes) != 1 {
		t.Fatalf("should have own votes")
	}
	if poll.OwnVotes[0] != 1 {
		t.Fatalf("want %d but %d", 1, poll.OwnVotes[0])
	}
	if len(poll.Options) != 2 {
		t.Fatalf("should have 2 options")
	}
	if poll.Options[0].Title != "accept" {
		t.Fatalf("want %q but %q", "accept", poll.Options[0].Title)
	}
	if poll.Options[0].VotesCount != 6 {
		t.Fatalf("want %q but %q", 6, poll.Options[0].VotesCount)
	}
	if poll.Options[1].Title != "deny" {
		t.Fatalf("want %q but %q", "deny", poll.Options[1].Title)
	}
	if poll.Options[1].VotesCount != 4 {
		t.Fatalf("want %q but %q", 4, poll.Options[1].VotesCount)
	}
}
