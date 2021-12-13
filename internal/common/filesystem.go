package common

import (
	"os"
	"os/user"
	"path/filepath"
	"runtime"
)

func DefaultDir(name string) string {
	home := os.Getenv("HOME")
	if user, err := user.Current(); err == nil {
		home = user.HomeDir
	}
	if runtime.GOOS == "windows" {
		return filepath.Join(home, "AppData", "PowCache")
	}
	return filepath.Join(home, name)
}
