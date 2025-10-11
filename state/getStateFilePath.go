package state

import (
	"os"
	"path/filepath"
	"runtime"
)

func getStateFilePath() (string, error) {
	var baseDir string
	appName := "P2Podium"

	if runtime.GOOS == "windows" {
		baseDir = os.Getenv("APPDATA")
		if baseDir == "" {
			return "", fmt.Errorf("APPDATA not set")
		}
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		baseDir = filepath.Join(home, ".local", "share")
	}

	dir := filepath.Join(baseDir, appName)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", err
	}

	return filepath.Join(dir, "state.json"), nil
}
