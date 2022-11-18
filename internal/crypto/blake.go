package crypto

import (
	"github.com/dchest/blake2b"
)

func Blake2b(data, personal []byte, size int) []byte {
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

	h.Write(data)

	return h.Sum(nil)
}

func Blake2b256(data []byte) []byte {
	out := blake2b.Sum256(data)
	return out[:]
}
