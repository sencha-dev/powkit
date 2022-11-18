// Copyright (c) 2018-2019 The kaspanet developers

package heavyhash

import (
	"encoding/binary"

	"github.com/sencha-dev/powkit/internal/crypto"
)

const (
	size       = 64
	iterations = size / 4
)

func heavyHash(hash []byte, timestamp int64, nonce uint64) []byte {
	// initialize the matrix
	s0 := binary.LittleEndian.Uint64(hash[0:8])
	s1 := binary.LittleEndian.Uint64(hash[8:16])
	s2 := binary.LittleEndian.Uint64(hash[16:24])
	s3 := binary.LittleEndian.Uint64(hash[24:32])
	mat := newMatrix(s0, s1, s2, s3)

	// build the header
	header := make([]byte, 32+8+32+8)
	copy(header[:32], hash)
	binary.LittleEndian.PutUint64(header[32:40], uint64(timestamp))
	binary.LittleEndian.PutUint64(header[72:80], nonce)
	header = crypto.CShake256(header, []byte("ProofOfWorkHash"), 32)

	// initialize the vector and product arrays
	var v, p [size]uint16
	for i := 0; i < size/2; i++ {
		v[i*2] = uint16(header[i] >> 4)
		v[i*2+1] = uint16(header[i] & 0x0f)
	}

	// build the product array
	for i := 0; i < size; i++ {
		var s uint16
		for j := 0; j < size; j++ {
			s += mat[i][j] * v[j]
		}
		p[i] = s >> 10
	}

	// calculate the digest
	digest := make([]byte, 32)
	for i := range digest {
		digest[i] = header[i] ^ (byte(p[i*2]<<4) | byte(p[i*2+1]))
	}

	// hash the digest a final time, reverse bytes
	digest = crypto.CShake256(digest, []byte("HeavyHash"), 32)
	for i, j := 0, len(digest)-1; i < j; i, j = i+1, j-1 {
		digest[i], digest[j] = digest[j], digest[i]
	}

	return digest
}
