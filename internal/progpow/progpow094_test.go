// ethash: C/C++ implementation of Ethash, the Ethereum Proof of Work algorithm.
// Copyright 2018-2019 Pawel Bylica.
// Licensed under the Apache License, Version 2.0.

package progpow

import (
	"bytes"
	"testing"

	"github.com/sencha-dev/powkit/internal/common"
	"github.com/sencha-dev/powkit/internal/common/testutil"
	"github.com/sencha-dev/powkit/internal/dag"
)

func TestProgpow094(t *testing.T) {
	tests := []struct {
		height uint64
		nonce  uint64
		hash   []byte
		mix    []byte
		digest []byte
	}{
		{
			height: 0,
			nonce:  0x123456789abcdef0,
			hash:   testutil.MustDecodeHex("0xffeeddccbbaa9988776655443322110000112233445566778899aabbccddeeff"),
			mix:    testutil.MustDecodeHex("0xc2e883b6876ec4cc514b9cea269f343095619faf9f2edcafb3fcf6928fa58141"),
			digest: testutil.MustDecodeHex("0xfa70fbf9979f80ec3db2c3f118a5e683fcf5f54ea7edc41b0b5d336508694cb8"),
		},
		{
			height: 0,
			nonce:  0x0000000000000001,
			hash:   testutil.MustDecodeHex("0xffeeddccbbaa9988776655443322110000112233445566778899aabbccddeeff"),
			mix:    testutil.MustDecodeHex("0xab94222a9736b2c93282fbf6e7217f792f87504033a83eb5501beb24f0d235e7"),
			digest: testutil.MustDecodeHex("0x1260d102572f6ab9840556e8766ba670b511bcf768767b0ab05af45ea9fbad8d"),
		},
		{
			height: 49,
			nonce:  0x123456789abcdef0,
			hash:   testutil.MustDecodeHex("0xffeeddccbbaa9988776655443322110000112233445566778899aabbccddeeff"),
			mix:    testutil.MustDecodeHex("0xa0e00c15ccff10aefeeef6ca28260807fdd7f2daaff7948b15857e3a65908f09"),
			digest: testutil.MustDecodeHex("0xa66465873e0674e95ac58efba116458342b3252abeb47874adaf139843ef79bb"),
		},
		{
			height: 49,
			nonce:  0x0000000000000001,
			hash:   testutil.MustDecodeHex("0xffeeddccbbaa9988776655443322110000112233445566778899aabbccddeeff"),
			mix:    testutil.MustDecodeHex("0x1704a993e5a8603615b964990253896681da83ddd10c0e6e8fee2f273fa2a961"),
			digest: testutil.MustDecodeHex("0x528dff2f543825030a8e0943013de7bc6a4b7c203c7398607811176b03ce70f4"),
		},
		{
			height: 14999,
			nonce:  0x123456789abcdef0,
			hash:   testutil.MustDecodeHex("0xffeeddccbbaa9988776655443322110000112233445566778899aabbccddeeff"),
			mix:    testutil.MustDecodeHex("0xfbbed3db6316658244eef0a897a901fdb40956de9439cf15a74582427443d3bc"),
			digest: testutil.MustDecodeHex("0xcaaa67746a4a26c102580851c4f8542f455cd97c6f2749de216c4425504d53c4"),
		},
		{
			height: 14999,
			nonce:  0x0000000000000001,
			hash:   testutil.MustDecodeHex("0xffeeddccbbaa9988776655443322110000112233445566778899aabbccddeeff"),
			mix:    testutil.MustDecodeHex("0x0a66b4d37962836650099ad914d2688ffb5dc8688424256cf177c3e7b3f85e88"),
			digest: testutil.MustDecodeHex("0x4bd39ef9155cfd42f0ebb486ee7097d08f793147a9d157027db3d188770ac29d"),
		},
		{
			height: 30000,
			nonce:  0x123456789abcdef0,
			hash:   testutil.MustDecodeHex("0xffeeddccbbaa9988776655443322110000112233445566778899aabbccddeeff"),
			mix:    testutil.MustDecodeHex("0x493c13e9807440571511b561132834bbd558dddaa3b70c09515080a6a1aff6d0"),
			digest: testutil.MustDecodeHex("0x46b72b75f238bea3fcfd227e0027dc173dceaa1fb71744bd3d5e030ed2fed053"),
		},
		{
			height: 30000,
			nonce:  0x0000000000000001,
			hash:   testutil.MustDecodeHex("0xffeeddccbbaa9988776655443322110000112233445566778899aabbccddeeff"),
			mix:    testutil.MustDecodeHex("0x67ed6cc54b89262d1ac9c3f24ad0d2362cb703c1e23881713d0350ca3035e2ae"),
			digest: testutil.MustDecodeHex("0xad8fa791bfb2f474a3487d27075bf339d73e0e69d62ecab14166add94c8d0f92"),
		},
	}

	dagClient := dag.New(
		dag.Config{
			Name:       "PROGPOW094",
			Revision:   23,
			StorageDir: common.DefaultDir(".powcache"),

			DatasetInitBytes:   1 << 30,
			DatasetGrowthBytes: 1 << 23,
			CacheInitBytes:     1 << 24,
			CacheGrowthBytes:   1 << 17,

			DatasetSizes: nil,
			CacheSizes:   nil,

			MixBytes:        128,
			DatasetParents:  512,
			EpochLength:     30000,
			SeedEpochLength: 30000,

			CacheRounds:    3,
			CachesCount:    3,
			CachesLockMmap: false,

			L1Enabled:       true,
			L1CacheSize:     4096 * 4,
			L1CacheNumItems: 4096,
		},
	)

	for i, tt := range tests {
		epoch := dagClient.CalcEpoch(tt.height)
		datasetSize := dagClient.DatasetSize(epoch)
		cache := dagClient.GetCache(epoch)
		lookup := dagClient.NewLookupFunc2048(cache, epoch)

		mix, digest := compute(tt.hash, tt.height, tt.nonce, datasetSize, lookup, cache.L1())

		if bytes.Compare(mix, tt.mix) != 0 {
			t.Errorf("failed on %d: mix mismatch: have %x, want %x", i, mix, tt.mix)
		}

		if bytes.Compare(digest, tt.digest) != 0 {
			t.Errorf("failed on %d: digest mismatch: have %x, want %x", i, digest, tt.digest)
		}
	}
}
