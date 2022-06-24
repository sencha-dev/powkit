package verthash

import (
	"bytes"
	"encoding/binary"

	"github.com/sencha-dev/powkit/internal/crypto"

	"golang.org/x/crypto/sha3"
)

func verthash(data, input []byte) []byte {
	p1 := [32]byte{}

	inputCopy := make([]byte, len(input))
	copy(inputCopy[:], input[:])
	sha3hash := sha3.Sum256(inputCopy)

	copy(p1[:], sha3hash[:])
	p0 := make([]byte, verthashSubset)
	for i := uint32(0); i < verthashIter; i++ {
		inputCopy[0] += 0x01
		digest64 := sha3.Sum512(inputCopy)
		copy(p0[i*verthashP0Size:], digest64[:])
	}

	buf := bytes.NewBuffer(p0)
	p0Index := make([]uint32, len(p0)/4)
	for i := 0; i < len(p0Index); i++ {
		binary.Read(buf, binary.LittleEndian, &p0Index[i])
	}

	seekIndexes := make([]uint32, verthashIndexes)
	for x := uint32(0); x < verthashRotations; x++ {
		copy(seekIndexes[x*verthashSubset/4:], p0Index)
		for y := 0; y < len(p0Index); y++ {
			p0Index[y] = (p0Index[y] << 1) | (1 & (p0Index[y] >> 31))
		}
	}

	mdiv := ((uint32(datasetSize) - verthashHashOutSize) / verthashByteAlignment) + 1
	valueAccumulator := uint32(0x811c9dc5)
	buf = bytes.NewBuffer(p1[:])
	p1Arr := make([]uint32, verthashHashOutSize/4)
	for i := 0; i < len(p1Arr); i++ {
		binary.Read(buf, binary.LittleEndian, &p1Arr[i])
	}
	for i := uint32(0); i < verthashIndexes; i++ {
		offset := (crypto.Fnv1a(seekIndexes[i], valueAccumulator) % mdiv) * verthashByteAlignment
		data := data[offset : offset+verthashHashOutSize]
		for i2 := uint32(0); i2 < verthashHashOutSize/4; i2++ {
			value := binary.LittleEndian.Uint32(data[i2*4 : ((i2 + 1) * 4)])
			p1Arr[i2] = crypto.Fnv1a(p1Arr[i2], value)
			valueAccumulator = crypto.Fnv1a(valueAccumulator, value)
		}
	}

	for i := uint32(0); i < verthashHashOutSize/4; i++ {
		binary.LittleEndian.PutUint32(p1[i*4:], p1Arr[i])
	}

	return p1[:]
}
