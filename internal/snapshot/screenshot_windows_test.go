package snapshot

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"
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
}
