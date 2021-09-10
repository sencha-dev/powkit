package equihash

import (
	"encoding/binary"
	"encoding/hex"
	"strings"
	"unsafe"
)

func mustDecodeHex(inp string) []byte {
	inp = strings.Replace(inp, "0x", "", -1)
	out, err := hex.DecodeString(inp)
	if err != nil {
		panic(err)
	}

	return out
}

func isLittleEndian() bool {
	n := uint32(0x01020304)
	return *(*byte)(unsafe.Pointer(&n)) == 0x04
}

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
