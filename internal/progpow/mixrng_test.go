// ethash: C/C++ implementation of Ethash, the Ethereum Proof of Work algorithm.
// Copyright 2018-2019 Pawel Bylica.
// Licensed under the Apache License, Version 2.0.

package progpow

import (
	"testing"
)

func TestInitMixRngState(t *testing.T) {
	tests := []struct {
		number uint64
		size   uint32
		src    []byte
		dst    []byte
		z      uint32
		w      uint32
		jsr    uint32
		jcong  uint32
	}{
		{
			number: 30000 / 50,
			size:   32,
			src: []byte{
				0x1A, 0x1E, 0x01, 0x13, 0x0B, 0x15, 0x0F, 0x12,
				0x03, 0x11, 0x1F, 0x10, 0x1C, 0x04, 0x16, 0x17,
				0x02, 0x0D, 0x1D, 0x18, 0x0A, 0x0C, 0x05, 0x14,
				0x07, 0x08, 0x0E, 0x1B, 0x06, 0x19, 0x09, 0x00,
			},
			dst: []byte{
				0x00, 0x04, 0x1B, 0x1A, 0x0D, 0x0F, 0x11, 0x07,
				0x0E, 0x08, 0x09, 0x0C, 0x03, 0x0A, 0x01, 0x0B,
				0x06, 0x10, 0x1C, 0x1F, 0x02, 0x13, 0x1E, 0x16,
				0x1D, 0x05, 0x18, 0x12, 0x19, 0x17, 0x15, 0x14,
			},
			z:     0x6535921C,
			w:     0x29345B16,
			jsr:   0xC0DD7F78,
			jcong: 0x1165D7EB,
		},
	}

	for i, tt := range tests {
		state := initMixRngState(tt.number, tt.size)

		if len(state.srcSequence) != len(tt.src) {
			t.Errorf("failed on %d: srcSequence length mismatch: have %d, want %d",
				i, len(state.srcSequence), len(tt.src))
		} else {
			for j := 0; j < len(state.srcSequence); j++ {
				if state.srcSequence[j] != uint32(tt.src[j]) {
					t.Errorf("failed on %d: srcSequence mismatch (%d): have %x, want %x",
						i, j, state.srcSequence[j], tt.src[j])
				}
			}
		}

		if len(state.dstSequence) != len(tt.dst) {
			t.Errorf("failed on %d: dstSequence length mismatch: have %d, want %d",
				i, len(state.dstSequence), len(tt.dst))
		} else {
			for j := 0; j < len(state.srcSequence); j++ {
				if state.dstSequence[j] != uint32(tt.dst[j]) {
					t.Errorf("failed on %d: dstSequence mismatch (%d): have %x, want %x",
						i, j, state.dstSequence[j], tt.dst[j])
				}
			}
		}

		if state.rng.z != tt.z {
			t.Errorf("failed on %d: z mismatch: have %d, want %d", i, state.rng.z, tt.z)
		}

		if state.rng.w != tt.w {
			t.Errorf("failed on %d: w mismatch: have %d, want %d", i, state.rng.w, tt.w)
		}

		if state.rng.jsr != tt.jsr {
			t.Errorf("failed on %d: jsr mismatch: have %d, want %d", i, state.rng.jsr, tt.jsr)
		}

		if state.rng.jcong != tt.jcong {
			t.Errorf("failed on %d: jcong mismatch: have %d, want %d", i, state.rng.jcong, tt.jcong)
		}
	}
}
