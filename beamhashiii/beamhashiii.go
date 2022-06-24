// Copyright (c) 2020 The Beam Team

package beamhashiii

import (
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/sencha-dev/powkit/internal/crypto"
)

const (
	workBitsSize      = 448
	collisionBitsSize = 24
	numRounds         = 5
)

var (
	indexMask     = new(big.Int).SetUint64(0x1ffffff)
	collisionMask = new(big.Int).SetUint64(0xffffff)
	mixMask       = new(big.Int).SetUint64(0xffffffffffffffff)
)

func rotl(a, b uint64) uint64 {
	return (a << b) | (a >> (64 - b))
}

func reverseBytes(b []byte) []byte {
	_b := make([]byte, len(b))
	copy(_b, b)

	for i, j := 0, len(_b)-1; i < j; i, j = i+1, j-1 {
		_b[i], _b[j] = _b[j], _b[i]
	}
	return _b
}

func indicesFromMinimal(soln []byte) []uint32 {
	streamBig := new(big.Int).SetBytes(reverseBytes(soln[:100]))

	indices := make([]uint32, 32)
	for i := 0; i < 32; i++ {
		indices[i] = uint32(new(big.Int).And(streamBig, indexMask).Uint64())
		streamBig.Rsh(streamBig, collisionBitsSize+1)
	}

	return indices
}

type node struct {
	bitset  *big.Int
	indices []uint32
}

func newNode(prePow []uint64, idx uint32) *node {
	bitset := new(big.Int)
	for i := 6; i >= 0; i-- {
		bitset.Lsh(bitset, 64)

		hasher := crypto.NewSipHasher(prePow[0], prePow[1], prePow[2], prePow[3])
		hasher.Hash24(uint64(idx)<<3 + uint64(i))
		value := hasher.XorLanes()

		bitset.Or(bitset, new(big.Int).SetUint64(value))
	}

	return &node{
		bitset:  bitset,
		indices: []uint32{idx},
	}
}

func newNodeFromChildrenRef(a, b *node, remLen uint32) *node {
	workBits := new(big.Int).Xor(a.bitset, b.bitset)
	workBits.Rsh(workBits, collisionBitsSize)

	rawMask := make([]byte, workBitsSize/8)
	for i := range rawMask {
		rawMask[i] = 0xFF
	}

	maskBig := new(big.Int).SetBytes(rawMask)
	maskBig.Rsh(maskBig, uint(workBitsSize-remLen))
	workBits.And(workBits, maskBig)

	indices := make([]uint32, 0)
	if a.indices[0] < b.indices[0] {
		indices = append(indices, a.indices...)
		indices = append(indices, b.indices...)
	} else {
		indices = append(indices, b.indices...)
		indices = append(indices, a.indices...)
	}

	n := &node{
		bitset:  workBits,
		indices: indices,
	}

	return n
}

func validateSubtrees(n, k uint32, a, b *node) error {
	maskedA := new(big.Int).And(a.bitset, collisionMask)
	maskedB := new(big.Int).And(b.bitset, collisionMask)

	// check hasCollision
	if maskedA.Cmp(maskedB) != 0 {
		return fmt.Errorf("collision")
	}

	// check indicesBefore
	if b.indices[0] < a.indices[0] {
		return fmt.Errorf("out of order")
	}

	// check distinctIndices
	for _, i := range a.indices {
		for _, j := range b.indices {
			if i == j {
				return fmt.Errorf("duplicate indices")
			}
		}
	}

	return nil
}

func applyMix(n *node, remLen uint32) {
	tempBits := new(big.Int).Set(n.bitset)

	padNum := ((512 - remLen) + collisionBitsSize) / (collisionBitsSize + 1)
	if uint32(len(n.indices)) < padNum {
		padNum = uint32(len(n.indices))
	}

	for i := uint32(0); i < padNum; i++ {
		tmp := new(big.Int).SetUint64(uint64(n.indices[i]))
		tmp.Lsh(tmp, uint((remLen + i*(collisionBitsSize+1))))
		tempBits.Or(tempBits, tmp)
	}

	var result uint64
	for i := uint64(0); i < 8; i++ {
		tmp := new(big.Int).And(tempBits, mixMask).Uint64()
		tempBits.Rsh(tempBits, 64)
		result += rotl(tmp, (29*(i+1))&0x3F)
	}

	resultBig := new(big.Int).SetUint64(rotl(result, 24))

	n.bitset.Rsh(n.bitset, 64)
	n.bitset.Lsh(n.bitset, 64)
	n.bitset.Or(n.bitset, resultBig)
}

func initializePrePow(personal, header, soln []byte) []uint64 {
	personalState := make([]byte, len(personal)+8)
	copy(personalState, personal)
	binary.LittleEndian.PutUint32(personalState[len(personal):], workBitsSize)
	binary.LittleEndian.PutUint32(personalState[len(personal)+4:], numRounds)

	state := make([]byte, len(header)+4)
	copy(state, header)
	binary.LittleEndian.PutUint32(state[len(header):], binary.LittleEndian.Uint32(soln[100:]))

	hash := crypto.Blake2b(state, personalState, 32)
	prePow := []uint64{
		binary.LittleEndian.Uint64(hash[0:8]),
		binary.LittleEndian.Uint64(hash[8:16]),
		binary.LittleEndian.Uint64(hash[16:24]),
		binary.LittleEndian.Uint64(hash[24:32]),
	}

	return prePow
}

func verify(n, k uint32, personal, header, soln []byte) (bool, error) {
	prePow := initializePrePow(personal, header, soln)
	indices := indicesFromMinimal(soln)
	rows := make([]*node, len(indices))
	for i := range indices {
		rows[i] = newNode(prePow, indices[i])
	}

	var round uint32 = 1
	for len(rows) > 1 {
		curRows := make([]*node, 0)
		for i := 0; i < len(rows); i += 2 {
			remLen := workBitsSize - (round-1)*collisionBitsSize
			if round == 5 {
				remLen -= 64
			}

			applyMix(rows[i], remLen)
			applyMix(rows[i+1], remLen)

			err := validateSubtrees(n, k, rows[i], rows[i+1])
			if err != nil {
				return false, err
			}

			remLen = workBitsSize - round*collisionBitsSize
			if round == 4 {
				remLen -= 64
			} else if round == 5 {
				remLen = collisionBitsSize
			}

			row := newNodeFromChildrenRef(rows[i], rows[i+1], remLen)
			curRows = append(curRows, row)
		}

		rows = curRows
		round++
	}

	if new(big.Int).Cmp(rows[0].bitset) != 0 {
		return false, nil
	}

	return true, nil
}
