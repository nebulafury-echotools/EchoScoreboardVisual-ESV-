package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"echo-scoreboard-visual/internal/board"
	"echo-scoreboard-visual/internal/echo"
	"echo-scoreboard-visual/internal/spark"
)

const listenAddr = ":8080"

func main() {
	poller := NewBoardPoller(5 * time.Second)
	go poller.Run()

	mux := http.NewServeMux()
	mux.HandleFunc("/board.json", poller.HandleBoardJSON)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		data, err := os.ReadFile("web/index.html")
		if err != nil {
			http.Error(w, "failed to load page", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(data)
	})

	log.Printf("scoreboard listening on http://localhost%s", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, mux))
}

type BoardPoller struct {
	mu      sync.RWMutex
	latest  board.Board
	staleAt time.Time
}

func NewBoardPoller(interval time.Duration) *BoardPoller {
	poller := &BoardPoller{}
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			poller.Refresh()
		}
	}()
	return poller
}

func (p *BoardPoller) Run() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		p.Refresh()
	}
}

func (p *BoardPoller) Refresh() {
	echoClient := echo.NewClient("http://127.0.0.1:6721/session")
	session, err := echoClient.Fetch()
	if err != nil {
		p.mu.Lock()
		p.latest = board.Board{WaitingForMatch: true, Message: "waiting for match"}
		p.staleAt = time.Now()
		p.mu.Unlock()
		return
	}

	sparkClient := spark.NewClient("http://localhost:6724")
	stats, err := sparkClient.FetchStats()
	if err != nil {
		stats = spark.Stats{}
	}

	merged := board.Build(session, stats)
	p.mu.Lock()
	p.latest = merged
	p.staleAt = time.Time{}
	p.mu.Unlock()
}

func (p *BoardPoller) HandleBoardJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	p.mu.RLock()
	response := p.latest
	staleSince := p.staleAt
	p.mu.RUnlock()

	if staleSince.IsZero() {
		response.StaleSince = ""
	} else {
		response.StaleSince = time.Since(staleSince).Round(time.Second).String()
	}

	if response.Message == "" && response.WaitingForMatch {
		response.Message = "waiting for match"
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("encode board response: %v", err)
		http.Error(w, fmt.Sprintf("encode response: %v", err), http.StatusInternalServerError)
	}
}
