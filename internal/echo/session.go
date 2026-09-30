package echo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	URL        string
	HTTPClient *http.Client
}

func NewClient(url string) *Client {
	return &Client{
		URL:        url,
		HTTPClient: &http.Client{Timeout: 3 * time.Second},
	}
}

type Session struct {
	SessionID        string   `json:"sessionid,omitempty"`
	MatchType        string   `json:"match_type,omitempty"`
	GameState        string   `json:"game_status,omitempty"`
	ClockDisplay     string   `json:"game_clock_display,omitempty"`
	BluePoints       int      `json:"blue_points,omitempty"`
	OrangePoints     int      `json:"orange_points,omitempty"`
	BlueRoundScore   int      `json:"blue_round_score,omitempty"`
	OrangeRoundScore int      `json:"orange_round_score,omitempty"`
	Teams            []Team   `json:"teams,omitempty"`
	Players          []Player `json:"players,omitempty"`
}

type Score struct {
	Blue   int `json:"blue_points,omitempty"`
	Orange int `json:"orange_points,omitempty"`
}

type Team struct {
	Name    string   `json:"team,omitempty"`
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

func ParseSession(data []byte) (Session, error) {
	data = bytes.TrimPrefix(data, []byte{0xEF, 0xBB, 0xBF})
	var session Session
	if err := json.Unmarshal(data, &session); err != nil {
		return Session{}, err
	}
	return session, nil
}

func (c *Client) Fetch() (Session, error) {
	session, _, err := c.FetchWithRaw()
	return session, err
}

func (c *Client) FetchWithRaw() (Session, []byte, error) {
	resp, err := c.HTTPClient.Get(c.URL)
	if err != nil {
		return Session{}, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Session{}, nil, fmt.Errorf("session endpoint returned %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return Session{}, nil, err
	}
	session, err := ParseSession(data)
	if err != nil {
		return Session{}, nil, err
	}
	return session, data, nil
}
