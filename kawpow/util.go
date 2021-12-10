// Copyright 2019 Victor Tran
// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package kawpow

import (
	"encoding/binary"
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
