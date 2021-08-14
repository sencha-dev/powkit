package pow

import (
	"encoding/binary"
)

func kawpow(l1 []uint32, hash []byte, height, nonce uint64, lookup func(index uint32) []uint32) ([]byte, []byte) {
	var state2 [8]uint32
	{
		var state [25]uint32
		for i := 0; i < 8; i += 1 {
			state[i] = binary.LittleEndian.Uint32(hash[i*4 : i*4+4])
		}

		state[8] = uint32(nonce)
		state[9] = uint32(nonce >> 32)

		for i := 10; i < 25; i++ {
			state[i] = ravencoinKawpow[i-10]
		}

		KeccakF800(&state)

		for i := 0; i < 8; i++ {
			state2[i] = state[i]
		}
	}

	seedHead := [2]uint32{state2[0], state2[1]}
	mixHash := hashMix(l1, height, seedHead, lookup)

	var state [25]uint32

	for i := 0; i < 8; i++ {
		state[i] = state2[i]
	}

	for i := 0; i < 8; i++ {
		state[i+8] = binary.LittleEndian.Uint32(mixHash[i*4 : i*4+4])
	}

	for i := 16; i < 25; i++ {
		state[i] = ravencoinKawpow[i-16]
	}

	KeccakF800(&state)

	digest := uint32ArrayToBytes(state[:8])

	return mixHash, digest
}

// hashimotoLight aggregates data from the full dataset (using only a small
// in-memory cache) in order to produce our final value for a particular header
// hash and nonce.
func kawpowLight(size, height, nonce uint64, cache []uint32, hash []byte) ([]byte, []byte) {
	keccak512Hasher := NewKeccak512Hasher()
	lookup := func(index uint32) []uint32 {
		return generateDatasetItem2048(cache, index, keccak512Hasher, datasetParentsRVN)
	}

	l1 := generateL1Cache(cache)

	return kawpow(l1, hash, height, nonce, lookup)
}
