package crypto

import (
	"testing"
)

func TestXoshiro256PlusPlusHasher(t *testing.T) {
	tests := []struct {
		state  [4]uint64
		values []uint64
	}{
		{
			state: [4]uint64{
				0,
				0,
				0,
				1,
			},
			values: []uint64{
				0x0000000000800000, 0x0000000000800011, 0x0002200000800011, 0x0002200044000020,
				0x00022000448008a1, 0x01102100802208a0, 0x0114002284200880, 0x0402896644a64095,
				0x8a14a0a244acc891, 0x8c88ba4680b20495, 0x8b48b35586484291, 0x09d8608086e66cd5,
			},
		},
		{
			state: [4]uint64{
				124123,
				591204,
				959691,
				959109,
			},
			values: []uint64{
				0x00000843b001e4db, 0x000003219d80c14a, 0x44be2003275edbad, 0x0e040b0682d9cf5b,
				0xae50d583e48813b7, 0x868e6c5267cb9e0e, 0x278a90e35d6235c9, 0xe0fbdb62d8981ac2,
				0xe9f929070cec215b, 0x933060e9112c0026, 0x8e3c312c503989f4, 0x0cb10bc24b1fbeca,
			},
		},
	}

	for i, tt := range tests {
		hasher := NewXoshiro256PlusPlusHasher(tt.state[0], tt.state[1], tt.state[2], tt.state[3])
		for j, expected := range tt.values {
			value := hasher.Next()
			if value != expected {
				t.Errorf("failed on %d: %d: have %d, want %d", i, j, value, expected)
			}
		}
	}
}
