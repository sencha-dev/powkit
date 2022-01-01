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

//go:generate ../.bin/gen-lookup -package kawpow -cacheInit 16777216 -cacheGrowth 131072 -datasetInit 1073741824 -datasetGrowth 8388608

package kawpow

import (
	"runtime"

	"github.com/sencha-dev/powkit/internal/common"
	"github.com/sencha-dev/powkit/internal/crypto"
	"github.com/sencha-dev/powkit/internal/dag"
)

type Kawpow struct {
	dag *dag.LightDAG
	cfg *dag.Config
}

func New(cfg *dag.Config) *Kawpow {
	client := &Kawpow{
		dag: dag.NewLightDAG(cfg),
		cfg: cfg,
	}

	return client
}

func NewRavencoin() *Kawpow {
	var cfg = &dag.Config{
		Name:       "RVN",
		Revision:   23,
		StorageDir: common.DefaultDir(".powcache"),

		DatasetInitBytes:   1 << 30,
		DatasetGrowthBytes: 1 << 23,
		CacheInitBytes:     1 << 24,
		CacheGrowthBytes:   1 << 17,

		CacheSizes:   dag.NewLookupTable(cacheSizes, 2048),
		DatasetSizes: dag.NewLookupTable(datasetSizes, 2048),

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

	return New(cfg)
}

func (e *Kawpow) Compute(height, nonce uint64, hash []byte) ([]byte, []byte) {
	epoch := e.cfg.CalcEpoch(height)
	datasetSize := e.cfg.DatasetSize(epoch)
	cache := e.dag.GetCache(epoch)

	keccak512Hasher := crypto.NewKeccak512Hasher()
	lookup := func(index uint32) []uint32 {
		return e.cfg.GenerateDatasetItem2048(cache.Cache(), index, keccak512Hasher)
	}

	mix, digest := compute(hash, height, nonce, datasetSize, lookup, cache.L1())
	runtime.KeepAlive(cache)

	return mix, digest
}
