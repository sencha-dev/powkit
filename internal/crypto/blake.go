package crypto

import (
	"github.com/codahale/blake2"
)

func Blake2b(inp, personal []byte, size int) []byte {
	h := blake2.New(&blake2.Config{
		Size:     uint8(size),
		Key:      nil,
		Salt:     nil,
		Personal: personal,
	})
	h.Write(inp)

	return h.Sum(nil)
}
