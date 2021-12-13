// ethash: C/C++ implementation of Ethash, the Ethereum Proof of Work algorithm.
// Copyright 2018-2019 Pawel Bylica.
// Licensed under the Apache License, Version 2.0.

package progpow

import (
	"github.com/sencha-dev/powkit/internal/crypto"
)

type mixRngState struct {
	size        uint32
	srcCounter  uint32
	dstCounter  uint32
	srcSequence []uint32
	dstSequence []uint32
	rng         *kiss99
}

func (s *mixRngState) nextSrc() uint32 {
	val := s.srcSequence[s.srcCounter%s.size]
	s.srcCounter++

	return val
}

func (s *mixRngState) nextDst() uint32 {
	val := s.dstSequence[s.dstCounter%s.size]
	s.dstCounter++

	return val
}

func (s *mixRngState) nextRng() uint32 {
	return s.rng.next()
}

func initMixRngState(seed uint64, size uint32) *mixRngState {
	var z, w, jsr, jcong uint32

	z = crypto.Fnv1a(fnvOffsetBasis, uint32(seed))
	w = crypto.Fnv1a(z, uint32(seed>>32))
	jsr = crypto.Fnv1a(w, uint32(seed))
	jcong = crypto.Fnv1a(jsr, uint32(seed>>32))

	rng := newKiss99(z, w, jsr, jcong)

	srcSeq := make([]uint32, size)
	dstSeq := make([]uint32, size)
	for i := uint32(0); i < size; i++ {
		dstSeq[i] = i
		srcSeq[i] = i
	}

	for i := size; i > 1; i-- {
		dstInd := rng.next() % i
		dstSeq[i-1], dstSeq[dstInd] = dstSeq[dstInd], dstSeq[i-1]

		srcInd := rng.next() % i
		srcSeq[i-1], srcSeq[srcInd] = srcSeq[srcInd], srcSeq[i-1]
	}

	state := &mixRngState{
		size:        size,
		srcCounter:  0,
		dstCounter:  0,
		srcSequence: srcSeq,
		dstSequence: dstSeq,
		rng:         rng,
	}

	return state
}
