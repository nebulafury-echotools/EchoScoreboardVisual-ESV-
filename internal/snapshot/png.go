package snapshot

import (
	"crypto/sha256"
	"encoding/hex"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"time"
)

func writePNG(imageData image.Image, eventID, directory string) (string, error) {
	digest := sha256.Sum256([]byte(eventID))
	name := "scoreboard_" + time.Now().Format("20060102_150405") + "_" + hex.EncodeToString(digest[:4]) + ".png"
	path := filepath.Join(directory, name)
	file, err := os.Create(path)
	if err != nil {
		return "", err
	}
	if err := png.Encode(file, imageData); err != nil {
		_ = file.Close()
		_ = os.Remove(path)
		return "", err
	}
	if err := file.Close(); err != nil {
		_ = os.Remove(path)
		return "", err
	}
	return path, nil
}
