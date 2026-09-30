package board

import (
	"echo-scoreboard-visual/internal/echo"
	"echo-scoreboard-visual/internal/spark"
	"strings"
)

type Board struct {
	WaitingForMatch bool        `json:"waitingForMatch"`
	Message         string      `json:"message,omitempty"`
	GameStatus      string      `json:"gameStatus,omitempty"`
	GameOver        bool        `json:"gameOver,omitempty"`
	TimeLeft        string      `json:"timeLeft,omitempty"`
	Score           Score       `json:"score,omitempty"`
	RoundScores     RoundScores `json:"roundScores,omitempty"`
	Teams           []Team      `json:"teams,omitempty"`
	Players         []Player    `json:"players,omitempty"`
	StaleSince      string      `json:"staleSince,omitempty"`
}

type Score struct {
	Blue   int `json:"blue,omitempty"`
	Orange int `json:"orange,omitempty"`
}

type RoundScores struct {
	Blue   int `json:"blue,omitempty"`
	Orange int `json:"orange,omitempty"`
}

type Team struct {
	Name    string   `json:"name,omitempty"`
	Color   string   `json:"color,omitempty"`
	Players []Player `json:"players,omitempty"`
}

type Player struct {
	Name     string             `json:"name,omitempty"`
	Username string             `json:"username,omitempty"`
	Team     string             `json:"team,omitempty"`
	Ping     int                `json:"ping,omitempty"`
	MVP      float64            `json:"mvp,omitempty"`
	Stats    map[string]float64 `json:"stats,omitempty"`
}

func Build(session echo.Session, stats spark.Stats) Board {
	if len(session.Teams) == 0 && session.BluePoints == 0 && session.OrangePoints == 0 && session.ClockDisplay == "" && session.GameState == "" {
		return Board{WaitingForMatch: true, Message: "waiting for match"}
	}

	board := Board{
		GameStatus: session.GameState,
		TimeLeft:   normalizeClock(session.ClockDisplay),
		Score: Score{
			Blue:   session.BluePoints,
			Orange: session.OrangePoints,
		},
		RoundScores: RoundScores{
			Blue:   session.BlueRoundScore,
			Orange: session.OrangeRoundScore,
		},
	}

	if board.TimeLeft == "00:00.00" || board.TimeLeft == "00:00.000" || strings.Contains(strings.ToLower(board.GameStatus), "post") {
		board.GameOver = true
		board.Message = "Game Over"
	}

	if len(session.Teams) == 0 {
		for _, player := range session.Players {
			board.Players = append(board.Players, normalizePlayer(player, stats))
		}
		return board
	}

	for _, team := range session.Teams {
		name := strings.TrimSpace(team.Name)
		if strings.Contains(strings.ToUpper(name), "SPECT") {
			continue
		}
		color := "blue"
		if strings.Contains(strings.ToUpper(name), "ORANGE") {
			color = "orange"
		}

		boardTeam := Team{Name: name, Color: color}
		players := team.Players
		if len(players) > 4 {
			players = players[:4]
		}
		for _, player := range players {
			normalized := normalizePlayer(player, stats)
			normalized.Team = color
			boardTeam.Players = append(boardTeam.Players, normalized)
			board.Players = append(board.Players, normalized)
		}
		board.Teams = append(board.Teams, boardTeam)
	}

	return board
}

func normalizePlayer(player echo.Player, stats spark.Stats) Player {
	normalized := Player{
		Name:     player.Name,
		Username: player.Username,
		Team:     player.Team,
		Ping:     player.Ping,
		MVP:      player.MVP,
		Stats:    map[string]float64{},
	}
	if player.Stats != nil {
		for k, v := range player.Stats {
			normalized.Stats[k] = v
		}
	}
	if _, exists := normalized.Stats["points"]; !exists {
		normalized.Stats["points"] = 0
	}
	if _, exists := normalized.Stats["assists"]; !exists {
		normalized.Stats["assists"] = 0
	}
	if _, exists := normalized.Stats["saves"]; !exists {
		normalized.Stats["saves"] = 0
	}
	if _, exists := normalized.Stats["stuns"]; !exists {
		normalized.Stats["stuns"] = 0
	}
	if _, exists := normalized.Stats["ping"]; !exists {
		normalized.Stats["ping"] = float64(player.Ping)
	}
	if normalized.MVP == 0 {
		if stat, ok := stats.Players[player.Username]; ok {
			normalized.MVP = stat.MVPPercentage
		}
		if normalized.MVP == 0 {
			if stat, ok := stats.Players[player.Name]; ok {
				normalized.MVP = stat.MVPPercentage
			}
		}
	}
	if normalized.MVP > 0 || normalized.Stats["mvp"] == 0 {
		normalized.Stats["mvp"] = normalized.MVP
	}
	return normalized
}

func normalizeClock(raw string) string {
	if raw == "" {
		return "00:00.00"
	}
	if strings.HasSuffix(raw, "000") {
		return raw[:len(raw)-3] + "." + raw[len(raw)-3:]
	}
	return raw
}
