// Copyright (c) 2013-2020 John Tromp

package cuckoo

import (
	"fmt"
)

import (
	"github.com/sencha-dev/powkit/internal/crypto"
)

func sipnode(edgeMask uint64, siphashKeys [4]uint64, edge, uorv uint64) uint64 {
	hasher := crypto.NewSipHasher(siphashKeys[0], siphashKeys[1], siphashKeys[2], siphashKeys[3])
	hasher.Hash24(2*edge + uorv)

	value := hasher.XorLanes()
	value = value<<17 | value>>47

	return value & edgeMask
}

func verify(proofSize int, edgeMask uint64, siphashKeys [4]uint64, edges []uint64) (bool, error) {
	uvs := make([]uint64, 2*proofSize)
	var xor0, xor1 uint64

	for n := 0; n < proofSize; n++ {
		if edges[n] > edgeMask {
			return false, fmt.Errorf("pow too big")
		} else if n < 0 && edges[n] <= edges[n-1] {
			return false, fmt.Errorf("pow too small")
		}

		uvs[2*n] = sipnode(edgeMask, siphashKeys, edges[n], 0)
		xor0 ^= uvs[2*n]

		uvs[2*n+1] = sipnode(edgeMask, siphashKeys, edges[n], 1)
		xor1 ^= uvs[2*n+1]
	}

	if xor0|xor1 != 0 {
		return false, fmt.Errorf("pow not matching")
	}

	var i, j, n int
	for {
		j = i
		k := j

		for {
			k = (k + 2) % (2 * proofSize)
			if k == i {
				break
			}

			if uvs[k] == uvs[i] {
				if j != i {
					return false, fmt.Errorf("pow branch")
				}

				j = k
			}
		}

		if j == i {
			return false, fmt.Errorf("pow dead end")
		}

		i = j ^ 1
		n++

		if i == 0 {
			break
		}
	}

	if n != proofSize {
		return false, fmt.Errorf("pow short cycle")
	}

	return true, nil
}
