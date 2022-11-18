// Copyright (c) 2018-2019 The kaspanet developers

package heavyhash

import (
	"math"

	"github.com/sencha-dev/powkit/internal/crypto"
)

const (
	epsilon = 1e-9
)

type matrix [size][size]uint16

func newMatrix(s0, s1, s2, s3 uint64) *matrix {
	hasher := crypto.NewXoshiro256PlusPlusHasher(s0, s1, s2, s3)

	var mat matrix
	for calculateRank(&mat) != size {
		for i := 0; i < size; i++ {
			for j := 0; j < size; j += iterations {
				value := hasher.Next()
				for k := 0; k < iterations; k++ {
					mat[i][j+k] = uint16(value >> (4 * k) & 0x0f)
				}
			}
		}
	}

	return &mat
}

func calculateRank(mat *matrix) int {
	var copied [size][size]float64
	for i := range mat {
		for j := range mat[i] {
			copied[i][j] = float64(mat[i][j])
		}
	}

	var rank int
	var rowsSelected [size]bool
	for i := 0; i < size; i++ {
		var j int
		for j = 0; j < size; j++ {
			if !rowsSelected[j] && math.Abs(copied[j][i]) > epsilon {
				break
			}
		}

		if j != size {
			rank++
			rowsSelected[j] = true
			for k := i + 1; k < size; k++ {
				copied[j][k] /= copied[j][i]
			}

			for k := 0; k < size; k++ {
				if k == j || math.Abs(copied[k][i]) <= epsilon {
					continue
				}

				for l := i + 1; l < size; l++ {
					copied[k][l] -= copied[j][l] * copied[k][i]
				}
			}
		}
	}

	return rank
}
