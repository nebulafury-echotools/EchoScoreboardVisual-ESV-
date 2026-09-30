package echo

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestSaveSnapshotStoresLatestPopulatedSessionPerMatch(t *testing.T) {
	root := t.TempDir()
	first := []byte(`{"sessionid":"match-123","match_type":"Echo_Arena","game_status":"playing","teams":[{"team":"BLUE TEAM","players":[{"name":"stardust","stats":{"points":6,"assists":1,"saves":0}}]}]}`)
	last := []byte(`{"sessionid":"match-123","match_type":"Echo_Arena","game_status":"post_match","blue_points":2,"teams":[{"team":"BLUE TEAM","players":[{"name":"stardust","stats":{"points":6,"assists":1,"saves":0}}]}]}`)

	if err := SaveSnapshot(root, first); err != nil {
		t.Fatalf("save first populated snapshot: %v", err)
	}
	if err := SaveSnapshot(root, last); err != nil {
		t.Fatalf("save last populated snapshot: %v", err)
	}

	latestPath := filepath.Join(root, "match-data", "echo_session.json")
	latest, err := os.ReadFile(latestPath)
	if err != nil {
		t.Fatalf("read latest snapshot: %v", err)
	}
	if !bytes.Equal(latest, last) {
		t.Fatal("latest snapshot was not updated with the last successful response")
	}

	matchPath := filepath.Join(root, "match-data", "matches", "match-123.json")
	match, err := os.ReadFile(matchPath)
	if err != nil {
		t.Fatalf("read per-match snapshot: %v", err)
	}
	if !bytes.Equal(match, last) {
		t.Fatal("per-match snapshot did not retain the last successful response")
	}

	if err := SaveSnapshot(root, []byte(`{}`)); err == nil {
		t.Fatal("expected an unpopulated response to be rejected")
	}
	latest, err = os.ReadFile(latestPath)
	if err != nil || !bytes.Equal(latest, last) {
		t.Fatal("invalid response replaced the last populated snapshot")
	}
}
