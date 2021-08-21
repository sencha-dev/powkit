package pow

import (
	"bytes"
	"testing"
)

func TestComputeETH(t *testing.T) {
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
			hash:   mustDecodeHex("0x69e71ffd37268b6cf7096cdd917c2c175eaaee8eb7afed4b5cf8521b09024818"),
			mix:    mustDecodeHex("0x37fde31175fe180346444d15b4dfc6a9da3b2b41ee2298ceeccaf888b2d45df4"),
			digest: mustDecodeHex("0x00000000000004418223e154cd70e083a4964684bf333f4654dac2fe76b999f6"),
		},
		{
			height: 12000000,
			nonce:  0xb62c052c3d4a3866,
			hash:   mustDecodeHex("0x1940ee93bb48f1982b9fc546ae69ca9a59de4e55f9944900b3abf04436eb1ee1"),
			mix:    mustDecodeHex("0x114f16f97d044682678844ea69212d6764998107ccb698fc6ae4fd2d71a33104"),
			digest: mustDecodeHex("0x000000000000089029050a54342c13f7614c95da4a471ce0ef10c009f84c762b"),
		},
		{
			height: 12965001,
			nonce:  0x956e895d988798e,
			hash:   mustDecodeHex("0xcf133ce0cccd4ad877d671b310c27f5ce19c28c14455dac45b90171bac5581c7"),
			mix:    mustDecodeHex("0xcb3166ebb1888430069b769145b20ba5e3a55f32fd2fa39f0ebdc08d60b4557e"),
			digest: mustDecodeHex("0x00000000000000012923a9ab2605573e0158adeb21c86b22d8ebd33b8ee08856"),
		},
	}

	dag, err := NewLightDag("ETH", false)
	if err != nil {
		t.Errorf("%v", err)
		return
	}

	for i, tt := range tests {
		mix, digest, err := dag.Compute(tt.hash, tt.height, tt.nonce)
		if err != nil {
			t.Errorf("compute - ETH: error computing digest for test %d", i)
			continue
		}

		if bytes.Compare(tt.mix, mix) != 0 {
			t.Errorf("compute - ETH: mixhash does not match for test %d", i)
		}

		if bytes.Compare(tt.digest, digest) != 0 {
			t.Errorf("compute - ETH: digest does not match for test %d", i)
		}
	}
}

func TestComputeETC(t *testing.T) {
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
			hash:   mustDecodeHex("0x9451157911b225460f07e5eac6f63e39a0f4a6952ba544302c6b2aae51f36064"),
			mix:    mustDecodeHex("0x929ce5aad1cfe51bfc88104b0ce6fdde605d21a6e1fa171cafbf7cccf310c626"),
			digest: mustDecodeHex("0x0000000000012c6470946acc24916d4c5d4c4f143fef2db2f40e1d7e2e47776f"),
		},
		{
			height: 12500000,
			nonce:  0x431d9c1839838ea1,
			hash:   mustDecodeHex("0x1d225541aeda9d1110b9e61c0c6c81f12376e9b4111472351196b9d557ba83d8"),
			mix:    mustDecodeHex("0xc280665934366eb1ac3fcf5e6f4e37945e05d764acf2a8106a982826078b613c"),
			digest: mustDecodeHex("0x00000000000069bdf3a3bf3bab6529869d8419e477e53898d6c123d023386b20"),
		},
		{
			height: 13344137,
			nonce:  0x9827862e22a92ff1,
			hash:   mustDecodeHex("0x27eaf677273c9147cd27b99c34b3783243255864a54b169af238750c39b3c167"),
			mix:    mustDecodeHex("0x6dd0879bfe248c4ac73160a3d2554ce12431d2033b5f4464559368d855795df7"),
			digest: mustDecodeHex("0x000000000000ca4bd2875398d73ab24e8467d5986bfe85ec7fca1b860a540d14"),
		},
	}

	dag, err := NewLightDag("ETC", false)
	if err != nil {
		t.Errorf("%v", err)
		return
	}

	for i, tt := range tests {
		mix, digest, err := dag.Compute(tt.hash, tt.height, tt.nonce)
		if err != nil {
			t.Errorf("compute - ETC: error computing digest for test %d", i)
			continue
		}

		if bytes.Compare(tt.mix, mix) != 0 {
			t.Errorf("compute - ETC: mixhash does not match for test %d", i)
		}

		if bytes.Compare(tt.digest, digest) != 0 {
			t.Errorf("compute - ETC: digest does not match for test %d", i)
		}
	}
}

func TestComputeRVN(t *testing.T) {
	tests := []struct {
		height uint64
		nonce  uint64
		hash   []byte
		mix    []byte
		digest []byte
	}{
		{
			height: 1880000,
			nonce:  0x25ca7a0109cf8f2d,
			hash:   mustDecodeHex("439fe77436016853df3ec5ca24d654da32845f389334cadad356a42ef62e19cd"),
			mix:    mustDecodeHex("7ff370c848f6f553fa7ca8d68c3515dca5833d330b048dd8a27ccd4b137dc1d7"),
			digest: mustDecodeHex("00000000000021c9af1f188fcf1c913b2583133b0fe58bb6ff0532dad895fdfc"),
		},
		{
			height: 1888000,
			nonce:  0xdc30900000a4c493,
			hash:   mustDecodeHex("911a676cf0e5077a24e4917483bcca4bbd461a679b1a780b9d15c8b6bf5bc1d7"),
			mix:    mustDecodeHex("28605c4c11c72f1d3af0baef8cca10ccf570aa0c89c93596b2bfd485e30bd9f7"),
			digest: mustDecodeHex("0000000000000d6ea9c4b4b759a157bc2ef74bd9a7ed04ca6ac9611e218ccddc"),
		},
		{
			height: 1888509,
			nonce:  0xf09b0e1342275f3f,
			hash:   mustDecodeHex("14f2c18d74d48abe637437458c10ff5283a9a5197e8b5e740a161f4411b97a43"),
			mix:    mustDecodeHex("3dedcc6fb28c6bf3f5d203f29188bef3ff86688be34c93f28bd227eced9226e4"),
			digest: mustDecodeHex("0000000000005e6585e5e6ab7e4d75a98810204341def05823ad3a5ca1fa0d83"),
		},
	}

	dag, err := NewLightDag("RVN", false)
	if err != nil {
		t.Errorf("%v", err)
		return
	}

	for i, tt := range tests {
		mix, digest, err := dag.Compute(tt.hash, tt.height, tt.nonce)
		if err != nil {
			t.Errorf("compute - RVN: error computing digest for test %d", i)
			continue
		}

		if bytes.Compare(tt.mix, mix) != 0 {
			t.Errorf("compute - RVN: mixhash does not match for test %d", i)
		}

		if bytes.Compare(tt.digest, digest) != 0 {
			t.Errorf("compute - RVN: digest does not match for test %d", i)
		}
	}
}
