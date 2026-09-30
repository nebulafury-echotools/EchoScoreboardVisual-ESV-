package main

import (
	"log"
	"sync"
	"time"

	"echo-scoreboard-visual/internal/board"
	"echo-scoreboard-visual/internal/echo"
	"echo-scoreboard-visual/internal/snapshot"
	scoreboardweb "echo-scoreboard-visual/web"

	webview2 "github.com/jchv/go-webview2"
)

type BoardPoller struct {
	mu          sync.RWMutex
	latest      board.Board
	staleAt     time.Time
	projectRoot string
}

func main() {
	projectRoot, err := snapshot.ProjectRoot()
	if err != nil {
		log.Fatalf("find project directory for match data: %v", err)
	}
	poller := NewBoardPoller(projectRoot)
	poller.Refresh()
	go poller.Run()

	window := webview2.NewWithOptions(webview2.WebViewOptions{
		AutoFocus: true,
		WindowOptions: webview2.WindowOptions{
			Title:  "Echo Scoreboard",
			Width:  1100,
			Height: 800,
			Center: true,
		},
	})
	if window == nil {
		log.Fatal("failed to create scoreboard window; verify the WebView2 runtime is installed")
	}
	defer window.Destroy()

	recorder := snapshot.NewRecorder()
	if err := window.Bind("getBoard", poller.Snapshot); err != nil {
		log.Fatalf("bind board data: %v", err)
	}
	if err := window.Bind("saveScoreboardScreenshot", func(eventID string) (string, error) {
		return recorder.Save(window.Window(), eventID)
	}); err != nil {
		log.Fatalf("bind screenshot callback: %v", err)
	}

	window.SetHtml(string(scoreboardweb.IndexHTML))
	window.Run()
}

func NewBoardPoller(projectRoot string) *BoardPoller {
	return &BoardPoller{
		latest:      board.Board{WaitingForMatch: true, Message: "waiting for match"},
		projectRoot: projectRoot,
	}
}

func (p *BoardPoller) Run() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		p.Refresh()
	}
}

func (p *BoardPoller) Snapshot() board.Board {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.latest
}

func (p *BoardPoller) Refresh() {
	echoClient := echo.NewClient("http://127.0.0.1:6721/session")
	session, rawSession, err := echoClient.FetchWithRaw()
	if err != nil {
		p.mu.Lock()
		p.latest = board.Board{WaitingForMatch: true, Message: "waiting for match"}
		p.staleAt = time.Now()
		p.mu.Unlock()
		return
	}
	if err := echo.SaveSnapshot(p.projectRoot, rawSession); err != nil {
		log.Printf("save populated Echo session snapshot: %v", err)
	}

	merged := board.Build(session)
	p.mu.Lock()
	merged.StaleSince = ""
	p.latest = merged
	p.staleAt = time.Time{}
	p.mu.Unlock()
}
