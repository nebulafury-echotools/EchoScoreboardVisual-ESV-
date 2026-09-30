package board

import (
	"echo-scoreboard-visual/internal/echo"
	"echo-scoreboard-visual/internal/spark"
	"os"
	"testing"
)

func TestBuildReturnsWaitingStateForEmptySession(t *testing.T) {
	board := Build(echo.Session{}, spark.Stats{})
	if !board.WaitingForMatch {
		t.Fatal("expected empty session to produce waiting state")
	}
}

func TestBuildUsesCapturedEchoSample(t *testing.T) {
	if _, err := os.Stat("../../testdata/echo_session.json"); err != nil {
		t.Skip("live Echo fixture not present")
	}

	client := echo.NewClient("http://127.0.0.1:1")
	session, err := client.Fetch()
	if err != nil {
		t.Fatalf("expected a valid captured session: %v", err)
	}

	built := Build(session, spark.Stats{})
	if built.Score.Blue == 0 && built.Score.Orange == 0 {
		t.Fatal("expected score fields to be populated from the live sample")
	}
	if len(built.Teams) < 2 {
		t.Fatalf("expected at least two teams in the board, got %d", len(built.Teams))
	}
	if built.TimeLeft == "" {
		t.Fatal("expected time remaining to be populated from the live sample")
	}
	for _, team := range built.Teams {
		if len(team.Players) != 4 {
			t.Fatalf("expected four players per team, got %d for %s", len(team.Players), team.Name)
		}
	}
}
