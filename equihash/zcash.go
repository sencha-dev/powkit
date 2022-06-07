// Copyright (c) 2021 Electric Coin Company

package equihash

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/sencha-dev/powkit/internal/common/convutil"
	"github.com/sencha-dev/powkit/internal/crypto"
)

const (
	uint32Size = 4
	wordSize   = 32
	wordMask   = (1 << wordSize) - 1
)

func expandArray(input []byte, bitLen, bytePad uint32) ([]byte, error) {
	if bitLen < 8 {
		return nil, fmt.Errorf("bitLen must be no less than 8")
	} else if wordSize < 7+bitLen {
		return nil, fmt.Errorf("bitLen must be no greater than %d", wordSize-7)
	}

	inputLen := uint32(len(input))
	outputWidth := (bitLen+7)/8 + bytePad
	outputLen := 8 * outputWidth * inputLen / bitLen

	if outputLen == inputLen {
		return input, nil
	}

	output := make([]byte, outputLen)
	var bitLenMask uint32 = (1 << bitLen) - 1

	var accBits, accValue, j uint32
	for i := range input {
		accValue = (accValue << 8) | uint32(input[i])
		accBits += 8

		if accBits >= bitLen {
			accBits -= bitLen
			for x := bytePad; x < outputWidth; x++ {
				p1 := accValue >> (accBits + (8 * (outputWidth - x - 1)))
				p2 := (bitLenMask >> (8 * (outputWidth - x - 1))) & 0xFF
				output[j+x] = uint8(p1 & p2)
			}

			j += outputWidth
		}
	}

	return output, nil
}

func collisionBitLength(n, k uint32) uint32 {
	return n / (k + 1)
}

func collisionByteLength(n, k uint32) uint32 {
	return (collisionBitLength(n, k) + 7) / 8
}

func indicesPerHashOutput(n uint32) uint32 {
	return 512 / n
}

func hashOutput(n uint32) uint8 {
	output := indicesPerHashOutput(n) * n / 8

	return uint8(output)
}

func hashLength(n, k uint32) uint32 {
	return (k + 1) * collisionByteLength(n, k)
}

func blakePersonal(personal []byte, n, k uint32) []byte {
	nBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(nBytes, n)

	kBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(kBytes, k)

	personalBytes := bytes.Join([][]byte{
		personal,
		nBytes,
		kBytes,
	}, nil)

	return personalBytes
}

func indicesFromMinimal(n, k uint32, minimal []byte) ([]uint32, error) {
	cBitLen := collisionBitLength(n, k)
	minimalLen := uint32(len(minimal))

	if minimalLen != ((1<<k)*(cBitLen+1))/8 {
		return nil, fmt.Errorf("invalid minimal for parameters")
	}

	if (((cBitLen + 1) + 7) / 8) > 4 {
		return nil, fmt.Errorf("invalid n, k parameters")
	}

	//lenIndices := 8 * uint32Size * minimalLen / (cBitLen + 1)
	bytePad := uint32Size - ((cBitLen+1)+7)/8

	indices, err := expandArray(minimal, cBitLen+1, bytePad)
	if err != nil {
		return nil, err
	}

	return convutil.BytesToUint32Array(indices, binary.BigEndian), nil
}

func hasCollision(a, b *node, len uint32) bool {
	for i := uint32(0); i < len; i++ {
		if a.hash[i] != b.hash[i] {
			return false
		}
	}

	return true
}

func distinctIndices(a, b *node) bool {
	for _, i := range a.indices {
		for _, j := range b.indices {
			if i == j {
				return false
			}
		}
	}

	return true
}

func validateSubtrees(n, k uint32, a, b *node) error {
	if !hasCollision(a, b, collisionByteLength(n, k)) {
		return fmt.Errorf("collision")
	} else if b.indicesBefore(a) {
		return fmt.Errorf("out of order")
	} else if !distinctIndices(a, b) {
		return fmt.Errorf("duplicate indices")
	}

	return nil
}

type node struct {
	hash    []byte
	indices []uint32
}

func (n *node) indicesBefore(b *node) bool {
	return n.indices[0] < b.indices[0]
}

func newNode(n, k uint32, personal, state []byte, i uint32) (*node, error) {
	iBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(iBytes, i/indicesPerHashOutput(n))
	state = bytes.Join([][]byte{state, iBytes}, nil)

	hashPersonal := blakePersonal(personal, n, k)
	hashSize := int(hashOutput(n))
	hash := crypto.Blake2b(state, hashPersonal, hashSize)

	start := (i % indicesPerHashOutput(n)) * n / 8
	end := start + n/8

	minimalHash, err := expandArray(hash[start:end], collisionBitLength(n, k), 0)
	if err != nil {
		return nil, err
	}

	return &node{
		hash:    minimalHash,
		indices: []uint32{i},
	}, nil
}

func newNodeFromChildrenRef(a, b *node, trim uint32) *node {
	len := uint32(len(a.hash))
	hash := make([]byte, len-trim)
	for i := trim; i < len; i++ {
		hash[i-trim] = a.hash[i] ^ b.hash[i]
	}

	indices := make([]uint32, 0)
	if a.indicesBefore(b) {
		indices = append(indices, a.indices...)
		indices = append(indices, b.indices...)
	} else {
		indices = append(indices, b.indices...)
		indices = append(indices, a.indices...)
	}

	n := &node{
		hash:    hash,
		indices: indices,
	}

	return n
}

func isValidSolutionIterative(n, k uint32, personal, state []byte, indices []uint32) (bool, error) {
	var err error
	rows := make([]*node, len(indices))
	for i := range indices {
		rows[i], err = newNode(n, k, personal, state, indices[i])
		if err != nil {
			return false, err
		}
	}

	hashLen := hashLength(n, k)
	for len(rows) > 1 {
		curRows := make([]*node, 0)
		for i := 0; i < len(rows); i += 2 {
			a := rows[i]
			b := rows[i+1]
			err := validateSubtrees(n, k, a, b)
			if err != nil {
				return false, err
			}

			row := newNodeFromChildrenRef(a, b, collisionByteLength(n, k))
			curRows = append(curRows, row)
		}

		rows = curRows
		hashLen -= collisionByteLength(n, k)
	}

	for i := uint32(0); i < hashLen; i++ {
		if rows[0].hash[i] != 0 {
			return false, nil
		}
	}

	return true, nil
}

func ZCashVerify(n, k uint32, personal, header, soln []byte) (bool, error) {
	indices, err := indicesFromMinimal(n, k, soln)
	if err != nil {
		return false, err
	}

	return isValidSolutionIterative(n, k, personal, header, indices)
}
