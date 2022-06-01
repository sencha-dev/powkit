// ethash: C/C++ implementation of Ethash, the Ethereum Proof of Work algorithm.
// Copyright 2018-2019 Pawel Bylica.
// Licensed under the Apache License, Version 2.0.

package kawpow

import (
	"encoding/binary"

	"github.com/sencha-dev/powkit/internal/common/convutil"
	"github.com/sencha-dev/powkit/internal/crypto"
	"github.com/sencha-dev/powkit/internal/progpow"
)

var ravencoinKawpow [15]uint32 = [15]uint32{
	0x00000072, //R
	0x00000041, //A
	0x00000056, //V
	0x00000045, //E
	0x0000004E, //N
	0x00000043, //C
	0x0000004F, //O
	0x00000049, //I
	0x0000004E, //N
	0x0000004B, //K
	0x00000041, //A
	0x00000057, //W
	0x00000050, //P
	0x0000004F, //O
	0x00000057, //W
}

func initialize(hash []byte, nonce uint64) ([25]uint32, uint64) {
	var seed [25]uint32
	for i := 0; i < 8; i++ {
		seed[i] = binary.LittleEndian.Uint32(hash[i*4 : i*4+4])
	}

	seed[8] = uint32(nonce)
	seed[9] = uint32(nonce >> 32)

	for i := 10; i < 25; i++ {
		seed[i] = ravencoinKawpow[i-10]
	}

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

	for i := 16; i < 25; i++ {
		state[i] = ravencoinKawpow[i-16]
	}

	crypto.KeccakF800(&state)

	return convutil.Uint32ArrayToBytes(state[:8], binary.LittleEndian)
}

func kawpow(hash []byte, height, nonce, datasetSize uint64, lookup func(index uint32) []uint32, l1 []uint32) ([]byte, []byte) {
	var cfg = &progpow.Config{
		PeriodLength:        3,
		DagLoads:            4,
		CacheBytes:          16 * 1024,
		LaneCount:           16,
		RegisterCount:       32,
		RoundCount:          64,
		RoundCacheAccesses:  11,
		RoundMathOperations: 18,
	}

	seed, seedHead := initialize(hash, nonce)
	mixHash := progpow.Hash(cfg, height, seedHead, datasetSize, lookup, l1)
	digest := finalize(seed, mixHash)

	return mixHash, digest
}
