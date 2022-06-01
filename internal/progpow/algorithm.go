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
)

type Config struct {
	PeriodLength        uint64
	DagLoads            int
	CacheBytes          uint32
	LaneCount           int
	RegisterCount       int
	RoundCount          int
	RoundCacheAccesses  int
	RoundMathOperations int
}

func initMix(seed uint64, numLanes, numRegs int) [][]uint32 {
	z := crypto.Fnv1a(fnvOffsetBasis, uint32(seed))
	w := crypto.Fnv1a(z, uint32(seed>>32))

	mix := make([][]uint32, numLanes)
	for lane := range mix {
		jsr := crypto.Fnv1a(w, uint32(lane))
		jcong := crypto.Fnv1a(jsr, uint32(lane))

		rng := newKiss99(z, w, jsr, jcong)

		mix[lane] = make([]uint32, numRegs)
		for reg := range mix[lane] {
			mix[lane][reg] = rng.next()
		}
	}

	return mix
}

func round(cfg *Config, seed uint64, r uint32, mix [][]uint32, datasetSize uint64, lookup dag.LookupFunc, l1 []uint32) [][]uint32 {
	state := initMixRngState(seed, uint32(cfg.RegisterCount))
	numItems := uint32(datasetSize / (2 * 128))
	itemIndex := mix[r%uint32(cfg.LaneCount)][0] % numItems

	item := lookup(itemIndex)

	numWordsPerLane := len(item) / cfg.LaneCount
	maxOperations := max(cfg.RoundCacheAccesses, cfg.RoundMathOperations)
	for i := 0; i < maxOperations; i++ {
		if i < cfg.RoundCacheAccesses {
			src := state.nextSrc()
			dst := state.nextDst()
			sel := state.nextRng()

			for l := 0; l < cfg.LaneCount; l++ {
				offset := mix[l][src] % (cfg.CacheBytes / 4)
				mix[l][dst] = randomMerge(mix[l][dst], l1[offset], sel)
			}
		}

		if i < cfg.RoundMathOperations {
			srcRand := state.nextRng() % (uint32(cfg.RegisterCount) * uint32(cfg.RegisterCount-1))
			src1 := srcRand % uint32(cfg.RegisterCount)
			src2 := srcRand / uint32(cfg.RegisterCount)
			if src2 >= src1 {
				src2 += 1
			}

			sel1 := state.nextRng()
			dst := state.nextDst()
			sel2 := state.nextRng()

			for l := 0; l < cfg.LaneCount; l++ {
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

		sels[i] = state.nextRng()
	}

	for l := 0; l < cfg.LaneCount; l++ {
		offset := ((uint32(l) ^ r) % uint32(cfg.LaneCount)) * uint32(numWordsPerLane)
		for i := 0; i < numWordsPerLane; i++ {
			word := item[offset+uint32(i)]
			mix[l][dsts[i]] = randomMerge(mix[l][dsts[i]], word, sels[i])
		}
	}

	return mix
}

func Hash(cfg *Config, height, seed, datasetSize uint64, lookup dag.LookupFunc, l1 []uint32) []byte {
	mix := initMix(seed, cfg.LaneCount, cfg.RegisterCount)

	number := height / cfg.PeriodLength
	for i := 0; i < cfg.RoundCount; i++ {
		mix = round(cfg, number, uint32(i), mix, datasetSize, lookup, l1)
	}

	laneHash := make([]uint32, cfg.LaneCount)
	for l := range laneHash {
		laneHash[l] = fnvOffsetBasis

		for i := 0; i < cfg.RegisterCount; i++ {
			laneHash[l] = crypto.Fnv1a(laneHash[l], mix[l][i])
		}
	}

	numWords := 8
	mixHash := make([]uint32, numWords)
	for i := 0; i < numWords; i++ {
		mixHash[i] = fnvOffsetBasis
	}

	for l := 0; l < cfg.LaneCount; l++ {
		mixHash[l%numWords] = crypto.Fnv1a(mixHash[l%numWords], laneHash[l])
	}

	return convutil.Uint32ArrayToBytes(mixHash, binary.LittleEndian)
}
