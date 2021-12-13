// ethash: C/C++ implementation of Ethash, the Ethereum Proof of Work algorithm.
// Copyright 2018-2019 Pawel Bylica.
// Licensed under the Apache License, Version 2.0.

package progpow

import (
	"testing"
)

func TestKiss99(t *testing.T) {
	tests := []struct {
		z     uint32
		w     uint32
		jsr   uint32
		jcong uint32
		index int
		value uint32
	}{
		{
			z:     362436069,
			w:     521288629,
			jsr:   123456789,
			jcong: 380116160,
			index: 1,
			value: 769445856,
		},
		{
			z:     362436069,
			w:     521288629,
			jsr:   123456789,
			jcong: 380116160,
			index: 2,
			value: 742012328,
		},
		{
			z:     362436069,
			w:     521288629,
			jsr:   123456789,
			jcong: 380116160,
			index: 3,
			value: 2121196314,
		},
		{
			z:     362436069,
			w:     521288629,
			jsr:   123456789,
			jcong: 380116160,
			index: 4,
			value: 2805620942,
		},
		{
			z:     362436069,
			w:     521288629,
			jsr:   123456789,
			jcong: 380116160,
			index: 100000,
			value: 941074834,
		},
	}

	for i, tt := range tests {
		rng := newKiss99(tt.z, tt.w, tt.jsr, tt.jcong)

		var value uint32
		for i := 0; i < tt.index; i++ {
			value = rng.next()
		}

		if value != tt.value {
			t.Errorf("failed on %d: have %d, want %d", i, value, tt.value)
		}
	}
}
