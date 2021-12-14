// ethash: C/C++ implementation of Ethash, the Ethereum Proof of Work algorithm.
// Copyright 2018-2019 Pawel Bylica.
// Licensed under the Apache License, Version 2.0.

package firopow

import (
	"encoding/binary"

	"github.com/sencha-dev/powkit/internal/common/convutil"
	"github.com/sencha-dev/powkit/internal/crypto"
	"github.com/sencha-dev/powkit/internal/progpow"
)

func firopow(hash []byte, height, nonce uint64, lookup func(index uint32) []uint32, l1 []uint32) ([]byte, []byte) {
	// temporary initialization state
	var tempState [25]uint32
	for i := 0; i < 8; i += 1 {
		tempState[i] = binary.LittleEndian.Uint32(hash[i*4 : i*4+4])
	}

	tempState[8] = uint32(nonce)
	tempState[9] = uint32(nonce >> 32)
	tempState[10] = 0x00000001
	tempState[18] = 0x80008081

	// mixhash
	crypto.KeccakF800(&tempState)
	seedHead := uint64(tempState[0]) + (uint64(tempState[1]) << 32)

	mixHash := progpow.HashMix(progpow.Firopow, height, seedHead, lookup, l1)

	// final hashed digest
	var state [25]uint32
	for i := 0; i < 8; i++ {
		state[i] = tempState[i]
		state[i+8] = binary.LittleEndian.Uint32(mixHash[i*4 : i*4+4])
	}

	state[17] = 0x00000001
	state[24] = 0x80008081

	crypto.KeccakF800(&state)

	digest := convutil.Uint32ArrayToBytesLE(state[:8])

	return mixHash, digest
}
