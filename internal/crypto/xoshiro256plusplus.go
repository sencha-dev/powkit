package crypto

func rotl(a, b uint64) uint64 {
	return (a << b) | (a >> (64 - b))
}

type Xoshiro256PlusPlusHasher struct {
	s0 uint64
	s1 uint64
	s2 uint64
	s3 uint64
}

func NewXoshiro256PlusPlusHasher(s0, s1, s2, s3 uint64) Xoshiro256PlusPlusHasher {
	hasher := Xoshiro256PlusPlusHasher{
		s0: s0,
		s1: s1,
		s2: s2,
		s3: s3,
	}

	return hasher
}

func (h *Xoshiro256PlusPlusHasher) Next() uint64 {
	value := rotl(h.s0+h.s3, 23) + h.s0
	state := h.s1 << 17

	h.s2 ^= h.s0
	h.s3 ^= h.s1
	h.s1 ^= h.s2
	h.s0 ^= h.s3

	h.s2 ^= state
	h.s3 = rotl(h.s3, 45)

	return value
}
