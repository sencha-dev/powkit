// ethash: C/C++ implementation of Ethash, the Ethereum Proof of Work algorithm.
// Copyright 2018-2019 Pawel Bylica.
// Licensed under the Apache License, Version 2.0.

package progpow

import (
	"testing"
)

func TestKiss99(t *testing.T) {
	var z, w, jsr, jcong uint32
	z, w = 362436069, 521288629
	jsr, jcong = 123456789, 380116160

	results := map[int]uint32{
		1:      769445856,
		2:      742012328,
		3:      2121196314,
		4:      2805620942,
		100000: 941074834,
	}

	kiss := newKiss99(z, w, jsr, jcong)

	for i := 1; i < 100001; i++ {
		actual := kiss.Next()

		if expected, ok := results[i]; ok {
			if actual != expected {
				t.Errorf("failed Kiss99 test on iteration %d", i)
			}
		}
	}
}
