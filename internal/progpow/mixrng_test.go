// ethash: C/C++ implementation of Ethash, the Ethereum Proof of Work algorithm.
// Copyright 2018-2019 Pawel Bylica.
// Licensed under the Apache License, Version 2.0.

package progpow

import (
	"testing"
)

func TestInitMixRngState(t *testing.T) {
	const period uint64 = 50
	const height uint64 = 30000

	expectedSrc := []byte{
		0x1A, 0x1E, 0x01, 0x13, 0x0B, 0x15, 0x0F, 0x12,
		0x03, 0x11, 0x1F, 0x10, 0x1C, 0x04, 0x16, 0x17,
		0x02, 0x0D, 0x1D, 0x18, 0x0A, 0x0C, 0x05, 0x14,
		0x07, 0x08, 0x0E, 0x1B, 0x06, 0x19, 0x09, 0x00,
	}

	expectedDst := []byte{
		0x00, 0x04, 0x1B, 0x1A, 0x0D, 0x0F, 0x11, 0x07,
		0x0E, 0x08, 0x09, 0x0C, 0x03, 0x0A, 0x01, 0x0B,
		0x06, 0x10, 0x1C, 0x1F, 0x02, 0x13, 0x1E, 0x16,
		0x1D, 0x05, 0x18, 0x12, 0x19, 0x17, 0x15, 0x14,
	}

	const expectedZ uint32 = 0x6535921C
	const expectedW uint32 = 0x29345B16
	const expectedJsr uint32 = 0xC0DD7F78
	const expectedJcong uint32 = 0x1165D7EB

	number := height / period
	state := initMixRngState(number)

	for i := 0; i < len(expectedSrc); i++ {
		if state.SrcSeq[i] != uint32(expectedSrc[i]) {
			t.Errorf("failed initMixRngState test for SrcSeq value at index %d", i)
		}

		if state.DstSeq[i] != uint32(expectedDst[i]) {
			t.Errorf("failed initMixRngState test for DstSeq value at index %d", i)
		}
	}

	if state.Rng.z != expectedZ {
		t.Errorf("failed initMixRngState test for Rng value z")
	}

	if state.Rng.w != expectedW {
		t.Errorf("failed initMixRngState test for Rng value w")
	}

	if state.Rng.jsr != expectedJsr {
		t.Errorf("failed initMixRngState test for Rng value jsr")
	}

	if state.Rng.jcong != expectedJcong {
		t.Errorf("failed initMixRngState test for Rng value jcong")
	}
}
