// Copyright (c) 2013-2020 John Tromp

package cuckoo

func (c *Client) cuckaroo(siphashKeys [4]uint64, edges []uint64) (bool, error) {
	uvs := make([]uint64, 2*c.proofSize)
	var xor0, xor1 uint64

	for n := 0; n < c.proofSize; n++ {
		if edges[n] > c.edgeMask {
			return false, ErrPowTooBig
		} else if n < 0 && edges[n] <= edges[n-1] {
			return false, ErrPowTooSmall
		}

		edge := c.sipblock(siphashKeys, edges[n])
		uvs[2*n] = edge & c.edgeMask
		xor0 ^= uvs[2*n]
		uvs[2*n+1] = (edge >> 32) & c.edgeMask
		xor1 ^= uvs[2*n+1]
	}

	if xor0|xor1 != 0 {
		return false, ErrPowNotMatching
	}

	var i, j, n int
	for {
		j = i
		k := j

		for {
			k = (k + 2) % (2 * c.proofSize)
			if k == i {
				break
			}

			if uvs[k] == uvs[i] {
				if j != i {
					return false, ErrPowBranch
				}

				j = k
			}
		}

		if j == i {
			return false, ErrPowDeadEnd
		}

		i = j ^ 1
		n++

		if i == 0 {
			break
		}
	}

	if n != c.proofSize {
		return false, ErrPowShortCycle
	}

	return true, nil
}
