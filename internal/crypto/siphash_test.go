package crypto

import (
	"testing"
)

func TestHash24(t *testing.T) {
	tests := []struct {
		v0     uint64
		v1     uint64
		v2     uint64
		v3     uint64
		nonce  uint64
		result uint64
	}{
		{
			v0:     1,
			v1:     2,
			v2:     3,
			v3:     4,
			nonce:  10,
			result: 928382149599306901,
		},
		{
			v0:     1,
			v1:     2,
			v2:     3,
			v3:     4,
			nonce:  111,
			result: 10524991083049122233,
		},
		{
			v0:     9,
			v1:     7,
			v2:     6,
			v3:     7,
			nonce:  12,
			result: 1305683875471634734,
		},
		{
			v0:     9,
			v1:     7,
			v2:     6,
			v3:     7,
			nonce:  10,
			result: 11589833042187638814,
		},
	}

	for i, tt := range tests {
		hasher := NewSipHasher(tt.v0, tt.v1, tt.v2, tt.v3)
		hasher.Hash24(tt.nonce)
		result := hasher.XorLanes()

		if result != tt.result {
			t.Errorf("failed on %d: result mismatch: have %x, want %x", i, result, tt.result)
		}
	}
}
