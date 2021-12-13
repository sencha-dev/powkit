// ethash: C/C++ implementation of Ethash, the Ethereum Proof of Work algorithm.
// Copyright 2018-2019 Pawel Bylica.
// Licensed under the Apache License, Version 2.0.

package progpow

import (
	"github.com/sencha-dev/powkit/internal/crypto"
)

type mixRngState struct {
	SrcCounter uint32
	DstCounter uint32
	SrcSeq     [progpowRegs]uint32
	DstSeq     [progpowRegs]uint32
	Rng        *kiss99
}

func (s *mixRngState) nextDst() uint32 {
	val := s.DstSeq[s.DstCounter%progpowRegs]
	s.DstCounter++

	return val
}

func (s *mixRngState) nextSrc() uint32 {
	val := s.SrcSeq[s.SrcCounter%progpowRegs]
	s.SrcCounter++

	return val
}

func initMixRngState(seed uint64) *mixRngState {
	var z, w, jsr, jcong uint32

	z = crypto.Fnv1a(fnvOffsetBasis, uint32(seed))
	w = crypto.Fnv1a(z, uint32(seed>>32))
	jsr = crypto.Fnv1a(w, uint32(seed))
	jcong = crypto.Fnv1a(jsr, uint32(seed>>32))

	rng := newKiss99(z, w, jsr, jcong)

	var srcSeq [progpowRegs]uint32
	var dstSeq [progpowRegs]uint32

	var i uint32
	for i = 0; i < progpowRegs; i++ {
		dstSeq[i] = i
		srcSeq[i] = i
	}

	for i = progpowRegs; i > 1; i-- {
		dstInd := rng.Next() % i
		dstSeq[i-1], dstSeq[dstInd] = dstSeq[dstInd], dstSeq[i-1]

		srcInd := rng.Next() % i
		srcSeq[i-1], srcSeq[srcInd] = srcSeq[srcInd], srcSeq[i-1]
	}

	return &mixRngState{0, 0, srcSeq, dstSeq, rng}
}
