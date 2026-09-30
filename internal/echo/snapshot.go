package echo

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

func SaveSnapshot(projectRoot string, data []byte) error {
	session, err := ParseSession(data)
	if err != nil {
		return err
	}
	if session.SessionID == "" || session.MatchType == "" || len(session.Teams) == 0 {
		return errors.New("refusing to save an unpopulated Echo session")
	}

	matchID := safeSessionID(session.SessionID)
	if matchID == "" {
		return errors.New("Echo session ID contains no filename-safe characters")
	}
	dataDirectory := filepath.Join(projectRoot, "match-data")
	matchDirectory := filepath.Join(dataDirectory, "matches")
	if err := os.MkdirAll(matchDirectory, 0o700); err != nil {
		return err
	}
	if err := atomicWrite(filepath.Join(dataDirectory, "echo_session.json"), data); err != nil {
		return err
	}
	return atomicWrite(filepath.Join(matchDirectory, matchID+".json"), data)
}

func safeSessionID(sessionID string) string {
	var cleaned strings.Builder
	for _, character := range sessionID {
		switch {
		case character >= 'a' && character <= 'z', character >= 'A' && character <= 'Z', character >= '0' && character <= '9', character == '-', character == '_':
			cleaned.WriteRune(character)
		}
	}
	return cleaned.String()
}

func atomicWrite(path string, data []byte) error {
	file, err := os.CreateTemp(filepath.Dir(path), ".echo-session-*.tmp")
	if err != nil {
		return err
	}
	temporaryPath := file.Name()
	defer os.Remove(temporaryPath)
	if err := file.Chmod(0o600); err != nil {
		_ = file.Close()
		return err
	}
	if _, err := file.Write(data); err != nil {
		_ = file.Close()
		return err
	}
	if err := file.Sync(); err != nil {
		_ = file.Close()
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	return os.Rename(temporaryPath, path)
}
