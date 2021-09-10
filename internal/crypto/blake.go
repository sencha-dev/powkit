package crypto

import (
	"github.com/dchest/blake2b"
)

func Blake2b(inp, personal []byte, size int) []byte {
	h, err := blake2b.New(&blake2b.Config{
		Size:   uint8(size),
		Key:    nil,
		Salt:   nil,
		Person: personal,
	})

	// this error only happens on invalid configs
	if err != nil {
		panic(err)
	}

	h.Write(inp)

	return h.Sum(nil)
}
