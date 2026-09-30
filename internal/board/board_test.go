package board

import (
	"echo-scoreboard-visual/internal/echo"
	"os"
	"testing"
)

func TestBuildReturnsWaitingStateForEmptySession(t *testing.T) {
	board := Build(echo.Session{})
	if !board.WaitingForMatch {
		t.Fatal("expected empty session to produce waiting state")
	}
}

func TestBuildMarksMVPPartialWhenCoreStatsAreMissing(t *testing.T) {
	session := echo.Session{
		MatchType: "Echo_Arena",
		Teams: []echo.Team{
			{Name: "BLUE TEAM", Players: []echo.Player{{Name: "BlueOne", Stats: map[string]float64{
				"points": 2, "assists": 1,
			}}}},
			{Name: "ORANGE TEAM", Players: []echo.Player{{Name: "OrangeOne", Stats: map[string]float64{
				"points": 1, "assists": 0,
			}}}},
		},
	}

	built := Build(session)
	blue := built.Teams[0].Players[0]
	orange := built.Teams[1].Players[0]
	if blue.MVPScore != 13 {
		t.Fatalf("expected BlueOne partial MVP score 13, got %v", blue.MVPScore)
	}
	if orange.MVPScore != 5 {
		t.Fatalf("expected OrangeOne partial MVP score 5, got %v", orange.MVPScore)
	}
	if !blue.IsMVP || orange.IsMVP {
		t.Fatalf("expected BlueOne to be MVP only; blue=%v orange=%v", blue.IsMVP, orange.IsMVP)
	}
	if !built.MVPAvailable || !blue.MVPAvailable || !orange.MVPAvailable {
		t.Fatal("expected echo_arena match to have MVP scoring available")
	}
	if built.MVPComplete || blue.MVPComplete || orange.MVPComplete {
		t.Fatal("expected MVP rating to be marked partial when reward stats are missing")
	}
}

func TestBuildUsesOwnerMVPWeightsForScreenshotMatch(t *testing.T) {
	session := echo.Session{
		MatchType: "Echo_Arena",
		Teams: []echo.Team{
			{Name: "BLUE TEAM", Players: []echo.Player{{Name: "stardust", Stats: map[string]float64{
				"points": 6, "assists": 1, "saves": 0,
			}}}},
			{Name: "ORANGE TEAM", Players: []echo.Player{
				{Name: "Coastal", Stats: map[string]float64{"points": 2, "assists": 5, "saves": 0}},
				{Name: "Aqua", Stats: map[string]float64{"points": 6, "assists": 0, "saves": 1}},
			}},
		},
	}

	built := Build(session)
	if built.Teams[0].Players[0].MVPScore != 33 {
		t.Fatalf("expected stardust score 33, got %v", built.Teams[0].Players[0].MVPScore)
	}
	if built.Teams[1].Players[0].MVPScore != 25 {
		t.Fatalf("expected Coastal score 25, got %v", built.Teams[1].Players[0].MVPScore)
	}
	if built.Teams[1].Players[1].MVPScore != 32 {
		t.Fatalf("expected Aqua score 32, got %v", built.Teams[1].Players[1].MVPScore)
	}
	if built.MVPName != "stardust" {
		t.Fatalf("expected stardust to win MVP, got %q", built.MVPName)
	}
}

func TestEchoArenaMVPScoreUsesSpecifiedWeights(t *testing.T) {
	score, complete := echoArenaMVPScore(map[string]float64{"points": 2, "assists": 3, "saves": 4})
	if score != 27 {
		t.Fatalf("expected weighted score 27, got %v", score)
	}
	if !complete {
		t.Fatal("expected complete stat set to produce a complete rating")
	}
}

func TestBuildDoesNotMarkMVPForOtherMatchTypes(t *testing.T) {
	built := Build(echo.Session{
		MatchType: "Echo_Combat",
		Teams:     []echo.Team{{Name: "BLUE TEAM", Players: []echo.Player{{Name: "Player"}}}},
	})
	if built.MVPAvailable || built.MVPName != "" || built.Teams[0].Players[0].IsMVP {
		t.Fatal("expected MVP scoring to be unavailable outside Echo_Arena")
	}
}

func TestSnapshotIDIsStableAcrossEndStates(t *testing.T) {
	session := echo.Session{
		SessionID:      "sample-session",
		MatchType:      "Echo_Arena",
		GameState:      "round_over",
		ClockDisplay:   "00:00.00",
		BluePoints:     5,
		OrangePoints:   3,
		BlueRoundScore: 1,
		Teams:          []echo.Team{{Name: "BLUE TEAM", Players: []echo.Player{{Name: "BlueOne"}}}},
	}
	first := Build(session)
	session.GameState = "post_match"
	second := Build(session)
	if first.SnapshotID == "" || first.SnapshotID != second.SnapshotID {
		t.Fatalf("expected stable non-empty snapshot ID, got %q and %q", first.SnapshotID, second.SnapshotID)
	}
}

func TestBuildUsesCapturedEchoSample(t *testing.T) {
	data, err := os.ReadFile("../../testdata/echo_session.json")
	if err != nil {
		t.Skip("live Echo fixture not present")
	}

	session, err := echo.ParseSession(data)
	if err != nil {
		t.Fatalf("expected a valid captured session: %v", err)
	}

	built := Build(session)
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
