package equihash

import (
	"encoding/hex"

	"github.com/sencha-dev/go-pow/internal/crypto"
)

const (
	SEED_LENGTH     = 16 //Length of seed in dwords ;
	NONCE_LENGTH    = 24 //Length of nonce in bytes;
	MAX_NONCE       = 0xFFFFF
	MAX_N           = 32 //Max length of n in bytes, should not exceed 32
	LIST_LENGTH     = 5
	FORK_MULTIPLIER = 3 //Maximum collision factor
)

type Equihash struct {
	n        uint64
	k        uint64
	personal []byte
}

func New(n, k uint64, personal string) (*Equihash, error) {
	personalBytes, err := hex.DecodeString(personal)
	if err != nil {
		return nil, err
	}

	equihash := &Equihash{
		n:        n,
		k:        k,
		personal: personalBytes,
	}

	return equihash, nil
}

func (e *Equihash) Verify(seedBytes, inputBytes []byte, nonce uint32) bool {
	seed := make([]uint32, SEED_LENGTH)
	partialSeed := bytesToUint32Array(seedBytes)
	for i, val := range partialSeed {
		seed[i] = val
	}

	inputs := bytesToUint32Array(inputBytes)

	input := make([]uint32, SEED_LENGTH+2)
	for i := 0; i < SEED_LENGTH; i++ {
		input[i] = seed[i]
	}

	input[SEED_LENGTH] = nonce
	input[SEED_LENGTH+1] = 0

	buf := make([]uint32, MAX_N/4)
	blocks := make([]uint32, e.k+1)

	for i := 0; i < len(inputs); i++ {
		input[SEED_LENGTH+1] = inputs[i]
		inputBytes := uint32ArrayToBytes(input)
		tempBuf := crypto.Blake2b(inputBytes, e.personal, MAX_N)
		buf = bytesToUint32Array(tempBuf)

		for j := 0; j < (int(e.k) + 1); j++ {
			blocks[j] = blocks[j] ^ (buf[j] >> (32 - e.n/(e.k+1)))
		}
	}

	for j := 0; j < (int(e.k) + 1); j++ {
		if blocks[j] != 0 {
			return false
		}
	}

	return true
}
