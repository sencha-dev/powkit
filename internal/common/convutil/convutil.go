package convutil

import (
	"encoding/binary"
)

func Uint32ArrayToBytesLE(arr []uint32) []byte {
	buf := make([]byte, len(arr)*4)

	for i, v := range arr {
		binary.LittleEndian.PutUint32(buf[i*4:], v)
	}

	return buf
}

func BytesToUint32ArrayLE(buf []byte) []uint32 {
	length := len(buf) / 4
	arr := make([]uint32, length)

	for i := 0; i < length; i++ {
		arr[i] = binary.LittleEndian.Uint32(buf[i*4 : (i+1)*4])
	}

	return arr
}

func Uint32ToBytesBE(val uint32) []byte {
	data := make([]byte, 4)
	binary.BigEndian.PutUint32(data, val)

	return data
}

func Uint64ToBytesBE(val uint64) []byte {
	data := make([]byte, 8)
	binary.BigEndian.PutUint64(data, val)

	return data
}
