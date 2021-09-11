/* ethash: C/C++ implementation of Ethash, the Ethereum Proof of Work algorithm.
 * Copyright 2018-2019 Pawel Bylica.
 * Licensed under the Apache License, Version 2.0.
 */

package crypto

import (
	"testing"
)

func TestSimple(t *testing.T) {
	var h1, d1, h2, d2, h3, d3 uint32
	var expected1, expected2, expected3 uint32

	h1, d1, expected1 = 0x811C9DC5, 0xDDD0A47B, 0xD37EE61A
	h2, d2, expected2 = 0xD37EE61A, 0xEE304846, 0xDEDC7AD4
	h3, d3, expected3 = 0xDEDC7AD4, 0x00000000, 0xA9155BBC

	actual1 := Fnv1a(h1, d1)
	if actual1 != expected1 {
		t.Errorf("failed fnv1a test 1")
	}

	actual2 := Fnv1a(h2, d2)
	if actual2 != expected2 {
		t.Errorf("failed fnv1a test 2")
	}

	actual3 := Fnv1a(h3, d3)
	if actual3 != expected3 {
		t.Errorf("failed fnv1a test 3")
	}
}
