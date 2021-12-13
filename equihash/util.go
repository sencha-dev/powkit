package equihash

import (
	"encoding/binary"
	"unsafe"
)

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

func bytesToUint32ArrayLE(arr []byte) []uint32 {
	length := len(arr) / 4
	data := make([]uint32, length)

	for i := 0; i < length; i++ {
		data[i] = binary.LittleEndian.Uint32(arr[i*4 : (i+1)*4])
	}

	return data
}

func bytesToUint32ArrayBE(arr []byte) []uint32 {
	length := len(arr) / 4
	data := make([]uint32, length)

	for i := 0; i < length; i++ {
		data[i] = binary.BigEndian.Uint32(arr[i*4 : (i+1)*4])
	}

	return data
}
