package autolykos2

import (
	"encoding/binary"
	"math/big"

	"github.com/sencha-dev/powkit/internal/common/convutil"
	"github.com/sencha-dev/powkit/internal/crypto"
)

const (
	increaseStart          = 600 * 1024
	increasePeriodForN     = 50 * 1024
	nIncreasementHeightMax = 4198400
)

func concatBytes(a, b []byte) []byte {
	c := make([]byte, len(a), len(a)+len(b))
	copy(c, a)
	c = append(c, b...)

	return c
}

func generateM(size uint64) []byte {
	m := make([]byte, size*8)
	for i := uint64(0); i < size; i++ {
		binary.BigEndian.PutUint64(m[i*8:i*8+8], i)
	}

	return m
}

func calcN(nBase uint32, height uint64) uint32 {
	if height > nIncreasementHeightMax {
		return 2143944600
	}

	var n uint32 = nBase
	if height >= increaseStart {
		iters := int((height-increaseStart)/increasePeriodForN + 1)
		for i := 0; i < iters; i++ {
			n = (n / 100) * 105
		}
	}

	return n
}

func genIndexes(seed []byte, k, n uint32) []uint32 {
	hash := crypto.Blake2b256(seed)
	extendedHash := append(hash, hash...)

	indexes := make([]uint32, k)
	for i := range indexes {
		indexes[i] = binary.BigEndian.Uint32(extendedHash[i:i+4]) % n
	}

	return indexes
}

func compute(k, nBase uint32, msg []byte, nonce, height uint64) []byte {
	m := generateM(1024)
	h := convutil.Uint32ToBytes(uint32(height), binary.BigEndian)
	nonceBytes := convutil.Uint64ToBytes(nonce, binary.BigEndian)

	n := calcN(nBase, height)
	bigN := new(big.Int).SetUint64(uint64(n))

	fullMsg := concatBytes(msg, nonceBytes)
	msgHash := crypto.Blake2b256(fullMsg)
	prei8 := new(big.Int).SetBytes(msgHash[24:32])

	i := convutil.Uint32ToBytes(uint32(prei8.Mod(prei8, bigN).Uint64()), binary.BigEndian)
	f := crypto.Blake2b256(concatBytes(i, concatBytes(h, m)))[1:32]

	seed := concatBytes(f, concatBytes(msg, nonceBytes))
	indexes := genIndexes(seed, k, n)

	f2 := new(big.Int)
	for _, index := range indexes {
		elem := concatBytes(convutil.Uint32ToBytes(index, binary.BigEndian), concatBytes(h, m))
		elemHash := crypto.Blake2b256(elem)[1:]
		f2.Add(f2, new(big.Int).SetBytes(elemHash))
	}

	ha := crypto.Blake2b256(f2.Bytes())

	return ha
}
