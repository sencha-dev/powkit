// Copyright (c) 2021 Electric Coin Company

package equihash

import (
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

func collisionBitLength(n, k uint32) uint32 {
	return n / (k + 1)
}

func collisionByteLength(n, k uint32) uint32 {
	return (collisionBitLength(n, k) + 7) / 8
}

func indicesPerHashOutput(n uint32) uint32 {
	return 512 / n
}

func hashOutput(n uint32) uint32 {
	return indicesPerHashOutput(n) * ((n + 7) / 8)
}

func hashLength(n, k uint32) uint32 {
	return (k + 1) * collisionByteLength(n, k)
}

func expandArray(input []byte, bitLen, bytePad uint32) ([]byte, error) {
	if bitLen < 8 {
		return nil, fmt.Errorf("bitLen must be no less than 8")
	} else if wordSize < bitLen {
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

func indicesFromMinimal(n, k uint32, minimal []byte) ([]uint32, error) {
	cBitLen := collisionBitLength(n, k)
	minimalLen := uint32(len(minimal))

	if minimalLen != ((1<<k)*(cBitLen+1))/8 {
		return nil, fmt.Errorf("invalid minimal for parameters")
	} else if (((cBitLen + 1) + 7) / 8) > 4 {
		return nil, fmt.Errorf("invalid n, k parameters")
	}

	bytePad := uint32Size - ((cBitLen+1)+7)/8
	indices, err := expandArray(minimal, cBitLen+1, bytePad)
	if err != nil {
		return nil, err
	}

	return convutil.BytesToUint32Array(indices, binary.BigEndian), nil
}

func hashBlakeWithOffset(initialState, personalState []byte, offset, hashLength uint32) []byte {
	newState := make([]byte, len(initialState)+4)
	copy(newState, initialState)
	binary.LittleEndian.PutUint32(newState[len(initialState):], offset)

	return crypto.Blake2b(newState, personalState, int(hashLength))
}

func generateHash(initialState, personalState []byte, g, hashLength uint32, twist bool) []byte {
	if !twist {
		return hashBlakeWithOffset(initialState, personalState, g, hashLength)
	}

	myHash := make([]uint32, 16)
	startIndex := g & 0xFFFFFFF0
	for g2 := startIndex; g2 <= g; g2++ {
		tmpHash := hashBlakeWithOffset(initialState, personalState, g2, hashLength)
		for i := 0; i < 16; i++ {
			myHash[i] += binary.LittleEndian.Uint32(tmpHash[i*4 : (i+1)*4])
		}
	}

	hash := convutil.Uint32ArrayToBytes(myHash, binary.LittleEndian)
	for j := 15; j < int(hashLength); j += 16 {
		hash[j] &= 0xF8
	}

	return hash
}

type node struct {
	hash    []byte
	indices []uint32
}

func newNode(n, k uint32, personalState []byte, twist bool, state []byte, i uint32) (*node, error) {
	g := i / indicesPerHashOutput(n)
	hash := generateHash(state, personalState, g, hashOutput(n), twist)

	start := (i % indicesPerHashOutput(n)) * ((n + 7) / 8)
	end := start + ((n + 7) / 8)

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
	if a.indices[0] < b.indices[0] {
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

func validateSubtrees(n, k uint32, a, b *node) error {
	// check hasCollision
	for i := uint32(0); i < collisionByteLength(n, k); i++ {
		if a.hash[i] != b.hash[i] {
			return fmt.Errorf("collision")
		}
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

func verify(n, k uint32, personal, header, soln []byte, twist bool) (bool, error) {
	personalState := make([]byte, len(personal)+8)
	copy(personalState, personal)
	binary.LittleEndian.PutUint32(personalState[len(personal):], n)
	binary.LittleEndian.PutUint32(personalState[len(personal)+4:], k)

	indices, err := indicesFromMinimal(n, k, soln)
	if err != nil {
		return false, err
	}

	rows := make([]*node, len(indices))
	for i := range indices {
		rows[i], err = newNode(n, k, personalState, twist, header, indices[i])
		if err != nil {
			return false, err
		}
	}

	for len(rows) > 1 {
		curRows := make([]*node, 0)
		for i := 0; i < len(rows); i += 2 {
			err := validateSubtrees(n, k, rows[i], rows[i+1])
			if err != nil {
				return false, err
			}

			row := newNodeFromChildrenRef(rows[i], rows[i+1], collisionByteLength(n, k))
			curRows = append(curRows, row)
		}

		rows = curRows
	}

	for _, v := range rows[0].hash {
		if v != 0 {
			return false, nil
		}
	}

	return true, nil
}
