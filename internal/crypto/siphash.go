package crypto

type SipHasher struct {
	v0 uint64
	v1 uint64
	v2 uint64
	v3 uint64
}

func NewSipHasher(v0, v1, v2, v3 uint64) *SipHasher {
	hasher := &SipHasher{
		v0: v0,
		v1: v1,
		v2: v2,
		v3: v3,
	}

	return hasher
}

func (h *SipHasher) XorLanes() uint64 {
	return h.v0 ^ h.v1 ^ h.v2 ^ h.v3
}

func (h *SipHasher) SipRound() {
	h.v0 += h.v1
	h.v1 = h.v1<<13 | h.v1>>51
	h.v1 ^= h.v0
	h.v0 = h.v0<<32 | h.v0>>32

	h.v2 += h.v3
	h.v3 = h.v3<<16 | h.v3>>48
	h.v3 ^= h.v2

	h.v0 += h.v3
	h.v3 = h.v3<<21 | h.v3>>43
	h.v3 ^= h.v0

	h.v2 += h.v1
	h.v1 = h.v1<<17 | h.v1>>47
	h.v1 ^= h.v2
	h.v2 = h.v2<<32 | h.v2>>32
}

func (h *SipHasher) Hash24(nonce uint64) {
	h.v3 ^= nonce
	h.SipRound()
	h.SipRound()
	h.v0 ^= nonce
	h.v2 ^= 0xff
	h.SipRound()
	h.SipRound()
	h.SipRound()
	h.SipRound()
}

func (h *SipHasher) Hash48(nonce uint64) {
	h.v3 ^= nonce
	h.SipRound()
	h.SipRound()
	h.SipRound()
	h.SipRound()
	h.v0 ^= nonce
	h.v2 ^= 0xff
	h.SipRound()
	h.SipRound()
	h.SipRound()
	h.SipRound()
	h.SipRound()
	h.SipRound()
	h.SipRound()
	h.SipRound()
}

type SipNodeFunc func(uint64, [4]uint64, uint64, uint64) uint64

func SipNode24(edgeMask uint64, siphashKeys [4]uint64, edge, uorv uint64) uint64 {
	hasher := NewSipHasher(siphashKeys[0], siphashKeys[1], siphashKeys[2], siphashKeys[3])
	hasher.Hash24(2*edge + uorv)
	value := hasher.XorLanes()

	return value & edgeMask
}

func SipNode24Legacy(edgeMask uint64, siphashKeys [4]uint64, edge, uorv uint64) uint64 {
	hasher := NewSipHasher(siphashKeys[0], siphashKeys[1], siphashKeys[2], siphashKeys[3])
	hasher.Hash24(2*edge + uorv)

	value := hasher.XorLanes()
	value = value<<17 | value>>47

	return value & edgeMask
}

type SipBlockFunc func([4]uint64, uint64) uint64

func SipBlock48(siphashKeys [4]uint64, edge uint64) uint64 {
	const edgeBlockBits uint64 = 6
	const edgeBlockSize uint64 = (1 << edgeBlockBits)
	const edgeBlockMask uint64 = (edgeBlockSize - 1)

	hasher := NewSipHasher(siphashKeys[0], siphashKeys[1], siphashKeys[2], siphashKeys[3])
	block := make([]uint64, edgeBlockSize)
	edge0 := edge & ^edgeBlockMask

	var i uint64
	for i = 0; i < edgeBlockSize; i++ {
		hasher.Hash48(edge0 + uint64(i))
		block[i] = hasher.XorLanes()
	}

	last := block[edgeBlockMask]
	for i = 0; i < edgeBlockMask; i++ {
		block[i] ^= last
	}

	return block[edge&edgeBlockMask]
}
