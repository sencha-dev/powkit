package kawpow

import (
	"encoding/binary"
	"math/bits"
	"unsafe"
)

func isLittleEndian() bool {
	n := uint32(0x01020304)
	return *(*byte)(unsafe.Pointer(&n)) == 0x04
}

/* Array utils */

func uint32ArrayToBytes(c []uint32) []byte {
	buf := make([]byte, len(c)*4)
	if isLittleEndian() {
		for i, v := range c {
			binary.LittleEndian.PutUint32(buf[i*4:], v)
		}
	} else {
		for i, v := range c {
			binary.BigEndian.PutUint32(buf[i*4:], v)
		}
	}
	return buf
}

func bytesToUint32Array(arr []byte) []uint32 {
	length := len(arr) / 4
	data := make([]uint32, length)

	for i := 0; i < length; i++ {
		data[i] = binary.LittleEndian.Uint32(arr[i*4 : (i+1)*4])
	}

	return data
}

/* Math utils */

func minUint32(a, b uint32) uint32 {
	if a > b {
		return b
	}

	return a
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

// following functionst taken from github.com/pkt-cash/pktd/

func rotl32(a, b uint32) uint32 {
	return a<<(b&31) | a>>((32-b)&31)
}

func rotr32(a, b uint32) uint32 {
	return a<<((32-b)&31) | a>>(b&31)
}

func clz32(a uint32) uint32 {
	return uint32(bits.LeadingZeros32(a))
}

func popcount32(a uint32) uint32 {
	return uint32(bits.OnesCount32(a))
}

func mul_hi32(a, b uint32) uint32 {
	return uint32((uint64(a) * uint64(b)) >> 32)
}

/* kawpow utils */

func randomMath(a, b, selector uint32) uint32 {
	switch selector % 11 {
	case 0:
		return a + b
	case 1:
		return a * b
	case 2:
		return mul_hi32(a, b)
	case 3:
		return minUint32(a, b)
	case 4:
		return rotl32(a, b)
	case 5:
		return rotr32(a, b)
	case 6:
		return a & b
	case 7:
		return a | b
	case 8:
		return a ^ b
	case 9:
		return clz32(a) + clz32(b)
	case 10:
		return popcount32(a) + popcount32(b)
	}

	return 0
}

func randomMerge(a, b, selector uint32) uint32 {
	x := ((selector >> 16) % 31) + 1

	switch selector % 4 {
	case 0:
		return (a * 33) + b
	case 1:
		return (a ^ b) * 33
	case 2:
		return rotl32(a, x) ^ b
	case 3:
		return rotr32(a, x) ^ b
	}

	return 0
}
