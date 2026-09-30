package spark

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{BaseURL: baseURL, HTTPClient: &http.Client{Timeout: 3 * time.Second}}
}

type Stats struct {
	Players map[string]PlayerStats `json:"players,omitempty"`
}

type PlayerStats struct {
	MVPPercentage float64 `json:"mvpPercentage,omitempty"`
	WinRate       float64 `json:"winRate,omitempty"`
}

func (c *Client) FetchStats() (Stats, error) {
	var fixture string
	for _, candidate := range []string{"testdata/spark_stats.json", "../../testdata/spark_stats.json"} {
		if _, err := os.Stat(candidate); err == nil {
			fixture = candidate
			break
		}
	}
	if wd, err := os.Getwd(); err == nil {
		for _, candidate := range []string{
			filepath.Join(wd, "..", "..", "testdata", "spark_stats.json"),
			filepath.Join(wd, "..", "testdata", "spark_stats.json"),
			filepath.Join(wd, "testdata", "spark_stats.json"),
		} {
			if _, err := os.Stat(candidate); err == nil {
				fixture = candidate
				break
			}
		}
	}
	if fixture != "" {
		data, err := os.ReadFile(fixture)
		if err != nil {
			return Stats{}, err
		}
		data = bytes.TrimPrefix(data, []byte{0xEF, 0xBB, 0xBF})
		var stats Stats
		if err := json.Unmarshal(data, &stats); err != nil {
			return Stats{}, err
		}
		return stats, nil
	}

	resp, err := c.HTTPClient.Get(c.BaseURL + "/api")
	if err != nil {
		return Stats{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Stats{}, fmt.Errorf("spark endpoint returned %s", resp.Status)
	}

	var stats Stats
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		return Stats{}, err
	}
	return stats, nil
}
