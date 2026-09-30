package echo

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestFetchUsesConfiguredEndpoint(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"match_type":"Echo_Arena","blue_points":7}`))
	}))
	defer server.Close()

	session, err := NewClient(server.URL).Fetch()
	if err != nil {
		t.Fatalf("expected endpoint response: %v", err)
	}
	if session.BluePoints != 7 {
		t.Fatalf("expected live score 7, got %d", session.BluePoints)
	}
}

func TestParseCapturedEchoSession(t *testing.T) {
	data, err := os.ReadFile("../../testdata/echo_session.json")
	if err != nil {
		t.Skip("live Echo fixture not present")
	}

	session, err := ParseSession(data)
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
