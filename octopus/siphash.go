package octopus

type SipHasher struct {
	v0 uint64
	v1 uint64
	v2 uint64
	v3 uint64
}

func rotl(x, b uint64) uint64 {
	return (x << b) | (x >> (64 - b))
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

func (h *SipHasher) xorLanes() uint64 {
	return h.v0 ^ h.v1 ^ h.v2 ^ h.v3
}

func (h *SipHasher) sipRound() {
	h.v0 = h.v0 + h.v1
	h.v2 = h.v2 + h.v3
	h.v1 = rotl(h.v1, 13)
	h.v3 = rotl(h.v3, 16)
	h.v1 ^= h.v0
	h.v3 ^= h.v2
	h.v0 = rotl(h.v0, 32)
	h.v2 = h.v2 + h.v1
	h.v0 = h.v0 + h.v3
	h.v1 = rotl(h.v1, 17)
	h.v3 = rotl(h.v3, 21)
	h.v1 ^= h.v2
	h.v3 ^= h.v0
	h.v2 = rotl(h.v2, 32)
}

func (h *SipHasher) hash24(nonce uint64) {
	h.v3 ^= nonce
	h.sipRound()
	h.sipRound()
	h.v0 ^= nonce
	h.v2 ^= 0xff
	h.sipRound()
	h.sipRound()
	h.sipRound()
	h.sipRound()
}
