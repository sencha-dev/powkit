// ethash: C/C++ implementation of Ethash, the Ethereum Proof of Work algorithm.
// Copyright 2018-2019 Pawel Bylica.
// Licensed under the Apache License, Version 2.0.

package progpow

type kiss99 struct {
	z     uint32
	w     uint32
	jsr   uint32
	jcong uint32
}

func newKiss99(z, w, jsr, jcong uint32) *kiss99 {
	return &kiss99{z, w, jsr, jcong}
}

func (k *kiss99) next() uint32 {
	k.z = 36969*(k.z&65535) + (k.z >> 16)
	k.w = 18000*(k.w&65535) + (k.w >> 16)

	k.jcong = 69069*k.jcong + 1234567

	k.jsr = k.jsr ^ (k.jsr << 17)
	k.jsr = k.jsr ^ (k.jsr >> 13)
	k.jsr = k.jsr ^ (k.jsr << 5)

	return (((k.z << 16) + k.w) ^ k.jcong) + k.jsr
}
