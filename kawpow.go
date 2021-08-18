package pow

import (
	"encoding/binary"
	"runtime"
)

func kawpow(l1 []uint32, hash []byte, height, nonce uint64, lookup func(index uint32) []uint32) ([]byte, []byte) {
	// temporary initialization state
	var tempState [25]uint32
	for i := 0; i < 8; i += 1 {
		tempState[i] = binary.LittleEndian.Uint32(hash[i*4 : i*4+4])
	}

	tempState[8] = uint32(nonce)
	tempState[9] = uint32(nonce >> 32)

	for i := 10; i < 25; i++ {
		tempState[i] = ravencoinKawpow[i-10]
	}

	keccakF800(&tempState)

	// mixhash
	seedHead := uint64(tempState[0]) + (uint64(tempState[1]) << 32)
	mixHash := hashProgpowMix(l1, height, seedHead, lookup)

	// final hashed digest
	var state [25]uint32
	for i := 0; i < 8; i++ {
		state[i] = tempState[i]
		state[i+8] = binary.LittleEndian.Uint32(mixHash[i*4 : i*4+4])

	}

	for i := 16; i < 25; i++ {
		state[i] = ravencoinKawpow[i-16]
	}

	keccakF800(&state)

	digest := uint32ArrayToBytes(state[:8])

	return mixHash, digest
}

func (dag *LightDag) kawpowLight(height, nonce uint64, hash []byte) ([]byte, []byte) {
	epoch := calcEpoch(height, dag.EpochLength)
	cache := dag.getCache(epoch)

	keccak512Hasher := newKeccak512Hasher()
	lookup := func(index uint32) []uint32 {
		return generateDatasetItem2048(cache.cache, index, keccak512Hasher, dag.DatasetParents)
	}

	mix, digest := kawpow(cache.l1, hash, height, nonce, lookup)
	runtime.KeepAlive(cache)

	return mix, digest
}
