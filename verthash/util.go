package verthash

import (
	"encoding/hex"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
)

func defaultDir() string {
	home := os.Getenv("HOME")
	if user, err := user.Current(); err == nil {
		home = user.HomeDir
	}
	if runtime.GOOS == "windows" {
		return filepath.Join(home, "AppData", "PowCache")
	}
	return filepath.Join(home, ".powcache")
}

// should only be used for tests
func mustDecodeHex(inp string) []byte {
	inp = strings.Replace(inp, "0x", "", -1)
	out, err := hex.DecodeString(inp)
	if err != nil {
		panic(err)
	}

	return out
}

func log2(x int64) int64 {
	r := int64(0)
	for ; x > 1; x >>= 1 {
		r++
	}

	return r
}
