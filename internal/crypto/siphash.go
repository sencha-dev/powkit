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
