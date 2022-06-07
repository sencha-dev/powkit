// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package dag

import (
	"reflect"
	"testing"

	"github.com/sencha-dev/powkit/internal/common"
	"github.com/sencha-dev/powkit/internal/common/testutil"
)

func TestEpochNumber(t *testing.T) {
	tests := []struct {
		height uint64
		epoch  uint64
	}{
		{
			height: 0,
			epoch:  0,
		},
		{
			height: 1,
			epoch:  0,
		},
		{
			height: 7499,
			epoch:  0,
		},
		{
			height: 7500,
			epoch:  1,
		},
		{
			height: 7501,
			epoch:  1,
		},
		{
			height: 7502,
			epoch:  1,
		},
		{
			height: 1245000,
			epoch:  166,
		},
	}

	var d = &DAG{
		Config: Config{
			Name:       "RVN",
			Revision:   23,
			StorageDir: common.DefaultDir(".powcache"),

			DatasetInitBytes:   1 << 30,
			DatasetGrowthBytes: 1 << 23,
			CacheInitBytes:     1 << 24,
			CacheGrowthBytes:   1 << 17,

			DatasetSizes: nil,
			CacheSizes:   nil,

			DatasetParents:  512,
			EpochLength:     7500,
			SeedEpochLength: 7500,

			CacheRounds:    3,
			CachesCount:    3,
			CachesLockMmap: false,

			L1Enabled:       true,
			L1CacheSize:     4096 * 4,
			L1CacheNumItems: 4096,
		},
	}

	for i, tt := range tests {
		epoch := d.CalcEpoch(tt.height)

		if epoch != tt.epoch {
			t.Errorf("failed on %d: epoch mismatch: have %d want %d", i, epoch, tt.epoch)
		}
	}
}

func TestSeedHash(t *testing.T) {
	tests := []struct {
		epoch uint64
		seed  []byte
	}{
		{
			epoch: 0,
			seed:  testutil.MustDecodeHex("0x0000000000000000000000000000000000000000000000000000000000000000"),
		},
		{
			epoch: 1,
			seed:  testutil.MustDecodeHex("0x290decd9548b62a8d60345a988386fc84ba6bc95484008f6362f93160ef3e563"),
		},
		{
			epoch: 171,
			seed:  testutil.MustDecodeHex("0xa9b0e0c9aca72c07ba06b5bbdae8b8f69e61878301508473379bb4f71807d707"),
		},
		{
			epoch: 2048,
			seed:  testutil.MustDecodeHex("0x20a7678ca7b50829183baac2e1e3c43fa3c4bcbc171b11cf5a9f30bebd172920"),
		},
		{
			epoch: 29998,
			seed:  testutil.MustDecodeHex("0x1222b1faed7f93098f8ae498621fb3479805a664b70186063861c46596c66164"),
		},
		{
			epoch: 29999,
			seed:  testutil.MustDecodeHex("0xee1d0f61b054dff0f3025ebba821d405c8dc19a983e582e9fa5436fc3e7a07d8"),
		},
	}

	var d = &DAG{
		Config: Config{
			Name:       "RVN",
			Revision:   23,
			StorageDir: common.DefaultDir(".powcache"),

			DatasetInitBytes:   1 << 30,
			DatasetGrowthBytes: 1 << 23,
			CacheInitBytes:     1 << 24,
			CacheGrowthBytes:   1 << 17,

			DatasetSizes: nil,
			CacheSizes:   nil,

			DatasetParents:  512,
			EpochLength:     7500,
			SeedEpochLength: 7500,

			CacheRounds:    3,
			CachesCount:    3,
			CachesLockMmap: false,

			L1Enabled:       true,
			L1CacheSize:     4096 * 4,
			L1CacheNumItems: 4096,
		},
	}

	for i, tt := range tests {
		seed := d.SeedHash(uint64(tt.epoch)*d.EpochLength + 1)
		if !reflect.DeepEqual(seed, tt.seed) {
			t.Errorf("failed on %d: seed mismatch: have %x, want %x", i, seed, tt.seed)
		}
	}
}

func TestCalcCacheSize(t *testing.T) {
	tests := []struct {
		epoch uint64
		size  uint64
	}{
		{
			epoch: 0,
			size:  16776896,
		},
		{
			epoch: 1,
			size:  16907456,
		},
		{
			epoch: 2,
			size:  17039296,
		},
		{
			epoch: 3,
			size:  17170112,
		},
		{
			epoch: 4,
			size:  17301056,
		},
		{
			epoch: 5,
			size:  17432512,
		},
		{
			epoch: 6,
			size:  17563072,
		},
		{
			epoch: 7,
			size:  17693888,
		},
		{
			epoch: 8,
			size:  17824192,
		},
		{
			epoch: 9,
			size:  17955904,
		},
	}

	var d = &DAG{
		Config: Config{
			Name:       "ETH",
			Revision:   23,
			StorageDir: common.DefaultDir(".powcache"),

			DatasetInitBytes:   1 << 30,
			DatasetGrowthBytes: 1 << 23,
			CacheInitBytes:     1 << 24,
			CacheGrowthBytes:   1 << 17,

			DatasetSizes: nil,
			CacheSizes:   nil,

			DatasetParents:  256,
			EpochLength:     30000,
			SeedEpochLength: 30000,

			CacheRounds:    3,
			CachesCount:    3,
			CachesLockMmap: false,

			L1Enabled: false,
		},
	}

	for i, tt := range tests {
		size := d.calcCacheSize(tt.epoch)
		if size != tt.size {
			t.Errorf("failed on %d: size mismatch: have %d, want %d", i, size, tt.size)
		}
	}
}

func TestCalcDatasetSize(t *testing.T) {
	tests := []struct {
		epoch uint64
		size  uint64
	}{
		{
			epoch: 0,
			size:  1073739904,
		},
		{
			epoch: 1,
			size:  1082130304,
		},
		{
			epoch: 2,
			size:  1090514816,
		},
		{
			epoch: 3,
			size:  1098906752,
		},
		{
			epoch: 4,
			size:  1107293056,
		},
		{
			epoch: 5,
			size:  1115684224,
		},
		{
			epoch: 6,
			size:  1124070016,
		},
		{
			epoch: 7,
			size:  1132461952,
		},
		{
			epoch: 8,
			size:  1140849536,
		},
		{
			epoch: 9,
			size:  1149232768,
		},
	}

	var d = &DAG{
		Config: Config{
			Name:       "ETH",
			Revision:   23,
			StorageDir: common.DefaultDir(".powcache"),

			DatasetInitBytes:   1 << 30,
			DatasetGrowthBytes: 1 << 23,
			CacheInitBytes:     1 << 24,
			CacheGrowthBytes:   1 << 17,

			DatasetSizes: nil,
			CacheSizes:   nil,

			MixBytes:        128,
			DatasetParents:  256,
			EpochLength:     30000,
			SeedEpochLength: 30000,

			CacheRounds:    3,
			CachesCount:    3,
			CachesLockMmap: false,

			L1Enabled: false,
		},
	}

	for i, tt := range tests {
		size := d.calcDatasetSize(tt.epoch)
		if size != tt.size {
			t.Errorf("failed on %d: size mismatch: have %d, want %d", i, size, tt.size)
		}
	}
}
