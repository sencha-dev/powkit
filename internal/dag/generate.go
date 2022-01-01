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

package dag

import (
	"encoding/binary"
	"reflect"
	"unsafe"

	"github.com/sencha-dev/powkit/internal/common/bitutil"
	"github.com/sencha-dev/powkit/internal/crypto"
)

// generateCache creates a verification cache of a given size for an input seed.
// The cache production process involves first sequentially filling up 32 MB of
// memory, then performing two passes of Sergio Demian Lerner's RandMemoHash
// algorithm from Strict Memory Hard Hashing Functions (2014). The output is a
// set of 524288 64-byte values.
// This method places the result into dest in machine byte order.
func (d *DAG) generateCache(dest []uint32, epoch uint64, seed []byte) {
	// Convert our destination slice to a byte buffer
	header := *(*reflect.SliceHeader)(unsafe.Pointer(&dest))
	header.Len *= 4
	header.Cap *= 4
	cache := *(*[]byte)(unsafe.Pointer(&header))

	// Calculate the number of theoretical rows (we'll store in one buffer nonetheless)
	size := uint64(len(cache))
	rows := int(size) / hashBytes

	// Create a hasher to reuse between invocations
	keccak512Hasher := crypto.NewKeccak512Hasher()

	// Sequentially produce the initial dataset
	keccak512Hasher(cache, seed)
	for offset := uint64(hashBytes); offset < size; offset += hashBytes {
		keccak512Hasher(cache[offset:], cache[offset-hashBytes:offset])
	}

	// Use a low-round version of randmemohash
	temp := make([]byte, hashBytes)

	for i := 0; i < d.CacheRounds; i++ {
		for j := 0; j < rows; j++ {
			var (
				srcOff = ((j - 1 + rows) % rows) * hashBytes
				dstOff = j * hashBytes
				xorOff = (binary.LittleEndian.Uint32(cache[dstOff:]) % uint32(rows)) * hashBytes
			)
			bitutil.XORBytes(temp, cache[srcOff:srcOff+hashBytes], cache[xorOff:xorOff+hashBytes])
			keccak512Hasher(cache[dstOff:], temp)
		}
	}
}

func (d *DAG) generateL1Cache(dest []uint32, cache []uint32) {
	keccak512Hasher := crypto.NewKeccak512Hasher()

	header := *(*reflect.SliceHeader)(unsafe.Pointer(&dest))
	header.Len *= 4
	header.Cap *= 4
	l1 := *(*[]byte)(unsafe.Pointer(&header))

	size := uint64(len(l1))
	rows := int(size) / hashBytes

	for i := 0; i < rows; i++ {
		item := d.generateDatasetItem(cache, uint32(i), keccak512Hasher)
		copy(l1[i*hashBytes:], item)
	}
}

// generateDatasetItem combines data from 256 pseudorandomly selected cache nodes,
// and hashes that to compute a single dataset node.
func (d *DAG) generateDatasetItem(cache []uint32, index uint32, keccak512Hasher crypto.Hasher) []byte {
	// Calculate the number of theoretical rows (we use one buffer nonetheless)
	rows := uint32(len(cache) / hashWords)

	// Initialize the mix
	mix := make([]byte, hashBytes)

	binary.LittleEndian.PutUint32(mix, cache[(index%rows)*hashWords]^index)
	for i := 1; i < hashWords; i++ {
		binary.LittleEndian.PutUint32(mix[i*4:], cache[(index%rows)*hashWords+uint32(i)])
	}

	keccak512Hasher(mix, mix)

	// Convert the mix to uint32s to avoid constant bit shifting
	intMix := make([]uint32, hashWords)
	for i := 0; i < len(intMix); i++ {
		intMix[i] = binary.LittleEndian.Uint32(mix[i*4:])
	}

	// fnv it with a lot of random cache nodes based on index
	for i := uint32(0); i < d.DatasetParents; i++ {
		parent := crypto.Fnv1(index^i, intMix[i%16]) % rows
		crypto.FnvHash(intMix, cache[parent*hashWords:])
	}

	// Flatten the uint32 mix into a binary one and return
	for i, val := range intMix {
		binary.LittleEndian.PutUint32(mix[i*4:], val)
	}

	keccak512Hasher(mix, mix)

	return mix
}

func (d *DAG) generateDatasetItemUint(cache []uint32, index, size uint32, keccak512Hasher crypto.Hasher) []uint32 {
	data := make([]uint32, hashWords*size)
	for n := 0; n < int(size); n++ {
		item := d.generateDatasetItem(cache, index*size+uint32(n), keccak512Hasher)

		for i := 0; i < hashWords; i++ {
			data[n*hashWords+i] = binary.LittleEndian.Uint32(item[i*4:])
		}
	}

	return data
}
