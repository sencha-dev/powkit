package testutil

import (
	"encoding/hex"
	"strings"
)

// should only be used for tests
func MustDecodeHex(inp string) []byte {
	inp = strings.Replace(inp, "0x", "", -1)
	out, err := hex.DecodeString(inp)
	if err != nil {
		panic(err)
	}

	return out
}
