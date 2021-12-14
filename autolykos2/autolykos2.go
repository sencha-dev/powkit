package autolykos2

import (
	"encoding/binary"
	"math/big"

	"github.com/sencha-dev/powkit/internal/crypto"
)

/* config */

type Config struct {
	k                      int
	n                      int
	nBase                  uint32
	increaseStart          uint64
	increasePeriodForN     uint64
	nIncreasementHeightMax uint64
}

// k = 32, n = 26
func New(k, n int) *Config {
	cfg := &Config{
		k:                      k,
		n:                      n,
		nBase:                  1 << n,
		increaseStart:          600 * 1024,
		increasePeriodForN:     50 * 1024,
		nIncreasementHeightMax: 4198400,
	}

	return cfg
}

/* helpers */

func generateM(size uint64) []byte {
	m := make([]byte, size*8)
	for i := uint64(0); i < size; i++ {
		binary.BigEndian.PutUint64(m[i*8:i*8+8], i)
	}

	return m
}

func uint32BEToBytes(i uint32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, i)
	return b
}

func concatBytes(a, b []byte) []byte {
	c := make([]byte, len(a), len(a)+len(b))
	copy(c, a)
	c = append(c, b...)

	return c
}

/* algorithm */

func (cfg *Config) calcN(height uint64) uint32 {
	if height > cfg.nIncreasementHeightMax {
		return 2143944600
	}

	n := cfg.nBase
	if height >= cfg.increaseStart {
		iters := int((height-cfg.increaseStart)/cfg.increasePeriodForN + 1)
		for i := 0; i < iters; i++ {
			n = (n / 100) * 105
		}
	}

	return n
}

func (cfg *Config) genIndexes(seed []byte, n uint32) []uint32 {
	hash := crypto.Blake2b256(seed)
	extendedHash := append(hash, hash...)

	indexes := make([]uint32, cfg.k)
	for i := range indexes {
		indexes[i] = binary.BigEndian.Uint32(extendedHash[i:i+4]) % n
	}

	return indexes
}

func (cfg *Config) Compute(msg []byte, nonce, height uint64) []byte {
	m := generateM(1024)
	h := uint32BEToBytes(uint32(height))

	n := cfg.calcN(height)
	bigN := new(big.Int).SetUint64(uint64(n))

	nonceBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(nonceBytes, nonce)

	fullMsg := concatBytes(msg, nonceBytes)
	msgHash := crypto.Blake2b256(fullMsg)
	prei8 := new(big.Int).SetBytes(msgHash[24:32])

	i := uint32BEToBytes(uint32(prei8.Mod(prei8, bigN).Uint64()))
	f := crypto.Blake2b256(concatBytes(i, concatBytes(h, m)))[1:32]

	seed := concatBytes(f, concatBytes(msg, nonceBytes))
	indexes := cfg.genIndexes(seed, n)

	f2 := new(big.Int)
	for _, index := range indexes {
		elem := concatBytes(uint32BEToBytes(index), concatBytes(h, m))
		elemHash := crypto.Blake2b256(elem)[1:]
		f2.Add(f2, new(big.Int).SetBytes(elemHash))
	}

	ha := crypto.Blake2b256(f2.Bytes())

	return ha
}
