// ethash: C/C++ implementation of Ethash, the Ethereum Proof of Work algorithm.
// Copyright 2018-2019 Pawel Bylica.
// Licensed under the Apache License, Version 2.0.

package progpow

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strings"
	"testing"

	"github.com/sencha-dev/go-pow/internal/crypto"
	"github.com/sencha-dev/go-pow/internal/dag"
)

// should only be used for tests
func mustDecodeHex(inp string) []byte {
	inp = strings.Replace(inp, "0x", "", -1)
	out, err := hex.DecodeString(inp)
	if err != nil {
		panic(err)
	}

	return out
}

func TestInitMix(t *testing.T) {
	seed := uint64(0xEE304846DDD0A47B)

	lanesExpected := map[int][32]uint32{
		0: [32]uint32{
			0x10C02F0D, 0x99891C9E, 0xC59649A0, 0x43F0394D,
			0x24D2BAE4, 0xC4E89D4C, 0x398AD25C, 0xF5C0E467,
			0x7A3302D6, 0xE6245C6C, 0x760726D3, 0x1F322EE7,
			0x85405811, 0xC2F1E765, 0xA0EB7045, 0xDA39E821,
			0x79FC6A48, 0x089E401F, 0x8488779F, 0xD79E414F,
			0x041A826B, 0x313C0D79, 0x10125A3C, 0x3F4BDFAC,
			0xA7352F36, 0x7E70CB54, 0x3B0BB37D, 0x74A3E24A,
			0xCC37236A, 0xA442B311, 0x955AB27A, 0x6D175B7E,
		},
		13: [32]uint32{
			0x4E46D05D, 0x2E77E734, 0x2C479399, 0x70712177,
			0xA75D7FF5, 0xBEF18D17, 0x8D42252E, 0x35B4FA0E,
			0x462C850A, 0x2DD2B5D5, 0x5F32B5EC, 0xED5D9EED,
			0xF9E2685E, 0x1F29DC8E, 0xA78F098B, 0x86A8687B,
			0xEA7A10E7, 0xBE732B9D, 0x4EEBCB60, 0x94DD7D97,
			0x39A425E9, 0xC0E782BF, 0xBA7B870F, 0x4823FF60,
			0xF97A5A1C, 0xB00BCAF4, 0x02D0F8C4, 0x28399214,
			0xB4CCB32D, 0x83A09132, 0x27EA8279, 0x3837DDA3,
		},
	}

	mix := initMix(seed)

	for lane := range mix {
		if expected, ok := lanesExpected[lane]; ok {
			actual := mix[lane]

			for reg := range expected {
				if actual[reg] != expected[reg] {
					t.Errorf("failed initMix test on iteration %d", lane)
					continue
				}
			}
		}
	}
}

func seedHash(block uint64, epochLength uint64) []byte {
	seed := make([]byte, 32)
	if block < epochLength {
		return seed
	}
	keccak256Hasher := crypto.NewKeccak256Hasher()
	for i := 0; i < int(block/epochLength); i++ {
		keccak256Hasher(seed, seed)
	}
	return seed
}

