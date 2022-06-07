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

const (
	hashBytes = 64 // Hash length in bytes
	hashWords = 16 // Number of 32 bit ints in a hash
)

type LookupTable struct {
	maxEpoch uint64
	table    []uint64
}

func NewLookupTable(table []uint64, maxEpoch uint64) *LookupTable {
	lookupTable := &LookupTable{
		maxEpoch: maxEpoch,
		table:    table,
	}

	return lookupTable
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
	MixBytes        uint64
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
