// ethash: C/C++ implementation of Ethash, the Ethereum Proof of Work algorithm.
// Copyright 2018-2019 Pawel Bylica.
// Licensed under the Apache License, Version 2.0.

package progpow

import (
	"encoding/binary"

	"github.com/sencha-dev/powkit/internal/common/convutil"
	"github.com/sencha-dev/powkit/internal/crypto"
	"github.com/sencha-dev/powkit/internal/dag"
)

const (
	fnvOffsetBasis uint32 = 0x811c9dc5
	progpowRegs    uint32 = 32
	progpowLanes   uint32 = 16
)

type mixArray [progpowLanes][progpowRegs]uint32

type lookupFunc func(index uint32) []uint32

func progpowSeed(hash []byte, nonce uint64) uint64 {
	var tempState [25]uint32
	for i := 0; i < 8; i += 1 {
		tempState[i] = binary.LittleEndian.Uint32(hash[i*4 : i*4+4])
	}

	tempState[8] = uint32(nonce)
	tempState[9] = uint32(nonce >> 32)

	crypto.KeccakF800(&tempState)

	seedHead := uint64(tempState[0]) + (uint64(tempState[1]) << 32)

	return seedHead
}

func initMix(seed uint64) mixArray {
	z := crypto.Fnv1a(fnvOffsetBasis, uint32(seed))
	w := crypto.Fnv1a(z, uint32(seed>>32))

	var mix mixArray
	for l := range mix {
		jsr := crypto.Fnv1a(w, uint32(l))
		jcong := crypto.Fnv1a(jsr, uint32(l))

		rng := newKiss99(z, w, jsr, jcong)

		for r := range mix[l] {
			mix[l][r] = rng.Next()
		}
	}

	return mix
}

func round(cfg *Config, seed uint64, r uint32, mix mixArray, datasetSize uint64, lookup lookupFunc, l1 []uint32) mixArray {
	state := initMixRngState(seed)
	numItems := uint32(datasetSize / (2 * 128))
	itemIndex := mix[r%progpowLanes][0] % numItems

	item := lookup(itemIndex)

	numWordsPerLane := len(item) / int(progpowLanes)
	maxOperations := max(cfg.RoundCacheAccesses, cfg.RoundMathOperations)
	for i := 0; i < maxOperations; i++ {
		if i < cfg.RoundCacheAccesses {
			src := state.nextSrc()
			dst := state.nextDst()
			sel := state.Rng.Next()

			for l := 0; l < int(progpowLanes); l++ {
				offset := mix[l][src] % (cfg.CacheBytes / 4)
				mix[l][dst] = randomMerge(mix[l][dst], l1[offset], sel)
			}
		}

		if i < cfg.RoundMathOperations {
			srcRand := state.Rng.Next() % (progpowRegs * (progpowRegs - 1))
			src1 := srcRand % progpowRegs
			src2 := srcRand / progpowRegs
			if src2 >= src1 {
				src2 += 1
			}

			sel1 := state.Rng.Next()
			dst := state.nextDst()
			sel2 := state.Rng.Next()

			for l := 0; l < int(progpowLanes); l++ {
				data := randomMath(mix[l][src1], mix[l][src2], sel1)
				mix[l][dst] = randomMerge(mix[l][dst], data, sel2)
			}
		}
	}

	// DAG access pattern.
	dsts := make([]uint32, numWordsPerLane)
	sels := make([]uint32, numWordsPerLane)
	for i := 0; i < numWordsPerLane; i++ {
		if i == 0 {
			dsts[i] = 0
		} else {
			dsts[i] = state.nextDst()
		}

		sels[i] = state.Rng.Next()
	}

	for l := 0; l < int(progpowLanes); l++ {
		offset := ((uint32(l) ^ r) % progpowLanes) * uint32(numWordsPerLane)
		for i := 0; i < numWordsPerLane; i++ {
			word := item[offset+uint32(i)]
			mix[l][dsts[i]] = randomMerge(mix[l][dsts[i]], word, sels[i])
		}
	}

	return mix
}

func HashMix(cfg *Config, height, seed uint64, lookup lookupFunc, l1 []uint32) []byte {
	mix := initMix(seed)

	number := height / cfg.PeriodLength
	epoch := dag.CalcEpoch(cfg.DagCfg, height)
	datasetSize := dag.DatasetSize(cfg.DagCfg, epoch)

	for i := 0; i < cfg.RoundCount; i++ {
		mix = round(cfg, number, uint32(i), mix, datasetSize, lookup, l1)
	}

	var laneHash [progpowLanes]uint32
	for l := 0; l < int(progpowLanes); l++ {
		laneHash[l] = fnvOffsetBasis

		for i := 0; i < int(progpowRegs); i++ {
			laneHash[l] = crypto.Fnv1a(laneHash[l], mix[l][i])
		}
	}

	numWords := 8
	mixHash := make([]uint32, numWords)
	for i := 0; i < numWords; i++ {
		mixHash[i] = fnvOffsetBasis
	}

	for l := 0; l < int(progpowLanes); l++ {
		mixHash[l%numWords] = crypto.Fnv1a(mixHash[l%numWords], laneHash[l])
	}

	return convutil.Uint32ArrayToBytes(mixHash)
}
