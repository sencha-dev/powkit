package ethash

import (
	"encoding/binary"

	"github.com/sencha-dev/go-pow/internal/crypto"
)

const (
	mixBytes  = 128 // Width of mix
	hashBytes = 64  // Hash length in bytes
	hashWords = 16  // Number of 32 bit ints in a hash

	// hashimoto constants
	hashimotoRounds = 64 // Number of accesses in hashimoto loop
)

// hashimoto aggregates data from the full dataset in order to produce our final
// value for a particular header hash and nonce.
func hashimoto(hash []byte, nonce uint64, size uint64, lookup func(index uint32) []uint32) ([]byte, []byte) {
	// Calculate the number of theoretical rows (we use one buffer nonetheless)
	rows := uint32(size / mixBytes)

	// Combine header+nonce into a 64 byte seed
	seed := make([]byte, 40)
	copy(seed, hash)
	binary.LittleEndian.PutUint64(seed[32:], nonce)

	seed = crypto.Keccak512(seed)
	seedHead := binary.LittleEndian.Uint32(seed)

	// Start the mix with replicated seed
	mix := make([]uint32, mixBytes/4)
	for i := 0; i < len(mix); i++ {
		mix[i] = binary.LittleEndian.Uint32(seed[i%16*4:])
	}
	// Mix in random dataset nodes
	temp := make([]uint32, len(mix))

	for i := 0; i < hashimotoRounds; i++ {
		parent := crypto.Fnv1(uint32(i)^seedHead, mix[i%len(mix)]) % rows
		for j := uint32(0); j < mixBytes/hashBytes; j++ {
			copy(temp[j*hashWords:], lookup(2*parent+j))
		}
		crypto.FnvHash(mix, temp)
	}
	// Compress mix
	for i := 0; i < len(mix); i += 4 {
		mix[i/4] = crypto.Fnv1(crypto.Fnv1(crypto.Fnv1(mix[i], mix[i+1]), mix[i+2]), mix[i+3])
	}
	mix = mix[:len(mix)/4]

	digest := make([]byte, 32)
	for i, val := range mix {
		binary.LittleEndian.PutUint32(digest[i*4:], val)
	}

	return digest, crypto.Keccak256(append(seed, digest...))
}
