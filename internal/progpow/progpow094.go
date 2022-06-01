// ethash: C/C++ implementation of Ethash, the Ethereum Proof of Work algorithm.
// Copyright 2018-2019 Pawel Bylica.
// Licensed under the Apache License, Version 2.0.

package progpow

import (
	"encoding/binary"

	"github.com/sencha-dev/powkit/internal/common/convutil"
	"github.com/sencha-dev/powkit/internal/crypto"
)

var progpow094Cfg = &Config{
	PeriodLength:        10,
	DagLoads:            4,
	CacheBytes:          16 * 1024,
	LaneCount:           16,
	RegisterCount:       32,
	RoundCount:          64,
	RoundCacheAccesses:  11,
	RoundMathOperations: 18,
}

func initialize(hash []byte, nonce uint64) ([25]uint32, uint64) {
	var seed [25]uint32
	for i := 0; i < 8; i += 1 {
		seed[i] = binary.LittleEndian.Uint32(hash[i*4 : i*4+4])
	}

	seed[8] = uint32(nonce)
	seed[9] = uint32(nonce >> 32)

	crypto.KeccakF800(&seed)

	seedHead := uint64(seed[0]) + (uint64(seed[1]) << 32)

	return seed, seedHead
}

func finalize(seed [25]uint32, mixHash []byte) []byte {
	var state [25]uint32
	for i := 0; i < 8; i++ {
		state[i] = seed[i]
		state[i+8] = binary.LittleEndian.Uint32(mixHash[i*4 : i*4+4])
	}

	crypto.KeccakF800(&state)

	return convutil.Uint32ArrayToBytes(state[:8], binary.LittleEndian)
}

func compute(hash []byte, height, nonce, datasetSize uint64, lookup func(index uint32) []uint32, l1 []uint32) ([]byte, []byte) {
	seed, seedHead := initialize(hash, nonce)
	mixHash := Hash(progpow094Cfg, height, seedHead, datasetSize, lookup, l1)
	digest := finalize(seed, mixHash)

	return mixHash, digest
}
