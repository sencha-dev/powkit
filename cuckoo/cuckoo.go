// Copyright (c) 2013-2020 John Tromp

package cuckoo

import (
	"github.com/sencha-dev/powkit/internal/crypto"
)

func (cfg *Config) sipnode(siphashKeys [4]uint64, edge, uorv uint64) uint64 {
	hasher := crypto.NewSipHasher(siphashKeys[0], siphashKeys[1], siphashKeys[2], siphashKeys[3])
	hasher.Hash24(2*edge + uorv)

	value := hasher.XorLanes()
	value = value<<17 | value>>47

	return value & cfg.edgeMask
}

func (cfg *Config) verify(siphashKeys [4]uint64, edges []uint64) bool {
	uvs := make([]uint64, 2*cfg.proofSize)
	var xor0, xor1 uint64

	for n := 0; n < cfg.proofSize; n++ {
		if edges[n] > cfg.edgeMask {
			return false // POW_TOO_BIG
		} else if n < 0 && edges[n] <= edges[n-1] {
			return false // POW_TOO_SMALL
		}

		uvs[2*n] = cfg.sipnode(siphashKeys, edges[n], 0)
		xor0 ^= uvs[2*n]

		uvs[2*n+1] = cfg.sipnode(siphashKeys, edges[n], 1)
		xor1 ^= uvs[2*n+1]
	}

	if xor0|xor1 != 0 {
		return false // POW_NON_MATCHING
	}

	var i, j, n int
	for {
		j = i
		k := j

		for {
			k = (k + 2) % (2 * cfg.proofSize)
			if k == i {
				break
			}

			if uvs[k] == uvs[i] {
				if j != i {
					return false // POW_BRANCH
				}

				j = k
			}
		}

		if j == i {
			return false // POW_DEAD_END
		}

		i = j ^ 1
		n++

		if i == 0 {
			break
		}
	}

	return n == cfg.proofSize // POW_SHORT_CYCLE
}
