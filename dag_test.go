package pow

import (
	"bytes"
	"testing"
)


func TestComputeETH(t *testing.T) {
	tests := []struct {
		height uint64
		nonce uint64
		hash []byte
		mix []byte
		digest []byte
	}{
		{
			height: 12000000,
			nonce: 0xb62c052c3d4a3866,
			hash: MustDecodeHex("0x1940ee93bb48f1982b9fc546ae69ca9a59de4e55f9944900b3abf04436eb1ee1"),
			mix: MustDecodeHex("0x114f16f97d044682678844ea69212d6764998107ccb698fc6ae4fd2d71a33104"),
		},
	}

	dag, err := NewLightDag("ETH")
	if err != nil {
		t.Errorf("%v", err)
		return
	}

	for i, tt := range tests {
		mix, _, err := dag.Compute(tt.hash, tt.height, tt.nonce)
		if err != nil {
			t.Errorf("compute - ETH: error computing digest for test %d", i)
			continue
		}

		if bytes.Compare(tt.mix, mix) != 0 {
			t.Errorf("compute - ETH: mixhash does not match for test %d", i)
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
			hash:   MustDecodeHex("0x9451157911b225460f07e5eac6f63e39a0f4a6952ba544302c6b2aae51f36064"),
			mix:    MustDecodeHex("0x929ce5aad1cfe51bfc88104b0ce6fdde605d21a6e1fa171cafbf7cccf310c626"),
		},
		{
			height: 13000000,
			nonce:  0x1f636f90ee9459db,
			hash:   MustDecodeHex("0xae3cb028681a4de56d3c2c994c2a054f1b1052232d100e5fc3ef8bc5e9abd562"),
			mix:    MustDecodeHex("0x55354b04d489599ed8a999667501a6d950bc511a16ddb9dbdaa81d66667dfb42"),
		},
		{
			height: 13344137,
			nonce:  0x2d9b9d79ff89caf2,
			hash:   MustDecodeHex("0xfed3eb74e36f2abc693168c7f96675a7302a38b1bea2b6468d69fcbde270541f"),
			mix:    MustDecodeHex("0x91776c8615a6ed101d80ebf31458aeb07e770d50d6b8eedef1615fcf6654261b"),
		},
	}

	dag, err := NewLightDag("ETC")
	if err != nil {
		t.Errorf("%v", err)
		return
	}

	for i, tt := range tests {
		// @TODO: add digest verification

		mix, _, err := dag.Compute(tt.hash, tt.height, tt.nonce)
		if err != nil {
			t.Errorf("compute - ETC: error computing digest for test %d", i)
			continue
		}

		if bytes.Compare(tt.mix, mix) != 0 {
			t.Errorf("compute - ETC: mixhash does not match for test %d", i)
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
			hash:   MustDecodeHex("439fe77436016853df3ec5ca24d654da32845f389334cadad356a42ef62e19cd"),
			mix:    MustDecodeHex("7ff370c848f6f553fa7ca8d68c3515dca5833d330b048dd8a27ccd4b137dc1d7"),
			digest: MustDecodeHex("00000000000021c9af1f188fcf1c913b2583133b0fe58bb6ff0532dad895fdfc"),
		},
		{
			height: 1888000,
			nonce:  0xdc30900000a4c493,
			hash:   MustDecodeHex("911a676cf0e5077a24e4917483bcca4bbd461a679b1a780b9d15c8b6bf5bc1d7"),
			mix:    MustDecodeHex("28605c4c11c72f1d3af0baef8cca10ccf570aa0c89c93596b2bfd485e30bd9f7"),
			digest: MustDecodeHex("0000000000000d6ea9c4b4b759a157bc2ef74bd9a7ed04ca6ac9611e218ccddc"),
		},
		{
			height: 1888509,
			nonce:  0xf09b0e1342275f3f,
			hash:   MustDecodeHex("14f2c18d74d48abe637437458c10ff5283a9a5197e8b5e740a161f4411b97a43"),
			mix:    MustDecodeHex("3dedcc6fb28c6bf3f5d203f29188bef3ff86688be34c93f28bd227eced9226e4"),
			digest: MustDecodeHex("0000000000005e6585e5e6ab7e4d75a98810204341def05823ad3a5ca1fa0d83"),
		},
	}

	dag, err := NewLightDag("RVN")
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
