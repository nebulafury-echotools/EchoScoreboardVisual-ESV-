package snapshot

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestWritePNGCreatesReadableImage(t *testing.T) {
	directory := t.TempDir()
	imageData := image.NewRGBA(image.Rect(0, 0, 2, 2))
	imageData.SetRGBA(0, 0, color.RGBA{R: 255, A: 255})

	path, err := writePNG(imageData, "test-event", directory)
	if err != nil {
		t.Fatalf("write PNG: %v", err)
	}
	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("open saved PNG: %v", err)
	}
	defer file.Close()

	decoded, err := png.Decode(file)
	if err != nil {
		t.Fatalf("decode saved PNG: %v", err)
	}
	if got := decoded.Bounds().Size(); got != image.Pt(2, 2) {
		t.Fatalf("expected 2x2 PNG, got %v", got)
	}
	if !strings.Contains(filepath.Base(path), "_EST_") {
		t.Fatalf("expected EST in screenshot filename, got %q", filepath.Base(path))
	}
}

func TestESTTimestampUsesFixedUTCMinusFive(t *testing.T) {
	instant := time.Date(2026, 9, 30, 12, 0, 0, 0, time.UTC)
	if got, want := estTimestamp(instant), "20260930_070000"; got != want {
		t.Fatalf("expected EST timestamp %q, got %q", want, got)
	}
}

func TestProjectRootUsesBuildDirectoryParent(t *testing.T) {
	root := filepath.Join(t.TempDir(), "EchoScoreboardVisual")
	executable := filepath.Join(root, "builds", "scoreboard.exe")
	if got := projectRoot(executable, t.TempDir()); got != root {
		t.Fatalf("expected project root %q, got %q", root, got)
	}
}
