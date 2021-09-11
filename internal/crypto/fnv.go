/* ethash: C/C++ implementation of Ethash, the Ethereum Proof of Work algorithm.
 * Copyright 2018-2019 Pawel Bylica.
 * Licensed under the Apache License, Version 2.0.
 */

package crypto

const (
	fnvPrime uint32 = 0x01000193
)

// See https://en.wikipedia.org/wiki/Fowler%E2%80%93Noll%E2%80%93Vo_hash_function#FNV-1_hash.
func Fnv1(u, v uint32) uint32 {
	return (u * fnvPrime) ^ v
}

// See https://en.wikipedia.org/wiki/Fowler%E2%80%93Noll%E2%80%93Vo_hash_function#FNV-1a_hash.
func Fnv1a(u, v uint32) uint32 {
	return (u ^ v) * fnvPrime
}

// fnvHash mixes in data into mix using the ethash fnv method.
func FnvHash(mix []uint32, data []uint32) {
	for i := 0; i < len(mix); i++ {
		mix[i] = Fnv1(mix[i], data[i])
	}
}
