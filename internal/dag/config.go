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
	"github.com/sencha-dev/go-pow/internal/common"
)

const (
	mixBytes  = 128 // Width of mix
	hashBytes = 64  // Hash length in bytes
	hashWords = 16  // Number of 32 bit ints in a hash
)

type LookupTable struct {
	maxEpoch uint64
	table    []uint64
}

type Config struct {
	Name       string
	Revision   int
	StorageDir string

	// size variables
	DatasetInitBytes   uint64 // Bytes in dataset at genesis
	DatasetGrowthBytes uint64 // Dataset growth per epoch
	CacheInitBytes     uint64 // Bytes in cache at genesis
	CacheGrowthBytes   uint64 // Cache growth per epoch

	// lookup tables
	DatasetSizes *LookupTable
	CacheSizes   *LookupTable

	// algorithm variables
	DatasetParents  uint32
	EpochLength     uint64
	SeedEpochLength uint64 // ETC uses a different seed epoch length

	// cache variables
	CacheRounds    int
	CachesCount    int // Maximum number of caches to keep before eviction (only init, don't modify)
	CachesLockMmap bool

	// L1 variables
	L1Enabled       bool
	L1CacheSize     uint64
	L1CacheNumItems uint
}

var (
	EthereumCfg = &Config{
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
	}

	EthereumClassicCfg = &Config{
		Name:       "ETC",
		Revision:   23,
		StorageDir: common.DefaultDir(".powcache"),

		DatasetInitBytes:   1 << 30,
		DatasetGrowthBytes: 1 << 23,
		CacheInitBytes:     1 << 24,
		CacheGrowthBytes:   1 << 17,

		DatasetSizes: nil,
		CacheSizes:   nil,

		DatasetParents:  256,
		EpochLength:     60000,
		SeedEpochLength: 30000,

		CacheRounds:    3,
		CachesCount:    3,
		CachesLockMmap: false,

		L1Enabled: false,
	}

	Progpow092Cfg = &Config{
		Name:       "PROGPOW092",
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

		L1Enabled:       true,
		L1CacheSize:     4096 * 4,
		L1CacheNumItems: 4096,
	}

	Progpow093Cfg = &Config{
		Name:       "PROGPOW093",
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

		L1Enabled:       true,
		L1CacheSize:     4096 * 4,
		L1CacheNumItems: 4096,
	}

	Progpow094Cfg = &Config{
		Name:       "PROGPOW094",
		Revision:   23,
		StorageDir: common.DefaultDir(".powcache"),

		DatasetInitBytes:   1 << 30,
		DatasetGrowthBytes: 1 << 23,
		CacheInitBytes:     1 << 24,
		CacheGrowthBytes:   1 << 17,

		DatasetSizes: nil,
		CacheSizes:   nil,

		DatasetParents:  512,
		EpochLength:     30000,
		SeedEpochLength: 30000,

		CacheRounds:    3,
		CachesCount:    3,
		CachesLockMmap: false,

		L1Enabled:       true,
		L1CacheSize:     4096 * 4,
		L1CacheNumItems: 4096,
	}

	RavencoinCfg = &Config{
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
	}

	FiroCfg = &Config{
		Name:       "FIRO",
		Revision:   23,
		StorageDir: common.DefaultDir(".powcache"),

		DatasetInitBytes:   (1 << 30) + (1 << 29),
		DatasetGrowthBytes: 1 << 23,
		CacheInitBytes:     1 << 24,
		CacheGrowthBytes:   1 << 17,

		DatasetSizes: nil,
		CacheSizes:   nil,

		DatasetParents:  512,
		EpochLength:     1300,
		SeedEpochLength: 1300,

		CacheRounds:    3,
		CachesCount:    3,
		CachesLockMmap: false,

		L1Enabled:       true,
		L1CacheSize:     4096 * 4,
		L1CacheNumItems: 4096,
	}
)
