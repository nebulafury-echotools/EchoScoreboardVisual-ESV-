package snapshot

import (
	"image"
	"os"
	"path/filepath"
	"sync"
	"syscall"
	"unsafe"

	"github.com/kbinani/screenshot"
)

type windowRect struct {
	Left   int32
	Top    int32
	Right  int32
	Bottom int32
}

var getWindowRect = syscall.NewLazyDLL("user32.dll").NewProc("GetWindowRect")

type Recorder struct {
	mu    sync.Mutex
	saved map[string]string
}

func NewRecorder() *Recorder {
	return &Recorder{saved: map[string]string{}}
}

func (r *Recorder) Save(window unsafe.Pointer, eventID string) (string, error) {
	if eventID == "" {
		return "", nil
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if path, exists := r.saved[eventID]; exists {
		return path, nil
	}

	var bounds windowRect
	result, _, callErr := getWindowRect.Call(uintptr(window), uintptr(unsafe.Pointer(&bounds)))
	if result == 0 {
		return "", callErr
	}
	area := image.Rect(int(bounds.Left), int(bounds.Top), int(bounds.Right), int(bounds.Bottom))
	if area.Empty() {
		return "", syscall.EINVAL
	}
	imageData, err := screenshot.CaptureRect(area)
	if err != nil {
		return "", err
	}

	baseDirectory := os.Getenv("LOCALAPPDATA")
	if baseDirectory == "" {
		baseDirectory, err = os.UserConfigDir()
		if err != nil {
			return "", err
		}
	}
	directory := filepath.Join(baseDirectory, "EchoScoreboardVisual", "screenshots")
	if err := os.MkdirAll(directory, 0o700); err != nil {
		return "", err
	}

	path, err := writePNG(imageData, eventID, directory)
	if err != nil {
		return "", err
	}
	r.saved[eventID] = path
	return path, nil
}
