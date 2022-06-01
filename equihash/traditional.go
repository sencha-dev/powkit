package equihash

import (
	"encoding/binary"

	"github.com/sencha-dev/powkit/internal/common/convutil"
	"github.com/sencha-dev/powkit/internal/crypto"
)

const (
	SeedLength     = 16 //Length of seed in dwords ;
	NonceLength    = 24 //Length of nonce in bytes;
	MaxNonce       = 0xFFFFF
	MaxN           = 32 //Max length of n in bytes, should not exceed 32
	ListLength     = 5
	ForkMultiplier = 3 //Maximum collision factor
)

func TraditionalVerify(n, k uint32, personalBytes, seedBytes, inputBytes []byte, nonce uint32) bool {
	seed := make([]uint32, SeedLength)
	partialSeed := convutil.BytesToUint32Array(seedBytes, binary.LittleEndian)
	for i, val := range partialSeed {
		seed[i] = val
	}

	inputs := convutil.BytesToUint32Array(inputBytes, binary.LittleEndian)

	input := make([]uint32, SeedLength+2)
	for i := 0; i < SeedLength; i++ {
		input[i] = seed[i]
	}

	input[SeedLength] = nonce
	input[SeedLength+1] = 0

	buf := make([]uint32, MaxN/4)
	blocks := make([]uint32, k+1)

	for i := 0; i < len(inputs); i++ {
		input[SeedLength+1] = inputs[i]
		inputBytes := convutil.Uint32ArrayToBytes(input, binary.LittleEndian)
		tempBuf := crypto.Blake2b(inputBytes, personalBytes, MaxN)
		buf = convutil.BytesToUint32Array(tempBuf, binary.LittleEndian)

		for j := 0; j < (int(k) + 1); j++ {
			blocks[j] = blocks[j] ^ (buf[j] >> (32 - n/(k+1)))
		}
	}

	for j := 0; j < (int(k) + 1); j++ {
		if blocks[j] != 0 {
			return false
		}
	}

	return true
}
