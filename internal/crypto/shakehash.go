package crypto

import (
	"golang.org/x/crypto/sha3"
)

func CShake256(data, personal []byte, size int) []byte {
	out := make([]byte, size)
	h := sha3.NewCShake256(nil, personal)
	h.Write(data)
	h.Read(out)

	return out
}
