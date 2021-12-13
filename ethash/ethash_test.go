package ethash

import (
	"bytes"
	"testing"

	"github.com/sencha-dev/powkit/internal/common/testutil"
)

func TestComputeEthereum(t *testing.T) {
	tests := []struct {
		height uint64
		nonce  uint64
		hash   []byte
		mix    []byte
		digest []byte
	}{
		{
			height: 10000000,
			nonce:  0x2f6923f80426f157,
			hash:   testutil.MustDecodeHex("0x69e71ffd37268b6cf7096cdd917c2c175eaaee8eb7afed4b5cf8521b09024818"),
			mix:    testutil.MustDecodeHex("0x37fde31175fe180346444d15b4dfc6a9da3b2b41ee2298ceeccaf888b2d45df4"),
			digest: testutil.MustDecodeHex("0x00000000000004418223e154cd70e083a4964684bf333f4654dac2fe76b999f6"),
		},
		{
			height: 12000000,
			nonce:  0xb62c052c3d4a3866,
			hash:   testutil.MustDecodeHex("0x1940ee93bb48f1982b9fc546ae69ca9a59de4e55f9944900b3abf04436eb1ee1"),
			mix:    testutil.MustDecodeHex("0x114f16f97d044682678844ea69212d6764998107ccb698fc6ae4fd2d71a33104"),
			digest: testutil.MustDecodeHex("0x000000000000089029050a54342c13f7614c95da4a471ce0ef10c009f84c762b"),
		},
		{
			height: 12965001,
			nonce:  0x956e895d988798e,
			hash:   testutil.MustDecodeHex("0xcf133ce0cccd4ad877d671b310c27f5ce19c28c14455dac45b90171bac5581c7"),
			mix:    testutil.MustDecodeHex("0xcb3166ebb1888430069b769145b20ba5e3a55f32fd2fa39f0ebdc08d60b4557e"),
			digest: testutil.MustDecodeHex("0x00000000000000012923a9ab2605573e0158adeb21c86b22d8ebd33b8ee08856"),
		},
	}

	client := NewEthereum()
	for i, tt := range tests {
		mix, digest := client.Compute(tt.height, tt.nonce, tt.hash)

		if bytes.Compare(mix, tt.mix) != 0 {
			t.Errorf("failed on %d: mix mismatch: have %x, want %x", i, mix, tt.mix)
		}

		if bytes.Compare(digest, tt.digest) != 0 {
			t.Errorf("failed on %d: digest mismatch: have %x, want %x", i, digest, tt.digest)
		}
	}
}

func TestComputeEthereumClassic(t *testing.T) {
	tests := []struct {
		height uint64
		nonce  uint64
		hash   []byte
		mix    []byte
		digest []byte
	}{
		{
			height: 12000000,
			nonce:  0x37850ed39b8fedee,
			hash:   testutil.MustDecodeHex("0x9451157911b225460f07e5eac6f63e39a0f4a6952ba544302c6b2aae51f36064"),
			mix:    testutil.MustDecodeHex("0x929ce5aad1cfe51bfc88104b0ce6fdde605d21a6e1fa171cafbf7cccf310c626"),
			digest: testutil.MustDecodeHex("0x0000000000012c6470946acc24916d4c5d4c4f143fef2db2f40e1d7e2e47776f"),
		},
		{
			height: 12500000,
			nonce:  0x431d9c1839838ea1,
			hash:   testutil.MustDecodeHex("0x1d225541aeda9d1110b9e61c0c6c81f12376e9b4111472351196b9d557ba83d8"),
			mix:    testutil.MustDecodeHex("0xc280665934366eb1ac3fcf5e6f4e37945e05d764acf2a8106a982826078b613c"),
			digest: testutil.MustDecodeHex("0x00000000000069bdf3a3bf3bab6529869d8419e477e53898d6c123d023386b20"),
		},
		{
			height: 13344137,
			nonce:  0x9827862e22a92ff1,
			hash:   testutil.MustDecodeHex("0x27eaf677273c9147cd27b99c34b3783243255864a54b169af238750c39b3c167"),
			mix:    testutil.MustDecodeHex("0x6dd0879bfe248c4ac73160a3d2554ce12431d2033b5f4464559368d855795df7"),
			digest: testutil.MustDecodeHex("0x000000000000ca4bd2875398d73ab24e8467d5986bfe85ec7fca1b860a540d14"),
		},
	}

	client := NewEthereumClassic()
	for i, tt := range tests {
		mix, digest := client.Compute(tt.height, tt.nonce, tt.hash)

		if bytes.Compare(mix, tt.mix) != 0 {
			t.Errorf("failed on %d: mix mismatch: have %x, want %x", i, mix, tt.mix)
		}

		if bytes.Compare(digest, tt.digest) != 0 {
			t.Errorf("failed on %d: digest mismatch: have %x, want %x", i, digest, tt.digest)
		}
	}
}