func TestProgpow(t *testing.T) {
	tests := []struct {
		config         *Config
		epochLength    uint64
		datasetParents uint32
		height         uint64
		nonce          uint64
		headerHash     []byte
		mixHash        []byte
		digest         []byte
	}{
		{
			config:         ProgPow093,
			epochLength:    30000,
			datasetParents: 256,
			height:         30001,
			nonce:          0x123456789abcdef0,
			headerHash:     mustDecodeHex("ffeeddccbbaa9988776655443322110000112233445566778899aabbccddeeff"),
			mixHash:        mustDecodeHex("11f19805c58ab46610ff9c719dcf0a5f18fa2f1605798eef770c47219274767d"),
			digest:         mustDecodeHex("5b7ccd472dbefdd95b895cac8ece67ff0deb5a6bd2ecc6e162383d00c3728ece"),
		},

		{
			config:         ProgPow094,
			epochLength:    30000,
			datasetParents: 512,
			height:         0,
			nonce:          0x123456789abcdef0,
			headerHash:     mustDecodeHex("ffeeddccbbaa9988776655443322110000112233445566778899aabbccddeeff"),
			mixHash:        mustDecodeHex("c2e883b6876ec4cc514b9cea269f343095619faf9f2edcafb3fcf6928fa58141"),
			digest:         mustDecodeHex("fa70fbf9979f80ec3db2c3f118a5e683fcf5f54ea7edc41b0b5d336508694cb8"),
		},
		{
			config:         ProgPow094,
			epochLength:    30000,
			datasetParents: 512,
			height:         0,
			nonce:          0x0000000000000001,
			headerHash:     mustDecodeHex("0xffeeddccbbaa9988776655443322110000112233445566778899aabbccddeeff"),
			mixHash:        mustDecodeHex("0xab94222a9736b2c93282fbf6e7217f792f87504033a83eb5501beb24f0d235e7"),
			digest:         mustDecodeHex("0x1260d102572f6ab9840556e8766ba670b511bcf768767b0ab05af45ea9fbad8d"),
		},
		{
			config:         ProgPow094,
			epochLength:    30000,
			datasetParents: 512,
			height:         49,
			nonce:          0x123456789abcdef0,
			headerHash:     mustDecodeHex("0xffeeddccbbaa9988776655443322110000112233445566778899aabbccddeeff"),
			mixHash:        mustDecodeHex("0xa0e00c15ccff10aefeeef6ca28260807fdd7f2daaff7948b15857e3a65908f09"),
			digest:         mustDecodeHex("0xa66465873e0674e95ac58efba116458342b3252abeb47874adaf139843ef79bb"),
		},
		{
			config:         ProgPow094,
			epochLength:    30000,
			datasetParents: 512,
			height:         49,
			nonce:          0x0000000000000001,
			headerHash:     mustDecodeHex("0xffeeddccbbaa9988776655443322110000112233445566778899aabbccddeeff"),
			mixHash:        mustDecodeHex("0x1704a993e5a8603615b964990253896681da83ddd10c0e6e8fee2f273fa2a961"),
			digest:         mustDecodeHex("0x528dff2f543825030a8e0943013de7bc6a4b7c203c7398607811176b03ce70f4"),
		},
		{
			config:         ProgPow094,
			epochLength:    30000,
			datasetParents: 512,
			height:         14999,
			nonce:          0x123456789abcdef0,
			headerHash:     mustDecodeHex("0xffeeddccbbaa9988776655443322110000112233445566778899aabbccddeeff"),
			mixHash:        mustDecodeHex("0xfbbed3db6316658244eef0a897a901fdb40956de9439cf15a74582427443d3bc"),
			digest:         mustDecodeHex("0xcaaa67746a4a26c102580851c4f8542f455cd97c6f2749de216c4425504d53c4"),
		},
		{
			config:         ProgPow094,
			epochLength:    30000,
			datasetParents: 512,
			height:         14999,
			nonce:          0x0000000000000001,
			headerHash:     mustDecodeHex("0xffeeddccbbaa9988776655443322110000112233445566778899aabbccddeeff"),
			mixHash:        mustDecodeHex("0x0a66b4d37962836650099ad914d2688ffb5dc8688424256cf177c3e7b3f85e88"),
			digest:         mustDecodeHex("0x4bd39ef9155cfd42f0ebb486ee7097d08f793147a9d157027db3d188770ac29d"),
		},
		{
			config:         ProgPow094,
			epochLength:    30000,
			datasetParents: 512,
			height:         30000,
			nonce:          0x123456789abcdef0,
			headerHash:     mustDecodeHex("0xffeeddccbbaa9988776655443322110000112233445566778899aabbccddeeff"),
			mixHash:        mustDecodeHex("0x493c13e9807440571511b561132834bbd558dddaa3b70c09515080a6a1aff6d0"),
			digest:         mustDecodeHex("0x46b72b75f238bea3fcfd227e0027dc173dceaa1fb71744bd3d5e030ed2fed053"),
		},
		{
			config:         ProgPow094,
			epochLength:    30000,
			datasetParents: 512,
			height:         30000,
			nonce:          0x0000000000000001,
			headerHash:     mustDecodeHex("0xffeeddccbbaa9988776655443322110000112233445566778899aabbccddeeff"),
			mixHash:        mustDecodeHex("0x67ed6cc54b89262d1ac9c3f24ad0d2362cb703c1e23881713d0350ca3035e2ae"),
			digest:         mustDecodeHex("0xad8fa791bfb2f474a3487d27075bf339d73e0e69d62ecab14166add94c8d0f92"),
		},
	}

	for i, tt := range tests {
		lightDag := &dag.LightDag{
			Name:            fmt.Sprintf("progpow-test-%d-%d", tt.epochLength, tt.datasetParents),
			EpochLength:     tt.epochLength,
			SeedEpochLength: tt.epochLength,
			DatasetParents:  tt.datasetParents,
			NumCaches:       1,
			NeedsL1:         true,
		}

		epoch := dag.CalcEpoch(tt.height, tt.epochLength)
		cache := lightDag.GetCache(epoch)

		keccak512Hasher := crypto.NewKeccak512Hasher()
		lookup := func(index uint32) []uint32 {
			return dag.GenerateDatasetItem2048(cache.Cache(), index, keccak512Hasher, tt.datasetParents)
		}

		seed := progpowSeed(tt.headerHash, tt.nonce)
		mix := HashMix(tt.config, tt.height, seed, cache.L1(), lookup)
		if bytes.Compare(tt.mixHash, mix) != 0 {
			t.Errorf("failed on %d: mix mismatch: have %s, want %s",
				i, hex.EncodeToString(mix), hex.EncodeToString(tt.mixHash))
		}
	}
}
