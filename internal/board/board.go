package board

import (
	"fmt"
	"strings"

	"echo-scoreboard-visual/internal/echo"
)

type Board struct {
	WaitingForMatch bool        `json:"waitingForMatch"`
	Message         string      `json:"message,omitempty"`
	MVPName         string      `json:"mvpName,omitempty"`
	MVPScore        float64     `json:"mvpScore,omitempty"`
	MVPAvailable    bool        `json:"mvpAvailable"`
	MVPComplete     bool        `json:"mvpComplete"`
	GameStatus      string      `json:"gameStatus,omitempty"`
	GameOver        bool        `json:"gameOver,omitempty"`
	RoundOver       bool        `json:"roundOver,omitempty"`
	SnapshotID      string      `json:"snapshotId,omitempty"`
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
	Name         string             `json:"name,omitempty"`
	Username     string             `json:"username,omitempty"`
	Team         string             `json:"team,omitempty"`
	Ping         int                `json:"ping,omitempty"`
	MVPScore     float64            `json:"mvpScore,omitempty"`
	MVPAvailable bool               `json:"mvpAvailable"`
	MVPComplete  bool               `json:"mvpComplete"`
	IsMVP        bool               `json:"isMVP,omitempty"`
	Stats        map[string]float64 `json:"stats,omitempty"`
}

func Build(session echo.Session) Board {
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
	board.RoundOver = board.GameOver || strings.EqualFold(board.GameStatus, "round_over")
	if board.RoundOver {
		board.SnapshotID = strings.Join([]string{
			session.SessionID,
			fmt.Sprintf("%d-%d", board.Score.Blue, board.Score.Orange),
			fmt.Sprintf("%d-%d", board.RoundScores.Blue, board.RoundScores.Orange),
		}, "-")
	}

	if len(session.Teams) == 0 {
		for _, player := range session.Players {
			board.Players = append(board.Players, normalizePlayer(player, session.MatchType))
		}
		setMVP(&board)
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
			normalized := normalizePlayer(player, session.MatchType)
			normalized.Team = color
			boardTeam.Players = append(boardTeam.Players, normalized)
			board.Players = append(board.Players, normalized)
		}
		board.Teams = append(board.Teams, boardTeam)
	}

	setMVP(&board)
	return board
}

func normalizePlayer(player echo.Player, matchType string) Player {
	mvpScore, mvpComplete := 0.0, false
	mvpAvailable := strings.EqualFold(matchType, "Echo_Arena")
	if mvpAvailable {
		mvpScore, mvpComplete = echoArenaMVPScore(player.Stats)
	}

	normalized := Player{
		Name:         player.Name,
		Username:     player.Username,
		Team:         player.Team,
		Ping:         player.Ping,
		MVPScore:     mvpScore,
		MVPAvailable: mvpAvailable,
		MVPComplete:  mvpComplete,
		Stats:        map[string]float64{},
	}
	for key, value := range player.Stats {
		normalized.Stats[key] = value
	}
	for _, key := range []string{"points", "assists", "saves", "stuns", "ping"} {
		if _, exists := normalized.Stats[key]; !exists {
			if key == "ping" {
				normalized.Stats[key] = float64(player.Ping)
			} else {
				normalized.Stats[key] = 0
			}
		}
	}
	return normalized
}

func setMVP(board *Board) {
	if len(board.Players) == 0 {
		return
	}
	board.MVPAvailable = true
	board.MVPComplete = true
	board.MVPScore = board.Players[0].MVPScore
	for index := range board.Players {
		player := &board.Players[index]
		if !player.MVPAvailable {
			board.MVPAvailable = false
		}
		if !player.MVPComplete {
			board.MVPComplete = false
		}
		if player.MVPScore > board.MVPScore {
			board.MVPScore = player.MVPScore
		}
	}
	if !board.MVPAvailable {
		board.MVPComplete = false
		return
	}

	winners := make([]string, 0, 1)
	for index := range board.Players {
		player := &board.Players[index]
		player.IsMVP = player.MVPScore == board.MVPScore
		if player.IsMVP {
			name := player.Username
			if name == "" {
				name = player.Name
			}
			winners = append(winners, name)
		}
	}
	board.MVPName = strings.Join(winners, ", ")

	for teamIndex := range board.Teams {
		for playerIndex := range board.Teams[teamIndex].Players {
			teamPlayer := &board.Teams[teamIndex].Players[playerIndex]
			for _, player := range board.Players {
				if player.Name == teamPlayer.Name && player.Team == teamPlayer.Team {
					teamPlayer.IsMVP = player.IsMVP
					break
				}
			}
		}
	}
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
