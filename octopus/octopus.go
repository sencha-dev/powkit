//  Copyright (C) 2019 Conflux Foundation

package octopus

import (
	"encoding/binary"
	"math"

	"github.com/sencha-dev/powkit/internal/common/convutil"
	"github.com/sencha-dev/powkit/internal/crypto"
)

const (
	powCacheRounds    = 3
	powMixBytes       = 256
	powAccesses       = 32
	powDatasetParents = 256
	powMod            = 1032193
	powModB           = 11
	powNK             = 10
	powN              = 1 << powNK
	powWarpSize       = 32
	powDataPerThread  = powN / powWarpSize

	nodeBytes  = 64
	nodeWords  = nodeBytes / 4
	nodeDwords = nodeWords / 2

	mixWords = powMixBytes / 4
	mixNodes = mixWords / nodeWords
)

func powerMod(a, n uint64) uint64 {
	var result uint64 = 1
	for n > 0 {
		if n%2 == 1 {
			result = result * a % powMod
		}

		a = a * a % powMod
		n >>= 1
	}

	return result
}

func gcd(a, b uint64) uint64 {
	if b == 0 {
		return a
	}

	return gcd(b, a%b)
}

func remap(num uint64) uint64 {
	e := num%(powMod-2) + 1
	for {
		g := gcd(e, powMod-1)
		if g == 1 {
			break
		}

		e /= g
	}

	return powerMod(powModB, e)
}

func computeC(a, b, h uint64) uint64 {
	for {
		c := remap(h)
		if b*b%powMod != 4*a*c%powMod {
			return c
		}
		h++
	}
}

func octopus(hash []byte, nonce, datasetSize uint64, lookup func(index uint32) []uint32) []byte {
	v0 := binary.LittleEndian.Uint64(hash[0:8])
	v1 := binary.LittleEndian.Uint64(hash[8:16])
	v2 := binary.LittleEndian.Uint64(hash[16:24])
	v3 := binary.LittleEndian.Uint64(hash[24:32])
	d := make([]uint32, powN)

	a := remap(v0)
	b := remap(v1)
	c := computeC(a, b, v2)
	w := remap(v3)

	warpID := nonce / powWarpSize
	for i := 0; i < powWarpSize; i++ {
		hasher := crypto.NewSipHasher(v0, v1, v2, v3)
		hasher.Hash24(warpID*powWarpSize + uint64(i))
		for j := 0; j < powDataPerThread; j++ {
			hasher.SipRound()
			d[(j*powWarpSize + i)] = uint32((hasher.XorLanes() & math.MaxUint32) % powMod)
		}
	}

	w2 := w * w % powMod
	var wpow, w2pow uint64 = 1, 1
	for i := 0; uint64(i) < nonce%powWarpSize; i++ {
		wpow = wpow * w % powMod
		w2pow = w2pow * w2 % powMod
	}

	wpowFull, w2powFull := wpow, w2pow
	for i := nonce % powWarpSize; i < powWarpSize; i++ {
		wpowFull = wpowFull * w % powMod
		w2powFull = w2powFull * w2 % powMod
	}

	var result uint64
	resBuf := make([]uint32, powDataPerThread)
	for i := 0; i < powDataPerThread; i++ {
		x := (a*w2pow + b*wpow + c) % powMod

		var pv uint64
		for j := 0; j < powN; j++ {
			pv = (pv*x + uint64(d[powN-j-1])) % powMod
		}

		resBuf[i] = uint32(pv)
		result = crypto.Fnv1Uint64(result, pv)
		if i+1 < powDataPerThread {
			wpow = wpow * wpowFull % powMod
			w2pow = w2pow * w2powFull % powMod
		}
	}

	halfMix := make([]byte, nodeBytes)
	copy(halfMix, hash)
	binary.LittleEndian.PutUint64(halfMix[len(hash):], result)
	halfMix = crypto.Keccak512(halfMix[:len(hash)+8])

	var mix []byte
	for i := 0; i < 4; i++ {
		mix = append(mix, halfMix...)
	}

	pageSize := 4 * mixWords
	numFullPages := uint32(datasetSize / uint64(pageSize))
	firstVal := binary.LittleEndian.Uint32(halfMix[:4])

	mixUWords := convutil.BytesToUint32Array(mix, binary.LittleEndian)
	for i := 0; i < powAccesses; i++ {
		idx := crypto.Fnv1(firstVal^uint32(i)^resBuf[i], mixUWords[i%mixWords]) % numFullPages
		for n := 0; n < mixNodes; n++ {
			tmpNode := lookup(idx*mixNodes + uint32(n))

			// @TODO: the x4 is weird
			for l, a := range mixUWords[n*mixNodes*4 : (n+1)*mixNodes*4] {
				// @TODO: the x4 is weird
				mixUWords[n*mixNodes*4+l] = crypto.Fnv1(a, tmpNode[l])
			}
		}
	}

	compressWords := make([]uint32, 8)
	for i := 0; i < 8; i++ {
		w := i * 4
		w2 := (8 + i) * 4

		reduction := mixUWords[w+0]
		reduction = reduction*crypto.FnvPrime ^ mixUWords[w+1]
		reduction = reduction*crypto.FnvPrime ^ mixUWords[w+2]
		reduction = reduction*crypto.FnvPrime ^ mixUWords[w+3]

		reduction2 := mixUWords[w2+0]
		reduction2 = reduction2*crypto.FnvPrime ^ mixUWords[w2+1]
		reduction2 = reduction2*crypto.FnvPrime ^ mixUWords[w2+2]
		reduction2 = reduction2*crypto.FnvPrime ^ mixUWords[w2+3]

		compressWords[i] = reduction*crypto.FnvPrime ^ reduction2
	}

	compressBytes := convutil.Uint32ArrayToBytes(compressWords, binary.LittleEndian)
	digest := crypto.Keccak256(append(mix[:nodeBytes], compressBytes...))

	return digest
}
