package convutil

import (
	"encoding/binary"
)

func Uint32ArrayToBytes(arr []uint32, byteOrder binary.ByteOrder) []byte {
	buf := make([]byte, len(arr)*4)
	for i, v := range arr {
		byteOrder.PutUint32(buf[i*4:], v)
	}

	return buf
}

func BytesToUint32Array(buf []byte, byteOrder binary.ByteOrder) []uint32 {
	length := len(buf) / 4
	arr := make([]uint32, length)
	for i := 0; i < length; i++ {
		arr[i] = byteOrder.Uint32(buf[i*4 : (i+1)*4])
	}

	return arr
}

func Uint32ToBytes(val uint32, byteOrder binary.ByteOrder) []byte {
	data := make([]byte, 4)
	byteOrder.PutUint32(data, val)

	return data
}

func Uint64ToBytes(val uint64, byteOrder binary.ByteOrder) []byte {
	data := make([]byte, 8)
	byteOrder.PutUint64(data, val)

	return data
}
