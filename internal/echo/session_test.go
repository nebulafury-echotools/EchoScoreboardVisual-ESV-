package echo

import (
	"os"
	"testing"
)

func TestFetchParsesCapturedEchoSession(t *testing.T) {
	if _, err := os.Stat("../../testdata/echo_session.json"); err != nil {
		t.Skip("live Echo fixture not present")
	}

	client := NewClient("http://127.0.0.1:1")
	session, err := client.Fetch()
	if err != nil {
		t.Fatalf("expected captured Echo session fixture to decode: %v", err)
	}
	if len(session.Teams) < 2 {
		t.Fatalf("expected at least two teams, got %d", len(session.Teams))
	}
	if len(session.Teams[0].Players) == 0 {
		t.Fatal("expected at least one player in the first team")
	}
}
